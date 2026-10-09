---
name: yorgut-architecture
description: Architecture layers and file-placement conventions for the Yorgut project (Go/Gin/GORM backend + React/MUI/TanStack Query frontend). Use when adding a feature or endpoint, deciding where a new file, type, query or component belongs, reviewing code for misplaced logic, or refactoring files that break the layering.
---

# Yorgut architecture

Yorgut is an Orkut clone: `backend/` (Go, Gin, GORM, Postgres in prod, SQLite in dev)
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
| DB | `internal/database/database.go` | Connection (sqlite when `DEV_MODE=true`, else postgres) and `AutoMigrate` list | — |
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
- **Swagger**: the spec (`backend/docs/`) is the contract the frontend client is generated from, so it must match what the handler really does. Every controller method has a `godoc` block with:
  - `@ID` (camelCase, unique; it becomes the generated function name: `findUserProfile` → `useFindUserProfile`, `getFindUserProfileQueryKey`),
  - `@Tags`, `@Summary`, `@Accept application/json` when it binds a JSON body, `@Produce application/json`,
  - one `@Param` per path param (`userId path int true "User ID"`, typed `int` because ids are uints), query param and body/formData field,
  - `@Success <status> {object} transport.<Response>` with the type the handler actually writes (`204` with no type for empty answers), `@Failure` only where the frontend cares,
  - `@Security BearerAuth` on everything behind `AuthMiddleware` (all but `/auth/google`),
  - `@Router` with the exact path from `routes.go` (without the `/api/v1` base path).
  Regenerate with `make gen-api` (spec **and** frontend client; needs `swag` v1.16.4: `go install github.com/swaggo/swag/cmd/swag@v1.16.4`). It runs swag with `--requiredByDefault`, so every field without `omitempty` is required in the spec. Commit `backend/docs/` and `frontend/src/api/generated/` together.
- **Nullable fields**: swag does not mark pointer fields (`*string`) as nullable, so they are typed non-null in the generated client even though the API sends `null`. Until the DTOs carry `extensions:"x-nullable"`, treat pointer-backed fields as possibly `null` in the UI (see known-deviations).
- **Schema changes**: add the model to `database.Migrate`; there are no SQL migrations in use.

### Recipe: adding a backend feature

Reference implementation: user reviews (`UserReview`).

1. `model/model.go`: struct + `New<Entity>` constructor.
2. `database/database.go`: append to `AutoMigrate`.
3. `transport/dto.go`: request/response types. `transport/mapper.go`: `From<Model>` methods that build the responses.
4. `repository/<entity>.go`: `type <Entity>Repository struct { repository[*model.<Entity>] }`, `New<Entity>Repository`, custom queries. Add it to `Repositories` / `NewRepositories` in `repository/repository.go`.
5. `service/<entity>.go`: struct with its repositories + `Validate` (+ `transactor` if a method writes more than one row), `New<Entity>Service`, methods taking primitives (or a `<Action>Input` struct when there are many). Permission checks go here. Add it to `Services` / `NewServices` in `service/service.go`. Tests go in `service/<entity>_test.go`, using `newTestRepositories(t)` (in-memory SQLite).
6. `controller/<entity>.go`: controller struct, `New<Entity>Controller`, handlers with godoc.
7. `server/server.go`: build the controller from `services`. `server/routes.go`: add the parameter to `NewRouter` and register routes.
8. `go build ./... && go vet ./... && go test ./...`, then `make gen-api` (regenerates the spec and the frontend client) and fix any frontend type errors it causes.

---

## Frontend (`frontend/src/`)

Imports are absolute from `src/` (`tsconfig.json` `baseUrl: "src"`), also inside a feature: `import { apiRequest } from "api/http"`, `import { MembersOnly } from "features/community/components/MembersOnly"`.

```
src/
  index.tsx, App.tsx, theme.ts, *.css   app shell
  routes/                               route table + ProtectedRoute
  api/
    http.ts                             apiRequest (orval mutator)
    generated/                          orval output, committed, never edited
    <domain>/queries.ts                 the hooks the UI uses (one file per backend resource)
  features/<feature>/
    pages/<Page>/index.tsx              route screens of the feature (parent layouts and nested routes)
    components/<Name>/index.tsx         sections only this feature uses
    <helper>.ts                         feature-only data/helpers (e.g. profile/profileFields.ts)
  components/<Name>/index.tsx           shared components (used by two or more features, or the app shell)
  contexts/, hooks/, helpers/, assets/
```

Features today: `auth` (login), `home`, `scrapbook` (`/scraps`), `profile` (`/profile/:userId/*`), `community` (`/community/:communityId/*`), `communities` (`/communities/*`), `friends`, `search`, `account`, `info` (help, what's new, media), `errors`.

| Layer | Location | Owns | Must not |
|---|---|---|---|
| App shell | `index.tsx`, `App.tsx`, `theme.ts`, `index.css`, `App.css` | Providers (QueryClient, MUI theme, Google OAuth, contexts), global styles, `palette` | — |
| Routing | `routes/index.tsx`, `routes/ProtectedRoute.tsx` | Route table, importing the `features/*/pages` screens; nested routes render via `<Outlet />` in the parent page | Fetch data |
| HTTP client | `api/http.ts` | `apiRequest`, the orval mutator every generated call goes through: base URL (`__API_BASE_URL__`), query string, auth header from `localStorage.token`, JSON vs `FormData` bodies, `401` → remove the token, error messages | Know about specific endpoints |
| Generated client | `api/generated/` (`endpoints/<tag>/<tag>.ts`, `model/*.ts`), produced by orval from `backend/docs/swagger.json` (`orval.config.ts`, `scripts/orval-transformer.mjs`) | Request functions (`findUserProfile`), query hooks (`useFindUserProfile`), query-key helpers (`getFindUserProfileQueryKey`) and every request/response type (`FindUserWithProfileByIdResponse`) | Be edited by hand: regenerate with `yarn gen:api` |
| API domain | `api/<domain>/queries.ts` (one per backend resource, shared by all features) | The hooks the UI uses: one `use<Thing>` per endpoint, wrapping a generated query hook with its options (`enabled`, `staleTime`, `placeholderData`) or a generated request function in `useMutation` with its cache invalidation | Declare request/response types (import them from `api/generated/model`); render anything |
| Feature pages | `features/<feature>/pages/<Page>/index.tsx` | Route-level screens and their layout, including parent pages with `<Outlet />` (`profile/pages/Profile`) and nested routes (`profile/pages/Scrapbook`) | Be imported by anything but `routes/` |
| Feature sections | `features/<feature>/components/<Name>/index.tsx` | Pieces used only inside that feature (`community/components/HotTopics`) | Be imported by another feature: move it to `components/` instead |
| Shared components | `components/<Name>/index.tsx` (flat) | Design-system pieces (`Button`, `BoxWithTitle`, `StripedGrid`, `PicAndText`, `ProfilePicture`), app-shell parts (`Wrapper`, `Navbar`, `UserBar`, `YorgutTitle`) and sections used by two or more features (`Scrapbook`, `FriendsList`, `HomeMenu`, `UserProfilePicture`). May call `api/*` hooks | Import from `features/` |
| Contexts | `contexts/<Name>Context.tsx` | Cross-cutting client state (auth, global loading) | Duplicate server state that TanStack Query already caches |
| Hooks | `hooks/use<Name>.ts` | Non-API reusable hooks (`useToken`, context accessors) | Wrap endpoints (those go in `api/`) |
| Helpers | `helpers/` (`helpers/formatting/*` for formatters) | Pure functions | Use React |
| Assets | `assets/{fonts,images}` | Static files imported by components | — |

Why this shape: a feature folder holds everything one area of the app needs, so most changes stay inside one folder and "is this used elsewhere?" has an answer (if it's in `features/x`, only `x` uses it). The old atoms/molecules/organisms/templates split was dropped: with ~20 shared components the level was a judgement call that didn't help find anything, and it let two different components both be called `ProfilePicture`. A flat `components/` with unique names is enough. API hooks stay in `api/` rather than in features because endpoints don't belong to one feature (`/me` feeds home, profile, communities and the shell).

### Conventions

- **Files**: one component per folder, `index.tsx`, named export (`export const ProfileHeader: React.FC<IProps>`); styled-components variants go in a sibling `styles.tsx`. Names in `components/` are unique; inside a feature they only need to be unique within the feature (`profile/components/ProfileHeader` and `community/components/ProfileHeader` coexist).
- **Promote, don't reach across**: when a second feature needs a section, move it to `components/` (with `git mv`) and, if the two uses differ only in data source, parameterise it (see `components/FriendsList`: `userId` omitted = the logged user).
- **API types** come from the generated client and keep the backend DTO names: `import { FindUserResponse } from "api/generated/model"`. Never hand-write a type for an API request or response; if the generated one is wrong, fix the controller's godoc (or the DTO) and regenerate. Local, non-API interfaces are prefixed with `I` (`IProps`).
- **API base**: `__API_BASE_URL__` (defined in `rsbuild.config.ts`, typed in `global.d.ts`) is added by `api/http.ts`; generated URLs are the spec paths (`/users/1/profile`).
- **Query keys** are orval's: `[path]` or `[path, params]` (e.g. `["/communities/search", { name, page, pageSize }]`). Mutations invalidate with the generated helpers in `onSettled`: `queryClient.invalidateQueries({ queryKey: getFindUserScrapsQueryKey(userId) })`. Calling a helper without its params (`getSearchCommunitiesQueryKey()`) matches every params variant, since keys match by prefix.
- **Who calls what**: components call the hooks in `api/<domain>/queries.ts`, never `api/generated/endpoints` directly and never `fetch`. The domain layer keeps mutation variables in the shape the UI uses (`mutate(payload)`, `mutate(photoId)`) instead of orval's `{ data }` objects, and is the one place to read a query's options or a mutation's invalidation.
- **Uploads**: multipart endpoints take a generated `*Body` type (`UploadPhotoBody`); the generated function builds the `FormData`, and `api/http.ts` drops the JSON `Content-Type` for it. Optional files are `undefined` (not `null`) when absent.
- **Generated code is committed** (`api/generated/`, plus `backend/docs/`), so the frontend builds and type-checks without Go or swag installed. Regenerate in the same change as the backend change; a regenerated client that doesn't compile is the point (the frontend must be updated too). ESLint ignores the folder.
- **Server state** lives in TanStack Query only.
- **Uploaded images**: the API returns upload paths relative to the backend (`storage/uploads/<uuid>`). Always render them through `uploadUrl()` from `helpers/uploadUrl`, which resolves them against the API origin. The shared image components (`PicAndText`, `PicAndTextHorizontal`, `ProfilePicture`) already do this; a raw `<img>` must call it itself.
- **Type-check stays clean**: `yarn typecheck` (`tsc --noEmit -p .`) has zero errors; keep it that way.
- **Dependencies**: React 18 with `@types/react` 18, TypeScript 5.9. Types, TypeScript, lint and codegen tooling are `devDependencies`. Yarn 1: change dependencies with `yarn add` / `yarn add -D` so `yarn.lock` stays in sync.
- **UI language** is Portuguese (pt-BR); code identifiers are English.

### Recipe: adding a frontend feature

Reference implementation: user reviews (`api/userReview/queries.ts`, `features/friends/components/FriendRating`).

1. Backend first: endpoint + godoc (see the backend recipe), then `make gen-api` in `backend/` (or `yarn gen:api` in `frontend/` when only the spec changed). This updates `api/generated/`; commit it.
2. `api/<domain>/queries.ts`: `use<Thing>` wrapping the generated query hook (add options here), and `use<Action><Thing>` = `useMutation` around the generated request function, invalidating the generated query keys it affects.
3. The UI goes in `features/<feature>/components/<Name>/` if one feature uses it, or in `components/<Name>/` if it is (or will be) shared. A new route screen goes in `features/<feature>/pages/<Page>/` and is registered in `routes/index.tsx`.
4. Mount it from the page or section that owns it.
5. Check with `yarn typecheck` (must stay at zero errors) and `npx rsbuild build`.

---

## Reviewing placement

When asked to audit or fix layering:

1. Read [known-deviations.md](known-deviations.md) for the current backlog, and the
   repo-root `ROADMAP.md` for larger planned changes (auth, transactions, tests, migrations).
2. Check a file against the "Must not" column of its layer.
3. Prefer small moves that keep behavior identical, one deviation per change, and update the backlog file.
