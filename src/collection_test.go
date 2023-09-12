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

	collectMetricsOfType("ha-frontend", HAProxyFrontendStats, from, i, "testhost")

	entityIDAttrs := integration.IDAttribute{Key: "clusterName", Value: args.HAProxyClusterName}
	e, err := i.Entity("testpx/testsv", "ha-frontend", entityIDAttrs)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 8, len(e.Metrics[0].Metrics))
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

	collectInventoryOfType("ha-frontend", from, i, "testClusterName")

	entityIDAttrs := integration.IDAttribute{Key: "clusterName", Value: args.HAProxyClusterName}
	e, err := i.Entity("testpx/testsv", "ha-frontend", entityIDAttrs)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, len(e.Inventory.Items()))
	assert.Equal(t, "4", e.Inventory.Items()["slim"]["value"])
	assert.Equal(t, nil, e.Inventory.Items()["empty"]["value"])
}

//nolint:paralleltest // integration being used to check metrics results
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
		"qtime":  "10",
		"ctime":  "11",
		"rtime":  "12",
		"ttime":  "13",
	}

	server := map[string]string{
		"pxname":         "testpx",
		"svname":         "testsv",
		"type":           "2",
		"scur":           "1",
		"qtime":          "100",
		"ctime":          "101",
		"rtime":          "102",
		"ttime":          "103",
		"agent_duration": "104",
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

	entityIDAttrs := integration.IDAttribute{Key: "clusterName", Value: args.HAProxyClusterName}
	frontendEntity, err := i.Entity("testpx/testsv", "ha-frontend", entityIDAttrs)
	if err != nil {
		t.Error(err)
	}
	backendEntity, err := i.Entity("testpx/testsv", "ha-backend", entityIDAttrs)
	if err != nil {
		t.Error(err)
	}
	serverEntity, err := i.Entity("testpx/testsv", "ha-server", entityIDAttrs)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, float64(1.0), frontendEntity.Metrics[0].Metrics["frontend.currentSessions"])

	assert.Equal(t, float64(1.0), backendEntity.Metrics[0].Metrics["backend.currentSessions"])
	assert.Equal(t, float64(0.01), backendEntity.Metrics[0].Metrics["backend.averageQueueTimeInSeconds"])
	assert.Equal(t, float64(0.011), backendEntity.Metrics[0].Metrics["backend.averageConnectTimeInSeconds"])
	assert.Equal(t, float64(0.012), backendEntity.Metrics[0].Metrics["backend.averageResponseTimeInSeconds"])
	assert.Equal(t, float64(0.013), backendEntity.Metrics[0].Metrics["backend.averageTotalSessionTimeInSeconds"])

	assert.Equal(t, float64(1.0), serverEntity.Metrics[0].Metrics["server.currentSessions"])
	assert.Equal(t, float64(0.1), serverEntity.Metrics[0].Metrics["server.averageQueueTimeInSeconds"])
	assert.Equal(t, float64(0.101), serverEntity.Metrics[0].Metrics["server.averageConnectTimeInSeconds"])
	assert.Equal(t, float64(0.102), serverEntity.Metrics[0].Metrics["server.averageResponseTimeInSeconds"])
	assert.Equal(t, float64(0.103), serverEntity.Metrics[0].Metrics["server.averageTotalSessionTimeInSeconds"])
	assert.Equal(t, float64(0.104), serverEntity.Metrics[0].Metrics["server.agentDurationSeconds"])
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

	collectInventory(frontend, i, "testClusterName")
	collectInventory(backend, i, "testClusterName")
	collectInventory(server, i, "testClusterName")
	collectInventory(listener, i, "testClusterName")
	collectInventory(invalid, i, "testClusterName")

	entityIDAttrs := integration.IDAttribute{Key: "clusterName", Value: args.HAProxyClusterName}
	frontendEntity, err := i.Entity("testpx/testsv", "ha-frontend", entityIDAttrs)
	if err != nil {
		t.Error(err)
	}
	backendEntity, err := i.Entity("testpx/testsv", "ha-backend", entityIDAttrs)
	if err != nil {
		t.Error(err)
	}
	serverEntity, err := i.Entity("testpx/testsv", "ha-server", entityIDAttrs)
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
	assert.Equal(t, "test1/test2", name)
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

	collectInventoryOfType("frontend", from, i, "testClusterName")

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
