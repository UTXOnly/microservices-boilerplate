use std::env;

pub struct Config {
    pub app_name: String,
    pub port: u16,
    pub database_url: String,
}

impl Config {
    pub fn load() -> Self {
        let _ = dotenvy::dotenv();

        let port = env::var("PORT")
            .ok()
            .and_then(|s| s.parse().ok())
            .unwrap_or(8081);

        let database_url = env::var("DATABASE_URL").unwrap_or_else(|_| {
            "postgres://postgres:postgres@localhost:5432/microservices".to_string()
        });

        let app_name = env::var("APP_NAME")
            .unwrap_or_else(|_| "rust-microservice".to_string());

        Self {
            app_name,
            port,
            database_url,
        }
    }
}
