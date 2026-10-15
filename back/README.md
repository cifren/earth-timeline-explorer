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
| `test`    | Unit tests: no network, Commons and Redis faked in-process |
| `test-integration` | Unit and integration tests: the real API on the Redis of the stack, Commons faked |
| `lint`    | Lint, check formatting and architecture boundaries         |
| `format`  | Format the code                                            |
| `api`     | Regenerate the server code from `api/openapi.yaml`         |
| `zones`   | Regenerate the zone outlines and country points from Natural Earth (needs network) |

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
| `OTEL_EXPORTER_OTLP_ENDPOINT` | no | OTLP/HTTP receiver of the traces, e.g. `http://jaeger:4318`; tracing is off when unset |
| `OTEL_SERVICE_NAME`    | no       | Service name shown in the traces                     |

## Traces and logs

Every request is traced with OpenTelemetry and sent to Jaeger
(http://localhost:16686, service `frise-chrono-back`). A trace shows the request with its
Redis commands and Wikimedia Commons calls, each with its duration. Logs written during
a request appear twice: in `./compose logs back`, with `trace_id` and `span_id`, and in
the trace, under the *Events* of the request step; an error log marks that step as
failed.

**Example — finding why a draw was slow:** open Jaeger, pick the service and the
operation `GET /api/v1/images`, sort by duration and open the slowest trace. Each Redis
command and Commons call is a step with its duration; a `serving stale cache entry`
event on the request step means Commons was down and the cache answered.

## How it works

**Example — drawing a photo of Europe in the 1950s** (`GET /api/v1/images?zone=europe&decade=1950`):

1. The request is checked: the decade must be on the timeline (1840 to the current
   decade, multiple of 10) and the zone must exist; otherwise the answer is a `400`
   Problem Details.
2. The back lists the countries having a `1950s photographs of <country>` category on
   Wikimedia Commons, keeps those placed in Europe (France, the Netherlands, the Soviet
   Union…) and draws one at random. This list is cached in Redis for 7 days and still
   served if Commons is down.
3. It asks Commons for a random file in that country's category tree. A missing file
   triggers a new draw, up to 3 attempts.
4. The file's free-text date ("Taken in 1953", "circa 1955", "03/04/1953") becomes a
   time span: an EDTF value (`1953-04-03`) and a plain English label ("3 April 1953");
   an unreadable date or one outside the decade falls back to the decade ("1950s").
5. The answer carries the image URLs on Wikimedia's servers, the zone, the source,
   license and author, and a "Learn more" link to Wikipedia. When Commons knows where
   the photo was taken, the location holds that point (GeoJSON) and links to it on
   GeoHack, OpenStreetMap and Google Maps; otherwise it holds a point inside the
   country, without links, and the granularity tells the difference.

**Zone outlines** are not computed at runtime: `make zones` downloads country borders
from Natural Earth, merges them zone by zone and generates a Go package compiled into
the binary, along with a point inside each country.

**Organisation**: the code follows a simplified DDD. Each business area (zones, images)
is a bounded context that never imports another; concepts shared by several contexts
(zones, decades) live apart; inside a context, the business model knows nothing of HTTP,
Redis or Commons. Generated code (server interfaces, zone outlines) lives in `var/` and
is never edited. The rules are in the [ADRs](#decisions) and checked by `make lint`.

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
