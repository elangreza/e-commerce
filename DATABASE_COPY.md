# Database Copy Utility

## Overview

The `copy-docker-data.sh` script allows bidirectional copying of SQLite databases between Docker (`data/`) and local development (`data-local/`) environments.

## Usage

### Via Makefile (Recommended)

```bash
# Copy from Docker to local development
make copy-to-local

# Copy from local development to Docker
make copy-to-docker
```

### Direct Script Usage

```bash
# Copy data/ → data-local/ (default)
./copy-docker-data.sh
./copy-docker-data.sh to-local

# Copy data-local/ → data/
./copy-docker-data.sh to-docker
```

## Use Cases

### 1. Copy Docker Data to Local Development

**When**: You have data in Docker and want to use it for local development

```bash
make copy-to-local
```

This will:
- Copy all `*/data/` directories to `*/data-local/`
- Change ownership to your user
- Allow you to run `make run` locally with Docker's data

### 2. Copy Local Development Data to Docker

**When**: You've been developing locally and want Docker to use your local data

```bash
make copy-to-docker
```

This will:
- Copy all `*/data-local/` directories to `*/data/`
- Allow Docker to use your local development data
- Useful for testing local changes in Docker environment

## Examples

### Scenario 1: Start with Docker, then develop locally

```bash
# 1. Run Docker first
docker compose up

# 2. Create some data (users, orders, etc.)

# 3. Stop Docker
docker compose down

# 4. Copy data to local
make copy-to-local

# 5. Continue development locally
cd api && make run
```

### Scenario 2: Develop locally, then test in Docker

```bash
# 1. Develop locally
cd api && make run

# 2. Create test data locally

# 3. Stop local service (Ctrl+C)

# 4. Copy to Docker
make copy-to-docker

# 5. Test in Docker
docker compose up
```

## What It Does

The script:
1. ✅ Removes the destination directory (to avoid conflicts)
2. ✅ Copies the source directory recursively
3. ✅ Changes ownership to your user (for local development)
4. ✅ Handles all 6 services automatically

## Safety

- Uses `sudo` to handle permission changes
- Removes destination before copying (clean slate)
- Skips services where source doesn't exist
- Shows clear progress messages
