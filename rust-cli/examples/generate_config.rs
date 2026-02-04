//! Generate config schema and example config files.
//!
//! Run with: cargo run --example generate_config
//! Or use: just generate-config

use std::path::PathBuf;

fn main() -> anyhow::Result<()> {
    let manifest_dir = std::env::var("CARGO_MANIFEST_DIR").expect("CARGO_MANIFEST_DIR not set");
    let output_dir = PathBuf::from(&manifest_dir).join("examples");

    const APP_NAME: &str = env!("CARGO_PKG_NAME");
    const REPO_URL: &str = "https://github.com/byteowlz/rust-cli";

    println!("Generating config files to {}", output_dir.display());

    rust_cli::write_generated_files(&output_dir, APP_NAME, REPO_URL)?;

    println!("Generated:");
    println!("  - {}/config.schema.json", output_dir.display());
    println!("  - {}/config.toml", output_dir.display());

    Ok(())
}
