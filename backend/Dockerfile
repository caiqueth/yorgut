# The API as a small static binary. Used by the root docker-compose.yml
# ("app" profile); configuration comes from the environment (README).

# The build stage runs on the builder's own architecture and cross-compiles
# for the target one (CI builds amd64 and arm64).
FROM --platform=$BUILDPLATFORM golang:1.27.1-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# No cgo: the container uses MySQL; the SQLite driver (DEV_MODE=true) needs
# cgo and is only for local dev.
RUN export CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
 && go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api \
 && go build -trimpath -ldflags="-s -w" -o /out/seed ./cmd/seed

# CA certificates (Google login) and time zones (America/Sao_Paulo) included.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/api /out/seed /app/
EXPOSE 8080
ENTRYPOINT ["/app/api"]
