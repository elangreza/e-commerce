FROM e-commerce/cgo-builder:latest AS builder
WORKDIR /app
COPY ./shop/ ./shop/
COPY ./pkg/ ./pkg/
COPY ./gen/ ./gen/
WORKDIR /app/shop
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -trimpath -o /app/bin/shop ./cmd/server

# ===== Runtime stage =====
FROM e-commerce/runtime-base:latest

WORKDIR /app
COPY --from=builder /app/bin/shop .
COPY --from=builder /app/shop/migrations ./migrations
RUN chmod +x ./shop
EXPOSE 50054
CMD ["./shop"]