#!/bin/bash
# Copy database files between Docker (data/) and local development (data-local/)
# Usage: ./copy-docker-data.sh [to-local|to-docker]

DIRECTION=${1:-to-local}

if [ "$DIRECTION" != "to-local" ] && [ "$DIRECTION" != "to-docker" ]; then
    echo "Usage: $0 [to-local|to-docker]"
    echo "  to-local  : Copy data/ → data-local/ (default)"
    echo "  to-docker : Copy data-local/ → data/"
    exit 1
fi

if [ "$DIRECTION" = "to-local" ]; then
    echo "Copying Docker database files (data/) to local development (data-local/)..."
    SOURCE="data"
    DEST="data-local"
else
    echo "Copying local development database files (data-local/) to Docker (data/)..."
    SOURCE="data-local"
    DEST="data"
fi

for service in api order product warehouse shop payment; do
    if [ -d "$service/$SOURCE" ]; then
        echo "Copying $service/$SOURCE/*.db* to $service/$DEST/..."
        
        # Create destination directory if it doesn't exist
        mkdir -p "$service/$DEST"
        
        # Copy only database files (*.db, *.db-shm, *.db-wal)
        if ls "$service/$SOURCE"/*.db* 1> /dev/null 2>&1; then
            sudo cp "$service/$SOURCE"/*.db* "$service/$DEST/"
            sudo chown -R $USER:$USER "$service/$DEST"
            echo "  ✓ Database files copied to $service/$DEST"
        else
            echo "  ⊘ No database files found in $service/$SOURCE"
        fi
    else
        echo "  ⊘ $service/$SOURCE not found, skipping"
    fi
done

echo ""
if [ "$DIRECTION" = "to-local" ]; then
    echo "Done! You can now run services locally with Docker data."
    echo "Run: cd api && make run"
else
    echo "Done! Docker will now use your local development data."
    echo "Run: docker compose up"
fi
