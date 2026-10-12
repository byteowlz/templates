# Changelog

All notable changes to this project will be documented in this file.

## Unreleased

### Fixed

- Resolve XDG base directories per the Base Directory Specification, like the rust-cli template: use `XDG_*_HOME` only when it is an absolute path (unset, empty and relative values are ignored, without `~`/`$VAR` expansion), else `~/.config`, `~/.local/share`, `~/.local/state` and `~/.cache` on Linux and macOS (previously `~/Library/Application Support`), or `%APPDATA%`/`%LOCALAPPDATA%` on Windows. Drop the `platformdirs` dependency.

