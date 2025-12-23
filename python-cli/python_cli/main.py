from __future__ import annotations

import os
from pathlib import Path
from typing import Optional

import typer
from typer import Context

from . import __version__
from .config import APP_NAME
from .handlers import (
    InitOptions,
    RunOptions,
    handle_config_path,
    handle_config_reset,
    handle_config_show,
    handle_init,
    handle_run,
)
from .runtime import CommonFlags, RuntimeContext, new_runtime

app = typer.Typer(help="Opinionated starting point for cross-platform Python CLIs.")
config_app = typer.Typer(help="Inspect and manage configuration.")
app.add_typer(config_app, name="config")

SHELL_CHOICES = ["bash", "zsh", "fish", "powershell", "elvish"]


def _version_callback(value: bool) -> None:
    if value:
        typer.echo(__version__)
        raise typer.Exit(code=0)


@app.callback()
def main(
    ctx: Context,
    version: bool = typer.Option(
        False,
        "--version",
        help="Show the python-cli version and exit.",
        callback=_version_callback,
        is_eager=True,
    ),
    config: Optional[Path] = typer.Option(
        None,
        "--config",
        help="Override the config file path.",
        exists=False,
        file_okay=True,
        dir_okay=False,
        writable=True,
    ),
    quiet: bool = typer.Option(False, "--quiet", "-q", help="Reduce output to only errors."),
    verbose: int = typer.Option(
        0,
        "--verbose",
        "-v",
        help="Increase logging verbosity (stackable).",
        count=True,
    ),
    debug: bool = typer.Option(False, "--debug", help="Enable debug logging (equivalent to -vv)."),
    trace: bool = typer.Option(False, "--trace", help="Enable trace logging (overrides other levels)."),
    json_output: bool = typer.Option(False, "--json", help="Output machine-readable JSON."),
    yaml_output: bool = typer.Option(False, "--yaml", help="Output machine-readable YAML."),
    no_color: bool = typer.Option(False, "--no-color", help="Disable ANSI colors in output."),
    color: str = typer.Option(
        "auto",
        "--color",
        help="Color output policy: auto, always, or never.",
        case_sensitive=False,
    ),
    dry_run: bool = typer.Option(False, "--dry-run", help="Do not change anything on disk."),
    assume_yes: bool = typer.Option(
        False,
        "--yes",
        "-y",
        help="Assume yes for interactive prompts (alias for --force).",
    ),
    no_progress: bool = typer.Option(False, "--no-progress", help="Disable progress indicators."),
    diagnostics: bool = typer.Option(
        False,
        "--diagnostics",
        help="Emit additional diagnostics for troubleshooting.",
    ),
    timeout_seconds: Optional[int] = typer.Option(
        None,
        "--timeout",
        min=1,
        help="Maximum seconds to allow an operation to run.",
    ),
    parallelism: Optional[int] = typer.Option(
        None,
        "--parallel",
        min=1,
        help="Override the degree of parallelism.",
    ),
) -> None:
    del version  # handled via callback

    flags = CommonFlags(
        config_path=config,
        quiet=quiet,
        verbose=verbose,
        debug=debug,
        trace=trace,
        json_output=json_output,
        yaml_output=yaml_output,
        no_color=no_color,
        color=color.lower(),
        dry_run=dry_run,
        assume_yes=assume_yes,
        no_progress=no_progress,
        diagnostics=diagnostics,
        timeout_seconds=timeout_seconds,
        parallelism=parallelism,
    )

    try:
        runtime = new_runtime(flags, env=os.environ)
    except ValueError as exc:
        raise typer.BadParameter(str(exc))

    ctx.obj = runtime


def _runtime(ctx: Context) -> RuntimeContext:
    runtime = ctx.obj
    if not isinstance(runtime, RuntimeContext):
        raise RuntimeError("internal error: runtime context not initialized")
    return runtime


@app.command("run")
def run_command(
    ctx: Context,
    task: Optional[str] = typer.Argument(None, help="Optional task name to execute."),
    profile: Optional[str] = typer.Option(None, "--profile", help="Override the profile to run under."),
) -> None:
    runtime = _runtime(ctx)
    options = RunOptions(task=task or runtime.config.run.default_task, profile=profile)
    handle_run(runtime, options)


@app.command("init")
def init_command(
    ctx: Context,
    force: bool = typer.Option(False, "--force", help="Recreate configuration even if it already exists."),
) -> None:
    runtime = _runtime(ctx)
    handle_init(runtime, InitOptions(force=force))


@config_app.command("show")
def config_show(ctx: Context) -> None:
    runtime = _runtime(ctx)
    handle_config_show(runtime)


@config_app.command("path")
def config_path(ctx: Context) -> None:
    runtime = _runtime(ctx)
    handle_config_path(runtime)


@config_app.command("reset")
def config_reset(ctx: Context) -> None:
    runtime = _runtime(ctx)
    handle_config_reset(runtime)


@app.command("completions")
def completions(
    shell: str = typer.Argument(..., metavar="SHELL", help=f"Shell to generate completions for ({', '.join(SHELL_CHOICES)})."),
) -> None:
    shell = shell.lower()
    if shell not in SHELL_CHOICES:
        raise typer.BadParameter(f"Unsupported shell: {shell}. Choose from {', '.join(SHELL_CHOICES)}.")
    from typer.main import get_command
    from click.shell_completion import get_completion_class

    command = get_command(app)
    completion_class = get_completion_class(shell)
    env_var = f"_{APP_NAME.upper().replace('-', '_')}_COMPLETE"
    completer = completion_class(command, {}, command.name or APP_NAME, env_var)
    typer.echo(completer.source())


def app_name() -> str:
    return APP_NAME
