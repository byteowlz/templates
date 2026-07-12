from __future__ import annotations

from datetime import UTC, datetime
from pathlib import Path

from rich.console import Console

from python_cli.config import AppConfig
from python_cli.handlers import handle_config_schema
from python_cli.runtime import CommonFlags, RuntimeContext


def test_handle_config_schema_reads_packaged_schema(tmp_path: Path) -> None:
    stdout = Console(file=None, record=True, force_terminal=False)
    stderr = Console(file=None, record=True, force_terminal=False, stderr=True)
    ctx = RuntimeContext(
        flags=CommonFlags(),
        config=AppConfig(),
        config_path=tmp_path / "config.toml",
        console=stdout,
        err_console=stderr,
        started_at=datetime.now(UTC),
        data_dir=tmp_path / "data",
        state_dir=tmp_path / "state",
    )

    handle_config_schema(ctx)

    output = stdout.export_text()
    assert '"title": "python-cli configuration"' in output
    assert '"logging"' in output
