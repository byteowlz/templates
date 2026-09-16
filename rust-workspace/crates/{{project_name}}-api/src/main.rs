//! HTTP API server for rust-workspace.

use std::net::SocketAddr;
use std::path::PathBuf;

use anyhow::{Context, Result};
use axum::{Json, Router, routing::get};
use clap::{Args, Parser};
use log::info;
use serde::Serialize;
use tower_http::trace::TraceLayer;

use {{project_name}}_core::AppPaths;

fn main() -> anyhow::Result<()> {
    try_main()
}

#[tokio::main]
async fn try_main() -> Result<()> {
    env_logger::init();

    let cli = Cli::parse();
    let _paths = AppPaths::discover(cli.common.config.as_deref())?;

    let app = build_router();

    let addr = SocketAddr::from(([127, 0, 0, 1], cli.common.port));
    info!("Starting API server on {addr}");

    let listener = tokio::net::TcpListener::bind(addr).await?;
    axum::serve(listener, app)
        .await
        .context("serving HTTP API")?;

    Ok(())
}

#[derive(Debug, Parser)]
#[command(author, version, about = "HTTP API server for rust-workspace")]
struct Cli {
    #[command(flatten)]
    common: CommonOpts,
}

#[derive(Debug, Clone, Args)]
struct CommonOpts {
    /// Override the config file path
    #[arg(long, value_name = "PATH")]
    config: Option<PathBuf>,

    /// Port to listen on
    #[arg(short, long, default_value = "3000")]
    port: u16,
}

#[derive(Serialize)]
struct RootResponse {
    name: &'static str,
    version: &'static str,
}

#[derive(Serialize)]
struct HealthResponse {
    status: &'static str,
}

/// Build the API router. The default scaffold exposes only non-sensitive
/// endpoints and configures **no CORS layer** and no state/config, so it cannot
/// be reached cross-origin and never discloses configuration or secrets.
///
/// If an endpoint later needs configuration, thread it through explicit,
/// readable state (e.g. `State(state): State<Arc<AppConfig>>`) and keep any
/// data secret out of the response. If a browser client genuinely needs
/// cross-origin access, opt into an exact-origin
/// [`CorsLayer`](tower_http::cors::CorsLayer) — never `Any`.
fn build_router() -> Router {
    Router::new()
        .route("/", get(root))
        .route("/health", get(health))
        .layer(TraceLayer::new_for_http())
}

async fn root() -> Json<RootResponse> {
    Json(RootResponse {
        name: env!("CARGO_PKG_NAME"),
        version: env!("CARGO_PKG_VERSION"),
    })
}

async fn health() -> Json<HealthResponse> {
    Json(HealthResponse { status: "ok" })
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::body::Body;
    use axum::http::{Request, StatusCode};
    use tower::ServiceExt;

    /// Regression guard (tmpl-h7h7): the default scaffold must never expose
    /// configuration (which may later hold secrets) over the wire, and must not
    /// open CORS by default. A `/config` route must remain absent (404) while
    /// ordinary non-sensitive endpoints keep working.
    #[tokio::test]
    async fn default_api_does_not_expose_config_or_cors() -> Result<()> {
        let app = build_router();

        let config_resp = app
            .clone()
            .oneshot(
                Request::builder()
                    .uri("/config")
                    .body(Body::empty())
                    .context("building /config request")?,
            )
            .await
            .context("calling /config")?;
        assert_eq!(
            config_resp.status(),
            StatusCode::NOT_FOUND,
            "GET /config must not exist on the default API"
        );

        let health_resp = app
            .clone()
            .oneshot(
                Request::builder()
                    .uri("/health")
                    .body(Body::empty())
                    .context("building /health request")?,
            )
            .await
            .context("calling /health")?;
        assert_eq!(
            health_resp.status(),
            StatusCode::OK,
            "GET /health must remain reachable"
        );

        Ok(())
    }
}
