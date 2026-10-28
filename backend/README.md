# Yorgut
O.G. Social Media Backend

Settings come from `backend/.env`; copy `.env.example` to start. The API
won't start without `SESSION_SECRET` and `GOOGLE_CLIENT_ID`.

# Slowing the API down (development)

To see loading states (buttons with a spinner, skeletons), set `API_DELAY` in
`backend/.env` and restart the API. Every `/api` response then waits that long:
a fixed duration (`API_DELAY=1s`) or a random one in a range
(`API_DELAY=300ms-2s`). Pictures and Swagger aren't delayed. Leave it unset
(the default) for normal speed, and never set it in production.

# Uploads

Pictures go through `storage.ObjectStore` (`internal/storage/`). The database
stores only each picture's key (`<uuid>.jpg`), and the response mappers turn it
into a URL. At login, a user without a picture gets a copy of their Google
profile picture in the store (`AuthService.copyGooglePicture`); users created
before that still have Google URLs, which pass through until their next login
replaces them.
`STORAGE` in `backend/.env` picks the store; the API and the seed refuse to
start if its settings are missing.

A folder on disk, served by the API itself (the default, for dev without Docker):

```sh
STORAGE=disk
UPLOADS_DIR=C:/path/to/yorgut/backend/storage/uploads   # absolute, so it doesn't depend on where the API starts
MEDIA_BASE_URL=http://localhost:8080/media              # the API serves UPLOADS_DIR at this URL's path
```

A bucket on any S3-compatible service (SeaweedFS, Cloudflare R2, Backblaze B2,
AWS). The bucket must exist and allow public reads; the API checks it exists
at startup and clients load pictures straight from `MEDIA_BASE_URL`:

```sh
STORAGE=s3
S3_ENDPOINT=127.0.0.1:8333          # host[:port], no scheme
S3_USE_SSL=false                    # default true
S3_REGION=us-east-1                 # optional for most services
S3_BUCKET=yorgut
S3_ACCESS_KEY=localdev
S3_SECRET_KEY=localdev-secret
MEDIA_BASE_URL=http://localhost:8333/yorgut   # the bucket's public URL (or a proxy/CDN in front of it)
```

The values above are for the local SeaweedFS in the repo-root `docker-compose.yml`:
`make storage_up` starts it with the `yorgut` bucket, the `localdev` key pair
from the repo-root `seaweedfs-s3.json`, and public reads (anonymous users can read pictures,
but not list or write). Pictures survive restarts in the `yorgut-dev-storage`
volume; `make storage_down` stops it.

- Use `127.0.0.1`, not `localhost`, in `S3_ENDPOINT` when Docker runs in WSL:
  Go resolves `localhost` to `::1`, and WSL only forwards IPv4. Browsers fall
  back to IPv4 by themselves, so `MEDIA_BASE_URL` can say `localhost`.
- With Docker in WSL, the containers stop when the WSL distro does; keep a WSL
  terminal open while you work.

`make uploads_to_s3` copies the pictures in `UPLOADS_DIR` into the bucket with
their keys, so switching `STORAGE` from `disk` to `s3` keeps them all. It
overwrites, so it can be rerun.

The seed empties the store and fills it with generated pictures
(`internal/seed/pictures.go`). It only empties an S3 bucket on localhost, and
checks that before resetting the database.

# Database

The schema comes from one of two places, depending on the database
(`internal/database/migrate.go`):

| Database | Used by | Schema from |
|---|---|---|
| MySQL 8.0.13+ (`DEV_MODE=false`) | production, local dev with Docker | versioned SQL migrations in `internal/database/migrations/`, applied by [goose](https://github.com/pressly/goose) when the API starts |
| SQLite (`DEV_MODE=true`) | local dev without Docker, the service tests | GORM `AutoMigrate` from the models |

The migrations are the source of truth for production. `make test_mysql`
checks that they produce the same MySQL schema as `AutoMigrate` from the
models, and runs the service tests on MySQL.

`DATABASE_DSN` for MySQL is either the Go driver's form,
`user:pass@tcp(host:3306)/db?tls=true`, or a URL,
`mysql://user:pass@host:3306/db?tls=true`. For a hosted database add
`tls=true`, or `tls=skip-verify` if the provider signs its certificate with
its own CA. The API always reads and writes times as UTC and keeps
microseconds (`internal/database/database.go`).

## Local MySQL

```sh
make db_up          # docker compose up: MySQL 8.4, data kept in a volume
```

`backend/.env`:

```sh
DEV_MODE=false
DATABASE_DSN=localdev:localdev@tcp(localhost:3306)/yorgut-dev
```

`make run` applies the pending migrations; `make seed` resets the database
(drops all its tables) and fills it with fake data. The seed
refuses to reset a database that is not on localhost. `make db_down` stops it
(the volume stays; `docker compose down -v` deletes it, and the SeaweedFS
volume too).

Without Docker, keep `DEV_MODE=true` and `DATABASE_DSN=<file>.db` (SQLite);
the seed then deletes that file to reset.

The seed's main user (most mock users befriend it, it owns 21+ communities)
is `SEED_MAIN_USER_EMAIL` from `backend/.env`, or `main@test.com` without it.
Set it to the Google account you sign in with.

## Changing the schema

1. Change the model in `internal/model/model.go`.
2. `make migration_new name=add_something` creates the next
   `internal/database/migrations/000NN_add_something.sql`. Write the MySQL
   SQL under `-- +goose Up` (and `-- +goose Down`), matching what GORM would
   create: same column types (`datetime(6)`, `longtext`, `bigint unsigned`,
   `varchar(255)` for an indexed string, which needs a `size` tag) and index
   names (`idx_<table>_<column>`). MySQL commits each DDL statement on its
   own, so a migration that fails halfway keeps what ran before the failing
   statement: keep them small. Never edit a migration that has already run in
   production; add a new one.
3. `make db_up`, then `make test_mysql` (compares the migrations with
   `AutoMigrate`, runs the service tests on MySQL) and `make test`.
4. `make migration_status` lists applied and pending migrations of the
   database at `DATABASE_DSN`.

## Production

`deploy/docker-compose.yml` runs the API image and sets its variables from
`deploy/.env` ([README.md, "Deploying"](../README.md#deploying)). There,
`GIN_MODE=release` and `DISABLE_SWAGGER=true` are set, and `API_DELAY` must
stay unset.

The API runs the pending migrations at startup (under a MySQL named lock, so
two instances starting at once don't both run them) and refuses to start if
one fails. On an empty database the first start creates every table.

### Running the API elsewhere

The compose file covers these; to run the API image some other way, set:

- Set `SESSION_SECRET` and `GOOGLE_CLIENT_ID`; the server won't start without them.
- Set `ADMIN_EMAIL` to the address people should write to about content that
  breaks the rules (served by `GET /info`, shown on the help pages).
- Set `DATABASE_DSN` to the MySQL (8.0.13 or later: the baseline uses
  functional index parts), with `tls=true` or `tls=skip-verify` for a hosted
  one. The user needs full rights on that database only.
- Set `STORAGE` and `MEDIA_BASE_URL` for the pictures (`UPLOADS_DIR` on disk,
  or the `S3_*` variables for a bucket with public reads).
- Moving to a bucket later: create it with public reads, set the `S3_*`
  variables, run `make uploads_to_s3` with `UPLOADS_DIR` still set, then set
  `STORAGE=s3` and `MEDIA_BASE_URL` to the bucket's public URL.
