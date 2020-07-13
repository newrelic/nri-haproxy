package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

var now = time.Now

func addCSVtoURL(statsURL string) string {
	return statsURL + `;csv;norefresh`
}

func createStatsRequest(username, password, url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)

	return req, nil
}

func processResponseToMap(body io.Reader) ([]map[string]string, error) {
	bufReader := bufio.NewReader(body)

	// Trim past the first space to remove the '# '
	comment, err := bufReader.ReadBytes(' ')
	if err != nil {
		return nil, err
	}

	// Ensure the result is well-formed
	if string(comment) != "# " {
		return nil, errors.New("cannot decode statistics: malformed response CSV")
	}

	// Read the metric names line
	metricNames, err := bufReader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	metricNamesSplit := strings.Split(string(metricNames), ",")
	metricNamesSplit = metricNamesSplit[0 : len(metricNamesSplit)-1] // Remove the newline

	// Create the CSV reader with strict entry length requirements
	r := csv.NewReader(bufReader)
	r.FieldsPerRecord = len(metricNamesSplit) + 1

	// For each line, parse results into a map from metric name to value
	maps := make([]map[string]string, 0, 10)
	for {
		recordMap := make(map[string]string)

		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to parse CSV line: %s", err.Error())
		}

		// For each except the last entry (trailing comma)
		for index, stat := range record[0 : len(record)-1] {
			recordMap[metricNamesSplit[index]] = stat
		}

		maps = append(maps, recordMap)
	}

	return maps, nil
}

func createEntity(stats map[string]string, i *integration.Integration, url string) (*integration.Entity, error) {
	entityType, err := entityType(stats)
	if err != nil {
		return nil, err
	}

	entityName, displayName, err := entityName(entityType, stats, url)
	if err != nil {
		return nil, err
	}

	return i.NewEntity(entityName, entityType, displayName)
}

func collectInventory(stats map[string]string, e *integration.Entity) {
	for metricName, metricValue := range stats {
		if metricValue == "" {
			continue
		}

		_, ok := HAProxyInventory[metricName]
		if ok {
			err := e.AddInventoryItem(metricName, "value", metricValue)
			if err != nil {
				log.Error("Failed to set inventory item for %s", metricName)
			}
		}
	}
}

func collectMetrics(stats map[string]string, e *integration.Entity) {
	switch stats["type"] {
	case "0":
		collectMetricsFromDefinition(HAProxyFrontendStats, stats, e)
	case "1":
		collectMetricsFromDefinition(HAProxyBackendStats, stats, e)
	case "2":
		collectMetricsFromDefinition(HAProxyServerStats, stats, e)
	case "3":
		log.Error("Cannot collect listener stats")
		return
	default:
		log.Error("Invalid type %s", stats["type"])
		return
	}
}

func collectMetricsFromDefinition(definitions map[string]metricDefinition, stats map[string]string, e *integration.Entity) {
	for metricName, metricValue := range stats {
		if metricValue == "" {
			continue
		}

		def, ok := definitions[metricName]
		if ok {
			if def.IsAttribute {
				err := e.AddMetadata(def.MetricName, metricValue)
				if err != nil {
					log.Error("Failed to add metadata to entity: %s", err)
				}
				continue
			}

			currentTime := now()
			var newMetric metric.Metric
			switch def.SourceType {
			case metric.GAUGE:
				val, err := strconv.ParseFloat(metricValue, 64)
				if err != nil {
					log.Error("Failed to parse metric as float: %s", err)
					continue
				}
				newMetric, err = integration.Gauge(currentTime, def.MetricName, val)
				if err != nil {
					log.Error("Failed to create metric: %s", err)
					continue
				}
			case metric.COUNT:
				val, err := strconv.ParseFloat(metricValue, 64)
				if err != nil {
					log.Error("Failed to parse metric as float: %s", err)
					continue
				}
				newMetric, err = integration.Count(currentTime, def.MetricName, val)
				if err != nil {
					log.Error("Failed to create metric: %s", err)
					continue
				}
			case metric.CUMULATIVE_COUNT:
				val, err := strconv.ParseFloat(metricValue, 64)
				if err != nil {
					log.Error("Failed to parse metric as float: %s", err)
					continue
				}
				newMetric, err = integration.CumulativeCount(currentTime, def.MetricName, val)
				if err != nil {
					log.Error("Failed to create metric: %s", err)
					continue
				}
			case metric.CUMULATIVE_RATE:
				val, err := strconv.ParseFloat(metricValue, 64)
				if err != nil {
					log.Error("Failed to parse metric as float: %s", err)
					continue
				}
				newMetric, err = integration.CumulativeRate(currentTime, def.MetricName, val)
				if err != nil {
					log.Error("Failed to create metric: %s", err)
					continue
				}
			case metric.RATE:
				val, err := strconv.ParseFloat(metricValue, 64)
				if err != nil {
					log.Error("Failed to parse metric as float: %s", err)
					continue
				}
				newMetric, err = integration.Rate(currentTime, def.MetricName, val)
				if err != nil {
					log.Error("Failed to create metric: %s", err)
					continue
				}
			}

			e.AddMetric(newMetric)
		}
	}
}

func entityName(entityType string, stats map[string]string, url string) (string, string, error) {
	pxname, ok := stats["pxname"]
	if !ok {
		return "", "", errors.New("proxy name (pxname) does not exist")
	}

	svname, ok := stats["svname"]
	if !ok {
		return "", "", errors.New("service name (svname) does not exist")
	}

	return fmt.Sprintf("haproxy|%s|%s|%s|%s", url, entityType, pxname, svname), fmt.Sprintf("%s:%s:%s", entityType, pxname, svname), nil
}

func entityType(stats map[string]string) (string, error) {
	switch stats["type"] {
	case "0":
		return "HAProxyFrontend", nil
	case "1":
		return "HAProxyBackend", nil
	case "2":
		return "HAProxyServer", nil
	case "3":
		return "", errors.New("cannot collect metrics for HAProxyListener")
	default:
		return "", fmt.Errorf("unknown type %s", stats["type"])
	}
}
