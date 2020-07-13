package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/newrelic/infra-integrations-sdk/data/metadata"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_createStatsRequest_Error(t *testing.T) {
	_, err := createStatsRequest("user", "password", "!@#$%^&*()")
	assert.Error(t, err)
}

func Test_addCSVtoURL(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{"http://haproxy/stats", "http://haproxy/stats;csv;norefresh"},
		{"http://haproxy/haproxy?stats", "http://haproxy/haproxy?stats;csv;norefresh"},
	}

	for _, tc := range testCases {
		if addCSVtoURL(tc.input) != tc.output {
			t.Error("Failed to produce correct url with csv")
		}
	}
}

func Test_processResponseToMap(t *testing.T) {
	csv := []byte("# a,b,c,\n1,2,3,\n4,5,6,\n")
	csvReader := bytes.NewReader(csv)

	output, err := processResponseToMap(csvReader)
	if err != nil {
		t.Error(err)
	}

	expected := []map[string]string{
		{
			"a": "1",
			"b": "2",
			"c": "3",
		},
		{
			"a": "4",
			"b": "5",
			"c": "6",
		},
	}

	assert.Equal(t, expected, output)

}

func Test_processResponseToMap_Error1(t *testing.T) {
	csv := []byte("a,b,c,\n1,2,3,\n4,5,6,\n")
	csvReader := bytes.NewReader(csv)
	_, err := processResponseToMap(csvReader)
	assert.Error(t, err)
}

func Test_processResponseToMap_Error2(t *testing.T) {
	csv := []byte("! a,b,c,\n1,2,3,\n4,5,6,\n")
	csvReader := bytes.NewReader(csv)
	_, err := processResponseToMap(csvReader)
	assert.Error(t, err)
}

func Test_processResponseToMap_Error3(t *testing.T) {
	csv := []byte("# a,b,c,")
	csvReader := bytes.NewReader(csv)
	_, err := processResponseToMap(csvReader)
	assert.Error(t, err)
}

func Test_processResponseToMap_Error4(t *testing.T) {
	csv := []byte("# a,b,c,\n1,2,\n")
	csvReader := bytes.NewReader(csv)
	_, err := processResponseToMap(csvReader)
	assert.Error(t, err)
}

func Test_collectMetrics(t *testing.T) {
	from := map[string]string{
		"pxname": "testpx",
		"svname": "testsv",
		"scur":   "3",
		"smax":   "6",
		"empty":  "",
		"type":   "0",
	}

	i, _ := integration.New("test", "test")
	e, _ := i.NewEntity("testname", "HAProxyFrontend", "testname")

	now = func() time.Time {
		t, _ := time.Parse(time.ANSIC, time.ANSIC)
		return t
	}
	collectMetrics(from, e)

	expectedMetrics := metric.Metrics{
		func() metric.Metric {
			m, _ := metric.NewGauge(now(), "currentSessions", 3.0)
			return m
		}(),
		func() metric.Metric {
			m, _ := metric.NewGauge(now(), "maxSessions", 6.0)
			return m
		}(),
	}

	require.Equalf(t, expectedMetrics, e.Metrics, "Expected: %#v\nActual: %#v\n", expectedMetrics, e.Metrics)

	expectedMetadata := metadata.Map{
		"proxyName":   "testpx",
		"serviceName": "testsv",
	}

	require.Equal(t, expectedMetadata, e.Metadata.Metadata)

}

func Test_collectInventory(t *testing.T) {
	i, _ := integration.New("test", "test")

	frontend := map[string]string{
		"pxname": "testpx",
		"svname": "testsv",
		"type":   "0",
		"slim":   "1",
	}

	backend := map[string]string{
		"pxname": "testpx",
		"svname": "testsv",
		"type":   "1",
		"slim":   "1",
	}

	server := map[string]string{
		"pxname": "testpx",
		"svname": "testsv",
		"type":   "2",
		"slim":   "1",
	}

	frontendEntity, _ := i.NewEntity("testfrontend", "HAProxyFrontend", "")
	backendEntity, _ := i.NewEntity("testfrontend", "HAProxyFrontend", "")
	serverEntity, _ := i.NewEntity("testfrontend", "HAProxyFrontend", "")

	collectInventory(frontend, frontendEntity)
	collectInventory(backend, backendEntity)
	collectInventory(server, serverEntity)

	assert.Equal(t, "1", frontendEntity.Inventory.Items()["slim"]["value"])
	assert.Equal(t, "1", backendEntity.Inventory.Items()["slim"]["value"])
	assert.Equal(t, "1", serverEntity.Inventory.Items()["slim"]["value"])
}

func Test_createStatsRequest(t *testing.T) {
	req, err := createStatsRequest("testuser", "testpass", "http://haproxy/stats")
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, req)
}
