# New Relic Infrastructure Integration for HAProxy 

Reports status and metrics for HAProxy service

## Requirements

HAProxy instance with the statistics page enabled and activated.

## Installation

* Download an archive file for the `HAProxy` Integration
* Extract `haproxy-definition.yml` and the `bin` directory into `/var/db/newrelic-infra/newrelic-integrations`
* Add execute permissions for the binary file `nr-haproxy` (if required)
* Extract `haproxy-config.yml.sample` into `/etc/newrelic-infra/integrations.d`

## Usage

To run the HAProxy integration, you must have the agent installed (see [agent installation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/installation/install-infrastructure-linux)).

To use the HAProxy integration, first rename `haproxy-config.yml.sample` to `haproxy-config.yml`, then configure the integration
by editing the fields in the file. 

You can view your data in Insights by creating your own NRQL queries. To do so, use the **HAProxyFrontendSample**, **HAProxyBackendSample**, **HAProxyListenerSample**, and **HAProxyServerSample** event types. 

## Compatibility

* Supported OS: No limitations
* Supported versions: HAProxy 1.3 - 1.8 

## Integration Development usage

Assuming you have the source code, you can build and run the HAProxy integration locally

* Go to the directory of the HAProxy Integration and build it
```
$ make
```

* The command above will execute tests for the HAProxy integration and build an executable file called `nr-haproxy` in the `bin` directory
```
$ ./bin/nr-haproxy --help
```

For managing external dependencies, the [govendor tool](https://github.com/kardianos/govendor) is used. It is required to lock all external dependencies to a specific version (if possible) in the vendor directory.
