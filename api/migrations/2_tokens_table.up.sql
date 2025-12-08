CREATE TABLE IF NOT EXISTS "tokens" (
    "id" TEXT PRIMARY KEY,
    "user_id" TEXT NOT NULL REFERENCES "users" ("id"),
    "token" TEXT NOT NULL,
    "token_type" TEXT NOT NULL,
    "issued_at" TIMESTAMP NOT NULL,
    "expired_at" TIMESTAMP NOT NULL,
    "duration" TEXT NOT NULL
);