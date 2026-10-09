-- Baseline: the schema GORM AutoMigrate created on Postgres for the models at
-- commit 41a3a88, which is what production was running before migrations
-- existed. The statements are the ones GORM itself issued (captured from a
-- real Postgres 17), with IF NOT EXISTS added: on the existing production
-- database every statement is a no-op and goose only records version 1.
--
-- Caveat of IF NOT EXISTS: an existing table is left exactly as it is, even if
-- it differs from this definition. See backend/README.md for how to compare.

-- +goose Up
CREATE TABLE IF NOT EXISTS "users" (
    "id" bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "email" text,
    "first_name" text,
    "last_name" text,
    "birth_date" timestamptz,
    "friend_count" bigint,
    PRIMARY KEY ("id")
);
CREATE INDEX IF NOT EXISTS "idx_users_deleted_at" ON "users" ("deleted_at");

CREATE TABLE IF NOT EXISTS "user_profiles" (
    "id" bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "user_id" bigint,
    "profile_picture" text,
    "description" text,
    "relationship_status" text,
    "interests" text,
    "about" text,
    "religion" text,
    "humor" text,
    "fashion" text,
    "lives_with" text,
    "birth_city" text,
    "website" text,
    "passions" text,
    "sports" text,
    "hobbies" text,
    "books" text,
    "music" text,
    "movies" text,
    "cuisine" text,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_users_profile" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE SET NULL ON UPDATE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_user_profiles_user_id" ON "user_profiles" ("user_id");
CREATE INDEX IF NOT EXISTS "idx_user_profiles_deleted_at" ON "user_profiles" ("deleted_at");

CREATE TABLE IF NOT EXISTS "communities" (
    "id" bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "name" text,
    "created_by_id" bigint,
    "category" text,
    "owner_id" bigint,
    "member_count" bigint,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_communities_created_by" FOREIGN KEY ("created_by_id") REFERENCES "users" ("id"),
    CONSTRAINT "fk_communities_owner" FOREIGN KEY ("owner_id") REFERENCES "users" ("id")
);
CREATE INDEX IF NOT EXISTS "idx_communities_deleted_at" ON "communities" ("deleted_at");

CREATE TABLE IF NOT EXISTS "community_profiles" (
    "id" bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "description" text,
    "profile_picture" text,
    "community_id" bigint,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_communities_profile" FOREIGN KEY ("community_id") REFERENCES "communities" ("id") ON DELETE SET NULL ON UPDATE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_community_profiles_community_id" ON "community_profiles" ("community_id");
CREATE INDEX IF NOT EXISTS "idx_community_profiles_deleted_at" ON "community_profiles" ("deleted_at");

CREATE TABLE IF NOT EXISTS "community_memberships" (
    "id" bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "community_id" bigint,
    "user_id" bigint,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_community_memberships_community" FOREIGN KEY ("community_id") REFERENCES "communities" ("id"),
    CONSTRAINT "fk_community_memberships_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
CREATE INDEX IF NOT EXISTS "idx_community_memberships_deleted_at" ON "community_memberships" ("deleted_at");

CREATE TABLE IF NOT EXISTS "friendships" (
    "id" bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "user_requested_id" bigint,
    "user_requestee_id" bigint,
    "accepted" boolean,
    "accepted_at" timestamptz,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_friendships_user_requested" FOREIGN KEY ("user_requested_id") REFERENCES "users" ("id"),
    CONSTRAINT "fk_friendships_user_requestee" FOREIGN KEY ("user_requestee_id") REFERENCES "users" ("id")
);
CREATE INDEX IF NOT EXISTS "idx_friendships_deleted_at" ON "friendships" ("deleted_at");

CREATE TABLE IF NOT EXISTS "photo_albums" (
    "id" bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "user_id" bigint,
    "name" text,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_photo_albums_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
CREATE INDEX IF NOT EXISTS "idx_photo_albums_deleted_at" ON "photo_albums" ("deleted_at");

CREATE TABLE IF NOT EXISTS "photos" (
    "id" bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "photo_album_id" bigint,
    "url" text,
    "caption" text,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_photos_photo_album" FOREIGN KEY ("photo_album_id") REFERENCES "photo_albums" ("id")
);
CREATE INDEX IF NOT EXISTS "idx_photos_deleted_at" ON "photos" ("deleted_at");

CREATE TABLE IF NOT EXISTS "scraps" (
    "id" bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "user_author_id" bigint,
    "user_recipient_id" bigint,
    "content" text,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_scraps_user_author" FOREIGN KEY ("user_author_id") REFERENCES "users" ("id"),
    CONSTRAINT "fk_scraps_user_recipient" FOREIGN KEY ("user_recipient_id") REFERENCES "users" ("id")
);
CREATE INDEX IF NOT EXISTS "idx_scraps_deleted_at" ON "scraps" ("deleted_at");

CREATE TABLE IF NOT EXISTS "testimonials" (
    "id" bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "user_author_id" bigint,
    "user_recipient_id" bigint,
    "content" text,
    "approved" boolean,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_testimonials_user_author" FOREIGN KEY ("user_author_id") REFERENCES "users" ("id"),
    CONSTRAINT "fk_testimonials_user_recipient" FOREIGN KEY ("user_recipient_id") REFERENCES "users" ("id")
);
CREATE INDEX IF NOT EXISTS "idx_testimonials_deleted_at" ON "testimonials" ("deleted_at");

CREATE TABLE IF NOT EXISTS "topics" (
    "id" bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "community_id" bigint,
    "author_id" bigint,
    "parent_topic_id" bigint,
    "title" text,
    "text" text,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_topics_community" FOREIGN KEY ("community_id") REFERENCES "communities" ("id"),
    CONSTRAINT "fk_topics_author" FOREIGN KEY ("author_id") REFERENCES "users" ("id"),
    CONSTRAINT "fk_topics_parent_topic" FOREIGN KEY ("parent_topic_id") REFERENCES "topics" ("id")
);
CREATE INDEX IF NOT EXISTS "idx_topics_deleted_at" ON "topics" ("deleted_at");

CREATE TABLE IF NOT EXISTS "user_reviews" (
    "id" bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "reviewer_id" bigint,
    "reviewed_id" bigint,
    "cool" smallint,
    "sexy" smallint,
    "trustworthy" smallint,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_user_reviews_reviewer" FOREIGN KEY ("reviewer_id") REFERENCES "users" ("id"),
    CONSTRAINT "fk_user_reviews_reviewed" FOREIGN KEY ("reviewed_id") REFERENCES "users" ("id")
);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_user_review_pair" ON "user_reviews" ("reviewer_id", "reviewed_id");
CREATE INDEX IF NOT EXISTS "idx_user_reviews_deleted_at" ON "user_reviews" ("deleted_at");

CREATE TABLE IF NOT EXISTS "user_fans" (
    "id" bigserial,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "fan_id" bigint,
    "idol_id" bigint,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_user_fans_idol" FOREIGN KEY ("idol_id") REFERENCES "users" ("id"),
    CONSTRAINT "fk_user_fans_fan" FOREIGN KEY ("fan_id") REFERENCES "users" ("id")
);
CREATE UNIQUE INDEX IF NOT EXISTS "idx_user_fan_pair" ON "user_fans" ("fan_id", "idol_id");
CREATE INDEX IF NOT EXISTS "idx_user_fans_deleted_at" ON "user_fans" ("deleted_at");

-- +goose Down
-- Not reversible on purpose: undoing the baseline would drop all data.
-- +goose StatementBegin
DO $$ BEGIN RAISE EXCEPTION 'the baseline migration cannot be rolled back'; END $$;
-- +goose StatementEnd
