FROM e-commerce/cgo-builder:latest AS builder
WORKDIR /app
COPY ./api/ ./api/
COPY ./gen/ ./gen/
COPY ./pkg/ ./pkg/
WORKDIR /app/api
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -trimpath -o /app/bin/api ./cmd/server

# ===== Runtime stage =====
FROM e-commerce/runtime-base:latest

WORKDIR /app
COPY --from=builder /app/bin/api .
COPY --from=builder /app/api/migrations ./migrations
RUN chmod +x ./api
EXPOSE 8080
CMD ["./api"]