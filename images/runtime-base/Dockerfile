# images/runtime-base/Dockerfile
FROM alpine:latest

# Install common runtime dependencies for all Go + SQLite services
RUN apk --no-cache add ca-certificates sqlite-libs

# Optional: create non-root user (recommended for production)
# RUN adduser -D -s /bin/sh appuser
# USER appuser

WORKDIR /app