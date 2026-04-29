# Go Clean API Scaffold

A production-ready Go API scaffold built on **Clean Architecture** principles. Clone it, define your API in OpenAPI, and start building domain logic — the structure, wiring, and testing harness are already in place.

## What you get out of the box

- Strict Clean Architecture layer separation enforced by directory structure
- **Contract-first** API development via OpenAPI 3.0 + code generation
- Dependency injection wired automatically with Uber FX
- Structured logging with Uber Zap
- BDD integration tests using Gherkin/Cucumber syntax (godog)
- Unit tests with auto-generated mocks (mockery)
- Hot reload in development (Air)
- Docker Compose + DevContainer support
- Migration system with Up/Down support

## Technology Stack

| Concern | Library |
|---|---|
| HTTP framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) + SQLite |
| Dependency injection | [Uber FX](https://github.com/uber-go/fx) |
| API contract | [OpenAPI 3.0](https://swagger.io/specification/) + [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) |
| Configuration | [Viper](https://github.com/spf13/viper) |
| Logging | [Uber Zap](https://github.com/uber-go/zap) |
| BDD tests | [godog](https://github.com/cucumber/godog) |
| Mocks | [mockery](https://github.com/vektra/mockery) |
| Hot reload | [Air](https://github.com/air-verse/air) |

## Architecture

Bounded contexts are organised under `internal/app/<domain>_domain/<bounded_context>/`. The `<domain>_domain` folder groups related bounded contexts together; each `<bounded_context>` folder owns its full layer stack independently.

```
internal/app/hello_domain/hello/
├── domain/
│   ├── entity/                 # Core business entity (private fields, business rules)
│   └── repository/             # Repository interface (port, not implementation)
├── application/
│   └── usecase/                # Use case: orchestrates domain logic
├── adapter/
│   ├── controller/             # HTTP input adapter (receives request, calls use case)
│   └── presenter/              # Output adapter (maps use case output → HTTP response)
├── infrastructure/
│   └── repository/             # Repository implementation (GORM/SQLite)
├── mocks/                      # Auto-generated mocks — run `make mocks` to regenerate
└── provider.go                 # FX DI bindings for this bounded context
```

Dependency direction: `adapter → application → domain ← infrastructure`

The domain layer has zero dependencies on frameworks or the database. The `infrastructure` layer depends on the `domain` interface, not the other way around.

**Request flow** (end to end):

```
HTTP request
  → Gin router
  → generated strict handler (api/api.gen.go)
  → StrictServerStub (dispatch by method name)
  → Controller (parse input, call application)
  → Application (facade over use cases)
  → UseCase (business logic, calls domain + repository interface)
  → Domain entity (pure business rules)
       ↑
  Infrastructure repository (implements domain interface, talks to GORM/SQLite)
```

**Application layer** (`application/`): a thin facade that aggregates the use cases belonging to a bounded context. Controllers inject `HelloApplication` rather than individual use cases directly, which keeps the controller signature stable as the context grows and provides a single place to compose or sequence use cases.

**GORM model vs domain entity**: the domain entity (`domain/entity/`) has no ORM annotations and knows nothing about the database. ORM-specific models live in `internal/app/shared/model/` (or a domain-local `model/` if preferred). The infrastructure repository is responsible for translating between them in both directions — domain entity in, GORM model to DB; GORM model from DB, domain entity out.

## Getting Started

### Prerequisites

- Go 1.23+
- `make` (optional, but recommended)

### 1. Copy the config

```bash
cp config.example.yaml config.yaml
```

### 2. Run migrations

Always run migrations before starting the server for the first time, or after adding new migrations:

```bash
make migrate
```

### 3. Start the server

```bash
go run cmd/main.go
# or with hot reload:
make dev
```

### Run with Docker Compose

```bash
make up    # docker-compose up -d
make down  # docker-compose down
```

Migrations are **not** run automatically in Docker. Run them explicitly after the container is up:

```bash
docker exec go-clean-api-scaffold-api go run ./cmd/migrate/main.go
```

In production, run migrations as a separate step before deploying the new app version — a one-off `docker run` command, a Kubernetes Job, or a CI/CD pipeline step. Never auto-run migrations inside the app container itself: if multiple replicas start simultaneously they will race, and a failed migration will block the app from starting.

### Run with DevContainer (VS Code)

Open the repository in VS Code and select **Reopen in Container**. The container includes Go, Air, oapi-codegen, and mockery pre-installed.

### Remove the example domain

The `hello_domain/` directory is a reference implementation included to demonstrate the structure. Once you are ready to build your own project, delete it:

```bash
rm -rf internal/app/hello_domain
rm -rf tests/features/hello_domain
```

Also remove the corresponding entry from `app/provider/provider.go` and clean up the migration in `database/migration/migrate.go`.

## Development Workflow

### Define your API (OpenAPI-first)

Edit `api/openapi.yaml` to add or modify endpoints, then regenerate the server interface:

```bash
make api
# runs: oapi-codegen --config=gen.config.yaml ./api/openapi.yaml
# output: api/api.gen.go
```

The generated file contains the strict server interface your controllers must implement. The compiler enforces that every endpoint defined in the spec has a corresponding handler.

### Start the dev server with hot reload

```bash
make dev
# runs: air
```

## Make Commands

| Command | Description |
|---|---|
| `make dev` | Start server with hot reload (Air) |
| `make build` | Compile to `bin/app` |
| `make api` | Regenerate server code from `api/openapi.yaml` |
| `make mocks` | Regenerate all mocks with mockery |
| `make migrate` | Run database migrations |
| `make test:unit` | Run unit tests (`./internal/...`) |
| `make test:feature` | Run BDD feature tests |
| `make up` | Start Docker Compose services |
| `make down` | Stop Docker Compose services |

## Key Design Decisions

### StrictServerStub — auto-discovery controller dispatch

`StrictServerStub` is generated from a custom template (`api/templates/strict-interface.tmpl`). When controllers are registered via FX, `InitControllers` uses reflection to scan every controller's exported methods and builds a name-to-function dispatch map. Each endpoint on the stub looks up the method by name and delegates to it:

```go
// In app/server.go — controllers are collected by FX and passed here
func NewServer(ctls []Controller) api.ServerInterface {
    server := &ServerImpl{}
    server.InitControllers(ctls)
    return api.NewStrictHandler(server, nil)
}
```

This means: adding a new controller just requires registering it in your domain's `provider.go`. The stub automatically discovers and routes to its methods. Endpoints with no registered controller return a `"not implemented"` error at runtime instead of failing to compile.

### Presenter — separate response format from business logic

Controllers handle only input: parse the request, call the use case, pass the result to the presenter. Presenters handle only output: map use case outputs and errors to HTTP response objects.

This separation means:
- **Error mapping lives in one place** — the presenter decides whether a domain error becomes a 400 or 500
- **Controllers stay thin** — no `if err == ErrX { return 400 }` chains
- **Response format is independently replaceable** — swap JSON for another format without touching controller or use case code

### Shared response format

`api/api_shared_response.go` is a handwritten file. The custom template (`api/templates/strict-interface.tmpl`) generates `VisitXxxResponse` methods on these two types, making them usable as valid response objects for any endpoint:

```go
// api/api_shared_response.go
type SharedAPIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
}

type SharedAPIErrorResponse struct {
    StatusCode int      `json:"-"`
    Details    []string `json:"details"`
    Message    string   `json:"message"`
    Success    bool     `json:"success"`
}
```

Per-endpoint response types (`SayHello200JSONResponse`, `SayHello400JSONResponse`, `SayHello500JSONResponse`) are **generated** by `make api` from the spec. The error variants are type aliases of the spec-defined shared schemas (`Shared.APIBadRequestResponse`, `Shared.APIInternalServerErrorResponse`) which also live in `api.gen.go`. Presenters use these generated types directly for full type safety:

```go
// adapter/presenter/json/say_hello_presenter.go
func (p *SayHelloPresenter) Output(output *usecase.SayHelloOutput) api.SayHelloResponseObject {
    return api.SayHello200JSONResponse{Success: true, Data: api.HelloGreeting{Message: output.Message}}
}

func (p *SayHelloPresenter) Error(err error) api.SayHelloResponseObject {
    switch {
    case errors.Is(err, usecase.ErrNameRequired):
        return api.SayHello400JSONResponse{Success: false, Message: "Bad Request", Details: []string{err.Error()}}
    default:
        return api.SayHello500JSONResponse{Success: false, Message: "Internal Server Error"}
    }
}
```

**To add or change response fields for a specific endpoint**, update `api/openapi.yaml` and run `make api`. The generated structs will reflect the new shape automatically, and the compiler will tell you exactly which presenter methods need updating.

**To change a shared error schema** (e.g. add a `code` field to all 400 responses), update both `components/schemas/Shared.APIBadRequestResponse` in the spec and its corresponding struct in `api/api_shared_response.go`, then run `make api`.

### Migration system

Migrations are an ordered slice of `Migrator` interface values, each with `Up` and `Down` methods:

```go
// database/migration/migrate.go
var Migrations = []Migrator{
    migrations.GreetingMigrator{},
}
```

This design is intentionally simple: no version-tracking table, no state. `RunMigrate` always runs every migration in order; `RunRollback` runs them in reverse. In feature tests, the DB is reset and re-migrated before each scenario via the `Before` hook.

**Trade-off**: this works well for development and test resets, but it does not track which migrations have been applied in production. If you need versioned, incremental migrations for a production database, consider replacing this with [golang-migrate](https://github.com/golang-migrate/migrate) or [goose](https://github.com/pressly/goose) — both support the same `Up`/`Down` pattern and can coexist with GORM. Or, if you prefer a Laravel-inspired fluent API with version tracking, the author of this scaffold also maintains [go-migrate](https://github.com/laijunbin/go-migrate) — though fair warning, it currently targets MySQL, so factor that in if SQLite is staying.

### CORS

CORS is configured in `app/server.go`. Allowed origins are read from config:

```yaml
# config.yaml
server:
  allowedOrigins:
    - http://localhost:5173
    - https://your-production-domain.com
```

Allowed methods (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`) and headers (`Content-Type`) are set in `app/server.go` and can be adjusted there.

### Custom oapi-codegen templates

`api/templates/` contains two custom templates that override oapi-codegen defaults:

| Template | Purpose |
|---|---|
| `strict-interface.tmpl` | Generates `StrictServerInterface`, `StrictServerStub`, and wires `SharedAPIResponse`/`SharedAPIErrorResponse` as valid response types for every endpoint |
| `strict-gin.tmpl` | Generates the Gin-specific strict handler that dispatches requests to the stub |

These templates are referenced in `gen.config.yaml`. Modify them if you need to change code generation behaviour across all endpoints at once.

## Adding a New Domain

Each bounded context owns its own entities, repository interfaces, use cases, and infrastructure implementations. Cross-context communication should go through use case interfaces — never by directly accessing another context's entities or repositories.

The example below adds a bounded context `order` inside `order_domain`. If your project has multiple related bounded contexts (e.g. `order` and `fulfillment`), they can share the same `<domain>_domain` folder while remaining independently structured inside.

The implementation order within the layers is up to you — feel free to apply TDD, start from the domain inward, or prototype from the controller down. The three things that must be in place for the server to wire up correctly are covered below.

### 1. Define your API contract first

Add the new endpoints to `api/openapi.yaml`, then regenerate:

```bash
make api
```

This gives you the typed request/response structs and the server interface your controller must satisfy before you write a single line of business logic.

### 2. Implement the four layers

Create the bounded context under `internal/app/order_domain/order/`. The directory layout below matches the `hello` example and is a starting point — adapt it as needed (add `domain/service/`, `application/service/`, etc.). What matters is that dependencies flow inward toward the domain.

```
internal/app/order_domain/order/
├── domain/
│   ├── entity/order/
│   └── repository/
├── application/
│   └── usecase/
├── adapter/
│   ├── controller/
│   └── presenter/
├── infrastructure/
│   └── repository/
└── provider.go
```

### 3. Register with FX

Declare the bounded context's DI bindings in `provider.go`:

```go
func NewProvider() *provider.Provider {
    return provider.New(provider.Config{
        Infrastructures: []interface{}{repository.NewOrderRepository},
        Controllers:     []interface{}{controller.NewOrderController},
        Presenters:      []interface{}{presenter.NewOrderPresenter},
        Usecases:        []interface{}{usecase.NewCreateOrderUseCase},
        Applications:    []interface{}{application.NewOrderApplication},
    })
}
```

Then add it to the app in `app/provider/provider.go`:

```go
Providers = []*provider.Provider{
    shared.NewProvider(),
    hello.NewProvider(),
    order.NewProvider(), // add this line
}
```

### Other things to add as needed

- **Migration**: add a migrator in `database/migration/migrations/` and register it in `database/migration/migrate.go`
- **Mocks**: run `make mocks` after defining or changing interfaces — it auto-discovers all `*_domain/*/` directories, so no Makefile changes needed when adding new domains
- **Feature tests**: add `.feature` files under `tests/features/order_domain/order/features/`

## Project Layout

```
.
├── api/                        # OpenAPI spec + generated server code + custom templates
├── app/                        # Server bootstrap and FX provider aggregation
├── cmd/
│   ├── main.go                 # Application entrypoint
│   └── migrate/main.go         # Migration entrypoint
├── database/
│   └── migration/              # Migration definitions (Up/Down per table)
├── internal/
│   └── app/
│       ├── shared/             # Shared infrastructure (logger, base types)
│       └── <domain>_domain/    # One directory per domain
└── tests/
    ├── features/               # BDD feature files and test entry points
    ├── setup/                  # Test server and infrastructure setup
    └── testutils/              # Step definitions and test helpers
```

The structure is a starting point, not a constraint. The `provider.Config` type supports additional slots beyond what the `hello` example uses:

```go
type Config struct {
    Infrastructures []interface{}
    Controllers     []interface{}
    Presenters      []interface{}
    Services        []interface{}      // domain services
    Usecases        []interface{}
    Applications    []interface{}      // application services
    Factories       []interface{}
    Handlers        []interface{}
    Providers       []interface{}
    Invokes         []interface{}
}
```

The `Services` and `Applications` slots are there when you need them — the `hello` example intentionally leaves them empty to keep the demo minimal. Where you put the corresponding files is up to you; common conventions include `domain/service/` for domain services and `application/service/` for application services, but the scaffold does not enforce a path. The `mocks/` directory is auto-generated — run `make mocks` whenever an interface changes.

`Handlers` and `Invokes` both become `fx.Invoke` calls at startup (side effects, not provided values). The split is semantic: use `Handlers` for functions that register HTTP routes or middleware, and `Invokes` for general startup logic such as background workers or seed data. Both are collected in `app/provider/provider.go` and passed to FX together.

## Testing

### Unit tests

Unit tests live alongside the code they test (e.g. `application/usecase/say_hello_usecase_test.go`). Mocks are auto-generated by mockery and committed to the repo. Regenerate them whenever an interface changes:

```bash
make mocks        # regenerate all mocks
make test:unit    # run unit tests
```

### BDD feature tests

Feature tests are written in Gherkin and live under `tests/features/`. The test harness:

1. Starts the real server (in-process) on a separate port
2. Resets and migrates the test SQLite database before each scenario
3. Sends actual HTTP requests and asserts against the full JSON response

Example feature file (`tests/features/hello_domain/hello/features/say_hello.feature`):

```gherkin
Feature: Say Hello API

  Scenario: Successful greeting
    When I send a "GET" request to "/api/v1/hello?name=Alice":
    Then the response code should be 200
    And the response should match json:
      """
      {
        "success": true,
        "data": { "message": "Hello, Alice!" }
      }
      """

  Scenario: Missing name
    When I send a "GET" request to "/api/v1/hello":
    Then the response code should be 400
```

Run feature tests:

```bash
make test:feature
```

The harness (`tests/test.go`) passes `-p 1` to `go test`, forcing all feature packages to run serially. This is intentional — all scenarios share the same running server instance, so parallel package execution would cause port conflicts and DB race conditions. Keep this in mind when adding new feature packages: each package runs sequentially, but scenarios within a package also run one at a time.

The project includes `testcontainers-go` as a dependency. It is unused for SQLite (which needs no container), but the infrastructure is already wired in `tests/setup/setup.go` via the `containers` map. When you swap SQLite for a containerized database such as PostgreSQL or MySQL, spin up the container in `SetupInfrastructure()` and register it there — the teardown in `TerminateInfrastructure()` will handle cleanup automatically.

## Configuration

Copy `config.example.yaml` to `config.yaml` and adjust as needed:

```yaml
env: development  # development | release | production

server:
  port: 8080
  allowedOrigins:
    - http://localhost:5173

database:
  sqlite:
    path: ./hello.db

log:
  level: info  # debug | info | warn | error

test:
  server:
    port: 18080
  database:
    sqlite:
      path: ./test.db
```

All keys can be overridden via environment variables using `_` as a separator (e.g. `SERVER_PORT=9090`, `DATABASE_SQLITE_PATH=./data.db`).

The `Config` struct in `internal/config/config.go` can be extended with new sections as your application grows. Add a new struct, embed it in `Config`, and Viper will automatically map the corresponding YAML keys.

---

*README drafted with [GitHub Copilot](https://github.com/features/copilot) and [OpenAI Codex](https://openai.com/index/openai-codex/), refined with [Claude Code](https://claude.ai/code)*
