"""Async HTTP client for service-to-service calls."""

from functools import lru_cache

import httpx


@lru_cache
def get_http_client() -> httpx.AsyncClient:
    """
    Reusable async HTTP client for outbound requests.
    Use for calling other microservices, external APIs, etc.
    """
    return httpx.AsyncClient(
        timeout=30.0,
        follow_redirects=True,
        headers={"User-Agent": "microservice-client/1.0"},
    )


async def close_http_client() -> None:
    """Close HTTP client on shutdown."""
    client = get_http_client()
    await client.aclose()
