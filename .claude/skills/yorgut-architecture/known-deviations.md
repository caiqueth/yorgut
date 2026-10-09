# Known deviations from the architecture

This backlog lists places where the code breaks the conventions in [SKILL.md](SKILL.md).
Each entry has a location, the problem and a suggested fix. Remove an entry once it is fixed.
Line numbers are from when the entry was written, so search for the symbol if they have drifted.

Severity:
- **bug**: wrong at runtime.
- **layer**: logic is in the wrong layer.
- **consistency**: naming or style differs from the rest of the code.

## Backend

- [ ] **bug**: `UserRepository.SearchUsers` always reports `count: 0`.
  - Where: `internal/repository/user.go`, the second query `r.db.Where(...).Count(&count)`
  - Problem: the count query has no `Model(&model.User{})`, so GORM doesn't know which table to count.
  - Fix: add `.Model(&model.User{})` to the count query.
- [ ] **consistency**: a missing record returns `500` (e.g. `GET /users/9999` → `{"error":"record not found"}`) instead of `404`.
  - Fix: have services translate `gorm.ErrRecordNotFound` into sentinel errors (like `service.ErrFriendshipNotFound`) and map them to `404` in controllers.
- [ ] **consistency**: the `DELETE` route for scraps is registered without an id (`scraps.DELETE("", ...)`); the swagger docs describe it as is (`DELETE /scraps`, `501`).
  - Fix: add `/:scrapId` when the handler is implemented (it returns `501` today), and update its godoc.
- [ ] **bug** (types): pointer fields in `transport` DTOs (`*string` profile fields, `FindCommunityProfileByIdResponse.Description`, `AmIFriendOfResponse.FriendshipId`, `UserReviewSummaryResponse.MyReview`, ...) are sent as `null`, but swag doesn't mark them nullable, so the generated frontend types say `string`/`number`/object. A smoke test against seeded data confirmed every other field matches the spec.
  - Fix: add `extensions:"x-nullable"` to those DTO fields, `make gen-api`, then handle `null` where the frontend now gets type errors.
- [ ] **consistency**: `make swagger` still runs `swag init` without `--requiredByDefault`, so it produces a spec where every response field is optional. Use `make gen-api`.
  - Fix: add `--requiredByDefault` to the `swagger` target (left alone because the Makefile was being edited in parallel).
- [ ] **consistency**: `ReplyFriendshipRequest.Accept` has no `json` tag (swag names it `accept` from the Go name; Go decodes it case-insensitively), and `transport.FindUserByIdResponse` is unused and untagged. `CommunityItem.CreatedBy` is tagged `json:"firstName"` (it holds a user id), so the generated type has `firstName: number`.
  - Fix: tag `Accept`, delete `FindUserByIdResponse`, rename the `CreatedBy` tag to `createdBy` (nothing in the frontend reads it today).
- [ ] **consistency**: `/me/join/:communityId` and `/me/leave/:communityId` put verbs in the URL, and `leave` answers `201`.
  - Fix: if it's worth a breaking change, use `PUT`/`DELETE /me/memberships/:communityId` with `204`.

## Frontend

None open. The structure entries (cross-page imports, duplicated sections, inconsistent section placement, the two `ProfilePicture`s) and the `useGlobalStateContext` import were fixed by the move to feature folders. For the API types, see the nullable-fields entry under Backend.
