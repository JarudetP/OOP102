// repository/sqlite/db.go
// Database initialization และ schema setup

package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB = สร้าง database connection และ tables
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// createTables = สร้าง schema ทั้งหมด
func createTables(db *sql.DB) error {
	queries := []string{
		// Coffees table - เก็บเมนูกาแฟ
		`CREATE TABLE IF NOT EXISTS coffees (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			price REAL NOT NULL,
			emoji TEXT NOT NULL
		);`,

		// Orders table - เก็บออเดอร์
		`CREATE TABLE IF NOT EXISTS orders (
			id TEXT PRIMARY KEY,
			total REAL NOT NULL,
			status TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);`,

		// OrderItems table - รายการกาแฟในแต่ละออเดอร์
		`CREATE TABLE IF NOT EXISTS order_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id TEXT NOT NULL,
			coffee_id TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			FOREIGN KEY (order_id) REFERENCES orders(id),
			FOREIGN KEY (coffee_id) REFERENCES coffees(id)
		);`,

		// Insert default coffees if not exist
		`INSERT OR IGNORE INTO coffees (id, name, price, emoji) VALUES 
			('1', 'Latte', 65, '☕'),
			('2', 'Mocha', 75, '☕'),
			('3', 'Americano', 55, '☕'),
			('4', 'Matcha Latte', 70, '🍵'),
			('5', 'Thai Tea', 50, '🧋');`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}
