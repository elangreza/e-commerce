CREATE TABLE carts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_carts_user_id ON carts(user_id);

-- CREATE TABLE cart_items (
--     id TEXT PRIMARY KEY,
--     cart_id TEXT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
--     product_id TEXT NOT NULL,
--     quantity INT NOT NULL CHECK (quantity > 0),
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
--     UNIQUE(cart_id, product_id)
-- );

