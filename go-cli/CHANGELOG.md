# Changelog

All notable changes to this project will be documented in this file.

## Unreleased

### Added

- Unified log output format. Logs now render as pretty colored text on a
  terminal and as JSON Lines (`{"time","level","msg"}`) when stderr is piped or
  redirected. Control with the new `--log-format auto|text|json` flag or the
  `logging.format` config key (default `auto`). See `../LOGGING.md` for the
  shared cross-language specification.

