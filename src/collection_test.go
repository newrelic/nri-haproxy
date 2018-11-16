package main

import (
	"bytes"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/stretchr/testify/assert"
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

func Test_collectMetricsOfType(t *testing.T) {
	from := map[string]string{
		"pxname": "testpx",
		"svname": "testsv",
		"scur":   "3",
		"empty":  "",
	}

	i, _ := integration.New("test", "test")

	collectMetricsOfType("frontend", HAProxyFrontendStats, from, i, "testhost")

	e, err := i.Entity("testpx:testsv", "frontend")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 7, len(e.Metrics[0].Metrics))
	assert.Equal(t, float64(3.0), e.Metrics[0].Metrics["frontend.currentSessions"])
	assert.Equal(t, nil, e.Metrics[0].Metrics["empty"])
}

func Test_collectInventoryOfType(t *testing.T) {
	from := map[string]string{
		"pxname": "testpx",
		"svname": "testsv",
		"slim":   "4",
		"empty":  "",
	}

	i, _ := integration.New("test", "test")

	collectInventoryOfType("frontend", from, i)

	e, err := i.Entity("testpx:testsv", "frontend")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, len(e.Inventory.Items()))
	assert.Equal(t, "4", e.Inventory.Items()["slim"]["value"])
	assert.Equal(t, nil, e.Inventory.Items()["empty"]["value"])
}

func Test_collectMetrics(t *testing.T) {
	i, _ := integration.New("test", "test")

	frontend := map[string]string{
		"pxname": "testpx",
		"svname": "testsv",
		"type":   "0",
		"scur":   "1",
	}

	backend := map[string]string{
		"pxname": "testpx",
		"svname": "testsv",
		"type":   "1",
		"scur":   "1",
	}

	server := map[string]string{
		"pxname": "testpx",
		"svname": "testsv",
		"type":   "2",
		"scur":   "1",
	}

	invalid := map[string]string{
		"pxname": "testpx",
		"svname": "testsv",
		"type":   "4",
		"scur":   "1",
	}

	collectMetrics(frontend, i, "testhost")
	collectMetrics(backend, i, "testhost")
	collectMetrics(server, i, "testhost")
	collectMetrics(invalid, i, "testhost")

	frontendEntity, err := i.Entity("testpx:testsv", "frontend")
	if err != nil {
		t.Error(err)
	}
	backendEntity, err := i.Entity("testpx:testsv", "backend")
	if err != nil {
		t.Error(err)
	}
	serverEntity, err := i.Entity("testpx:testsv", "server")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, float64(1.0), frontendEntity.Metrics[0].Metrics["frontend.currentSessions"])
	assert.Equal(t, float64(1.0), backendEntity.Metrics[0].Metrics["backend.currentSessions"])
	assert.Equal(t, float64(1.0), serverEntity.Metrics[0].Metrics["server.currentSessions"])
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

	listener := map[string]string{
		"pxname": "testpx",
		"svname": "testsv",
		"type":   "3",
		"slim":   "1",
	}

	invalid := map[string]string{
		"pxname": "testpx",
		"svname": "testsv",
		"type":   "4",
		"slim":   "1",
	}

	collectInventory(frontend, i)
	collectInventory(backend, i)
	collectInventory(server, i)
	collectInventory(listener, i)
	collectInventory(invalid, i)

	frontendEntity, err := i.Entity("testpx:testsv", "frontend")
	if err != nil {
		t.Error(err)
	}
	backendEntity, err := i.Entity("testpx:testsv", "backend")
	if err != nil {
		t.Error(err)
	}
	serverEntity, err := i.Entity("testpx:testsv", "server")
	if err != nil {
		t.Error(err)
	}

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

func Test_entityName(t *testing.T) {
	in := map[string]string{
		"pxname": "test1",
		"svname": "test2",
	}

	name, err := entityName(in)
	assert.Nil(t, err)
	assert.Equal(t, "test1:test2", name)
}

func Test_entityName_Error(t *testing.T) {
	in1 := map[string]string{
		"pxname": "test1",
	}

	_, err := entityName(in1)
	assert.Error(t, err)

	in2 := map[string]string{
		"svname": "test1",
	}

	_, err = entityName(in2)
	assert.Error(t, err)
}

func Test_collectInventoryOfType_Error(t *testing.T) {
	i, _ := integration.New("test", "test")

	from := map[string]string{
		"pxname": "test",
	}

	collectInventoryOfType("frontend", from, i)

	assert.Equal(t, 0, len(i.Entities))
}

func Test_collectMetricsOfType_Error(t *testing.T) {
	i, _ := integration.New("test", "test")

	from := map[string]string{
		"pxname": "test",
	}

	collectMetricsOfType("frontend", HAProxyFrontendStats, from, i, "testhost")

	assert.Equal(t, 0, len(i.Entities))
}
