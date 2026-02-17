use axum::{
    extract::{Path, Query, State},
    http::StatusCode,
    response::IntoResponse,
    Json,
};
use sqlx::PgPool;

use crate::models::{Item, ItemCreate, ItemUpdate};

#[derive(Clone)]
pub struct AppState {
    pub pool: PgPool,
    pub app_name: String,
}

// Health
pub async fn liveness(State(state): State<AppState>) -> impl IntoResponse {
    Json(serde_json::json!({
        "status": "alive",
        "service": state.app_name
    }))
}

pub async fn readiness(State(state): State<AppState>) -> impl IntoResponse {
    Json(serde_json::json!({
        "status": "ready",
        "service": state.app_name
    }))
}

// Root
pub async fn root(State(state): State<AppState>) -> impl IntoResponse {
    Json(serde_json::json!({
        "service": state.app_name,
        "docs": "See README for API docs",
        "health": "/health"
    }))
}

// Items
#[derive(Debug, serde::Deserialize)]
pub struct ListParams {
    #[serde(default)]
    pub skip: Option<i32>,
    #[serde(default)]
    pub limit: Option<i32>,
}

pub async fn list_items(
    State(state): State<AppState>,
    Query(params): Query<ListParams>,
) -> impl IntoResponse {
    let skip = params.skip.unwrap_or(0).max(0);
    let limit = params.limit.unwrap_or(100).clamp(1, 100);

    match sqlx::query_as::<_, Item>(
        r#"
        SELECT id, name, description, created_at, updated_at
        FROM items ORDER BY created_at DESC LIMIT $1 OFFSET $2
        "#,
    )
    .bind(limit)
    .bind(skip)
    .fetch_all(&state.pool)
    .await
    {
        Ok(items) => Json(items).into_response(),
        Err(e) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(serde_json::json!({"error": e.to_string()})),
        )
            .into_response(),
    }
}

pub async fn get_item(
    State(state): State<AppState>,
    Path(id): Path<i32>,
) -> impl IntoResponse {
    match sqlx::query_as::<_, Item>(
        r#"
        SELECT id, name, description, created_at, updated_at
        FROM items WHERE id = $1
        "#,
    )
    .bind(id)
    .fetch_optional(&state.pool)
    .await
    {
        Ok(Some(item)) => Json(item).into_response(),
        Ok(None) => (
            StatusCode::NOT_FOUND,
            Json(serde_json::json!({"error": "item not found"})),
        )
            .into_response(),
        Err(e) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(serde_json::json!({"error": e.to_string()})),
        )
            .into_response(),
    }
}

pub async fn create_item(
    State(state): State<AppState>,
    Json(data): Json<ItemCreate>,
) -> impl IntoResponse {
    if data.name.is_empty() {
        return (
            StatusCode::BAD_REQUEST,
            Json(serde_json::json!({"error": "name is required"})),
        )
            .into_response();
    }

    match sqlx::query_as::<_, Item>(
        r#"
        INSERT INTO items (name, description)
        VALUES ($1, $2)
        RETURNING id, name, description, created_at, updated_at
        "#,
    )
    .bind(&data.name)
    .bind(&data.description)
    .fetch_one(&state.pool)
    .await
    {
        Ok(item) => (StatusCode::CREATED, Json(item)).into_response(),
        Err(e) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(serde_json::json!({"error": e.to_string()})),
        )
            .into_response(),
    }
}

pub async fn update_item(
    State(state): State<AppState>,
    Path(id): Path<i32>,
    Json(data): Json<ItemUpdate>,
) -> impl IntoResponse {
    let existing = sqlx::query_as::<_, Item>(
        r#"SELECT id, name, description, created_at, updated_at FROM items WHERE id = $1"#,
    )
    .bind(id)
    .fetch_optional(&state.pool)
    .await;

    let existing = match existing {
        Ok(Some(item)) => item,
        Ok(None) => {
            return (
                StatusCode::NOT_FOUND,
                Json(serde_json::json!({"error": "item not found"})),
            )
                .into_response()
        }
        Err(e) => {
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({"error": e.to_string()})),
            )
                .into_response()
        }
    };

    let name = data.name.as_deref().filter(|s| !s.is_empty()).unwrap_or(&existing.name);
    let description = data.description.or(existing.description);

    match sqlx::query_as::<_, Item>(
        r#"
        UPDATE items SET name = $2, description = $3, updated_at = NOW()
        WHERE id = $1
        RETURNING id, name, description, created_at, updated_at
        "#,
    )
    .bind(id)
    .bind(name)
    .bind(&description)
    .fetch_one(&state.pool)
    .await
    {
        Ok(item) => Json(item).into_response(),
        Err(e) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(serde_json::json!({"error": e.to_string()})),
        )
            .into_response(),
    }
}

pub async fn delete_item(
    State(state): State<AppState>,
    Path(id): Path<i32>,
) -> impl IntoResponse {
    match sqlx::query(r#"DELETE FROM items WHERE id = $1"#)
        .bind(id)
        .execute(&state.pool)
        .await
    {
        Ok(_) => StatusCode::NO_CONTENT.into_response(),
        Err(e) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(serde_json::json!({"error": e.to_string()})),
        )
            .into_response(),
    }
}
