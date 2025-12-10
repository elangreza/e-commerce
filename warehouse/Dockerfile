FROM golang:1.24.3-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY . .
WORKDIR /app/warehouse
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -o /app/bin/warehouse ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /app
COPY --from=builder /app/bin/warehouse .
COPY --from=builder /app/warehouse/migrations ./migrations
RUN chmod +x ./warehouse
EXPOSE 50053
CMD ["./warehouse"]