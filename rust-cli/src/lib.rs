//! rust-cli library for config management and schema generation.
//!
//! This module provides functions to generate JSON schemas and example TOML
//! configurations from the config struct definitions.

use std::fs;
use std::path::Path;

use anyhow::{Context, Result};
use schemars::generate::SchemaSettings;
use schemars::JsonSchema;
use schemars::Schema;
use serde::{Deserialize, Serialize};
use serde_json::json;

/// Generated schema filename.
pub const SCHEMA_FILENAME: &str = "config.schema.json";

/// Generated config filename.
pub const CONFIG_FILENAME: &str = "config.toml";

/// Application configuration.
#[derive(Debug, Clone, Serialize, Deserialize, JsonSchema)]
#[serde(default)]
pub struct AppConfig {
    /// Active configuration profile name
    pub profile: String,
    /// Logging configuration
    pub logging: LoggingConfig,
    /// Runtime behavior configuration
    pub runtime: RuntimeConfig,
    /// Directory path overrides
    pub paths: PathsConfig,
}

impl Default for AppConfig {
    fn default() -> Self {
        Self {
            profile: "default".to_string(),
            logging: LoggingConfig::default(),
            runtime: RuntimeConfig::default(),
            paths: PathsConfig::default(),
        }
    }
}

/// Logging configuration.
#[derive(Debug, Clone, Serialize, Deserialize, JsonSchema)]
#[serde(default)]
pub struct LoggingConfig {
    /// Log level (trace, debug, info, warn, error)
    pub level: String,
    /// Optional file path to write logs to
    pub file: Option<String>,
}

impl Default for LoggingConfig {
    fn default() -> Self {
        Self {
            level: "info".to_string(),
            file: None,
        }
    }
}

/// Runtime behavior configuration.
#[derive(Debug, Clone, Serialize, Deserialize, JsonSchema)]
#[serde(default)]
pub struct RuntimeConfig {
    /// Number of parallel tasks (defaults to CPU count)
    pub parallelism: Option<usize>,
    /// Default timeout in seconds for operations
    pub timeout: Option<u64>,
    /// Stop on first error instead of continuing
    pub fail_fast: bool,
}

impl Default for RuntimeConfig {
    fn default() -> Self {
        Self {
            parallelism: None,
            timeout: Some(60),
            fail_fast: true,
        }
    }
}

/// Directory path overrides.
#[derive(Debug, Default, Clone, Serialize, Deserialize, JsonSchema)]
#[serde(default)]
pub struct PathsConfig {
    /// Override the data directory path
    pub data_dir: Option<String>,
    /// Override the state directory path
    pub state_dir: Option<String>,
}

/// Generate the JSON schema for AppConfig using schemars.
pub fn generate_schema(project_name: &str, repo_url: &str) -> Result<String> {
    // Use draft-07 for better TOML editor support
    let settings = SchemaSettings::draft07();
    let generator = settings.into_generator();
    let mut schema: Schema = generator.into_root_schema_for::<AppConfig>();

    // Set schema metadata
    schema.insert(
        "$id".to_string(),
        json!(format!("{repo_url}/schemas/config.schema.json")),
    );
    schema.insert(
        "title".to_string(),
        json!(format!("{project_name} configuration")),
    );
    schema.insert(
        "description".to_string(),
        json!(format!("Configuration schema for {project_name}")),
    );

    // Add $schema property for LSP/editor support
    if let Some(props) = schema.get_mut("properties") {
        if let Some(props_obj) = props.as_object_mut() {
            props_obj.insert(
                "$schema".to_string(),
                json!({
                    "type": "string",
                    "description": "JSON Schema reference for editor support"
                }),
            );
        }
    }

    serde_json::to_string_pretty(&schema).context("serializing JSON schema")
}

/// Generate the example TOML configuration from the default AppConfig.
pub fn generate_example_config(project_name: &str) -> Result<String> {
    let schema_url = format!(
        "https://raw.githubusercontent.com/byteowlz/schemas/refs/heads/main/{project_name}/{project_name}.config.schema.json"
    );

    // Serialize the default config to TOML
    let config = AppConfig::default();
    let toml_body =
        toml::to_string_pretty(&config).context("serializing default config to TOML")?;

    // Build output with schema reference and header
    let mut output = String::new();
    output.push_str(&format!(
        r#""$schema" = "{schema_url}"

# Configuration for {project_name}.
# Copy this file to $XDG_CONFIG_HOME/{project_name}/config.toml and adjust as needed.

"#
    ));
    output.push_str(&toml_body);

    Ok(output)
}

/// Write generated files to a directory.
pub fn write_generated_files(output_dir: &Path, project_name: &str, repo_url: &str) -> Result<()> {
    fs::create_dir_all(output_dir)
        .with_context(|| format!("creating output directory: {}", output_dir.display()))?;

    let schema = generate_schema(project_name, repo_url)?;
    let schema_path = output_dir.join(SCHEMA_FILENAME);
    fs::write(&schema_path, &schema)
        .with_context(|| format!("writing schema to {}", schema_path.display()))?;

    let config = generate_example_config(project_name)?;
    let config_path = output_dir.join(CONFIG_FILENAME);
    fs::write(&config_path, &config)
        .with_context(|| format!("writing config to {}", config_path.display()))?;

    Ok(())
}

/// Compare generated files against existing files in a directory.
/// Returns Ok(()) if they match, Err with diff details if they differ.
pub fn validate_against_examples(
    examples_dir: &Path,
    project_name: &str,
    repo_url: &str,
) -> Result<()> {
    let schema = generate_schema(project_name, repo_url)?;
    let config = generate_example_config(project_name)?;

    let schema_path = examples_dir.join(SCHEMA_FILENAME);
    let config_path = examples_dir.join(CONFIG_FILENAME);

    let mut errors = Vec::new();

    // Check schema
    if schema_path.exists() {
        let existing = fs::read_to_string(&schema_path)
            .with_context(|| format!("reading {}", schema_path.display()))?;
        if existing != schema {
            errors.push(format!(
                "{} is out of date. Run 'just generate-config' to update.",
                schema_path.display()
            ));
        }
    } else {
        errors.push(format!(
            "{} does not exist. Run 'just generate-config' to create.",
            schema_path.display()
        ));
    }

    // Check config
    if config_path.exists() {
        let existing = fs::read_to_string(&config_path)
            .with_context(|| format!("reading {}", config_path.display()))?;
        if existing != config {
            errors.push(format!(
                "{} is out of date. Run 'just generate-config' to update.",
                config_path.display()
            ));
        }
    } else {
        errors.push(format!(
            "{} does not exist. Run 'just generate-config' to create.",
            config_path.display()
        ));
    }

    if errors.is_empty() {
        Ok(())
    } else {
        anyhow::bail!(
            "Generated config/schema validation failed:\n  - {}",
            errors.join("\n  - ")
        )
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    const APP_NAME: &str = env!("CARGO_PKG_NAME");
    const REPO_URL: &str = "https://github.com/byteowlz/rust-cli";

    #[test]
    fn test_schema_generation() {
        let schema = generate_schema(APP_NAME, REPO_URL).expect("schema generation failed");
        assert!(schema.contains("\"title\""));
        assert!(schema.contains("rust-cli configuration"));
        assert!(schema.contains("\"$schema\""));
    }

    #[test]
    fn test_config_generation() {
        let config = generate_example_config(APP_NAME).expect("config generation failed");
        assert!(config.contains("[logging]"));
        assert!(config.contains("[runtime]"));
        assert!(config.contains("$schema"));
    }

    #[test]
    fn validate_examples_are_up_to_date() {
        let manifest_dir = std::env::var("CARGO_MANIFEST_DIR").expect("CARGO_MANIFEST_DIR not set");
        let crate_root = Path::new(&manifest_dir);
        let examples_dir = crate_root.join("examples");

        if !examples_dir.exists() {
            panic!(
                "examples/ directory not found at {}. Create it and run 'just generate-config'.",
                examples_dir.display()
            );
        }

        validate_against_examples(&examples_dir, APP_NAME, REPO_URL)
            .expect("examples/ files are out of date");
    }
}
