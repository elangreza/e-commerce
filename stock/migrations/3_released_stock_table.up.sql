CREATE TABLE released_stocks (
    id INTEGER PRIMARY KEY,
    stock_id INTEGER NOT NULL REFERENCES stocks(id),
    reserved_stock_id INTEGER NOT NULL REFERENCES reserved_stocks(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    user_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);