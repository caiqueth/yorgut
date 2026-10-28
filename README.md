# Yorgut
O.G. Social Media

Monorepo to keep it all glued together for now.

# Running

## Local development

### TL;DR

Needs Go 1.25 and Node with Yarn 1; no Docker.

1. Create a Google OAuth client id (Web application, origin
   `http://localhost:3000`).
2. `echo "GOOGLE_CLIENT_ID=<id>" > frontend/.env.local`
3. `cp backend/.env.example backend/.env` and fill in `SESSION_SECRET`
   (`openssl rand -hex 32`), `GOOGLE_CLIENT_ID`, `UPLOADS_DIR` (absolute path
   to `backend/storage/uploads`) and `SEED_MAIN_USER_EMAIL` (your Google
   e-mail).
4. Optional fake data: `cd backend && make seed`
5. In two terminals: `cd backend && make run` and
   `cd frontend && yarn install && yarn start`
6. Open http://localhost:3000 and sign in.

Details below.

### Detailed

Day to day, with no Docker at all: SQLite in a file and the pictures in a
local folder that the API serves itself. Needs Go 1.25 and Node with Yarn 1.

1. A Google OAuth client for the login: in the Google Cloud console,
   **APIs & Services > Credentials > Create credentials > OAuth client ID**,
   type **Web application**, with `http://localhost:3000` (and
   `http://localhost:8000` for Docker) in **Authorized JavaScript origins**.
   Put its client id in `frontend/.env.local` (ignored by git):

   ```sh
   echo "GOOGLE_CLIENT_ID=<your client id>" > frontend/.env.local
   ```

2. Backend settings: `.env.example` is already set up for this.

   ```sh
   cd backend
   cp .env.example .env
   ```

   In `backend/.env`, fill in:
   - `SESSION_SECRET`: any random value (`openssl rand -hex 32`).
   - `GOOGLE_CLIENT_ID`: the same client id.
   - `UPLOADS_DIR`: the absolute path of `backend/storage/uploads` in your
     clone. The API creates it, and git ignores what's in it.
   - `SEED_MAIN_USER_EMAIL`: the Google account you sign in with, so the seed
     gives it friends, communities and scraps.

   Keep `DEV_MODE=true`, `DATABASE_DSN=local.db`, `STORAGE=disk` and
   `MEDIA_BASE_URL=http://localhost:8080/media`.

3. Fake data (optional): `make seed` in `backend/`. It deletes `local.db` and
   the uploads folder and fills both again, so stop the API first.

4. Run the API and the frontend in two terminals:

   ```sh
   cd backend && make run                      # http://localhost:8080
   cd frontend && yarn install && yarn start   # http://localhost:3000
   ```

   The API creates the SQLite tables when it starts.

To use MySQL and an S3 bucket (SeaweedFS) like production instead, start them
from `docker-compose.yml` with `make db_up` and `make storage_up`, and change
`backend/.env` as [backend/README.md](backend/README.md) describes.

## The whole app in Docker (optional)

To try the production setup locally (nginx serving the built SPA and proxying
`/api` and `/media`, with MySQL and SeaweedFS), after the setup above and with
`http://localhost:8000` added to the OAuth client's origins:

```sh
docker compose --env-file backend/.env --profile app up -d --build   # http://localhost:8000
docker compose --profile app down                                    # -v also deletes the data
```

To seed it, run from `backend/` with `go run` (`make seed` would use the
`backend/.env` values instead):

```sh
DEV_MODE=false DATABASE_DSN="localdev:localdev@tcp(127.0.0.1:3306)/yorgut-dev" \
STORAGE=s3 S3_ENDPOINT=127.0.0.1:8333 S3_USE_SSL=false S3_BUCKET=yorgut \
S3_ACCESS_KEY=localdev S3_SECRET_KEY=localdev-secret MEDIA_BASE_URL=/media \
go run ./cmd/seed
```

# Deploying

Production is [`deploy/docker-compose.yml`](deploy/docker-compose.yml), on any
machine with Docker:

- **Images**: GitHub Actions ([`images.yml`](.github/workflows/images.yml))
  builds `yorgut-api` and `yorgut-web` on every push to `main` and pushes them
  to ghcr.io as `latest` and `<commit sha>`, for amd64 and arm64. The web image
  needs the repository variable `GOOGLE_CLIENT_ID` (**Settings > Secrets and
  variables > Actions > Variables**). The server only pulls them
  (`docker login ghcr.io` first if the packages are private).
- **Services**: the API, nginx serving the SPA, and a SeaweedFS bucket for the
  pictures (in the `storage-data` volume). The database is an external MySQL
  (8.0.13 or later).
- **Settings**: copy `deploy/.env.example` to `deploy/.env` and fill it in;
  every value is explained there. It holds the secrets, so it stays on the
  server (git ignores it). `backend/.env` isn't used.
- **Reaching it**: the app listens on `127.0.0.1:8000` only. The `cloudflared`
  service publishes it through a Cloudflare Tunnel (`TUNNEL_TOKEN`); to use a
  reverse proxy or load balancer instead, remove that service and point yours
  at port 8000. Either way, add the public origin to the Google OAuth client's
  authorized JavaScript origins, or login fails.

Start, and update after the images workflow finishes:

```sh
docker compose -f deploy/docker-compose.yml --env-file deploy/.env pull
docker compose -f deploy/docker-compose.yml --env-file deploy/.env up -d
```

The API runs the pending migrations when it starts. To roll back, set
`YORGUT_TAG` in `deploy/.env` to an earlier commit sha and run the same two
commands.

[`deploy/backup.sh`](deploy/backup.sh) dumps the database and archives the
pictures into `deploy/backups/`, keeping 14 days (`KEEP_DAYS`). Run it on a
schedule and copy the backups off the machine.

Never run the seed against production (it wipes the database), and never
`down -v` (it deletes the pictures).

# License

[MIT](LICENSE)
