from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime, timezone
import json
import logging
import os
from pathlib import Path
from typing import Any, Mapping

import yaml
from rich.console import Console

from .config import AppConfig, config_as_dict, default_data_dir, default_state_dir, load_config
from .logging_config import TRACE_LEVEL, configure_logging


@dataclass(slots=True)
class CommonFlags:
    config_path: Path | None = None
    quiet: bool = False
    verbose: int = 0
    debug: bool = False
    trace: bool = False
    json_output: bool = False
    yaml_output: bool = False
    no_color: bool = False
    color: str = "auto"
    dry_run: bool = False
    assume_yes: bool = False
    no_progress: bool = False
    diagnostics: bool = False
    timeout_seconds: int | None = None
    parallelism: int | None = None

    def resolve_log_level(self, config: AppConfig) -> int:
        level = config.logging.level or "info"
        level = level.lower()

        if self.trace:
            return TRACE_LEVEL
        if self.debug or self.verbose >= 2:
            return logging.DEBUG
        if self.verbose == 1 and level in {"warning", "error", "critical"}:
            return logging.INFO
        if self.quiet:
            return logging.ERROR

        return _level_from_string(level)

    def resolve_timeout(self, config: AppConfig) -> int | None:
        return self.timeout_seconds if self.timeout_seconds else config.run.timeout_seconds

    def resolve_parallelism(self, config: AppConfig) -> int | None:
        return self.parallelism if self.parallelism else config.run.parallelism

    def resolve_color_policy(self, config: AppConfig, env: Mapping[str, str]) -> str:
        if self.no_color:
            return "never"
        if self.color != "auto":
            return self.color
        if "NO_COLOR" in env:
            return "never"
        if env.get("FORCE_COLOR"):
            return "always"

        cfg_color = (config.output.color or "auto").lower()
        if cfg_color in {"auto", "always", "never"}:
            return cfg_color
        return "auto"

    def resolve_output_format(self) -> str:
        if self.json_output and self.yaml_output:
            raise ValueError("--json and --yaml cannot be used together")
        if self.json_output:
            return "json"
        if self.yaml_output:
            return "yaml"
        return "human"


@dataclass(slots=True)
class RuntimeContext:
    flags: CommonFlags
    config: AppConfig
    config_path: Path
    console: Console
    err_console: Console
    started_at: datetime
    data_dir: Path
    state_dir: Path
    output_format: str = "human"
    color_policy: str = "auto"
    env: Mapping[str, str] = field(default_factory=dict)

    def render(self, message: str, payload: Mapping[str, Any] | None = None) -> None:
        if self.output_format == "json":
            data = payload or {"message": message}
            self.console.print_json(json.dumps(data))
            return

        if self.output_format == "yaml":
            data = payload or {"message": message}
            text = yaml.safe_dump(data, sort_keys=False)
            self.console.print(text.rstrip())
            return

        if self.flags.quiet:
            if payload:
                self.console.print(payload)
            return

        self.console.print(f"[bold green]✔[/] {message}")
        if payload:
            from rich.table import Table

            table = Table(show_header=False, box=None, padding=(0, 1))
            for key, value in payload.items():
                table.add_row(f"[dim]{key}[/]", str(value))
            self.console.print(table)

    def render_config(self) -> None:
        payload = config_as_dict(self.config)
        self.render("Effective configuration", payload)


def new_runtime(flags: CommonFlags, env: Mapping[str, str] | None = None) -> RuntimeContext:
    env = env or os.environ

    config, config_path = load_config(flags.config_path, env)
    data_dir = default_data_dir(env)
    state_dir = default_state_dir(env)
    color_policy = flags.resolve_color_policy(config, env)
    output_format = flags.resolve_output_format()

    console = Console(
        stderr=False,
        no_color=color_policy == "never",
        force_terminal=color_policy == "always",
    )
    err_console = Console(
        stderr=True,
        no_color=color_policy == "never",
        force_terminal=color_policy == "always",
    )

    configure_logging(console=err_console, level=flags.resolve_log_level(config), diagnostics=flags.diagnostics)

    timeout = flags.resolve_timeout(config)
    if timeout and timeout < 0:
        raise ValueError("timeout must be positive")

    return RuntimeContext(
        flags=flags,
        config=config,
        config_path=config_path,
        console=console,
        err_console=err_console,
        started_at=datetime.now(timezone.utc),
        data_dir=data_dir,
        state_dir=state_dir,
        output_format=output_format,
        color_policy=color_policy,
        env=dict(env),
    )


def _level_from_string(level: str) -> int:
    mapping = {
        "trace": TRACE_LEVEL,
        "debug": logging.DEBUG,
        "info": logging.INFO,
        "warning": logging.WARNING,
        "warn": logging.WARNING,
        "error": logging.ERROR,
        "critical": logging.CRITICAL,
    }
    return mapping.get(level.lower(), logging.INFO)
