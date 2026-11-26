--  psql -d <your_db_name> -f scripts/seed.sql

--  products
INSERT INTO products (sku, name, price, stock) VALUES
	('RO-001', 'Rose', 25000, 52),
	('TU-002', 'Tulip', 35000, 65),
	('LI-003', 'Lily', 40000, 19),
	('OR-004', 'Orchid', 150000, 30),
	('SU-005', 'Sunflower', 50000, 25),
	('DA-006', 'Daisy', 20000, 40),
	('CB-007', 'Carnation', 28000, 200),
	('JA-008', 'Jasmine', 30000, 35),
	('LA-009', 'Lavender', 45000, 12),
	('PE-010', 'Peony', 120000, 7),
	('DH-011', 'Dahlia', 60000, 22),
	('MA-012', 'Marigold', 15000, 15)
ON CONFLICT (sku) DO NOTHING;

-- Stock movements (gives Rose history, so DELETE on it returns 409)
INSERT INTO stock_movements (product_id, type, qty, note)
SELECT id, 'in', 100, 'restock' FROM products WHERE sku = 'RO-001';

INSERT INTO stock_movements (product_id, type, qty, note)
SELECT id, 'out', 5, 'sale' FROM products WHERE sku = 'RO-001';

INSERT INTO stock_movements (product_id, type, qty, note)
SELECT id, 'in', 65, 'restock' FROM products WHERE sku = 'TU-002';

INSERT INTO stock_movements (product_id, type, qty, note)
SELECT id, 'out', 3, NULL FROM products WHERE sku = 'TU-002';

-- Order 
WITH o AS (
	INSERT INTO sales_orders (customer, total)
	VALUES ('John Doe', 85000)
	RETURNING id
)
INSERT INTO sales_order_items (order_id, sku, qty, unit_price)
SELECT o.id, 'RO-001', 2, 25000 FROM o
UNION ALL
SELECT o.id, 'TU-002', 1, 35000 FROM o;

WITH o AS (
	INSERT INTO sales_orders (customer, status, total, cancelled_at)
	VALUES ('Jane Doe', 'CANCELLED', 40000, now())
	RETURNING id
)
INSERT INTO sales_order_items (order_id, sku, qty, unit_price)
SELECT o.id, 'LI-003', 1, 40000 FROM o;
