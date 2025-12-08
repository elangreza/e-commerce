-- {
--         "id": "0198a7e8-5e30-714d-9a1d-b198f451f59c",
--         "name": "smartphone",
--         "description": "A handheld device that combines mobile phone and computing functions.",
--         "price": 699,
--         "image_url": "http://example.com/smartphone.png"
-- }

-- shop_id is used for grouping
-- can easily grow if multiple shop_id in one commerce

CREATE TABLE products (
    id TEXT PRIMARY KEY,
    shop_id INTEGER NOT NULL, 
    name TEXT NOT NULL,
    description TEXT,
    price INTEGER NOT NULL CHECK (price >= 0),
    currency TEXT,
    image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
