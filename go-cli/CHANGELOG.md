# Changelog

All notable changes to this project will be documented in this file.

## Unreleased

### Fixed

- Resolve XDG base directories per the Base Directory Specification, like the rust-cli template: use `XDG_*_HOME` only when it is an absolute path (unset, empty and relative values are ignored, without `~`/`$VAR` expansion), else `~/.config`, `~/.local/share`, `~/.local/state` and `~/.cache` on Linux and macOS (previously `~/Library/Application Support`, with data in `~/Library/Caches`), or `%APPDATA%`/`%LOCALAPPDATA%` on Windows.

### Added

- Unified log output format. Logs now render as pretty colored text on a
  terminal and as JSON Lines (`{"time","level","msg"}`) when stderr is piped or
  redirected. Control with the new `--log-format auto|text|json` flag or the
  `logging.format` config key (default `auto`). See `../LOGGING.md` for the
  shared cross-language specification.

