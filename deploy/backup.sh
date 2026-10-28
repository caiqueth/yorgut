#!/usr/bin/env bash
# Backs up Yorgut's database and pictures into deploy/backups/, and deletes
# backups older than KEEP_DAYS (default 14). Run on the server, e.g. from cron
# (README.md, "Deploying"). Copy the backups to another machine too: a backup on the
# same disk doesn't survive the disk.
set -euo pipefail

cd "$(dirname "$0")/.."
compose=(docker compose -f deploy/docker-compose.yml --env-file deploy/.env)
dir="deploy/backups"
keep_days="${KEEP_DAYS:-14}"
stamp="$(date +%Y-%m-%d_%H%M)"
mkdir -p "$dir"

# The database (the external MySQL in deploy/.env), as SQL that recreates
# its tables (each is dropped first), dumped by a throwaway MySQL client
# container. --single-transaction gives a consistent copy without locking.
# docker reads deploy/.env literally, so any password works.
docker run --rm --env-file deploy/.env mysql:8.4 sh -c '
  ssl_mode=REQUIRED
  if [ "${DB_TLS:-true}" = false ]; then ssl_mode=DISABLED; fi
  MYSQL_PWD="$DB_PASSWORD" exec mysqldump -h "$DB_HOST" -P "${DB_PORT:-3306}" -u "$DB_USER" \
    --ssl-mode="$ssl_mode" --single-transaction --no-tablespaces --set-gtid-purged=OFF "$DB_NAME"' \
  | gzip > "$dir/db_$stamp.sql.gz"

# The pictures: the bucket's volume, read by a throwaway container. The bucket
# is paused meanwhile (a few seconds without pictures) so the copy is
# consistent; the trap resumes it even if tar fails.
"${compose[@]}" pause storage
trap '"${compose[@]}" unpause storage' EXIT
docker run --rm \
  -v yorgut-prod_storage-data:/data:ro \
  -v "$PWD/$dir":/backup \
  alpine:3.22 sh -c "tar czf /backup/pictures_$stamp.tar.gz -C /data . && chown $(id -u):$(id -g) /backup/pictures_$stamp.tar.gz"
"${compose[@]}" unpause storage
trap - EXIT

find "$dir" -type f \( -name 'db_*.sql.gz' -o -name 'pictures_*.tar.gz' \) -mtime +"$keep_days" -delete

echo "$(date -Iseconds) backed up: $dir/db_$stamp.sql.gz, $dir/pictures_$stamp.tar.gz"
