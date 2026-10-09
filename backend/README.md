# Yorgut
O.G. Social Media Backend

# Database

The schema comes from one of two places, depending on the database
(`internal/database/migrate.go`):

| Database | Used by | Schema from |
|---|---|---|
| Postgres (`DEV_MODE=false`) | production, local dev with Docker | versioned SQL migrations in `internal/database/migrations/`, applied by [goose](https://github.com/pressly/goose) when the API starts |
| SQLite (`DEV_MODE=true`) | local dev without Docker, the service tests | GORM `AutoMigrate` from the models |

The migrations are the source of truth for production. `make test_pg` checks
that they produce the same Postgres schema as `AutoMigrate` from the models.

## Local Postgres

```sh
make db_up          # docker compose up: Postgres 17, data kept in a volume
```

`backend/.env`:

```sh
DEV_MODE=false
DATABASE_DSN=postgres://localdev:localdev@localhost:5432/yorgut-dev?sslmode=disable
```

`make run` applies the pending migrations; `make seed` resets the database
(drops and recreates the `public` schema) and fills it with fake data. The seed
refuses to reset a database that is not on localhost. `make db_down` stops it
(the volume stays; `docker compose down -v` deletes it).

Without Docker, keep `DEV_MODE=true` and `DATABASE_DSN=<file>.db` (SQLite);
the seed then deletes that file to reset.

## Changing the schema

1. Change the model in `internal/model/model.go`.
2. `make migration_new name=add_something` creates the next
   `internal/database/migrations/000NN_add_something.sql`. Write the Postgres
   SQL under `-- +goose Up` (and `-- +goose Down`), matching what GORM would
   create: same column types and index names (`idx_<table>_<column>`).
   Each file runs in a transaction. Never edit a migration that has already
   run in production; add a new one.
3. `make test_pg` (compares the migrations with `AutoMigrate`) and `make test`.
4. `make migration_status` lists applied and pending migrations of the
   database at `DATABASE_DSN`.

## Production

The API runs the pending migrations at startup (under a Postgres advisory
lock) and refuses to start if one fails. Migration 00001 is the baseline: the
schema `AutoMigrate` created before migrations existed, written with
`IF NOT EXISTS`, so on the existing database it only records version 1.
`IF NOT EXISTS` leaves an existing table as it is, even if it differs. To check
production matches, run these on it and on a freshly migrated local Postgres
(same migration version) and diff the output:

```sql
SELECT table_name, column_name, data_type, is_nullable, column_default
FROM information_schema.columns WHERE table_schema = 'public' ORDER BY 1, 2;
SELECT indexdef FROM pg_indexes WHERE schemaname = 'public' ORDER BY 1;
SELECT conrelid::regclass, conname, pg_get_constraintdef(oid)
FROM pg_constraint WHERE connamespace = 'public'::regnamespace ORDER BY 1, 2;
```

# TODO
- Authentication
  - Login endpoint
  - JWT configuration and validation pattern
- Communities:
  - Create Community endpoint.
  - Update Community endpoint.
  - Split communities in two types:
    - Public or Private:
      - Public communities can be joined by anyone. Also, their forums can be viewed by anyone, be it a member or not.
      - Private communities can be joined by anyone, but the owner must accept any new memberships. Also, their forums can only be read by community members.
  - Community Members endpoint:
    - Split in two: one that shows the last 8 logged in and another that paginates over all members.
  - Forum:
    - Create Forum model and endpoints:
      - A Forum is FK of a Community
      - A Forum has its own children: "Topics"
      - A Topic has its own children: "Posts"
- User:
  - My Friends:
    - Split in two endpoints: one that shows the last 8 logged in and another that paginates over all friends.
  - My Communities:
    - Split in two endpoints: one that shows random 8 and another that paginates over all communities.
  - Scrapbook:
    - List all user's scraps.
    - Create scrap.
  - General Search endpoint:
    - Search for other users and communities.
- Profile:
  - "Friends" page:
    - Paginate over all friends of a given user.
  - "Communities" page:
    - List communities joined paginated.
    - List communities owned paginated.
    - Create Testimonial endpoint.
    - List last N Testimonials endpoint.
    - List Testimonials paginated.
    - Voting in Trustworthy (😀), Cool (🧊), Sexy (❤️) endpoint.
    - Create fanship (?) model and crud.
- Others:
  - Home screen data endpoint.
    - Indirect Connections
    - Direct Connections (Friends)
    - Number of Fans
    - Number of Scraps
    - Profile Views - Total
    - Profile Views - Last 7 days
    - Last 10 visitors
    - Today's Luck
- Add unittests.
- Change Sqlite to Postgres for production.
