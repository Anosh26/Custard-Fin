# httpapi Adapter (`apps/api/internal/adapters/httpapi`)

The `httpapi` package provides the HTTP server, middlewares, and REST API handlers for the CustardFin Cloud Control API (`apid`).

## Public Endpoints
- `GET /healthz`: Liveness check (`{"status":"ok"}`).
- `POST /v1/auth/register`: User registration with bcrypt hashing, returns JWT access token.
- `POST /v1/auth/login`: User login, returns JWT access token.
- `POST /v1/nodes`: Pair/register an edge node (requires JWT `BearerAuth`), returns node credentials and FRP tunnel auth token.
- `PUT /v1/nodes/{id}/media/{mid}/manifest`: Upsert media manifest ($\le 3$ KB cap strictly enforced, requires `X-Node-Token`).
- `DELETE /v1/nodes/{id}/media/{mid}/manifest`: Unpublish media item (requires `X-Node-Token`).
- `GET /v1/catalog`: Public catalog search with query filtering (`q`) and pagination (`limit`, `offset`).

## Middlewares
- `LoggingMiddleware`: Logs request method, path, remote IP, HTTP status, and duration in ms.
- `CORSMiddleware`: Handles CORS headers and OPTIONS preflights.
- `RecoverMiddleware`: Catches panics and returns HTTP 500 JSON error without crashing.
- `RequireUserAuth`: Validates JWT token from `Authorization: Bearer <token>` and injects user claims into context.
- `RequireNodeAuth`: Validates `X-Node-Token` against node token hash and route path ID.

## How to Test
```bash
go test -v ./apps/api/internal/adapters/httpapi/...
```
