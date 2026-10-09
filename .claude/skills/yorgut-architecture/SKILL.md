---
name: yorgut-architecture
description: Architecture layers and file-placement conventions for the Yorgut project (Go/Gin/GORM backend + React/MUI/TanStack Query frontend). Use when adding a feature or endpoint, deciding where a new file, type, query or component belongs, reviewing code for misplaced logic, or refactoring files that break the layering.
---

# Yorgut architecture

Yorgut is an Orkut clone: `backend/` (Go, Gin, GORM, Postgres in prod and optionally in dev, SQLite for dev without Docker and tests)
and `frontend/` (React, MUI, TanStack Query, rsbuild). This skill describes the
**intended** layering. The code does not always follow it: see
[known-deviations.md](known-deviations.md) before treating an existing file as a
pattern to copy.

When placing new code, follow the rules below even if a neighbouring file
does otherwise. When fixing a deviation, update `known-deviations.md`.

---

## Backend (`backend/`)

### Request flow

```
HTTP → AuthMiddleware → routes.go → Controller → Service → Repository → GORM/DB
```

What crosses each boundary:

```
            binds Request                 primitives / service input            model
HTTP ──────────────────────→ Controller ─────────────────────────────→ Service ───────→ Repository
HTTP ←────────────────────── Controller ←───────────────────────────── Service ←─────── Repository
            Response (JSON)                    Response struct                   model
```

- **Requests** (`transport.*Request`) stop at the controller. The controller binds them,
  validates them, and passes their fields to the service as primitives or a
  service-owned input struct.
- **Services** never take `transport.*Request` types or anything from `net/http`,
  `multipart` or `gin`. They return response structs built with the `transport` mappers.
  Separate service result types are deliberately not used yet; introduce them
  only when one service result has to feed two different responses.
- **Models** only travel between services and repositories, except for the
  logged user (see below).

| Layer | Location | Owns | Must not |
|---|---|---|---|
| Entry points | `cmd/api/main.go`, `cmd/seed/main.go` | Process start, timezone, graceful shutdown, seeding | Contain business logic |
| Wiring | `repository.NewRepositories(db)` and `service.NewServices(repos, validate)` build everything; `internal/server/server.go` builds the controllers from them | Manual DI. `cmd/seed` uses the same two aggregates, so the wiring can't drift | Contain logic |
| Routing | `internal/server/routes.go` | URL → controller method mapping, CORS, swagger | Contain handlers inline (except `/api` health) |
| Middleware | `internal/server/middleware.go` | `AuthMiddleware`: reads the `Authorization: Bearer` Yorgut session token, resolves the logged user via `AuthService.UserFromSession` into `ctx.Set("user", *model.User)` | Do anything that is not request-wide; call Google |
| Controller | `internal/controller/<entity>.go` (one file per entity, snake_case); shared handler helpers (`notImplemented`) in `controller/controller.go` | Bind and validate the request, read path params, get logged user via `ctx.MustGet("user").(*model.User)`, turn request fields into service arguments (including reading uploaded files into `[]byte`), call the services it needs, choose HTTP status, swagger `godoc` comments | Touch repositories or GORM; contain business rules; pass a `Request` struct to a service |
| Service | `internal/service/<entity>.go`: struct, `New<Entity>Service`, and methods in one file. `service.go` only has the `Services` aggregate | Business rules, permission checks (e.g. "only friends"), orchestration across repositories, input structs for methods with many parameters (e.g. `service.UpdateUserInput`, no `json`/`binding` tags), returning response structs (via `transport` mappers), sentinel errors (`var ErrXNotFound = errors.New(...)`) | Import `gin`, `net/http` or `mime/multipart`; take `transport.*Request` parameters; write SQL; return models to controllers |
| Repository | `internal/repository/<entity>.go`: struct, `New<Entity>Repository`, custom queries and query-projection structs (e.g. `UserReviewAverages`). `repository.go` has the generic `repository[T]` (`Create`, `FindById`, …), the `Repositories` aggregate and the `Transactor` | All GORM queries, preloads, upserts, aggregates, and starting transactions (`Transactor`). Lookups by id use `First` so a missing row returns `gorm.ErrRecordNotFound` | Contain business rules or return DTOs |
| Model | `internal/model/model.go` | GORM entities (`gorm.Model` embedded) + `New<Entity>(...)` constructors | Have JSON tags or HTTP concerns |
| HTTP contract (`transport`) | `internal/transport/dto.go` (types), `internal/transport/mapper.go` (all model → response mapping, as `From<X>` methods) | The API's request and response shapes, with camelCase `json` tags and `binding` validation. Types ending in `Item` (`ScrapItem`, `UserItem`) are nested pieces of responses | Import services/repositories; be used as service inputs |
| Helpers | `internal/helpers/` | Stateless utilities: parsing (`MustParseUint`), file upload, session-token signing/parsing (`NewSessionToken`, `ParseSessionToken`) | Hold domain logic |
| DB | `internal/database/`: `database.go` (connection: sqlite when `DEV_MODE=true`, else postgres), `migrate.go` (`Migrate` picks the schema path: `MigratePostgres` runs the goose SQL files in `migrations/`, `AutoMigrateSQLite` runs GORM AutoMigrate over `Models()`), `reset.go` (localhost-only reset for the seed) | Run AutoMigrate on Postgres; edit a migration that already ran in production |
| Seed | `internal/seed/` + `cmd/seed/main.go` | Fake data and images for local dev | Be imported by the API |

### Authentication

- Google is the only identity provider (no passwords, no email flows, on purpose).
- `POST /api/v1/auth/google` (public) takes the Google ID token, `AuthService.LoginWithGoogle` validates it once against `GOOGLE_CLIENT_ID`, finds the user by `GoogleSub` (or links an existing one by email, or creates one) and returns a Yorgut JWT (HS256, `SESSION_SECRET`, 7 days).
- Every other route is in the `/api/v1` group behind `AuthMiddleware`, which only verifies that JWT and loads the user by id.
- The server refuses to start without `SESSION_SECRET` or `GOOGLE_CLIENT_ID`.
- Frontend: `api/auth` stores the token in `localStorage.token`; `api/http.ts` sends it and removes it on `401`.

### Conventions

- **Logged user** is always taken from the context (`ctx.MustGet("user")`), never from the request body.
- **Naming**:
  - Model: `Scrap`. Repository: `ScrapRepository`. Service: `ScrapService`. Controller: `ScrapController`.
  - Transport types: `Create<X>Request`, `Update<X>Request`, `Find<X>Response`, `<X>Item` for nested pieces of a response.
  - Service inputs: plain arguments, or `<Action>Input` structs in the service package.
  - Repository query methods: `Find…Of<Owner>`, `Count…`, `…By<Field>`.
- **Constructors**: every repository, service and controller has `New<Name>(deps...)`. Constructors must assign **every** field, including `Validate`.
- **Fields**: repository and service fields on structs are unexported (`scrapRepository`, `friendshipService`); only `Validate` is exported.
- **Dependencies are passed by value** (`*services.User` dereferenced in `server.go`); repositories only hold a `*gorm.DB`, so copies are cheap.
- **Transactions**: a service method that writes more than one row does it inside `s.transactor.Transaction(func(tx *repository.Repositories) error { ... })`, using only `tx.<Repo>` inside the callback. Returning an error rolls everything back. Services that need it take `repository.Transactor` as their first constructor argument (`r.Transactor` in `NewServices`). Reads and validation can happen before the transaction; file uploads can't be rolled back, so do them before it.
- **Denormalized counters** (`User.FriendCount`, `Community.MemberCount`) are never incremented in Go. They are recomputed in SQL (`UserRepository.UpdateFriendCount`, `CommunityRepository.UpdateMemberCount`) inside the same transaction as the change.
- **One service per entity**: a method lives in the service of the entity it changes or reads (`JoinCommunity` → `CommunityMembershipService`). Controllers depend on as many services as they need (e.g. `MeController` takes five).
- **Errors in controllers** use `ctx.JSON(status, gin.H{"error": err.Error()})` and **always `return`** after writing an error. Map service sentinel errors with `errors.Is` (e.g. `service.ErrFriendshipNotFound` → `404`).
- **Success statuses**: `200` with a body for reads, `201` for creations returning `transport.Response`, `204` with no body for actions, `501` via `notImplemented(ctx)` for routed-but-unwritten handlers (never `panic`).
- **Validation**: `transport.*Request` types use Gin `binding:"..."` tags, checked by `ctx.ShouldBindJSON`; use pointer fields (`*uint8`) when `0` is a valid value with `required`.
- **Routes**: groups are plural kebab-case (`/communities`, `/community-memberships`, `/photo-albums`), and every path starts with `/`. Resources nested under the owning user use `/users/:userId/<thing>`; "my" views of the logged user go under `/me/...`. Renaming a route means updating `frontend/src/api` in the same change.
- **Logged user exception**: `AuthService.UserFromSession` returns a `*model.User` because the middleware stores the model in the context; no other service method returns models.
- **Swagger**: every controller method has a `godoc` block; regenerate with `make swagger` (needs `swag`).
- **Schema changes** need two things that must agree: the model change (SQLite and the tests use `AutoMigrate`) and a new goose migration in `internal/database/migrations/` (`make migration_new name=...`), written in Postgres SQL that matches what GORM would create (types, `idx_<table>_<column>` index names). Migrations are the source of truth for production and run at API startup. A new model also goes in `database.Models()`. Check with `make test_pg`, which compares the migrated schema with AutoMigrate on a real (embedded) Postgres. Never edit a migration that already ran in production. Details in `backend/README.md` ("Database").

### Recipe: adding a backend feature

Reference implementation: user reviews (`UserReview`).

1. `model/model.go`: struct + `New<Entity>` constructor.
2. Schema: a new model goes in `Models()` in `database/migrate.go`; then `make migration_new name=<change>` and write the Postgres SQL for it in the new `database/migrations/000NN_<change>.sql`. Run `make test_pg`.
3. `transport/dto.go`: request/response types. `transport/mapper.go`: `From<Model>` methods that build the responses.
4. `repository/<entity>.go`: `type <Entity>Repository struct { repository[*model.<Entity>] }`, `New<Entity>Repository`, custom queries. Add it to `Repositories` / `NewRepositories` in `repository/repository.go`.
5. `service/<entity>.go`: struct with its repositories + `Validate` (+ `transactor` if a method writes more than one row), `New<Entity>Service`, methods taking primitives (or a `<Action>Input` struct when there are many). Permission checks go here. Add it to `Services` / `NewServices` in `service/service.go`. Tests go in `service/<entity>_test.go`, using `newTestRepositories(t)` (in-memory SQLite).
6. `controller/<entity>.go`: controller struct, `New<Entity>Controller`, handlers with godoc.
7. `server/server.go`: build the controller from `services`. `server/routes.go`: add the parameter to `NewRouter` and register routes.
8. `go build ./... && go vet ./... && go test ./...`, then `make swagger`.

---

## Frontend (`frontend/src/`)

Imports are absolute from `src/` (`tsconfig.json` `baseUrl: "src"`): `import { doGet } from "api/http"`.

| Layer | Location | Owns | Must not |
|---|---|---|---|
| App shell | `index.tsx`, `App.tsx`, `theme.ts`, `index.css`, `App.css` | Providers (QueryClient, MUI theme, Google OAuth, contexts), global styles, `palette` | — |
| Routing | `routes/index.tsx`, `routes/ProtectedRoute.tsx` | Route table; nested routes render via `<Outlet />` in the parent page | Fetch data |
| HTTP client | `api/http.ts` | `doGet`/`doPost`/`doPut`/`doDelete`, auth header from `localStorage.token`, error mapping | Know about specific endpoints |
| API domain | `api/<domain>/queries.ts` + `api/<domain>/types.ts` | One `use<Thing>` hook per endpoint (TanStack `useQuery`/`useMutation`), cache invalidation, request/response interfaces | Render anything; be called with raw `fetch` from elsewhere |
| Design system | `components/{atoms,molecules,organisms,templates}/<Name>/index.tsx` | Atomic-design components reused across **two or more pages** | Import from `pages/` |
| Pages | `pages/<Page>/index.tsx` | Route-level screen; layout of its sections | Be imported by another page |
| Sub-routes | `pages/<Page>/<SubRoute>/index.tsx` | Nested route screens (e.g. `pages/Profile/Scrapbook`) | — |
| Page sections | `pages/<Page>/components/<Name>/index.tsx` | Pieces used only by that page (and its sub-routes) | Be imported by other pages; promote to `components/` instead |
| Contexts | `contexts/<Name>Context.tsx` | Cross-cutting client state (auth, global loading) | Duplicate server state that TanStack Query already caches |
| Hooks | `hooks/use<Name>.ts` | Non-API reusable hooks (`useToken`, context accessors) | Wrap endpoints (those go in `api/`) |
| Helpers | `helpers/` (`helpers/formatting/*` for formatters) | Pure functions | Use React |
| Assets | `assets/{fonts,images}` | Static files imported by components | — |

### Conventions

- **Files**: one component per folder, `index.tsx`, named export (`export const ProfileHeader: React.FC<IProps>`); styled-components variants go in a sibling `styles.tsx`.
- **Types**: interfaces prefixed with `I` (`IScrapsResponse`, `ICreateScrapPayload`); `Response` for API reads and `Payload`/`Request` for request bodies.
- **API base**: `const basePath = \`${__API_BASE_URL__}/<resource>\`` (the global is defined in `rsbuild.config.ts`, typed in `global.d.ts`).
- **Query keys** are the full URL string: `queryKey: [\`${basePath}/${id}\`]`. Mutations invalidate exactly the URLs they affect in `onSettled`.
- **Server state** lives in TanStack Query only; components call `api/*` hooks, never `fetch`.
- **Uploaded images**: the API returns upload paths relative to the backend (`storage/uploads/<uuid>`). Always render them through `uploadUrl()` from `helpers/uploadUrl`, which resolves them against the API origin. The image atoms (`PicAndText`, `PicAndTextHorizontal`, `ProfilePicture`) already do this; a raw `<img>` must call it itself.
- **Atomic levels**:
  - atoms: no data fetching, purely presentational.
  - molecules: small compositions, may use one query hook.
  - organisms: sections like `Navbar` or `Scrapbook`.
  - templates: page frames like `Wrapper`.
- **UI language** is Portuguese (pt-BR); code identifiers are English.

### Recipe: adding a frontend feature

Reference implementation: user reviews.

1. `api/<domain>/types.ts`: response/payload interfaces.
2. `api/<domain>/queries.ts`: `use<Thing>` query + `use<Action><Thing>` mutation with invalidation.
3. The UI goes in `pages/<Page>/components/<Name>/` if one page uses it, or in `components/<level>/<Name>/` if it is (or will be) shared.
4. Mount it from the page or section that owns it.
5. Check with `npx tsc --noEmit -p .` (filter for your files; some errors already exist) and `npx rsbuild build`.

---

## Reviewing placement

When asked to audit or fix layering:

1. Read [known-deviations.md](known-deviations.md) for the current backlog, and the
   repo-root `ROADMAP.md` for larger planned changes (auth, transactions, tests, migrations).
2. Check a file against the "Must not" column of its layer.
3. Prefer small moves that keep behavior identical, one deviation per change, and update the backlog file.
