package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

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

func collectMetrics(stats map[string]string, i *integration.Integration) {
	switch stats["type"] {
	case "0":
		collectMetricsOfType("frontend", HAProxyFrontendStats, stats, i)
	case "1":
		collectMetricsOfType("backend", HAProxyBackendStats, stats, i)
	case "2":
		collectMetricsOfType("server", HAProxyServerStats, stats, i)
	case "3":
		log.Error("Cannot collect listener stats")
		return
	default:
		log.Error("Invalid type %s", stats["type"])
		return
	}
}

func collectInventory(stats map[string]string, i *integration.Integration) {
	switch stats["type"] {
	case "0":
		collectInventoryOfType("frontend", stats, i)
	case "1":
		collectInventoryOfType("backend", stats, i)
	case "2":
		collectInventoryOfType("server", stats, i)
	case "3":
		log.Error("Cannot collect listener stats")
		return
	default:
		log.Error("Invalid type %s", stats["type"])
		return
	}
}

func collectInventoryOfType(entityType string, stats map[string]string, i *integration.Integration) {
	entityName, err := entityName(stats)
	if err != nil {
		log.Error("Failed to determine entity name: %s", err.Error())
		return
	}
	e, err := i.Entity(entityName, entityType)
	if err != nil {
		log.Error("Failed to create entity for %s: %s", entityName, err.Error())
		return
	}

	for metricName, metricValue := range stats {
		if metricValue == "" {
			continue
		}

		_, ok := HAProxyInventory[metricName]
		if ok {
			e.SetInventoryItem(metricName, "value", metricValue)
		}
	}
}

func collectMetricsOfType(entityType string, definitions map[string]metricDefinition, stats map[string]string, i *integration.Integration) {
	entityName, err := entityName(stats)
	if err != nil {
		log.Error("Failed to determine entity name: %s", err.Error())
		return
	}

	e, err := i.Entity(entityName, entityType)
	if err != nil {
		log.Error("Failed to create entity for %s: %s", stats["pxname"], err.Error())
		return
	}

	ms := e.NewMetricSet(fmt.Sprintf("HAProxy%sSample", strings.Title(entityType)),
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "displayName", Value: entityType + ":" + e.Metadata.Name},
	)

	for metricName, metricValue := range stats {
		if metricValue == "" {
			continue
		}

		def, ok := definitions[metricName]
		if ok {
			err := ms.SetMetric(def.MetricName, metricValue, def.SourceType)
			if err != nil {
				log.Error("Failed to set metric %s for entity %s: %s", metricName, stats["pxname"], err.Error())
			}
		}
	}
}

func entityName(metrics map[string]string) (string, error) {
	pxname, ok := metrics["pxname"]
	if !ok {
		return "", errors.New("proxy name (pxname) does not exist")
	}

	svname, ok := metrics["svname"]
	if !ok {
		return "", errors.New("service name (svname) does not exist")
	}

	return fmt.Sprintf("%s:%s", pxname, svname), nil
}
