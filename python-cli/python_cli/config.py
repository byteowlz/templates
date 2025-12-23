from __future__ import annotations

import json
import os
from pathlib import Path
from typing import Any, Mapping, MutableMapping

from dataclasses import dataclass, field

try:  # Python 3.11+
    import tomllib as tomli  # type: ignore[no-redef]
except ModuleNotFoundError:  # pragma: no cover
    import tomli  # type: ignore[no-redef]

import tomli_w
from platformdirs import PlatformDirs

APP_NAME = "python-cli"
ENV_PREFIX = "PYTHON_CLI__"


@dataclass(slots=True)
class LoggingConfig:
    level: str = "info"
    file: Path | None = None
    json: bool = False


@dataclass(slots=True)
class OutputConfig:
    color: str = "auto"
    no_progress: bool = False


@dataclass(slots=True)
class RunConfig:
    default_task: str = "default"
    profile: str = "default"
    timeout_seconds: int | None = None
    parallelism: int | None = None


@dataclass(slots=True)
class AppConfig:
    logging: LoggingConfig = field(default_factory=LoggingConfig)
    output: OutputConfig = field(default_factory=OutputConfig)
    run: RunConfig = field(default_factory=RunConfig)


DEFAULT_CONFIG_TOML = """\
# Default configuration for python-cli.
# Copy this file to {config_path} to customize.

[logging]
# Supported levels: trace, debug, info, warning, error, critical
level = "info"
file = ""
json = false

[output]
# Color policies: auto, always, never
color = "auto"
no_progress = false

[run]
default_task = "default"
profile = "default"
# timeout_seconds = 60
# parallelism = 4
"""


def default_dirs() -> PlatformDirs:
    return PlatformDirs(appname=APP_NAME, appauthor=False)


def default_config_path(env: Mapping[str, str] | None = None) -> Path:
    env = env or os.environ
    base = env.get("XDG_CONFIG_HOME")
    if base:
        return _expand_path(base, env) / APP_NAME / "config.toml"
    return Path(default_dirs().user_config_dir) / "config.toml"


def default_data_dir(env: Mapping[str, str] | None = None) -> Path:
    env = env or os.environ
    base = env.get("XDG_DATA_HOME")
    if base:
        return _expand_path(base, env) / APP_NAME
    return Path(default_dirs().user_data_dir)


def default_state_dir(env: Mapping[str, str] | None = None) -> Path:
    env = env or os.environ
    base = env.get("XDG_STATE_HOME")
    if base:
        return _expand_path(base, env) / APP_NAME
    return Path(default_dirs().user_state_dir)


def load_config(path: Path | None, env: Mapping[str, str]) -> tuple[AppConfig, Path]:
    resolved_path = path if path is not None else default_config_path(env)

    merged: dict[str, Any] = _defaults_dict()
    if resolved_path.exists():
        data = tomli.loads(resolved_path.read_text())
        _deep_update(merged, data)

    env_overrides = _environment_overrides(env)
    if env_overrides:
        _deep_update(merged, env_overrides)

    config = _decode_app_config(merged)
    return config, resolved_path


def ensure_default_config(
    path: Path | None = None,
    *,
    force: bool = False,
    env: Mapping[str, str] | None = None,
) -> Path:
    env = env or os.environ
    target = path if path is not None else default_config_path(env)
    target.parent.mkdir(parents=True, exist_ok=True)

    if target.exists() and not force:
        return target

    content = DEFAULT_CONFIG_TOML.replace("{config_path}", str(target))
    target.write_text(content.strip() + "\n")
    return target


def config_as_dict(config: AppConfig) -> dict[str, Any]:
    return {
        "logging": {
            "level": config.logging.level,
            "file": str(config.logging.file) if config.logging.file else "",
            "json": config.logging.json,
        },
        "output": {
            "color": config.output.color,
            "no_progress": config.output.no_progress,
        },
        "run": {
            "default_task": config.run.default_task,
            "profile": config.run.profile,
            "timeout_seconds": config.run.timeout_seconds,
            "parallelism": config.run.parallelism,
        },
    }


def dump_config(config: AppConfig) -> str:
    data = config_as_dict(config)
    return tomli_w.dumps(data)


def _defaults_dict() -> dict[str, Any]:
    return {
        "logging": {
            "level": "info",
            "file": "",
            "json": False,
        },
        "output": {
            "color": "auto",
            "no_progress": False,
        },
        "run": {
            "default_task": "default",
            "profile": "default",
            "timeout_seconds": None,
            "parallelism": None,
        },
    }


def _deep_update(base: MutableMapping[str, Any], updates: Mapping[str, Any]) -> None:
    for key, value in updates.items():
        if (
            isinstance(value, Mapping)
            and key in base
            and isinstance(base[key], MutableMapping)
        ):
            _deep_update(base[key], value)
        else:
            base[key] = value


def _environment_overrides(env: Mapping[str, str]) -> dict[str, Any]:
    overrides: dict[str, Any] = {}
    prefix = ENV_PREFIX
    for key, raw in env.items():
        if not key.startswith(prefix):
            continue
        parts = key[len(prefix) :].strip("_").split("__")
        if not parts:
            continue
        target = overrides
        for segment in parts[:-1]:
            segment_key = segment.lower()
            target = target.setdefault(segment_key, {})  # type: ignore[assignment]
        target[parts[-1].lower()] = _parse_env_value(raw)
    return overrides


def _parse_env_value(raw: str) -> Any:
    lowered = raw.lower()
    if lowered in {"true", "false"}:
        return lowered == "true"
    if lowered in {"null", "none"}:
        return None
    try:
        return int(raw)
    except ValueError:
        pass
    try:
        return float(raw)
    except ValueError:
        pass
    try:
        return json.loads(raw)
    except (json.JSONDecodeError, TypeError):
        return raw


def _decode_app_config(data: Mapping[str, Any]) -> AppConfig:
    logging_data = data.get("logging", {})
    output_data = data.get("output", {})
    run_data = data.get("run", {})

    logging = LoggingConfig(
        level=str(logging_data.get("level", "info")).lower(),
        file=_coerce_path(logging_data.get("file")),
        json=bool(logging_data.get("json", False)),
    )
    output = OutputConfig(
        color=str(output_data.get("color", "auto")).lower(),
        no_progress=bool(output_data.get("no_progress", False)),
    )
    run = RunConfig(
        default_task=str(run_data.get("default_task", "default")),
        profile=str(run_data.get("profile", "default")),
        timeout_seconds=_coerce_optional_int(run_data.get("timeout_seconds")),
        parallelism=_coerce_optional_int(run_data.get("parallelism")),
    )
    return AppConfig(logging=logging, output=output, run=run)


def _coerce_path(value: Any) -> Path | None:
    if value in ("", None):
        return None
    return Path(os.path.expanduser(str(value)))


def _coerce_optional_int(value: Any) -> int | None:
    if value in ("", None):
        return None
    try:
        return int(value)
    except (TypeError, ValueError):
        return None


def _expand_path(raw: str, env: Mapping[str, str]) -> Path:
    expanded = os.path.expanduser(raw)
    for key, value in env.items():
        expanded = expanded.replace(f"${key}", value).replace(f"${{{key}}}", value)
    return Path(expanded)
