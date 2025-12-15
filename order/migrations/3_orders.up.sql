CREATE TABLE orders (
    id TEXT PRIMARY KEY,
    idempotency_key TEXT UNIQUE NOT NULL,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL,
    transaction_id TEXT,
    total_amount INTEGER NOT NULL,
    -- e.g., 125000 (IDR), 12500 (USD $125.00)
    currency TEXT NOT NULL DEFAULT 'IDR',
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE order_items (
    id TEXT PRIMARY KEY,
    order_id TEXT NOT NULL REFERENCES orders(id) ON
    DELETE
        CASCADE,
        product_id TEXT NOT NULL,
        name TEXT NOT NULL,
        price_per_unit_units INTEGER NOT NULL,
        -- minor units (e.g., cents, rupiah)
        currency TEXT NOT NULL,
        quantity INTEGER NOT NULL CHECK (quantity > 0),
        total_price_units INTEGER NOT NULL,
        -- = price_per_unit_units * quantity
        CHECK (
            total_price_units = price_per_unit_units * quantity
        )
);