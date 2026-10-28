-- Community polls ("enquetes"): the poll, its options and one vote per user.

-- +goose Up
CREATE TABLE `polls` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `community_id` bigint unsigned,
    `author_id` bigint unsigned,
    `question` longtext,
    `description` longtext,
    `closes_at` datetime(6) NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_polls_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_polls_community` FOREIGN KEY (`community_id`) REFERENCES `communities` (`id`),
    CONSTRAINT `fk_polls_author` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `poll_options` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `poll_id` bigint unsigned,
    `text` longtext,
    PRIMARY KEY (`id`),
    INDEX `idx_poll_options_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_poll_options_poll` FOREIGN KEY (`poll_id`) REFERENCES `polls` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `poll_votes` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `poll_id` bigint unsigned,
    `poll_option_id` bigint unsigned,
    `user_id` bigint unsigned,
    PRIMARY KEY (`id`),
    INDEX `idx_poll_votes_deleted_at` (`deleted_at`),
    UNIQUE INDEX `idx_poll_vote_pair` (`poll_id`, `user_id`),
    CONSTRAINT `fk_poll_votes_poll` FOREIGN KEY (`poll_id`) REFERENCES `polls` (`id`),
    CONSTRAINT `fk_poll_votes_poll_option` FOREIGN KEY (`poll_option_id`) REFERENCES `poll_options` (`id`),
    CONSTRAINT `fk_poll_votes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE `poll_votes`;
DROP TABLE `poll_options`;
DROP TABLE `polls`;
