FROM --platform=linux/amd64 golang:1.22.6-bookworm AS builder

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    gcc \
    librdkafka-dev \
    pkg-config && \
    rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/app ./cmd/app

# Final stage
FROM --platform=linux/amd64 debian:bookworm-slim

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    librdkafka1 \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/app .
COPY --from=builder /app/config/config.yaml ./config/

RUN useradd -m appuser && chown -R appuser:appuser /app
USER appuser

COPY --from=builder /app/static ./static

EXPOSE 8080

ENTRYPOINT ["./app"]