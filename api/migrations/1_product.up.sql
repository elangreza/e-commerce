-- {
--         "id": "0198a7e8-5e30-714d-9a1d-b198f451f59c",
--         "name": "smartphone",
--         "description": "A handheld device that combines mobile phone and computing functions.",
--         "price": 699,
--         "image_url": "http://example.com/smartphone.png"
-- }

CREATE TABLE products (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price INTEGER NOT NULL CHECK (price >= 0),
    currency TEXT,
    image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
