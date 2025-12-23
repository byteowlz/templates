//! Shared core library for the rust-workspace template.
//!
//! This crate provides:
//! - Configuration loading and management
//! - XDG-compliant path resolution
//! - Common types and error handling

pub mod config;
pub mod error;
pub mod paths;

pub use config::{AppConfig, LoggingConfig, PathsConfig, RuntimeConfig};
pub use error::{CoreError, Result};
pub use paths::AppPaths;

/// Application name used for config directories and environment prefix.
/// Override this constant when scaffolding a new project.
pub const APP_NAME: &str = "rust-workspace";

/// Returns the environment variable prefix for this application.
pub fn env_prefix() -> String {
    APP_NAME
        .chars()
        .map(|c| {
            if c.is_ascii_alphanumeric() {
                c.to_ascii_uppercase()
            } else {
                '_'
            }
        })
        .collect()
}

/// Returns the default parallelism based on available CPU cores.
pub fn default_parallelism() -> usize {
    std::thread::available_parallelism()
        .map(|n| n.get())
        .unwrap_or(1)
}
