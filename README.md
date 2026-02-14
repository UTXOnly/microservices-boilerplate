# Microservices Boilerplate

Production-ready boilerplate for building async microservices in **Python** (FastAPI + Pydantic) and **Go**. Both templates include async PostgreSQL, health checks, Docker support, and common microservice patterns.

Use this to rapidly scaffold new microservices—copy a template, rename it, and start building.

---

## Quick Start

### Run Everything with Docker

```bash
docker compose up --build
```

- **Python service**: http://localhost:8000  
- **Go service**: http://localhost:8080  
- **PostgreSQL**: localhost:5432 (postgres/postgres)

### API Documentation

| Service | Swagger UI | ReDoc |
|---------|------------|-------|
| Python | http://localhost:8000/docs | http://localhost:8000/redoc |

---

## Project Structure

```
microservices-boilerplate/
├── python/                 # FastAPI + Pydantic + async Postgres
│   ├── app/
│   │   ├── main.py         # Application entrypoint
│   │   ├── config.py      # Pydantic Settings
│   │   ├── database.py    # Async SQLAlchemy + asyncpg
│   │   ├── models/        # SQLAlchemy ORM models
│   │   ├── schemas/       # Pydantic request/response models
│   │   ├── routers/       # API route handlers
│   │   ├── services/      # Business logic layer
│   │   └── http_client.py # Service-to-service HTTP client
│   ├── requirements.txt
│   └── Dockerfile
├── go/                     # Echo + pgx + async/concurrent Postgres
│   ├── cmd/server/
│   │   └── main.go        # Application entrypoint
│   ├── internal/
│   │   ├── config/        # Configuration
│   │   ├── database/      # Connection pool + migrations
│   │   ├── models/        # Domain structs
│   │   ├── handlers/      # HTTP handlers
│   │   └── services/      # Business logic
│   ├── go.mod
│   └── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## Python Microservice

### Stack

- **FastAPI** – Async web framework
- **Pydantic** – Validation and serialization
- **SQLAlchemy 2.0** – Async ORM with asyncpg
- **httpx** – Async HTTP client for inter-service calls

### Run Locally

```bash
cd python
python -m venv .venv
source .venv/bin/activate   # Windows: .venv\Scripts\activate
pip install -r requirements.txt
cp .env.example .env        # Edit DATABASE_URL if needed
uvicorn app.main:app --reload
```

### Key Files

| File | Purpose |
|------|---------|
| `app/config.py` | Pydantic Settings from env vars |
| `app/database.py` | Async session, `get_db()` dependency |
| `app/models/` | SQLAlchemy ORM models |
| `app/schemas/` | Pydantic schemas (Create, Read, Update) |
| `app/services/` | Reusable business logic |
| `app/http_client.py` | Shared async HTTP client |

### Adding a New Resource

1. Add SQLAlchemy model in `app/models/`
2. Add Pydantic schemas in `app/schemas/`
3. Add service methods in `app/services/`
4. Add router in `app/routers/`
5. Include router in `app/main.py`

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_NAME` | python-microservice | Service name |
| `PORT` | 8000 | Server port |
| `DATABASE_URL` | postgresql+asyncpg://... | Async Postgres URL |
| `DEBUG` | false | Enable debug logging |

---

## Go Microservice

### Stack

- **Echo** – Lightweight HTTP framework
- **pgx v5** – Concurrent Postgres driver with connection pool
- **Standard library** – Config, JSON, errors

### Run Locally

```bash
cd go
cp .env.example .env        # Edit if needed
go run ./cmd/server
```

### Key Files

| File | Purpose |
|------|---------|
| `internal/config/` | Env-based config |
| `internal/database/` | Pool, migrations |
| `internal/models/` | Request/response structs |
| `internal/handlers/` | HTTP handlers |
| `internal/services/` | Business logic |

### Adding a New Resource

1. Add model structs in `internal/models/`
2. Add service in `internal/services/`
3. Add handler in `internal/handlers/`
4. Register routes in `cmd/server/main.go`

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_NAME` | go-microservice | Service name |
| `PORT` | 8080 | Server port |
| `DATABASE_URL` | postgres://... | Postgres connection string |
| `DEBUG` | false | Debug mode |

---

## API Endpoints

Both services expose the same API surface (shared `items` table):

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Liveness probe |
| GET | `/health/ready` | Readiness probe |
| GET | `/items` | List items (paginated) |
| GET | `/items/{id}` | Get item by ID |
| POST | `/items` | Create item |
| PATCH | `/items/{id}` | Partial update |
| DELETE | `/items/{id}` | Delete item |

### Example: Create an Item

```bash
curl -X POST http://localhost:8000/items \
  -H "Content-Type: application/json" \
  -d '{"name": "My Item", "description": "A test item"}'
```

---

## Docker

### Build & Run

```bash
# All services
docker compose up -d

# Specific service
docker compose up -d python-service

# Rebuild after changes
docker compose up --build -d
```

### Individual Dockerfiles

Each service can be built standalone:

```bash
# Python
docker build -t my-python-service ./python

# Go
docker build -t my-go-service ./go
```

---

## Using as a Template

1. **Clone or copy** the repo or the specific `python/` or `go/` directory.
2. **Rename** the service (app name, package names, Docker image).
3. **Replace** the example `Item` model with your domain entities.
4. **Add** your business logic in `services/`.
5. **Extend** routers/handlers for new endpoints.

---

## Roadmap

- [ ] Java template (planned)

---

## License

MIT
