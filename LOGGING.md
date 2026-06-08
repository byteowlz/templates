# Unified logging format

This is the shared logging contract for all byteowlz projects. Every tool — no
matter the language — produces logs that look the same to humans and parse the
same for machines. Keep implementations lean: the goal is one small, predictable
surface, not a logging framework.

## The two modes

Each tool emits logs in one of two formats, chosen automatically and overridable:

| Mode   | When                                   | For                          |
| ------ | -------------------------------------- | ---------------------------- |
| `text` | stderr is an interactive terminal      | a human reading the terminal |
| `json` | stderr is piped or redirected          | machines, aggregators, `jq`  |

**Auto-detection is the default.** A tool checks whether stderr is a TTY:
terminal → `text`, otherwise → `json`. So `mytool run` is pretty and colored,
while `mytool run 2> run.log` (or `| jq`) is JSON Lines — no flags needed.

## The JSON schema

One JSON object per line (JSON Lines / NDJSON). Field order is fixed so output
is byte-stable and diff-friendly:

```json
{"time":"2026-05-28T18:35:39Z","level":"info","msg":"running task demo with profile default"}
```

Required fields, always in this order:

| Field   | Type   | Notes                                              |
| ------- | ------ | -------------------------------------------------- |
| `time`  | string | RFC 3339 / ISO 8601                                |
| `level` | string | lowercase: `trace debug info warn error`           |
| `msg`   | string | the human-readable message                         |

Tools **may** append extra key/value fields after `msg` (e.g. `"task":"demo"`).
Consumers must ignore unknown fields. The three required fields are the contract;
everything else is per-tool.

## Levels

Five levels, ascending severity. Lowercase everywhere (config, flags, JSON):

```
trace  debug  info  warn  error
```

## Standard controls

Every tool exposes the same vocabulary:

| Control                  | Effect                                                      |
| ------------------------ | ---------------------------------------------------------- |
| `--log-format auto\|text\|json` | force a format; `auto` (default) detects the TTY    |
| `-v` / `-vv`             | raise verbosity (`debug` / `trace`)                        |
| `--debug`                | shortcut for `-vv` → `debug`                               |
| `--trace`                | force `trace` (most verbose)                               |
| `-q` / `--quiet`         | suppress all but errors                                    |
| `--no-color`             | disable ANSI color in `text` mode                          |
| config key `logging.level`  | default level                                          |
| config key `logging.format` | default format (`auto`/`text`/`json`)                 |
| config key `logging.file`   | optional file to mirror logs into                     |

**Precedence for format:** explicit `--log-format text|json` > `logging.format`
config > auto-detection. **All logs go to stderr** (stdout is reserved for the
tool's actual data output). Color applies only to `text` mode; JSON is never
colorized.

## Why these choices

- **JSON Lines, not pretty-printed JSON** — every aggregator, `jq -c`, and
  stream parser ingests it; one object per line survives truncation.
- **Auto-switch on TTY** — humans get readable output, pipelines get structured
  output, and nobody has to remember a flag.
- **`{time, level, msg}` as the floor** — the common denominator across Go's
  `slog`, Python's `logging`, and Rust's `log`/`tracing`. Each language keeps its
  native backend; only the wire format is standardized. Richer structured fields
  are allowed but optional, so no project is forced to adopt a heavier logger.
- **stderr only** — keeps machine-readable program output on stdout clean.

## Implementation status

| Template          | Status      | Backend                          |
| ----------------- | ----------- | -------------------------------- |
| `go-cli`          | done        | custom logger + JSON branch      |
| `python-cli`      | not started | stdlib `logging` + Rich + JSON formatter |
| `rust-cli`        | not started | `env_logger` `.format()` → JSON  |

When adding a new language template, conform to this document.
