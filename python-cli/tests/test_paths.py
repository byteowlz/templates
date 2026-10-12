from __future__ import annotations

from pathlib import Path

import pytest

from python_cli.config import APP_NAME, default_cache_dir, default_config_path, default_data_dir, default_state_dir


@pytest.mark.parametrize("value", [None, "", "relative/config", "~/config", "$HOME/config"])
def test_invalid_or_missing_xdg_values_use_spec_defaults(value: str | None) -> None:
    env = {"HOME": "/home/ada"}
    if value is not None:
        env.update(XDG_CONFIG_HOME=value, XDG_DATA_HOME=value, XDG_STATE_HOME=value, XDG_CACHE_HOME=value)
    assert [default_config_path(env), default_data_dir(env), default_state_dir(env), default_cache_dir(env)] == [
        Path("/home/ada/.config") / APP_NAME / "config.toml",
        Path("/home/ada/.local/share") / APP_NAME,
        Path("/home/ada/.local/state") / APP_NAME,
        Path("/home/ada/.cache") / APP_NAME,
    ]


def test_absolute_xdg_values_are_used() -> None:
    env = {"HOME": "/home/ada", "XDG_CONFIG_HOME": "/c", "XDG_DATA_HOME": "/d", "XDG_STATE_HOME": "/s", "XDG_CACHE_HOME": "/k"}
    assert [default_config_path(env), default_data_dir(env), default_state_dir(env), default_cache_dir(env)] == [
        Path("/c") / APP_NAME / "config.toml",
        Path("/d") / APP_NAME,
        Path("/s") / APP_NAME,
        Path("/k") / APP_NAME,
    ]
