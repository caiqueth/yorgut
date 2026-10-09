-- Login with own sessions (commit 0212a76): users are keyed by Google's
-- account id, and e-mails are unique. Same names GORM AutoMigrate uses for
-- `uniqueIndex` on User.GoogleSub and User.Email.
--
-- Fails if production has duplicate e-mails; check first with
--   SELECT email, COUNT(*) FROM users GROUP BY email HAVING COUNT(*) > 1;

-- +goose Up
ALTER TABLE "users" ADD COLUMN IF NOT EXISTS "google_sub" text;
CREATE UNIQUE INDEX IF NOT EXISTS "idx_users_google_sub" ON "users" ("google_sub");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_users_email" ON "users" ("email");

-- +goose Down
DROP INDEX IF EXISTS "idx_users_email";
DROP INDEX IF EXISTS "idx_users_google_sub";
ALTER TABLE "users" DROP COLUMN IF EXISTS "google_sub";
