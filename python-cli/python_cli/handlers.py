from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
import time
from typing import Any

import typer

from .config import (
    config_as_dict,
    default_cache_dir,
    ensure_default_config,
    load_config,
)
from .runtime import RuntimeContext


@dataclass(slots=True)
class InitOptions:
    force: bool = False


@dataclass(slots=True)
class RunOptions:
    task: str = "default"
    profile: str | None = None


def handle_init(ctx: RuntimeContext, opts: InitOptions) -> None:
    config_path = ctx.config_path
    if config_path.exists() and not (opts.force or ctx.flags.assume_yes):
        confirmed = typer.confirm(
            f"{config_path} already exists. Overwrite with defaults?", default=False
        )
        if not confirmed:
            ctx.err_console.print("[yellow]Skipped[/] existing configuration.")
            return

    ensure_default_config(config_path, force=True, env=ctx.env)
    ctx.data_dir.mkdir(parents=True, exist_ok=True)
    ctx.state_dir.mkdir(parents=True, exist_ok=True)

    ctx.render(
        "Initialized configuration and data directories",
        {
            "config_path": str(config_path),
            "data_dir": str(ctx.data_dir),
            "state_dir": str(ctx.state_dir),
        },
    )


def handle_config_show(ctx: RuntimeContext) -> None:
    ctx.render("Effective configuration", config_as_dict(ctx.config))


def handle_config_path(ctx: RuntimeContext) -> None:
    if ctx.output_format in {"json", "yaml"}:
        ctx.render("Config path", {"config_path": str(ctx.config_path)})
    else:
        ctx.console.print(str(ctx.config_path))


def handle_config_reset(ctx: RuntimeContext) -> None:
    ensure_default_config(ctx.config_path, force=True, env=ctx.env)
    refreshed, _ = load_config(ctx.config_path, ctx.env)
    ctx.config = refreshed
    ctx.render(
        "Regenerated default configuration", {"config_path": str(ctx.config_path)}
    )


def handle_config_paths(ctx: RuntimeContext) -> None:
    cache_dir = default_cache_dir(ctx.env)
    paths = {
        "config": str(ctx.config_path),
        "data": str(ctx.data_dir),
        "state": str(ctx.state_dir),
        "cache": str(cache_dir),
    }
    if ctx.output_format in {"json", "yaml"}:
        ctx.render("Paths", paths)
    else:
        ctx.console.print(f"config: {ctx.config_path}")
        ctx.console.print(f"data:   {ctx.data_dir}")
        ctx.console.print(f"state:  {ctx.state_dir}")
        ctx.console.print(f"cache:  {cache_dir}")


def handle_config_schema(ctx: RuntimeContext) -> None:
    schema_path = Path(__file__).parent.parent / "examples" / "config.schema.json"
    if schema_path.exists():
        ctx.console.print(schema_path.read_text())
    else:
        # Fall back to importlib.resources for installed packages
        try:
            from importlib.resources import files

            schema = (
                files("python_cli")
                .joinpath("../examples/config.schema.json")
                .read_text()
            )
            ctx.console.print(schema)
        except Exception:
            ctx.err_console.print("[red]Error:[/] config.schema.json not found")


def handle_run(ctx: RuntimeContext, opts: RunOptions) -> None:
    payload: dict[str, Any] = {
        "task": opts.task,
        "profile": opts.profile or ctx.config.run.profile,
        "dry_run": ctx.flags.dry_run,
        "timeout_seconds": ctx.flags.resolve_timeout(ctx.config),
        "parallelism": ctx.flags.resolve_parallelism(ctx.config),
    }

    if ctx.flags.dry_run:
        ctx.render("Dry run complete (no changes made)", payload)
        return

    # Simulate work with timing markers for the template.
    start = time.perf_counter()
    time.sleep(0.2)
    duration = time.perf_counter() - start
    payload["duration_seconds"] = round(duration, 3)

    ctx.render("Task completed successfully", payload)
