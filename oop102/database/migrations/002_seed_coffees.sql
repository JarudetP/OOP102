-- 002_seed_coffees.sql
-- เพิ่มข้อมูลเมนูกาแฟเริ่มต้น

INSERT OR IGNORE INTO coffees (id, name, price, emoji) VALUES 
    ('1', 'Latte', 65, '☕'),
    ('2', 'Mocha', 75, '☕'),
    ('3', 'Americano', 55, '☕'),
    ('4', 'Matcha Latte', 70, '🍵'),
    ('5', 'Thai Tea', 50, '🧋');
