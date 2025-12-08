insert into warehouses (id, name, is_active)
values 
    (1, 'w-JKT', TRUE),
    (2, 'w-BDG', TRUE);

-- simulate the stock is spread across multiple warehouses and multiple shop

insert into stocks (product_id, warehouse_id, shop_id, quantity)
values 
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5c9a', 1, 1, 10),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5c9b', 1, 1, 130),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5c9c', 1, 1, 16),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5c9d', 1, 1, 16),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5c9e', 1, 1, 243),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5c9f', 1, 1, 72),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca1', 1, 1, 657),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca2', 1, 1, 77),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca3', 1, 1, 53),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca4', 1, 1, 73),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca5', 1, 1, 65),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca6', 1, 1, 13),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca7', 1, 2, 2),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca8', 1, 2, 7),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca9', 1, 2, 2),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5caa', 1, 2, 0),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5cab', 1, 2, 7),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5cac', 1, 2, 2),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5cad', 1, 2, 0);

insert into stocks (product_id, warehouse_id, shop_id, quantity)
values 
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5c9a', 2, 1, 10),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5c9b', 2, 1, 130),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5c9c', 2, 1, 16),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5c9d', 2, 1, 16),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5c9e', 2, 1, 243),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5c9f', 2, 1, 72),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca1', 2, 1, 657),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca2', 2, 1, 77),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca3', 2, 1, 53),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca4', 2, 1, 73),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca5', 2, 1, 65),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca6', 2, 1, 13),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca7', 2, 2, 2),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca8', 2, 2, 7),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5ca9', 2, 2, 2),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5caa', 2, 2, 0),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5cab', 2, 2, 7),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5cac', 2, 2, 2),
    ('019394d0-4d5e-7d6a-9c4b-8a3f2e1d5cad', 2, 2, 0);