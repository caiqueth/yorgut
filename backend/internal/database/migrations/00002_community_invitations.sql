-- Invitations from a community's members for their friends to join it.
-- idx_community_invitations_active_pair is unique among the live (not
-- soft-deleted) rows only, like the active-pair indexes of the baseline.

-- +goose Up
CREATE TABLE `community_invitations` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(6) NULL,
    `updated_at` datetime(6) NULL,
    `deleted_at` datetime(6) NULL,
    `community_id` bigint unsigned,
    `inviter_id` bigint unsigned,
    `invitee_id` bigint unsigned,
    PRIMARY KEY (`id`),
    INDEX `idx_community_invitations_deleted_at` (`deleted_at`),
    UNIQUE INDEX `idx_community_invitations_active_pair` (`invitee_id`, `community_id`, (IF(`deleted_at` IS NULL, 1, NULL))),
    CONSTRAINT `fk_community_invitations_community` FOREIGN KEY (`community_id`) REFERENCES `communities` (`id`),
    CONSTRAINT `fk_community_invitations_inviter` FOREIGN KEY (`inviter_id`) REFERENCES `users` (`id`),
    CONSTRAINT `fk_community_invitations_invitee` FOREIGN KEY (`invitee_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE `community_invitations`;
