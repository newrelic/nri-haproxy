# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

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
