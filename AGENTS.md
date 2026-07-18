# AGENTS.md

## Quick Reference

- **Run**: `go run ./cmd/main.go` (needs `.env` with DSN, PORT, JWT_SECRET, FRONTEND_URL, Stripe keys, CLOUDINARY_URL)
- **Hot reload**: `air` (uses `.air.toml`, builds to `tmp/main`)
- **Build**: `go build -o ./tmp/main ./cmd/main.go`
- **No tests, linting, or CI configured**

## Module & Framework

- Module name is `gotickets` (not `go-ticket`). All imports use `gotickets/...`.
- **Echo v5** (`github.com/labstack/echo/v5`), not v4. Key differences: `echo.Context` is a pointer (`*echo.Context`), middleware signatures differ.
- **GORM** with PostgreSQL. Auto-migration runs on startup in `internal/server/http.go:31`.

## Architecture

Domain-driven, each domain in `internal/domain/<name>/` follows:

```
entity.go     — GORM model + ToResponse()
dto/          — request.go, response.go
repository.go — interface + gorm implementation
service.go    — business logic
handler.go    — HTTP handlers
register.go   — RegisterRoutes() wires DI and mounts routes
```

**Domains**: `user`, `event`, `bookings`
**Shared**: `internal/auth` (JWT), `internal/config` (env + DB + Cloudinary), `internal/middlewares`, `internal/httpResponse`, `internal/utils`, `internal/payment` (Stripe)

Entrypoint: `cmd/main.go` → loads config → connects DB → starts server.

## API Routes

All under `/api/v1/`. Auth via `Authorization: Bearer <jwt>`.

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| POST | `/api/v1/auth/register` | No | JSON body |
| POST | `/api/v1/auth/login` | No | JSON body, returns JWT |
| GET | `/api/v1/auth/me` | Yes | Returns current user |
| GET | `/api/v1/events` | No | List all |
| GET | `/api/v1/events/:id` | No | Get by ID |
| GET | `/api/v1/events/my-events` | Yes | Current user's events |
| POST | `/api/v1/events/create` | Yes | **multipart/form-data** (photo field) |
| PATCH | `/api/v1/events/:id` | Yes | **multipart/form-data**, owner only |
| POST | `/api/v1/bookings` | Yes | JSON body, creates Stripe checkout |
| GET | `/api/v1/bookings` | Yes | Current user's bookings |
| GET | `/api/v1/bookings/:id` | Yes | Get booking by ID |
| POST | `/webhook/stripe` | No | Stripe webhook (no /api/v1 prefix) |

## Gotchas

- **Event create/update use `multipart/form-data`**, not JSON — because of the photo upload field. Do not use `c.Bind()` for these; use `c.FormValue()` + `c.FormFile()`.
- **Dates are parsed in `Asia/Dhaka` timezone** (hardcoded in event handler). Format: `"2006-01-02 15:04:05"`.
- **JWT falls back to hardcoded `"secret_key"`** if `JWT_SECRET` is empty (`internal/auth/jwt.go:34`). Never rely on this in production.
- **Stripe webhook route** is at `/webhook/stripe` (no `/api/v1` prefix, no auth middleware).
- **Bookings flow**: create booking → Stripe checkout session → webhook confirms → tickets deducted. Expired webhooks delete the booking.
- **`BookingCode`** is a UUID prefixed with `"GT-"`. Generated via `uuid.New()`.
- **User password** is bcrypt-hashed on the entity itself (`entity.go`), not in the handler.
- **Response struct field `user/dto/response.go`**: `ID` serializes as `"Id"` (capital D) — inconsistent with other domains.

## Environment Variables

Required in `.env` (or system env):

| Variable | Purpose |
|----------|---------|
| `PORT` | Server port |
| `DSN` | PostgreSQL connection string |
| `JWT_SECRET` | JWT signing key |
| `FRONTEND_URL` | CORS allowed origin |
| `STRIPE_SECRET_KEY` | Stripe API key |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook verification |
| `STRIPE_SUCCESS_URL` | Redirect after successful payment |
| `STRIPE_CANCEL_URL` | Redirect after cancelled payment |
| `CLOUDINARY_URL` | Auto-read by Cloudinary SDK |
