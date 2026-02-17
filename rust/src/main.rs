
mod config;
mod database;
mod handlers;
mod models;

use axum::{
    routing::{delete, get, patch, post},
    Router,
};
use std::net::SocketAddr;
use tokio::net::TcpListener;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::registry()
        .with(tracing_subscriber::EnvFilter::new(
            std::env::var("RUST_LOG").unwrap_or_else(|_| "info".into()),
        ))
        .with(tracing_subscriber::fmt::layer())
        .init();

    let config = config::Config::load();

    let pool = database::create_pool(&config).await?;
    database::init_schema(&pool).await?;

    let state = handlers::AppState {
        pool,
        app_name: config.app_name.clone(),
    };

    let app = Router::new()
        .route("/", get(handlers::root))
        .route("/health", get(handlers::liveness))
        .route("/health/ready", get(handlers::readiness))
        .route("/items", get(handlers::list_items).post(handlers::create_item))
        .route(
            "/items/:id",
            get(handlers::get_item)
                .patch(handlers::update_item)
                .delete(handlers::delete_item),
        )
        .with_state(state);

    let addr = SocketAddr::from(([0, 0, 0, 0], config.port));
    tracing::info!("listening on {}", addr);

    let listener = TcpListener::bind(addr).await?;
    axum::serve(listener, app).await?;

    Ok(())
}
