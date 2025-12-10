FROM e-commerce/cgo-builder:latest AS builder
WORKDIR /app
COPY ./product/ ./product/
COPY ./pkg/ ./pkg/
COPY ./gen/ ./gen/
WORKDIR /app/product
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -trimpath -o /app/bin/product ./cmd/server

# ===== Runtime stage =====
FROM e-commerce/runtime-base:latest

WORKDIR /app
COPY --from=builder /app/bin/product .
COPY --from=builder /app/product/migrations ./migrations
RUN chmod +x ./product
EXPOSE 50052
CMD ["./product"]