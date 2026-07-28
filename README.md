# Bookmark Management Service

A Go REST API for creating temporary short links and redirecting them to their original URLs. Link data is stored in Redis and may expire after at most seven days.

## Features

- Create a short link with a caller-provided expiration time
- Redirect a short code to its original URL
- Redis-backed health check
- Swagger UI and OpenAPI specifications
- Local and reverse-proxy Docker Compose deployments

## Prerequisites

- Go 1.24+
- Docker and Docker Compose (recommended for local Redis and deployment)
- `make` for Makefile commands

## Configuration

The service reads environment variables directly. Defaults are used when a variable is omitted.

| Variable | Default | Description |
| --- | --- | --- |
| `APP_PORT` | `8080` | HTTP port exposed by the API |
| `SERVICE_NAME` | `bookmark_service` | Name returned by the health check |
| `INSTANCE_ID` | generated UUID | Service instance identifier |
| `BASE_PATH` | `/` | Swagger API base path |
| `LOG_LEVEL` | `info` | Zerolog log level |
| `REDIS_ADDR` | `localhost:6379` | Redis host and port |
| `REDIS_PWD` | empty | Redis password |
| `REDIS_DB` | `0` | Redis database index |

Do **not** commit environment files. For the Nginx deployment, copy the tracked template and adjust it locally:

```bash
cp deployment/bookmark-service/.env.example deployment/bookmark-service/.env
```

## Run locally

### With Docker Compose

This is the quickest option; it starts Redis and exposes the API at port `8080`.

```bash
docker compose up --build
```

Stop the stack with:

```bash
docker compose down
```

### With Go

Start Redis first, then run the API:

```bash
docker run --rm -p 6379:6379 redis:alpine
go mod download
make dev-run
```

`make dev-run` regenerates Swagger documentation before starting the API. To run without regeneration, use `make run-app`.

## API and Swagger

When running locally, Swagger UI is available at:

```text
http://localhost:8080/swagger/index.html
```

Main endpoints:

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health-check` | Checks that Redis is available |
| `POST` | `/v1/links/shorten` | Creates a short link |
| `GET` | `/v1/links/redirect/:code` | Redirects to the original URL |

Example request:

```bash
curl --request POST http://localhost:8080/v1/links/shorten \
  --header 'Content-Type: application/json' \
  --data '{"url":"https://example.com","exp":3600}'
```

`exp` is required, expressed in seconds, and cannot exceed `604800` (seven days).

## Testing and generated docs

```bash
make test
make swagger
```

`make test` enforces the configured coverage threshold and creates local coverage artifacts, all ignored by Git.

## Deployment with Nginx

The deployment stack starts Redis, the API, and Nginx. Nginx exposes the service on port `80` under `/api/bookmark_service/`.

```bash
cp deployment/bookmark-service/.env.example deployment/bookmark-service/.env
docker compose -f deployment/docker-compose.yaml up --build -d
```

Useful deployment commands:

```bash
docker compose -f deployment/docker-compose.yaml logs -f
docker compose -f deployment/docker-compose.yaml down
```

After deployment, use the proxied API path, for example:

```text
http://localhost/api/bookmark_service/swagger/index.html
```

## Repository hygiene

- Keep actual `.env` files local; use `.env.example` for safe defaults.
- Do not commit generated `repomix-output.xml`, binaries, coverage reports, or editor files.
- Regenerate Swagger files with `make swagger` after changing handler annotations.
