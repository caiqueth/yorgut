-- At most one live (not soft-deleted) membership per user and community, and
-- one live friendship per pair of users, in either direction. The services
-- check "already exists" before inserting, but two concurrent requests could
-- both pass that check; these indexes make the second insert fail instead.
-- Partial (WHERE deleted_at IS NULL) so leaving and rejoining, or unfriending
-- and asking again, still works: the old row stays, soft-deleted.
--
-- Existing duplicates would make the index creation fail, so they are
-- soft-deleted first:
--   * memberships: the oldest row (lowest id) is kept;
--   * friendships: an accepted row is kept over a pending one, then the oldest.
-- Then friend and member counts are recomputed for everyone (same rule as
-- UserRepository.UpdateFriendCount / CommunityRepository.UpdateMemberCount),
-- which also fixes counts that were already wrong before transactions.
--
-- The membership index is mirrored in model.CommunityMembership's tags. The
-- friendship one is an expression index GORM can't describe portably; SQLite
-- gets an equivalent one in database.AutoMigrateSQLite.

-- +goose Up
UPDATE "community_memberships" m
SET "deleted_at" = NOW()
WHERE m."deleted_at" IS NULL
  AND EXISTS (
    SELECT 1 FROM "community_memberships" o
    WHERE o."deleted_at" IS NULL
      AND o."user_id" = m."user_id"
      AND o."community_id" = m."community_id"
      AND o."id" < m."id"
  );

UPDATE "friendships" f
SET "deleted_at" = NOW()
WHERE f."deleted_at" IS NULL
  AND EXISTS (
    SELECT 1 FROM "friendships" o
    WHERE o."deleted_at" IS NULL
      AND o."id" <> f."id"
      AND LEAST(o."user_requested_id", o."user_requestee_id") = LEAST(f."user_requested_id", f."user_requestee_id")
      AND GREATEST(o."user_requested_id", o."user_requestee_id") = GREATEST(f."user_requested_id", f."user_requestee_id")
      AND (
        (COALESCE(o."accepted", false) AND NOT COALESCE(f."accepted", false))
        OR (COALESCE(o."accepted", false) = COALESCE(f."accepted", false) AND o."id" < f."id")
      )
  );

CREATE UNIQUE INDEX IF NOT EXISTS "idx_community_memberships_active_pair"
  ON "community_memberships" ("user_id", "community_id")
  WHERE "deleted_at" IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_friendships_active_pair"
  ON "friendships" (LEAST("user_requested_id", "user_requestee_id"), GREATEST("user_requested_id", "user_requestee_id"))
  WHERE "deleted_at" IS NULL;

UPDATE "users" u
SET "friend_count" = (
  SELECT COUNT(*) FROM "friendships" f
  WHERE f."deleted_at" IS NULL
    AND f."accepted"
    AND (f."user_requested_id" = u."id" OR f."user_requestee_id" = u."id")
);

UPDATE "communities" c
SET "member_count" = (
  SELECT COUNT(*) FROM "community_memberships" m
  WHERE m."deleted_at" IS NULL AND m."community_id" = c."id"
);

-- +goose Down
-- The soft-deleted duplicates are not restored.
DROP INDEX IF EXISTS "idx_friendships_active_pair";
DROP INDEX IF EXISTS "idx_community_memberships_active_pair";
