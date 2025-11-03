CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,                     -- Who placed the order?
    status TEXT NOT NULL CHECK (status IN (
        'pending', 'confirmed', 'shipped', 'delivered', 'cancelled', 'failed'
    )),
    total_amount DECIMAL(12, 2) NOT NULL,      -- Total price (cached for reporting)
    currency CHAR(3) NOT NULL DEFAULT 'IDR',   -- e.g., IDR, USD, EUR
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,   
    shipped_at TIMESTAMP DEFAULT NULL,
    cancelled_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,                 -- ID from Product Service (denormalized)
    name TEXT NOT NULL,                       -- Product name at time of order (denormalized)
    price_per_unit DECIMAL(10, 2) NOT NULL,   -- Snapshot of price (critical!)
    quantity INT NOT NULL CHECK (quantity > 0),
    total_price DECIMAL(12, 2) NOT NULL       -- = price_per_unit * quantity
);