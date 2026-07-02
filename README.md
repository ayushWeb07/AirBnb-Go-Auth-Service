# 🔐 Auth Service & API Gateway

> Identity, access control, and the front door of **AirBnB-Node** — one Go service plays both roles: it's where users register, log in, and get authorized, and it's the single entry point that fronts every other service in the system.

The module name (`AirBnb-Go-Api-Gateway`) gives away the real shape of this service: it started as "the auth service" but has grown into the system's **API Gateway**, since authentication naturally has to sit in front of everything anyway. Requests for `/api/v1/users`, `/api/v1/roles`, and `/api/v1/user-roles` are handled directly here; requests for `/api/v1/hotels/*` and `/api/v1/bookings/*` are transparently reverse-proxied through to `HotelService` and `BookingService`.

---

## Where this fits in the system

```
                        Client
                          │
                          ▼
        ┌───────────────────────────────────┐
        │   Auth Service / API Gateway (Go)   │   ← this repo
        │                                       │
        │   /users, /roles, /user-roles         │→ handled locally
        │   /hotels/*        ─────────────────┐│
        │   /bookings/*      ─────────────────┼┼→ reverse-proxied
        └───────────────────────────────────┘│
                          │                     │
                 ┌────────▼────────┐   ┌────────▼─────────┐   ┌────────────────┐
                 │  MySQL            │   │  HotelService      │   │  BookingService  │
                 │  users/roles/      │   │  :3000              │   │  :3010            │
                 │  sessions/otps     │   └────────────────────┘   └────────────────────┘
                 └───────────────────┘
```

Every other service in AirBnB-Node treats this service's JWTs as the source of truth for "who is this" — `HotelService` and `BookingService` call back into it (or validate its tokens) rather than managing their own user tables.

---

## What it does

### Authentication
- **Registration & login** with bcrypt-hashed passwords.
- **Email verification via OTP.** New accounts can't log in until verified — a 10-character random OTP is generated, SHA-256 hashed before storage (the raw code is never persisted), emailed via SMTP using a Go `text/template`, and consumed once on `/verify-otp`.
- **Short-lived access tokens, long-lived refresh tokens.** Login issues a 15-minute HS256 JWT access token and a 7-day HS256 refresh token, signed with separate secrets so a leaked access token can't be used to mint new ones.
- **Session tracking, not just stateless JWTs.** Every refresh token is hashed (SHA-256) and stored server-side as a `Session` row. This means refresh tokens can be **revoked** — `/logout` revokes the specific session, `/logout-from-all-sessions` revokes every session for that user — something pure stateless JWTs can't do on their own.
- **Token refresh flow.** `/refresh-access-token` verifies the refresh token against its stored session hash (and checks it hasn't been revoked) before issuing a new short-lived access token.

### Authorization (RBAC)
- **Roles** (`admin`, `user`, `moderator`, etc.) are first-class resources with their own CRUD endpoints.
- **User–role assignment** is a many-to-many join, with endpoints to assign, remove, and query roles per user.
- **Middleware-enforced RBAC** — `RequireUserAllRoles` and `RequireUserAnyRoles` middlewares gate routes based on the authenticated user's roles, layered on top of `AuthMiddleware` (which decodes the JWT and puts the user ID on the request context for downstream handlers).

### Gateway
- **Reverse proxy to downstream services.** `httputil.NewSingleHostReverseProxy` forwards `/api/v1/hotels/*` to `HotelService` and `/api/v1/bookings/*` to `BookingService`, rewriting the `Host` so the upstream sees a normal request.
- **Global rate limiting.** Every route (including proxied ones) is rate-limited per IP + endpoint via `go-chi/httprate`, protecting downstream services from abuse without each of them needing their own limiter.

---

## Auth flow at a glance

```
POST /users/register ──► bcrypt-hash password, insert user (unverified)
                              │
POST /send-otp-for-verification ──► generate OTP, hash+store, email it
                              │
POST /verify-otp ──► compare hash, mark user verified, delete OTPs
                              │
POST /users/login ──► issue access token (15m) + refresh token (7d)
                              │                  hash+store session
              ┌───────────────┴────────────────┐
              ▼                                 ▼
   POST /refresh-access-token          POST /logout
   (validate session, reissue          (revoke this session)
    access token)
                                        POST /logout-from-all-sessions
                                        (revoke every session)
```

---

## Tech stack

| Concern | Choice |
|---|---|
| Language | Go 1.25 |
| Router | `go-chi/chi` |
| Auth | `golang-jwt/jwt` (HS256), `golang.org/x/crypto/bcrypt` |
| Rate limiting | `go-chi/httprate` |
| Validation | `go-playground/validator` |
| Database | MySQL (`database/sql` + `go-sql-driver/mysql`) |
| Email | `net/smtp` with Go `text/template` email bodies |
| Logging | `go.uber.org/zap` |
| Reverse proxy | `net/http/httputil` |
| Migrations | Goose-style timestamped SQL files |
| Live reload (dev) | [Air](https://github.com/air-verse/air) |

---

## Project structure

```
cmd/app/               # App composition: wires config, DB, router into an http.Server
main.go                # Process entrypoint — loads config, builds the App, runs it
internal/
├── config/               # Env loading, server config, DB config, zap logger setup
├── controllers/           # HTTP handlers for users, roles, user-roles, health
├── services/              # Business logic: auth flows, RBAC, mail sending
├── repositories/           # SQL queries against MySQL
├── database/
│   ├── models/               # Row structs: User, Role, UserRole, Session, Otp
│   └── migrations/            # Timestamped .sql migrations
├── dtos/                  # Request/response payload structs
├── middlewares/            # JWT auth, RBAC, rate limiting, body/param validation
├── routers/                # Route registration per resource + the proxy wiring
├── utils/                  # Errors, JSON responses, OTP generation, reverse proxy, email templates
```

The layering mirrors the Node services in this system for consistency: **router → controller → service → repository → model**, with interfaces (`UserServiceInterface`, `UserRepositoryInterface`, etc.) at each layer boundary so components are swappable/mockable in tests.

---

## API surface (v1)

### Users & auth — `/api/v1/users`

| Method | Route | Auth | Purpose |
|---|---|---|---|
| `POST` | `/register` | — | Create a new (unverified) user |
| `POST` | `/login` | — | Log in, get access + refresh tokens |
| `GET` | `/` | — | List users |
| `GET` | `/profile` | JWT + role (`admin`/`user`/`moderator`) | Get the authenticated user's profile |
| `GET` | `/{id}` | — | Get a user by ID |
| `PUT` | `/{id}` | — | Update a user |
| `DELETE` | `/{id}` | — | Delete a user |
| `POST` | `/send-otp-for-verification` | — | Send an email verification OTP |
| `POST` | `/verify-otp` | — | Verify a user's email with the OTP |
| `POST` | `/refresh-access-token` | — | Exchange a valid refresh token for a new access token |
| `POST` | `/logout` | — | Revoke the current session |
| `POST` | `/logout-from-all-sessions` | — | Revoke every session for the user |

### Roles — `/api/v1/roles`

| Method | Route | Purpose |
|---|---|---|
| `POST` | `/` | Create a role |
| `GET` | `/` | List roles |
| `GET` | `/{id}` | Get a role |
| `PUT` | `/{id}` | Update a role |
| `DELETE` | `/{id}` | Delete a role |

### User–role assignment — `/api/v1/user-roles`

| Method | Route | Auth | Purpose |
|---|---|---|---|
| `GET` | `/user/{user_id}` | — | List a user's roles |
| `POST` | `/assign` | JWT + `admin` | Assign a role to a user |
| `POST` | `/remove` | JWT + `admin` | Remove a role from a user |
| `POST` | `/check-single` | — | Check if a user has a specific role |
| `POST` | `/check-all` | — | Check if a user has *all* of a set of roles |
| `POST` | `/check-any` | — | Check if a user has *any* of a set of roles |

### Gateway / misc

| Route | Purpose |
|---|---|
| `GET /health` | Health check |
| `ANY /api/v1/hotels/*` | Reverse-proxied to `HotelService` |
| `ANY /api/v1/bookings/*` | Reverse-proxied to `BookingService` |

---

## Getting started

### Prerequisites

- Go 1.25+
- MySQL running locally or reachable
- [Air](https://github.com/air-verse/air) (optional, for live reload — config already in `.air.toml`)
- SMTP credentials for sending OTP verification emails
- `HotelService` on `:3000` and `BookingService` on `:3010` reachable, if you want proxying to actually work

### Configure environment

Create a `.env` file in the project root:

```bash
# Server
ADDR=:3020
READ_TIMEOUT=15
WRITE_TIMEOUT=15
IDLE_TIMEOUT=180
APP_ENV=development
REQUESTS_PER_MINUTE=100

# Auth secrets
ACCESS_TOKEN_SECRET_KEY=change-me
REFRESH_TOKEN_SECRET_KEY=change-me-too

# Database
DB_USERNAME=admin
DB_PASSWORD=admin
DB_NET=tcp
DB_ADDRESS=127.0.0.1:3306
DB_NAME=dev_db

# Goose (migrations)
GOOSE_DRIVER=mysql
GOOSE_DBSTRING=admin:admin@tcp(127.0.0.1:3306)/dev_db
GOOSE_MIGRATION_DIR=internal/database/migrations

# Mail (SMTP, used for OTP verification emails)
MAIL_FROM=your-address@example.com
MAIL_PASSWORD=your-smtp-password
MAIL_SMTP_HOST=smtp.example.com
MAIL_SMTP_PORT=587
```

> Note: `ADDR` defaults to `:3020`, while the reverse proxy targets in `routers/main.go` currently point at `http://localhost:3000` (HotelService) and `http://localhost:3010` (BookingService) — update those if your local ports differ. Also note `GOOSE_DBSTRING` needs a `:` between user and password (the code's own default value is missing it — `admin@admin@...` — so don't copy that literally).

### Run migrations

Migrations are plain timestamped `.sql` files under `internal/database/migrations`, managed with [Goose](https://github.com/pressly/goose). With Goose installed and `GOOSE_DRIVER`/`GOOSE_DBSTRING` set (or exported in your shell), either call it directly or use the bundled Taskfile:

```bash
# check current migration status
goose status
# or: task goose:status

# apply all pending migrations
goose up
# or: task goose:up-all

# roll back the most recent migration
goose down
# or: task goose:down-by-one
```

### Start the server

```bash
go run main.go
```

Or via the Taskfile:

```bash
task server:run
```

For live reload during development (uses the [Air](https://github.com/air-verse/air) config already checked into `.air.toml`):

```bash
air
# or: task server:dev
```

By default the server listens on `:3020` (override with `ADDR`). Once it's up, `GET /health` should return a 200 — that's the quickest way to confirm it started cleanly and connected to MySQL.

---

## Design notes worth knowing

- **Two secrets, two lifetimes.** Access and refresh tokens are signed with *different* secrets (`ACCESS_TOKEN_SECRET_KEY` / `REFRESH_TOKEN_SECRET_KEY`), so compromising one doesn't compromise the other, and they're validated against different code paths (`AuthMiddleware` for access tokens, `RefreshAccessToken`/`LogoutUser` for refresh tokens).
- **Refresh tokens are revocable because they're stateful.** The JWT itself is never stored — only its SHA-256 hash, alongside a `revoked` flag. This is what makes "logout" and "logout everywhere" possible despite using JWTs, which are normally stateless-by-default.
- **OTPs are hashed at rest, single-use.** Like refresh tokens, the raw OTP is never stored — only its hash — and successful verification deletes all outstanding OTPs for that user, closing the verification window rather than leaving old codes valid.
- **Auth and gateway concerns share one codebase for now.** Because JWT validation has to happen at the edge anyway, folding the gateway into the auth service avoids a second network hop just to check a token — the tradeoff is that this service now has two responsibilities (identity + routing) instead of one, worth watching as the system grows and more services need proxying.
- **RBAC checks hit the DB per request.** `RequireUserAllRoles`/`RequireUserAnyRoles` query the `user_roles` join table live rather than trusting roles embedded in the JWT — slightly more DB load, but role changes take effect immediately instead of waiting for a token to expire.