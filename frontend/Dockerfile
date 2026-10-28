# The built SPA served by nginx, which also proxies /api to the API and
# /media to the bucket (nginx.conf). Used by the root docker-compose.yml
# ("app" profile).

# The output is plain files, so the build stage always runs on the builder's
# own architecture (CI builds amd64 and arm64 images).
FROM --platform=$BUILDPLATFORM node:24.20.0-alpine AS build
WORKDIR /src
COPY package.json yarn.lock ./
# Yarn 1 gives up on a tarball after 30s; big ones (@mui/material) can take
# longer while the api image builds in parallel on the same connection.
RUN yarn install --frozen-lockfile --network-timeout 600000
COPY . .
# The Google OAuth client id is baked into the bundle.
ARG GOOGLE_CLIENT_ID
# .env.docker: same-origin API URL (/api/v1).
RUN yarn build --env-mode docker

FROM nginx:1.30.5-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY security-headers.conf /etc/nginx/security-headers.conf
COPY --from=build /src/dist /usr/share/nginx/html
