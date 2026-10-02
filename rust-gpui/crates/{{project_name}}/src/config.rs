//! Configuration for a desktop app, loaded from a TOML file at
//! `$XDG_CONFIG_HOME/{{project_name}}/config.toml` (or
//! `~/.config/{{project_name}}/config.toml`). A default is written on first
//! run. See `examples/config.schema.json` for the editor/LSP schema.
//!
//! Resolution order is: command-line args > explicit config file > env vars >
//! this file. With no args/env, we use the local config file.

use std::fs;
use std::path::{Path, PathBuf};

use anyhow::{Context as _, Result};
use serde::{Deserialize, Serialize};

/// The app name, used as the config directory name and the env-var prefix.
pub const APP_NAME: &str = "{{project_name}}";

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Config {
    /// Name of the active design-system scheme (a scheme in the theme set).
    #[serde(default = "default_theme")]
    pub theme: String,
    /// Display face used for a few accents (section headings, status bar).
    #[serde(default = "default_accent_font")]
    pub accent_font: String,
    #[serde(default)]
    pub window: WindowConfig,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WindowConfig {
    #[serde(default = "default_width")]
    pub width: f32,
    #[serde(default = "default_height")]
    pub height: f32,
}

impl Default for WindowConfig {
    fn default() -> Self {
        Self {
            width: default_width(),
            height: default_height(),
        }
    }
}

fn default_theme() -> String {
    crate::theme::DARK.to_string()
}
fn default_accent_font() -> String {
    "SystemUIFont".to_string()
}
fn default_width() -> f32 {
    1200.0
}
fn default_height() -> f32 {
    800.0
}

impl Default for Config {
    fn default() -> Self {
        Self {
            theme: default_theme(),
            accent_font: default_accent_font(),
            window: WindowConfig {
                width: default_width(),
                height: default_height(),
            },
        }
    }
}

/// `$XDG_CONFIG_HOME` or `~/.config`.
fn config_root() -> PathBuf {
    if let Some(dir) = std::env::var_os("XDG_CONFIG_HOME") {
        return PathBuf::from(dir);
    }
    std::env::var_os("HOME")
        .map(PathBuf::from)
        .unwrap_or_else(std::env::temp_dir)
        .join(".config")
}

/// Path to this app's config file.
pub fn config_path() -> PathBuf {
    config_root().join(APP_NAME).join("config.toml")
}

const DEFAULT_TOML: &str = concat!(
    "theme = \"Example Dark\"\n",
    "accent_font = \"SystemUIFont\"\n",
    "\n",
    "[window]\n",
    "width = 1200.0\n",
    "height = 800.0\n",
);

/// Load config, writing a default file on first run.
pub fn load_or_create() -> Result<Config> {
    let path = config_path();
    if !path.exists() {
        let parent = path.parent().unwrap_or(Path::new("."));
        fs::create_dir_all(parent)
            .with_context(|| format!("create config dir {}", parent.display()))?;
        fs::write(&path, DEFAULT_TOML)
            .with_context(|| format!("write default config {}", path.display()))?;
        log::info!("wrote default config to {}", path.display());
    }
    load(&path)
}

/// Load config from a specific file.
pub fn load(path: &Path) -> Result<Config> {
    let builder = config::Config::builder().add_source(config::File::from(path));
    let source = builder.build().with_context(|| format!("load {}", path.display()))?;
    let cfg: Config = source.try_deserialize().context("parse config")?;
    Ok(cfg)
}

/// Validate that the default TOML matches the `Config` struct.
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn default_toml_deserializes() {
        let cfg = config::Config::builder()
            .add_source(config::File::from_str(DEFAULT_TOML, config::FileFormat::Toml))
            .build()
            .unwrap();
        let cfg: Config = cfg.try_deserialize().unwrap();
        assert_eq!(cfg.theme, "Example Dark");
        assert_eq!(cfg.window.width, 1200.0);
    }

    #[test]
    fn config_path_is_under_xdg() {
        assert!(config_path().ends_with("config.toml"));
    }
}