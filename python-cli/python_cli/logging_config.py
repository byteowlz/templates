from __future__ import annotations

import logging
from typing import Any

from rich.console import Console
from rich.logging import RichHandler

TRACE_LEVEL = 5
logging.addLevelName(TRACE_LEVEL, "TRACE")


def _trace(self: logging.Logger, message: str, *args: Any, **kwargs: Any) -> None:
    if self.isEnabledFor(TRACE_LEVEL):
        self._log(TRACE_LEVEL, message, args, **kwargs)


if not hasattr(logging.Logger, "trace"):
    setattr(logging.Logger, "trace", _trace)


def configure_logging(*, console: Console, level: int, diagnostics: bool) -> None:
    handler = RichHandler(
        console=console,
        show_time=True,
        show_path=diagnostics,
        markup=True,
        rich_tracebacks=True,
    )

    root = logging.getLogger()

    # Remove any existing RichHandler to avoid duplicate logs during tests.
    root.handlers = [h for h in root.handlers if not isinstance(h, RichHandler)]
    root.addHandler(handler)
    root.setLevel(level)
