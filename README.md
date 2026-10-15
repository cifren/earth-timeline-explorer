# Frise Chrono

Frise Chrono is a web application to **travel through time and space on Earth**. The
user picks a zone of the world on a map and a decade on a timeline, and discovers
photographs of that place and period, fetched from public sources (Wikimedia Commons
first). Every image shows where it comes from and how reliable its date and location
are.

The repository holds every application of the system:

| Folder     | Content                                                         | Documentation |
|------------|-----------------------------------------------------------------|---------------|
| `back/`    | Go service: business rules, API contract, source gateways       | [back/README.md](back/README.md) |
| `web/`     | React web application                                           | [web/README.md](web/README.md) |
| `mobile/`  | Mobile application (later)                                      | — |
| `docs/`    | Architecture decision records shared by every application       | [docs/](docs/) |

## Requirements

- Docker Engine with the Compose v2 plugin (`docker compose`).
- Bash.

Nothing else is installed on the host: Go and Node run in containers.

All containers join the external Docker network **`jacquet`** (created automatically
when missing). Docker's default address range collides with the VPN; set
`JACQUET_NETWORK` to use another network name.

## Getting started

```bash
./compose init          # create .env and the network, build the images, install dependencies
./compose up            # start back, web, redis and jaeger in the background
```

`./compose init` creates `.env` from `.env.example` when it does not exist yet; set a
contact in its `COMMONS_USER_AGENT`.

| Service       | URL                                  |
|---------------|--------------------------------------|
| Web           | http://localhost:5173                |
| API           | http://localhost:8080/api/v1         |
| API reference | http://localhost:8080/api/docs       |
| Traces        | http://localhost:16686               |

Stop everything with `./compose down`.

## Everyday commands

Every Docker interaction goes through the `./compose` script
([ADR devops 0001](docs/adr-devops/0001-single-compose-entrypoint-for-docker.md)).
Application commands are listed in each application's Makefile and run inside its
container.

| Command                          | Action                                             |
|----------------------------------|----------------------------------------------------|
| `./compose help`                 | List every `./compose` command                     |
| `./compose up` / `down`          | Start / stop the stack                             |
| `./compose logs [service]`       | Follow the logs                                    |
| `./compose make back help`       | List the back commands                             |
| `./compose make web help`        | List the web commands                              |
| `./compose bash <service>`       | Open a shell in a container (`back`, `web`)        |

## Configuration

Variables read from `.env` by Docker Compose (created by `./compose init` from `.env.example`):

| Variable             | Purpose                                                              |
|----------------------|----------------------------------------------------------------------|
| `COMMONS_USER_AGENT` | Identifying User-Agent sent to Wikimedia Commons, with a contact, as required by the Wikimedia User-Agent policy |
| `JACQUET_NETWORK`    | Name of the external Docker network (default `jacquet`)              |

## Architecture

The system is organised in simplified DDD: bounded contexts shared by back and web
(`geography`, `chronology`, `imagery`), all business rules on the back, and one API
contract (`back/api/openapi.yaml`) from which server and clients are generated.

Decisions are recorded as short ADRs:

| Folder                                         | Scope                                          |
|------------------------------------------------|------------------------------------------------|
| [docs/adr-common/](docs/adr-common/)           | Conventions shared by every application        |
| [docs/adr-functional/](docs/adr-functional/)   | Whole-system decisions driven by the business  |
| [docs/adr-devops/](docs/adr-devops/)           | Whole-system build and run                     |
| [back/docs/](back/docs/)                       | Back decisions (core, functional, devops)      |
| [web/docs/](web/docs/)                         | Web decisions (core, functional, devops)       |

## License

The code is released under the [GNU General Public License v3.0](LICENSE). Images shown
by the application are not part of this repository: each keeps its own license, displayed
with its attribution.
