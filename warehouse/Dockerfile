FROM e-commerce/cgo-builder:latest AS builder
WORKDIR /app
COPY ./warehouse/ ./warehouse/
COPY ./pkg/ ./pkg/
COPY ./gen/ ./gen/
WORKDIR /app/warehouse
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -trimpath -o /app/bin/warehouse ./cmd/server

# ===== Runtime stage =====
FROM e-commerce/runtime-base:latest

WORKDIR /app
COPY --from=builder /app/bin/warehouse .
COPY --from=builder /app/warehouse/migrations ./migrations
RUN chmod +x ./warehouse
EXPOSE 50053
CMD ["./warehouse"]