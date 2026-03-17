-- 001_init.sql
-- สร้างตารางสำหรับ Coffee Shop

-- Coffees table - เก็บเมนูกาแฟ
CREATE TABLE IF NOT EXISTS coffees (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    price REAL NOT NULL,
    emoji TEXT NOT NULL
);

-- Orders table - เก็บออเดอร์
CREATE TABLE IF NOT EXISTS orders (
    id TEXT PRIMARY KEY,
    total REAL NOT NULL,
    status TEXT NOT NULL,
    created_at DATETIME NOT NULL
);

-- OrderItems table - รายการกาแฟในแต่ละออเดอร์
CREATE TABLE IF NOT EXISTS order_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_id TEXT NOT NULL,
    coffee_id TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (coffee_id) REFERENCES coffees(id)
);

-- Create indices for better query performance
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_coffee_id ON order_items(coffee_id);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
