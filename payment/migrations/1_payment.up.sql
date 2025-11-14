CREATE TABLE payments (
    id TEXT PRIMARY KEY,
    status TEXT NOT NULL,
    total_amount INTEGER NOT NULL,
    currency TEXT NOT NULL,
    transaction_id TEXT UNIQUE NOT NULL,
    order_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payments_transaction_id ON payments(transaction_id);