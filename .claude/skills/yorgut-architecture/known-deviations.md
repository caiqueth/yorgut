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
- [ ] **consistency**: the `DELETE` route for scraps is registered without an id (`scraps.DELETE("", ...)`), while their swagger docs say `/{id}`.
  - Fix: add `/:<entity>Id` when the handlers are implemented (they return `501` today).
- [ ] **consistency**: `/me/join/:communityId` and `/me/leave/:communityId` put verbs in the URL, and `leave` answers `201`.
  - Fix: if it's worth a breaking change, use `PUT`/`DELETE /me/memberships/:communityId` with `204`.

## Frontend

### Bugs caused by structure

- [ ] **bug**: `hooks/useGlobalStateContext.ts` imports from `contexts/GlobalState`, but the file is `contexts/GlobalStateContext.tsx`, so this is a TS error.
- [ ] **bug**: `routes/ProtectedRoute.tsx` returns `ReactNode` from a `React.FC`, so this is a TS error.
  - Fix: wrap it in a fragment.

### Components in the wrong place

- [ ] **layer**: one page imports a component from another.
  - Where: `pages/MyScrapbook/index.tsx` imports `pages/Home/components/Menu`.
  - Fix: promote `Menu` to `components/organisms/`.
- [ ] **layer**: page sections are duplicated.
  - Where: `FriendsList` and `CommunitiesList` exist in both `pages/Home/components/` and `pages/Profile/components/`.
  - Fix: if they differ only in data source, promote one parameterised organism.
- [ ] **consistency**: route-level sections are placed differently across pages.
  - Profile keeps `ProfileHeader`, `ProfileDetails` and `ProfileMenu` at `pages/Profile/<Name>/`, as if they were sub-routes, while the page-only `FriendsList`/`CommunitiesList`/`UserReview` sit in `pages/Profile/components/`.
  - Community keeps `ProfileHeader`/`ProfileDetails` in `pages/Community/components/` but `ProfileMenu` at `pages/Community/ProfileMenu/`.
  - Fix: anything that is not a route belongs in `components/`.
- [ ] **consistency**: two different components are both named `ProfilePicture` (`components/atoms/ProfilePicture` and `components/molecules/ProfilePicture`).
  - Fix: rename one (for example, the molecule to `UserProfilePicture`, which is already its export name).
