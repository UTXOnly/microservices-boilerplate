"""
FastAPI async microservice - boilerplate entrypoint.

Run with: uvicorn app.main:app --reload
"""

from contextlib import asynccontextmanager

from fastapi import FastAPI

from app.config import get_settings
from app.database import close_db, init_db
from app.http_client import close_http_client
from app.routers import health, items


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Manage startup and shutdown lifecycle."""
    # Startup
    await init_db()
    yield
    # Shutdown
    await close_db()
    await close_http_client()


app = FastAPI(
    title="Python Microservice",
    description="""
    Async Python microservice boilerplate using FastAPI, Pydantic, and async Postgres.

    ## Features
    - **Async throughout** - Non-blocking I/O with SQLAlchemy 2.0 and asyncpg
    - **Pydantic models** - Request/response validation and serialization
    - **Health checks** - `/health` and `/health/ready` for orchestration
    - **Service layer** - Reusable business logic in `services/`
    - **OpenAPI docs** - Auto-generated at `/docs` and `/redoc`
    """,
    version="1.0.0",
    lifespan=lifespan,
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json",
)

settings = get_settings()

# Include routers
app.include_router(health.router)
app.include_router(items.router)


@app.get("/")
async def root() -> dict:
    """Root endpoint - service info."""
    return {
        "service": settings.app_name,
        "docs": "/docs",
        "health": "/health",
    }


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "app.main:app",
        host=settings.host,
        port=settings.port,
        reload=settings.debug,
    )
