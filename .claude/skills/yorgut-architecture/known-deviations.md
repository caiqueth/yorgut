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
- [ ] **consistency**: the `DELETE` route for scraps is registered without an id (`scraps.DELETE("", ...)`); the swagger docs describe it as is (`DELETE /scraps`, `501`).
  - Fix: add `/:scrapId` when the handler is implemented (it returns `501` today), and update its godoc.
- [ ] **consistency**: `transport.FindUserByIdResponse` is unused and untagged. `CommunityItem.CreatedBy` is tagged `json:"firstName"` (it holds a user id), so the generated type has `firstName: number`.
  - Fix: delete `FindUserByIdResponse`, rename the `CreatedBy` tag to `createdBy` (nothing in the frontend reads it today).
- [ ] **consistency**: `/me/join/:communityId` and `/me/leave/:communityId` put verbs in the URL, and `leave` answers `201`.
  - Fix: if it's worth a breaking change, use `PUT`/`DELETE /me/memberships/:communityId` with `204`.

## Frontend

None open. The structure entries (cross-page imports, duplicated sections, inconsistent section placement, the two `ProfilePicture`s) and the `useGlobalStateContext` import were fixed by the move to feature folders.
