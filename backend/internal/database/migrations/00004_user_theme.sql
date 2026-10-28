-- The look each user picked for the app (see model.User.Theme). Everyone
-- starts on the 2005 one, the only look before themes existed.

-- +goose Up
ALTER TABLE `users` ADD COLUMN `theme` varchar(32) NOT NULL DEFAULT '2005';

-- +goose Down
ALTER TABLE `users` DROP COLUMN `theme`;
