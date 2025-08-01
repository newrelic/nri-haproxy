# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

Unreleased section should follow [Release Toolkit](https://github.com/newrelic/release-toolkit#render-markdown-and-update-markdown)
## Unreleased

## v3.2.1 - 2025-07-01

### ⛓️ Dependencies
- Updated golang version to v1.24.4

## v3.2.0 - 2025-03-12

### 🚀 Enhancements
- Add FIPS compliant packages

## v3.1.2 - 2025-03-11

### ⛓️ Dependencies
- Updated golang patch version to v1.23.6

## v3.1.1 - 2025-01-21

### ⛓️ Dependencies
- Updated golang patch version to v1.23.5

## v3.1.0 - 2024-10-14

### dependency
- Upgrade go to 1.23.2

### 🚀 Enhancements
- Upgrade integrations SDK so the interval is variable and allows intervals up to 5 minutes

## v3.0.3 - 2024-09-10

### ⛓️ Dependencies
- Updated golang version to v1.23.1

## v3.0.2 - 2024-07-09

### ⛓️ Dependencies
- Updated golang version to v1.22.5

## v3.0.1 - 2024-04-16

### ⛓️ Dependencies
- Updated golang version to v1.22.2

## v3.0.0 - 2024-02-27

### ⚠️️ Breaking changes ⚠️
- `*InSeconds` metrics are now reported in seconds. The affected metrics used to be reported in milliseconds. Affected metrics: `backend.averageQueueTimeInSeconds`, `backend.averageResponseTimeInSeconds`, `backend.averageTotalSessionTimeInSeconds`, `server.agentDurationSeconds`, `server.averageConnectTimeInSeconds`, `server.averageQueueTimeInSeconds`, `server.averageResponseTimeInSeconds`, `server.averageTotalSessionTimeInSeconds`.

### ⛓️ Dependencies
- Updated github.com/newrelic/infra-integrations-sdk to v3.8.2+incompatible

## v2.5.1 - 2023-10-31

### ⛓️ Dependencies
- Updated golang version to 1.21

## 2.5.0 (2023-06-06)
### Changed
- Update Go version to 1.20

## 2.4.0  (2023-03-08)
### Changed
- Upgrade Go to 1.19 and bump dependencies

## 2.3.3  (2022-06-29)
### Changed
- Bump dependencies
### Added
Added support for more distributions:
- RHEL(EL) 9
- Ubuntu 22.04

## 2.3.2 (2021-10-20)
### Added
Added support for more distributions:
- Debian 11
- Ubuntu 20.10
- Ubuntu 21.04
- SUSE 12.15
- SUSE 15.1
- SUSE 15.2
- SUSE 15.3
- Oracle Linux 7
- Oracle Linux 8

## 2.3.1
Added infra agent dependency

## 2.3.0 (2021-08-27)
### Added

Moved default config.sample to [V4](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-newer-configuration-format/), added a dependency for infra-agent version 1.20.0

Please notice that old [V3](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-standard-configuration-format/) configuration format is deprecated, but still supported.

## 2.2.2 (2021-06-10)
### Changed
- Add ARM support

## 2.2.1 (2021-04-14)
### Changed
- Upgraded github.com/newrelic/infra-integrations-sdk to v3.6.7
- Switched to go modules
- Upgraded pipeline to go 1.16
- Replaced gometalinter with golangci-lint

## 2.2.0 (2021-03-03)
### Changed
- Decorate metrics with `haproxyClusterName`
  `clusterName` (deprecated) attribute might get overwritten by the Kubernetes integration when running inside kubernetes.
  This change allows HAProxy entities to be uniquely identified when multiple HAProxy clusters are located in the same Kubernetes cluster.

## 2.1.2 (2020-11-06)
### Fixed
- Print error when CSV parsing fails

## 2.1.1 (2020-09-25)
### Fixed
- Updated the SDK to fix bug with calculating positive rates

## 2.1.0 (2019-11-18)
### Changed
- Renamed the integration executable from nr-haproxy to nri-haproxy in order to be consistent with the package naming. **Important Note:** if you have any security module rules (eg. SELinux), alerts or automation that depends on the name of this binary, these will have to be updated.

## 2.0.3 - 2019-10-23
### Added
- Added resources for windows MSI packaging

## 2.0.2 - 2019-07-29
### Changed
- Updated sample config with cluster_name

## 2.0.1 - 2019-07-23
- Tarball release

## 2.0.0 - 2019-04-22
### Changed
- Added namespace prefixes for more unique entity keys
- Added cluster name argument and made it required

## 1.0.1 - 2019-03-19
### Changed
- GA Release

## 0.1.0 - 2018-11-12
### Added
- Initial version: Includes Metrics and Inventory data
