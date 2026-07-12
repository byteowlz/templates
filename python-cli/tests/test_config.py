from __future__ import annotations

from pathlib import Path

from python_cli.config import CONFIG_SCHEMA_URL, ensure_default_config, load_config


def test_load_config_applies_file_before_environment(tmp_path: Path) -> None:
    config_path = tmp_path / "config.toml"
    config_path.write_text(
        """
[logging]
level = "warning"

[run]
default_task = "from-file"
parallelism = 2
""".strip()
        + "\n",
        encoding="utf-8",
    )

    env = {
        "PYTHON_CLI__RUN__PARALLELISM": "8",
        "PYTHON_CLI__RUN__PROFILE": "from-env",
    }

    config, resolved = load_config(config_path, env)

    assert resolved == config_path
    assert config.logging.level == "warning"
    assert config.run.default_task == "from-file"
    assert config.run.parallelism == 8
    assert config.run.profile == "from-env"


def test_ensure_default_config_writes_schema_reference(tmp_path: Path) -> None:
    config_path = tmp_path / "config.toml"

    written = ensure_default_config(config_path, env={})
    content = written.read_text(encoding="utf-8")

    assert written == config_path
    assert content.startswith(f'"$schema" = "{CONFIG_SCHEMA_URL}"')
    assert "[logging]" in content
