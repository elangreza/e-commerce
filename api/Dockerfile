FROM golang:1.24.3-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY . .
COPY ./api/ ./api/
COPY ./gen/ ./gen/
COPY ./pkg/ ./pkg/
WORKDIR /app/api
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -o /app/bin/api ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /app
COPY --from=builder /app/bin/api .
COPY --from=builder /app/api/migrations ./migrations
RUN chmod +x ./api
EXPOSE 8080
CMD ["./api"]