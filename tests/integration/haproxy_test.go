// +build integration

package integration

import (
	"flag"
	"github.com/newrelic/infra-integrations-sdk/v3/log"
	"github.com/newrelic/nri-haproxy/tests/integration/helpers"
	"github.com/newrelic/nri-haproxy/tests/integration/jsonschema"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	iName = "haproxy"

	defaultContainer = "integration_nri-haproxy_1"

	defaultBinPath            = "/nri-haproxy"
	defaultStatsURL           = "http://haproxy:8404/stats"
	defaultHAProxyClusterName = "haproxy-cluster"

	// cli flags
	container = flag.String("container", defaultContainer, "container where the integration is installed")
	binPath   = flag.String("bin", defaultBinPath, "Integration binary path")

	statsURL           = flag.String("stats_url", defaultStatsURL, "haproxy stats url")
	haProxyClusterName = flag.String("ha_proxy_cluster_name", defaultHAProxyClusterName, "haproxy cluster name")
)

// Returns the standard output, or fails testing if the command returned an error
func runIntegration(t *testing.T, envVars ...string) (string, string, error) {
	t.Helper()

	command := make([]string, 0)
	command = append(command, *binPath)

	var (
		hasEnvStatsURL           bool
		hasEnvHAProxyClusterName bool
	)

	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "STATS_URL") {
			hasEnvStatsURL = true
		}
		if strings.HasPrefix(envVar, "HA_PROXY_CLUSTER_NAME") {
			hasEnvHAProxyClusterName = true
		}
	}

	if !hasEnvStatsURL && statsURL != nil {
		command = append(command, "--stats_url", *statsURL)
	}
	if !hasEnvHAProxyClusterName && haProxyClusterName != nil {
		command = append(command, "--ha_proxy_cluster_name", *haProxyClusterName)
	}

	stdout, stderr, err := helpers.ExecInContainer(*container, command, envVars...)

	if stderr != "" {
		log.Debug("Integration command Standard Error: ", stderr)
	}

	return stdout, stderr, err
}

func TestMain(m *testing.M) {
	flag.Parse()

	result := m.Run()
	os.Exit(result)
}

func TestHAProxyIntegration(t *testing.T) {
	stdout, stderr, err := runIntegration(t)

	assert.NotNil(t, stderr, "unexpected stderr")
	assert.NoError(t, err, "Unexpected error")

	schemaPath := filepath.Join("json-schema-files", "haproxy-schema.json")
	err = jsonschema.Validate(schemaPath, stdout)
	assert.NoError(t, err, "The output of HAProxy integration doesn't have expected format.")
}

func TestHAProxyIntegrationOnlyMetrics(t *testing.T) {
	stdout, stderr, err := runIntegration(t, "METRICS=true")

	assert.NotNil(t, stderr, "unexpected stderr")
	assert.NoError(t, err, "Unexpected error")

	schemaPath := filepath.Join("json-schema-files", "haproxy-schema-metrics.json")

	err = jsonschema.Validate(schemaPath, stdout)
	assert.NoError(t, err, "The output of HAProxy integration doesn't have expected format.")
}

func TestHAProxyIntegrationOnlyInventory(t *testing.T) {
	stdout, stderr, err := runIntegration(t, "INVENTORY=true")

	assert.NotNil(t, stderr, "unexpected stderr")
	assert.NoError(t, err, "Unexpected error")

	schemaPath := filepath.Join("json-schema-files", "haproxy-schema-inventory.json")

	err = jsonschema.Validate(schemaPath, stdout)
	assert.NoError(t, err, "The output of HAProxy integration doesn't have expected format.")
}

func TestHAProxyIntegrationInvalidStatsURL_EmptyStatsURL(t *testing.T) {
	stdout, stderr, err := runIntegration(t, "STATS_URL=")

	expectedErrorMessage := "Must supply a URL pointing to the HAProxy stats page"

	errMatch, _ := regexp.MatchString(expectedErrorMessage, stderr)
	assert.Error(t, err, "Expected error")
	assert.Truef(t, errMatch, "Expected error message: '%s', got: '%s'", expectedErrorMessage, stderr)

	assert.NotNil(t, stdout, "unexpected stdout")
}

func TestHAProxyIntegrationInvalidStatsURL(t *testing.T) {
	stdout, stderr, err := runIntegration(t, "STATS_URL=http://localhost/", "VERBOSE=true")

	expectedErrorMessage := "connection refused"

	errMatch, _ := regexp.MatchString(expectedErrorMessage, stderr)
	assert.Error(t, err, "Expected error")
	assert.Truef(t, errMatch, "Expected error message: '%s', got: '%s'", expectedErrorMessage, stderr)

	assert.NotNil(t, stdout, "unexpected stdout")
}

func TestHAProxyIntegrationInvalidStatsURL_NoExistingHost(t *testing.T) {
	stdout, stderr, err := runIntegration(t, "STATS_URL=http://nonExistingHost/stats")

	expectedErrorMessage := "no such host"

	errMatch, _ := regexp.MatchString(expectedErrorMessage, stderr)
	assert.Error(t, err, "Expected error")
	assert.Truef(t, errMatch, "Expected error message: '%s', got: '%s'", expectedErrorMessage, stderr)

	assert.NotNil(t, stdout, "unexpected stdout")
}

func TestHAProxyIntegrationNotValidURL_NoHttp(t *testing.T) {
	stdout, stderr, err := runIntegration(t, "STATS_URL=localhost/stats")

	expectedErrorMessage := "unsupported protocol scheme"

	errMatch, _ := regexp.MatchString(expectedErrorMessage, stderr)
	assert.Error(t, err, "Expected error")
	assert.Truef(t, errMatch, "Expected error message: '%s', got: '%s'", expectedErrorMessage, stderr)

	assert.NotNil(t, stdout, "unexpected stdout")
}

func TestHAProxyIntegrationNotValidURL_OnlyHttp(t *testing.T) {
	stdout, stderr, err := runIntegration(t, "STATS_URL=http://")

	expectedErrorMessage := "no such host"

	errMatch, _ := regexp.MatchString(expectedErrorMessage, stderr)
	assert.Error(t, err, "Expected error")
	assert.Truef(t, errMatch, "Expected error message: '%s', got: '%s'", expectedErrorMessage, stderr)

	assert.NotNil(t, stdout, "unexpected stdout")
}

func TestHAProxyIntegrationEmptyHAProxyClusterName(t *testing.T) {
	stdout, stderr, err := runIntegration(t, "HA_PROXY_CLUSTER_NAME=")

	expectedErrorMessage := "Must supply a cluster name to identify this HAProxy instance"

	errMatch, _ := regexp.MatchString(expectedErrorMessage, stderr)
	assert.Error(t, err, "Expected error")
	assert.Truef(t, errMatch, "Expected error message: '%s', got: '%s'", expectedErrorMessage, stderr)

	assert.NotNil(t, stdout, "unexpected stdout")
}
