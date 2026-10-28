-- Baseline: the whole schema on MySQL (8.0.13 or later, for the functional
-- index parts below), which replaced Postgres with an empty database. The
-- Postgres migrations it replaces are in the git history before this file.
-- Types and names are the ones GORM AutoMigrate uses for the models, so
-- `make test_mysql` can compare the two.
--
-- What GORM tags can't describe, only here (database.migrationOnly):
--   * idx_community_memberships_active_pair, idx_community_join_requests_active_pair
--     and idx_friendships_active_pair: unique among the live (not
--     soft-deleted) rows only, so leaving and rejoining, or unfriending and
--     asking again, still works. MySQL has no partial indexes; instead the
--     last key part is NULL for soft-deleted rows, and a unique index never
--     compares NULLs. Friendships are unique per pair in either direction.
--     SQLite gets equivalents from the models' tags and AutoMigrateSQLite.

-- +goose Up
CREATE TABLE `users` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `google_sub` varchar(255),
    `email` varchar(255),
    `first_name` longtext,
    `last_name` longtext,
    `birth_date` datetime(6) NULL,
    `friend_count` bigint unsigned,
    `last_seen_at` datetime(6) NULL,
    `record_visits` boolean NOT NULL DEFAULT true,
    `public_album` boolean NOT NULL DEFAULT true,
    `public_communities` boolean NOT NULL DEFAULT true,
    `public_scraps` boolean NOT NULL DEFAULT true,
    `public_testimonials` boolean NOT NULL DEFAULT true,
    `public_birthday` boolean NOT NULL DEFAULT false,
    `hidden_profile_fields` varchar(191) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    INDEX `idx_users_deleted_at` (`deleted_at`),
    UNIQUE INDEX `idx_users_google_sub` (`google_sub`),
    UNIQUE INDEX `idx_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user_profiles` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `user_id` bigint unsigned,
    `profile_picture` longtext,
    `description` longtext,
    `relationship_status` longtext,
    `interests` longtext,
    `about` longtext,
    `religion` longtext,
    `humor` longtext,
    `fashion` longtext,
    `lives_with` longtext,
    `birth_city` longtext,
    `website` longtext,
    `passions` longtext,
    `sports` longtext,
    `hobbies` longtext,
    `books` longtext,
    `music` longtext,
    `movies` longtext,
    `cuisine` longtext,
    PRIMARY KEY (`id`),
    INDEX `idx_user_profiles_deleted_at` (`deleted_at`),
    UNIQUE INDEX `idx_user_profiles_user_id` (`user_id`),
    CONSTRAINT `fk_users_profile` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `communities` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `name` longtext,
    `created_by_id` bigint unsigned,
    `category` longtext,
    `owner_id` bigint unsigned,
    `member_count` bigint unsigned,
    `private` boolean NOT NULL DEFAULT false,
    PRIMARY KEY (`id`),
    INDEX `idx_communities_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_communities_owner` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`),
    CONSTRAINT `fk_communities_created_by` FOREIGN KEY (`created_by_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `community_profiles` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `description` longtext,
    `profile_picture` longtext,
    `community_id` bigint unsigned,
    PRIMARY KEY (`id`),
    INDEX `idx_community_profiles_deleted_at` (`deleted_at`),
    UNIQUE INDEX `idx_community_profiles_community_id` (`community_id`),
    CONSTRAINT `fk_communities_profile` FOREIGN KEY (`community_id`) REFERENCES `communities` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `community_memberships` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `community_id` bigint unsigned,
    `user_id` bigint unsigned,
    PRIMARY KEY (`id`),
    INDEX `idx_community_memberships_deleted_at` (`deleted_at`),
    UNIQUE INDEX `idx_community_memberships_active_pair` (`user_id`, `community_id`, (IF(`deleted_at` IS NULL, 1, NULL))),
    CONSTRAINT `fk_community_memberships_community` FOREIGN KEY (`community_id`) REFERENCES `communities` (`id`),
    CONSTRAINT `fk_community_memberships_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `community_join_requests` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `community_id` bigint unsigned,
    `user_id` bigint unsigned,
    PRIMARY KEY (`id`),
    INDEX `idx_community_join_requests_deleted_at` (`deleted_at`),
    UNIQUE INDEX `idx_community_join_requests_active_pair` (`user_id`, `community_id`, (IF(`deleted_at` IS NULL, 1, NULL))),
    CONSTRAINT `fk_community_join_requests_community` FOREIGN KEY (`community_id`) REFERENCES `communities` (`id`),
    CONSTRAINT `fk_community_join_requests_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `profile_visits` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `visited_id` bigint unsigned,
    `visitor_id` bigint unsigned,
    `day` datetime(6) NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_profile_visits_deleted_at` (`deleted_at`),
    UNIQUE INDEX `idx_profile_visits_daily` (`visited_id`, `visitor_id`, `day`),
    CONSTRAINT `fk_profile_visits_visitor` FOREIGN KEY (`visitor_id`) REFERENCES `users` (`id`),
    CONSTRAINT `fk_profile_visits_visited` FOREIGN KEY (`visited_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `friendships` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `user_requested_id` bigint unsigned,
    `user_requestee_id` bigint unsigned,
    `accepted` boolean,
    `accepted_at` datetime(6) NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_friendships_deleted_at` (`deleted_at`),
    UNIQUE INDEX `idx_friendships_active_pair` (
        (LEAST(`user_requested_id`, `user_requestee_id`)),
        (GREATEST(`user_requested_id`, `user_requestee_id`)),
        (IF(`deleted_at` IS NULL, 1, NULL))
    ),
    CONSTRAINT `fk_friendships_user_requested` FOREIGN KEY (`user_requested_id`) REFERENCES `users` (`id`),
    CONSTRAINT `fk_friendships_user_requestee` FOREIGN KEY (`user_requestee_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `photo_albums` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `user_id` bigint unsigned,
    `name` longtext,
    PRIMARY KEY (`id`),
    INDEX `idx_photo_albums_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_photo_albums_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `photos` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `photo_album_id` bigint unsigned,
    `url` longtext,
    `caption` longtext,
    PRIMARY KEY (`id`),
    INDEX `idx_photos_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_photos_photo_album` FOREIGN KEY (`photo_album_id`) REFERENCES `photo_albums` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `scraps` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `user_author_id` bigint unsigned,
    `user_recipient_id` bigint unsigned,
    `content` longtext,
    PRIMARY KEY (`id`),
    INDEX `idx_scraps_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_scraps_user_author` FOREIGN KEY (`user_author_id`) REFERENCES `users` (`id`),
    CONSTRAINT `fk_scraps_user_recipient` FOREIGN KEY (`user_recipient_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `testimonials` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `user_author_id` bigint unsigned,
    `user_recipient_id` bigint unsigned,
    `content` longtext,
    `approved` boolean,
    PRIMARY KEY (`id`),
    INDEX `idx_testimonials_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_testimonials_user_author` FOREIGN KEY (`user_author_id`) REFERENCES `users` (`id`),
    CONSTRAINT `fk_testimonials_user_recipient` FOREIGN KEY (`user_recipient_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `topics` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `community_id` bigint unsigned,
    `author_id` bigint unsigned,
    `parent_topic_id` bigint unsigned,
    `title` longtext,
    `text` longtext,
    PRIMARY KEY (`id`),
    INDEX `idx_topics_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_topics_community` FOREIGN KEY (`community_id`) REFERENCES `communities` (`id`),
    CONSTRAINT `fk_topics_author` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`),
    CONSTRAINT `fk_topics_parent_topic` FOREIGN KEY (`parent_topic_id`) REFERENCES `topics` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user_reviews` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `reviewer_id` bigint unsigned,
    `reviewed_id` bigint unsigned,
    `cool` tinyint unsigned,
    `sexy` tinyint unsigned,
    `trustworthy` tinyint unsigned,
    PRIMARY KEY (`id`),
    INDEX `idx_user_reviews_deleted_at` (`deleted_at`),
    UNIQUE INDEX `idx_user_review_pair` (`reviewer_id`, `reviewed_id`),
    CONSTRAINT `fk_user_reviews_reviewer` FOREIGN KEY (`reviewer_id`) REFERENCES `users` (`id`),
    CONSTRAINT `fk_user_reviews_reviewed` FOREIGN KEY (`reviewed_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user_fans` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `fan_id` bigint unsigned,
    `idol_id` bigint unsigned,
    PRIMARY KEY (`id`),
    INDEX `idx_user_fans_deleted_at` (`deleted_at`),
    UNIQUE INDEX `idx_user_fan_pair` (`fan_id`, `idol_id`),
    CONSTRAINT `fk_user_fans_fan` FOREIGN KEY (`fan_id`) REFERENCES `users` (`id`),
    CONSTRAINT `fk_user_fans_idol` FOREIGN KEY (`idol_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
-- Not reversible on purpose: undoing the baseline would drop all data. The
-- table below doesn't exist, so a rollback fails with its name as the reason.
SELECT 1 FROM `the_baseline_migration_cannot_be_rolled_back`;
