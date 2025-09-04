CREATE TABLE stocks (
    id INTEGER PRIMARY KEY,
    product_id TEXT NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reserved_stock (
    id INTEGER PRIMARY KEY,
    stock_id INTEGER NOT NULL REFERENCES stocks(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    user_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE released_stocks (
    id INTEGER PRIMARY KEY,
    stock_id INTEGER NOT NULL REFERENCES stocks(id),
    reserved_stock_id INTEGER NOT NULL REFERENCES reserved_stock(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    user_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);