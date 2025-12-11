FROM e-commerce/cgo-builder:latest AS builder
WORKDIR /app
COPY ./payment/ ./payment/
COPY ./pkg/ ./pkg/
COPY ./gen/ ./gen/
WORKDIR /app/payment
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -trimpath -o /app/bin/payment ./cmd/server

# ===== Runtime stage =====
FROM e-commerce/runtime-base:latest

WORKDIR /app
COPY --from=builder /app/bin/payment .
COPY --from=builder /app/payment/migrations ./migrations
COPY --from=builder /app/payment/templates ./templates
RUN chmod +x ./payment
EXPOSE 50051 8081
CMD ["./payment"]