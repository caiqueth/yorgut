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

### 2. Transactions for multi-write operations ✅ (done 2026-10-09)

`repository.Transactor` runs a callback with repositories bound to one transaction. Used by
`CreateCommunity`, `UpdateCommunity`, `JoinCommunity`, `LeaveCommunity`, `ReplyFriendship` and
`DeleteFriendship`. Friend and member counts are recomputed in a single SQL `UPDATE ... (SELECT COUNT(*))`
inside the transaction instead of read-modify-write increments. This also fixed `DeleteFriendship`
decrementing both users' counts when cancelling a *pending* request. Pattern documented in the architecture skill.

Left to do:

- [ ] Counts already wrong in production heal on the next friendship/membership change of that user or
  community. To fix them all at once, run after deploying:
  `UPDATE users u SET friend_count = (SELECT COUNT(*) FROM friendships f WHERE f.deleted_at IS NULL AND f.accepted AND (f.user_requested_id = u.id OR f.user_requestee_id = u.id));`
  `UPDATE communities c SET member_count = (SELECT COUNT(*) FROM community_memberships m WHERE m.deleted_at IS NULL AND m.community_id = c.id);`
- [ ] Duplicate rows are still possible under concurrency: `JoinCommunity` and `CreateFriendship` check
  "already exists" and then insert. A transaction doesn't prevent that; a unique index would, but soft deletes
  (`deleted_at`) make it tricky (leaving and rejoining keeps the old row). Options: partial unique index
  `WHERE deleted_at IS NULL` (Postgres and SQLite support it; needs item 5's migrations) or hard deletes for these tables.
- [ ] An uploaded picture stays on disk if `CreateCommunity` then fails. Minor; clean up when there's an upload cleanup job.

---

## Next

### 3. Tests for the service layer

- [x] Test helper that opens in-memory SQLite, runs `database.Migrate`, and builds the repositories (`service/helpers_test.go`, `newTestRepositories`).
- [x] Started: login/session (`auth_test.go`), friend and member counts, transaction rollback (`counters_test.go`).
- [ ] Next targets: friendship permission rules (`ShouldThisUserBeReplying`/`Deleting`, duplicate requests), testimonials, user reviews/fans (friends-only rules).
- [x] `go test ./...` in the Makefile (`make test`, was already there).
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
