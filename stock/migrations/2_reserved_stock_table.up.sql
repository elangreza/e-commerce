CREATE TABLE reserved_stocks (
    id INTEGER PRIMARY KEY,
    stock_id INTEGER NOT NULL REFERENCES stocks(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    user_id TEXT NOT NULL,
    -- status can be "reserved", "released"
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);