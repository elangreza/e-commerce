# docker/cgo-builder/Dockerfile
FROM golang:1.24.3-alpine

ENV CGO_ENABLED=1

RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

# Optional: Set default working dir
WORKDIR /app