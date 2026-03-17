// repository/sqlite/order_repo.go
// SQLite implementation ของ domain.OrderRepository

package sqlite

import (
	"database/sql"
	"fmt"

	"coffee-shop/oop102/backend/domain"
)

// SQLiteOrderRepo = เก็บข้อมูลออเดอร์ใน SQLite
type SQLiteOrderRepo struct {
	db *sql.DB
}

// NewSQLiteOrderRepo = สร้าง order repository ใหม่
func NewSQLiteOrderRepo(db *sql.DB) domain.OrderRepository {
	return &SQLiteOrderRepo{
		db: db,
	}
}

// Save = บันทึกออเดอร์และ items ลง database
func (r *SQLiteOrderRepo) Save(order *domain.Order) error {
	// เริ่มต้น transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. บันทึก order หลัก
	orderQuery := `INSERT INTO orders (id, total, status, created_at) VALUES (?, ?, ?, ?)`
	_, err = tx.Exec(orderQuery, order.ID, order.Total, order.Status, order.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// 2. บันทึก order items ทีละรายการ
	itemQuery := `INSERT INTO order_items (order_id, coffee_id, quantity) VALUES (?, ?, ?)`
	for _, item := range order.Items {
		_, err := tx.Exec(itemQuery, order.ID, item.Coffee.ID, item.Quantity)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	// 3. Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// FindByID = ค้นหาออเดอร์เฉพาะ ID พร้อม items
func (r *SQLiteOrderRepo) FindByID(id string) (*domain.Order, error) {
	// Query order หลัก
	order := &domain.Order{}
	query := `SELECT id, total, status, created_at FROM orders WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(&order.ID, &order.Total, &order.Status, &order.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	// Query order items
	items, err := r.queryOrderItems(id)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

// FindAll = ดึงออเดอร์ทั้งหมด พร้อม items
func (r *SQLiteOrderRepo) FindAll() ([]domain.Order, error) {
	query := `SELECT id, total, status, created_at FROM orders ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	orders := []domain.Order{}
	for rows.Next() {
		order := domain.Order{}
		err := rows.Scan(&order.ID, &order.Total, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		// Query items สำหรับออเดอร์นี้
		items, err := r.queryOrderItems(order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return orders, nil
}

// ============ Helper Functions ============

// queryOrderItems = ค้นหา items ของออเดอร์เฉพาะ พร้อม coffee details
// ใช้ SQL JOIN เพื่อดึงข้อมูล coffee และ order_items ในคำสั่ง 1 ครั้ง (ไม่ N+1 queries)
func (r *SQLiteOrderRepo) queryOrderItems(orderID string) ([]domain.OrderItem, error) {
	// JOIN query - ดึง coffee details พร้อม order_items ในคำสั่ง 1 ครั้ง
	query := `
		SELECT 
			oi.coffee_id, 
			oi.quantity,
			c.id,
			c.name,
			c.price,
			c.emoji
		FROM order_items oi
		JOIN coffees c ON oi.coffee_id = c.id
		WHERE oi.order_id = ?
		ORDER BY oi.id
	`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to query order items: %w", err)
	}
	defer rows.Close()

	items := []domain.OrderItem{}
	for rows.Next() {
		var coffeeID, coffeeName, coffeeEmoji string
		var quantity int
		var price float64

		err := rows.Scan(&coffeeID, &quantity, &coffeeID, &coffeeName, &price, &coffeeEmoji)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		// สร้าง Coffee object จากผลลัพธ์ query
		coffee := &domain.Coffee{
			ID:    coffeeID,
			Name:  coffeeName,
			Price: price,
			Emoji: coffeeEmoji,
		}

		items = append(items, domain.OrderItem{
			Coffee:   *coffee,
			Quantity: quantity,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate order items: %w", err)
	}

	return items, nil
}
