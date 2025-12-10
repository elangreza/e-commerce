FROM golang:1.24.3-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY . .
WORKDIR /app/order
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -o /app/bin/order ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /app
COPY --from=builder /app/bin/order .
COPY --from=builder /app/order/migrations ./migrations
RUN chmod +x ./order
EXPOSE 50051
CMD ["./order"]