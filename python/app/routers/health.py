"""Health check endpoints for readiness and liveness probes."""

from fastapi import APIRouter

from app.config import get_settings

router = APIRouter(tags=["Health"])
settings = get_settings()


@router.get("/health", summary="Liveness probe")
async def liveness() -> dict:
    """
    Simple liveness check - returns 200 if the service is running.
    Use for Kubernetes/Docker health checks.
    """
    return {"status": "alive", "service": settings.app_name}


@router.get("/health/ready", summary="Readiness probe")
async def readiness() -> dict:
    """
    Readiness check - verifies the service can accept traffic.
    Extend to check DB connectivity, external services, etc.
    """
    return {"status": "ready", "service": settings.app_name}
