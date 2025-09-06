run:
	go run ./cmd/server/main.go
	@echo "Server is running."

gen:
	go generate ./...

seed:
	sqlite3 product.db < migrations/seed/seed.sql

.PHONY: run gen seed
.DEFAULT_GOAL := run