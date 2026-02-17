use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;

#[derive(Debug, Clone, Serialize, FromRow)]
pub struct Item {
    pub id: i32,
    pub name: String,
    pub description: Option<String>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Deserialize)]
pub struct ItemCreate {
    pub name: String,
    pub description: Option<String>,
}

#[derive(Debug, Deserialize, Default)]
pub struct ItemUpdate {
    pub name: Option<String>,
    pub description: Option<String>,
}
