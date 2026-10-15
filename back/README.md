# Frise Chrono — back

Go service holding every business rule of Frise Chrono: zones, decades and dates, the
random draw of images, and the gateway to Wikimedia Commons. It exposes the HTTP API
used by the web and, later, the mobile application. Image files are never served:
responses point to the source URLs.

## Start

From the repository root (see the [root README](../README.md) for the first setup):

```bash
./compose up back       # starts the back and Redis, with live reload
./compose logs back     # follows its logs
```

The API listens on http://localhost:8080/api/v1; its reference is on
http://localhost:8080/api/docs.

## Commands

Run with `./compose make back <target>`:

| Target    | Action                                                     |
|-----------|------------------------------------------------------------|
| `help`    | List the commands                                          |
| `install` | Download the Go modules                                    |
| `run`     | Run the server with live reload (default command of the container) |
| `build`   | Build the server binary into `bin/`                        |
| `lint`    | Lint, check formatting and architecture boundaries         |
| `format`  | Format the code                                            |
| `api`     | Regenerate the server code from `api/openapi.yaml`         |

## API

The contract `api/openapi.yaml` is the source of truth
([ADR functional 0003](../docs/adr-functional/0003-contract-first-with-openapi.md)).

| Method | Path                                | Purpose                                  |
|--------|-------------------------------------|------------------------------------------|
| GET    | `/api/v1/zones`                     | The 9 zones                              |
| GET    | `/api/v1/zones/geometry`            | Zone outlines (GeoJSON)                  |
| GET    | `/api/v1/decades/{decade}/zones`    | Zones having photos for a decade         |
| GET    | `/api/v1/images?zone=&decade=&limit=` | Random images of a zone and a decade   |
| GET    | `/api/health`                       | Health check                             |

Errors are Problem Details (RFC 9457).

**Changing the API**: edit `api/openapi.yaml`, run `./compose make back api`, then
implement the new or changed methods reported by the compiler in the `presentation`
layer of the right context.

## Configuration

Environment variables, set in `docker-compose.yml`:

| Variable               | Required | Purpose                                              |
|------------------------|----------|------------------------------------------------------|
| `HTTP_ADDR`            | no       | Listening address (default `:8080`)                  |
| `REDIS_ADDR`           | yes      | Redis address, e.g. `redis:6379`                     |
| `COMMONS_USER_AGENT`   | yes      | Identifying User-Agent sent to Wikimedia Commons     |
| `CORS_ALLOWED_ORIGINS` | no       | Comma-separated browser origins allowed to call the API |

## Code layout

Simplified DDD ([ADR common 0004](../docs/adr-common/0004-bounded-contexts-and-clean-layers.md),
[ADR back core 0001](docs/adr-core/0001-go-layout-for-contexts-and-layers.md)):

```
api/                    API contract and code generation configuration
cmd/server/             HTTP service entry point
internal/
├── app/                composition root: configuration, wiring, HTTP server
├── bc/
│   ├── geography/      zones and their outlines
│   └── imagery/        images, Commons gateway, random draw
├── core/
│   ├── timespan/       decades, EDTF dates and plain English labels
│   └── zoning/         zones and the country → zone table
└── shared/
    ├── apigen/         code generated from the contract (never edited)
    ├── cache/          Redis cache with stale-on-error
    └── problem/        Problem Details responses
```

Inside each context: `domain` → `infrastructure` → `application` → `presentation`.
Contexts never import each other; `make lint` enforces these boundaries.

## Production image

```bash
docker build --target prod -t frise-chrono-back back/
```

A static binary on a distroless base image, run as a non-root user
([ADR back devops 0001](docs/adr-devops/0001-container-image-and-dev-reload.md)).

## Decisions

[docs/adr-core/](docs/adr-core/) (Go baseline), [docs/adr-functional/](docs/adr-functional/)
(Commons access, cache, zone geometry), [docs/adr-devops/](docs/adr-devops/) (image and
development reload).
