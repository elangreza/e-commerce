FROM golang:1.24.3-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY . .
WORKDIR /app/product
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -o /app/bin/product ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /app
COPY --from=builder /app/bin/product .
COPY --from=builder /app/product/migrations ./migrations
RUN chmod +x ./product
EXPOSE 50052
CMD ["./product"]