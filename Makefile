.PHONY: up down build logs python go clean

up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build

logs:
	docker compose logs -f

# Run Python service locally (requires Postgres)
python:
	cd python && uvicorn app.main:app --reload

# Run Go service locally (requires Postgres)
go:
	cd go && go run ./cmd/server

# Remove containers and volumes
clean:
	docker compose down -v
