FROM golang:1.24.3-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY . .
WORKDIR /app/shop
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -o /app/bin/shop ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /app
COPY --from=builder /app/bin/shop .
COPY --from=builder /app/shop/migrations ./migrations
RUN chmod +x ./shop
EXPOSE 50054
CMD ["./shop"]