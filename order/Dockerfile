FROM e-commerce/cgo-builder:latest AS builder
WORKDIR /app
COPY ./order/ ./order/
COPY ./pkg/ ./pkg/
COPY ./gen/ ./gen/
WORKDIR /app/order
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -trimpath -o /app/bin/order ./cmd/server

# ===== Runtime stage =====
FROM e-commerce/runtime-base:latest

WORKDIR /app
COPY --from=builder /app/bin/order .
COPY --from=builder /app/order/migrations ./migrations
RUN chmod +x ./order
EXPOSE 50051
CMD ["./order"]