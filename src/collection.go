package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"net/http"
	"strings"
  "errors"
  "fmt"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
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
  if string( comment ) != "# " {
    return nil, errors.New("cannot decode statistics: malformed response CSV")
  }

  // Read the metric names line
	metricNames, err := bufReader.ReadBytes('\n') 
	if err != nil {
		return nil, err
	}
	metricNamesSplit := strings.Split(string(metricNames), ",")

  // Create the CSV reader with strict entry length requirements
	r := csv.NewReader(bufReader)
	r.FieldsPerRecord = len(metricNamesSplit)

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

		for index, stat := range record {
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
		collectMetricsOfType("listener", HAProxyListenerStats, stats, i)
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
		collectInventoryOfType("listener", stats, i)
	default:
		log.Error("Invalid type %s", stats["type"])
    return
	}
}

func collectInventoryOfType(entityType string, from map[string]string, i *integration.Integration) {
	e, err := i.Entity(from["pxname"], entityType)
	if err != nil {
		log.Error("Failed to create entity for %s: %s", from["pxname"], err.Error())
		return
	}

	for metricName, metricValue := range from {
		if metricValue == "" {
			continue
		}

		_, ok := HAProxyInventory[metricName]
		if ok {
			e.SetInventoryItem(metricName, "value", metricValue)
		}
	}
}

// TODO naming
func collectMetricsOfType(entityType string, collect map[string]metricDefinition, from map[string]string, i *integration.Integration) {
  e, err := i.Entity(from["pxname"]+":"+from["svname"], entityType)
	if err != nil {
		log.Error("Failed to create entity for %s: %s", from["pxname"], err.Error())
		return
	}

	ms := e.NewMetricSet(fmt.Sprintf("HAProxy%sSample", strings.Title(entityType)),
    metric.Attribute{Key:"displayName", Value:e.Metadata.Name},
    metric.Attribute{Key:"displayName", Value: entityType +":"+ e.Metadata.Name},
  )

	for metricName, metricValue := range from {
		if metricValue == "" {
			continue
		}

		def, ok := collect[metricName]
		if ok {
      err := ms.SetMetric(def.MetricName, metricValue, def.SourceType)
      if err != nil {
        log.Error("Failed to set metric %s for entity %s: %s", metricName, from["pxname"], err.Error())
      }
		}
	}
}
