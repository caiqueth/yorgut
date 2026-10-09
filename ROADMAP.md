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
- [x] `/auth/google` is in the Swagger docs and the `@securityDefinitions` comment is fixed (done with item 8).
- [ ] Keep in mind: never render user content (scraps, testimonials, topics) as raw HTML (`dangerouslySetInnerHTML`) without sanitizing. With the token in `localStorage`, XSS is the main risk. None today.

### 2. Transactions for multi-write operations ✅ (done 2026-10-09)

`repository.Transactor` runs a callback with repositories bound to one transaction. Used by
`CreateCommunity`, `UpdateCommunity`, `JoinCommunity`, `LeaveCommunity`, `ReplyFriendship` and
`DeleteFriendship`. Friend and member counts are recomputed in a single SQL `UPDATE ... (SELECT COUNT(*))`
inside the transaction instead of read-modify-write increments. This also fixed `DeleteFriendship`
decrementing both users' counts when cancelling a *pending* request. Pattern documented in the architecture skill.

Left to do:

- [x] Counts already wrong in production: migration `00003` recomputes every friend and member count on deploy.
- [x] Duplicate memberships/friendships under concurrency: migration `00003` soft-deletes existing duplicates and adds
  partial unique indexes (`WHERE deleted_at IS NULL`, friendships on the unordered pair). A concurrent duplicate now
  fails at the database, but the API answers 500 for it; mapping it to 409 belongs to item 4.
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

### 5. Versioned migrations and Postgres in dev ✅ (done 2026-10-09)

- [x] **goose** with embedded SQL migrations (`backend/internal/database/migrations/`), run at API startup under an
  advisory lock. `00001` is the baseline, taken from the DDL GORM issued for the `main` (`41a3a88`) models on a real
  Postgres and made idempotent, so the existing production DB just records it. `00002` adds `google_sub` and the
  unique email index; `00003` adds the duplicate protection (see item 2).
- [x] `backend/docker-compose.yml` (Postgres 17), `make db_up` / `db_down` / `migration_new name=...` /
  `migration_status`. The seed can reset a local Postgres (refuses non-localhost DSNs).
- SQLite stays for dev without Docker and for the tests, through `AutoMigrate`. A schema change needs a migration
  **and** the model change; `make test_pg` (embedded Postgres, `pgtest` build tag) checks both give the same schema.

Left to do:

- [ ] Before deploying: check there are no duplicate emails in production (`00002` would fail and stop startup):
  `SELECT email, COUNT(*) FROM users GROUP BY email HAVING COUNT(*) > 1;`
- [ ] After deploying: compare production's schema with a freshly migrated DB using the queries in `backend/README.md`
  ("Production"). `main` ignored AutoMigrate errors, so production may differ in ways `IF NOT EXISTS` leaves alone.
- [ ] `docker-compose.yml` is untested (no Docker on the machine it was written on).
- [ ] Run the service tests against Postgres too (item 3); the `pgtest` setup can be reused.

---

## Later

### 6. Smaller backend gotchas ✅ (done 2026-10-09)

- [x] Generic `repository[T].Update` used `Updates(struct)`, which silently skips zero values. Nothing called it, so it's deleted; every partial update already uses `Update`, a map, or `Select(...).Updates(struct)`. The rule is in the architecture skill ("Partial updates").
- [x] Services return `transport` DTOs: kept on purpose (documented in the skill); revisit if one service result needs to feed two responses.

### 7. Frontend structure and dependencies ✅ (done 2026-10-09)

- [x] Feature folders: `src/features/<feature>/pages/<Page>` (route screens) and `components/<Name>` (sections only
  that feature uses); no feature imports another. Shared pieces live in a flat `src/components/<Name>` (the
  atoms/molecules/organisms split is gone). `Menu` became the shared `HomeMenu`; the duplicated `FriendsList` /
  `CommunitiesList` are one component each (optional `userId`); the two `ProfilePicture`s are now `ProfilePicture`
  and `UserProfilePicture`.
- [x] React 18 with matching `@types/react` 18; TypeScript 5.9; types and tooling in `devDependencies`.
- [x] `npx tsc --noEmit -p .` at zero errors (`yarn typecheck`).

Left to do:

- [ ] Click through the app in a browser. The move was mechanical and type-checks, but nobody has looked at it.

### 8. Generate frontend types from Swagger ✅ (done 2026-10-09)

- [x] **orval** 8 generates TanStack Query hooks and types into `frontend/src/api/generated/` (committed) from
  `backend/docs/swagger.json`; its mutator is `api/http.ts`, so auth, 401 handling and error messages are unchanged.
  `api/<domain>/queries.ts` stays as a thin layer (same hook names; query options and cache invalidation live there,
  using orval's query-key helpers).
- [x] Hand-written `api/*/types.ts` deleted; pages import from `api/generated/model`.
- [x] `make gen-api` regenerates the spec (swag `--requiredByDefault`) and then the client (`yarn gen:api`).
- [x] Controller godoc fixed so the spec matches the handlers (`@ID`, params, response types, `BearerAuth`), and pointer
  fields in response DTOs are marked `x-nullable`, so the client types them `| null`. That surfaced 13 places that
  assumed a value could never be `null` (image props, the topic page); they're fixed.

Left to do:

- [ ] Query keys are now spec paths (`["/me/friends"]`) instead of full URLs. Worth a browser check that
  scraps, friendships, memberships and testimonials still refresh after changes.
