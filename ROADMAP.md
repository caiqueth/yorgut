# Architecture roadmap

Improvements that came out of an architecture review (2026-10-09). The stack
(Go API + React SPA) stays; this is about fixing what is underneath it.

Smaller layering and naming fixes live in
[`.claude/skills/yorgut-architecture/known-deviations.md`](.claude/skills/yorgut-architecture/known-deviations.md).
Tick items off here as they land, and remove a section once it's fully done.

Decisions already made:

- **No migration to Next.js.** The app is behind a login, so SSR/SEO buys little, and it would add a Node server and force an auth rework.
- **No rewrite in Django.** The "front and back aren't glued together" pain is mostly duplicated types; solved by item 8 below.
- **The frontend stays a React SPA.** No server-side rendering (Next.js, templ/htmx or otherwise).

---

## Now

### 1. Own sessions instead of using the Google ID token as a bearer ✅ (code done 2026-10-09, not yet tested with a real Google login or deployed)

Google validates the user once at `POST /api/v1/auth/google`; from then on the frontend sends a Yorgut JWT
(HS256, `SESSION_SECRET`, 7 days) as `Authorization: Bearer`. `AuthMiddleware` only verifies it and loads the
user by id. Users are keyed by `GoogleSub` (unique), with `Email` unique too; existing users get linked by
email on their next login. Details in the architecture skill ("Authentication").

Decisions: header + `localStorage` instead of httpOnly cookies (familiar pattern; the cookie's XSS gain
isn't worth the switch now). No refresh token: after 7 days the user logs in again. Google stays the only
identity provider: no passwords, no reset flow, no email service, on purpose.

Left to do:

- [ ] Log in with Google in the browser end to end (new user, existing user, logout, expired token → login).
- [ ] Before deploying: set `SESSION_SECRET` and `GOOGLE_CLIENT_ID` in production (the server won't start without them) and check prod has no duplicate emails (the new unique index would fail the migration, and a failed migration now stops startup): `SELECT email, COUNT(*) FROM users GROUP BY email HAVING COUNT(*) > 1;`
- [ ] `make swagger` to add `/auth/google` to the docs (needs `swag` installed). The `@securityDefinitions` comment in `routes.go` is also malformed.
- [ ] Keep in mind: never render user content (scraps, testimonials, topics) as raw HTML (`dangerouslySetInnerHTML`) without sanitizing. With the token in `localStorage`, XSS is the main risk. None today.

### 2. Transactions for multi-write operations

**Problem.** No `db.Transaction` anywhere. Flows that write several rows can half-succeed and leave counters wrong:

- `CommunityMembershipService.JoinCommunity` / `LeaveCommunity`: membership row + `IncrementCommunityMembers`.
- `FriendshipService.ReplyFriendship`: accept + `UpdateFriendCount` for both users (recounts).
- `FriendshipService.DeleteFriendship`: delete + `IncrementFriendCount(-1)` twice, **errors ignored** (inconsistent with `ReplyFriendship`, which recounts).

**Plan.**

- [ ] Repositories: add `WithTx(tx *gorm.DB)` returning a copy bound to the transaction (cheap, since repositories only hold a `*gorm.DB`). Add it to the generic `repository[T]` so all inherit it.
- [ ] Services: give the services that need it access to the `*gorm.DB` (or a small `TxRunner`) and wrap the flows above in `db.Transaction(func(tx *gorm.DB) error { ... })`.
- [ ] Pick one counter strategy and use it everywhere: either atomic `UPDATE ... SET count = count + ?` inside the transaction, or recount inside the transaction. Stop ignoring errors in `DeleteFriendship`.
- [ ] Grep services for any other place that calls two or more repository writes in a row and wrap those too.
- [ ] Document the pattern in the architecture skill (Service row + recipe).

**Done when:** every service method that writes more than one row runs in a single transaction, and counter updates never ignore errors.

---

## Next

### 3. Tests for the service layer

- [ ] Test helper that opens in-memory SQLite (`file::memory:?cache=shared`), runs `database.Migrate`, and builds `repository.NewRepositories` + `service.NewServices`.
- [ ] First targets: friendship rules (`ShouldThisUserBeReplying`/`Deleting`, duplicate requests, counts after accept/delete), join/leave community (owner can't leave, member counts).
- [ ] `go test ./...` in the Makefile.
- [ ] Later: run the same tests against Postgres (see item 5).

### 4. Error handling: domain errors → HTTP status in one place

- [ ] Shared service errors (`ErrNotFound`, `ErrForbidden`, `ErrConflict`, `ErrInvalid`), wrapped with context via `fmt.Errorf("...: %w", ErrNotFound)`.
- [ ] Repositories or services translate `gorm.ErrRecordNotFound`; services stop importing `gorm` (`service/user.go` does today).
- [ ] One controller helper (e.g. `respondError(ctx, err)`) that maps those errors to 404/403/409/400, else 500 (and logs it).
- [ ] Fixes the "record not found → 500" entry in known-deviations.

### 5. Versioned migrations and Postgres in dev

- [ ] Replace `AutoMigrate` with **goose** or **atlas** before there is production data to protect. Generate the baseline migration from the current schema.
- [ ] `docker-compose.yml` with Postgres for local dev, so dev and prod run the same database; keep SQLite only for fast unit tests, if at all.

---

## Later

### 6. Smaller backend gotchas

- [ ] Generic `repository[T].Update` uses `Updates(struct)`, which silently skips zero values (`false`, `0`, `""`). Use `Select("*")` / explicit columns, or a map, where clearing a field must work.
- [ ] Services return `transport` DTOs. Fine for now (documented in the skill); revisit if one service result needs to feed two responses.

### 7. Frontend structure and dependencies

- [ ] Consider feature folders (`features/profile`, `features/community`, `features/scrapbook`) instead of atomic design + `pages/*/components`, which causes most of the placement entries in known-deviations. Keep a small shared `components/` for real design-system pieces.
- [ ] Align versions: React 18 runtime with React 19 `@types/*`; TypeScript 4.9 is old. Move `@types/*`, `typescript`, eslint/prettier bits to `devDependencies`.
- [ ] Get `npx tsc --noEmit -p .` to zero errors and keep it there.

### 8. Generate frontend types from Swagger

- [ ] Use **orval** (types + TanStack Query hooks) on `backend/docs/swagger.json`. Chosen because it's already familiar from other projects.
  - Swag emits Swagger 2.0; check orval handles it directly, otherwise convert to OpenAPI 3 first.
  - Point orval's mutator at `api/http.ts` so auth and error handling stay in one place.
- [ ] Replace the hand-written `api/*/types.ts`; a backend DTO change then becomes a frontend compile error.
- [ ] Add a script (`make swagger && yarn gen:api`) so it's one command.
