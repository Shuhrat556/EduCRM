# syntax=docker/dockerfile:1

FROM golang:1.26-bookworm AS builder

WORKDIR /src

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates git \
	&& rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api \
	&& go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate \
	&& go build -trimpath -ldflags="-s -w" -o /out/seed ./cmd/seed

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /out/api /app/api
COPY --from=builder /out/migrate /app/migrate
COPY --from=builder /out/seed /app/seed
COPY migrations /app/migrations

RUN groupadd --system app && useradd --system --gid app --no-create-home app \
	&& chown -R app:app /app

USER app

ENV MIGRATIONS_PATH=/app/migrations

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
	CMD curl -fsS http://127.0.0.1:8080/health || exit 1

ENTRYPOINT ["/app/api"]
