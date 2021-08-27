# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

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
