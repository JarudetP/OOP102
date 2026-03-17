// repository/sqlite/coffee_repo.go
// SQLite implementation ของ domain.CoffeeRepository

package sqlite

import (
	"database/sql"
	"fmt"

	"coffee-shop/oop102/backend/domain"
)

// SQLiteCoffeeRepo = เก็บข้อมูลกาแฟใน SQLite
type SQLiteCoffeeRepo struct {
	db *sql.DB
}

// NewSQLiteCoffeeRepo = สร้าง coffee repository ใหม่
func NewSQLiteCoffeeRepo(db *sql.DB) domain.CoffeeRepository {
	return &SQLiteCoffeeRepo{db: db}
}

// FindByID = ค้นหากาแฟด้วย ID
func (r *SQLiteCoffeeRepo) FindByID(id string) (*domain.Coffee, error) {
	coffee := &domain.Coffee{}

	query := `SELECT id, name, price, emoji FROM coffees WHERE id = ?`
	err := r.db.QueryRow(query, id).Scan(
		&coffee.ID,
		&coffee.Name,
		&coffee.Price,
		&coffee.Emoji,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("coffee not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return coffee, nil
}

// FindAll = ดึงกาแฟทั้งหมด เรียงตาม ID
func (r *SQLiteCoffeeRepo) FindAll() ([]domain.Coffee, error) {
	query := `SELECT id, name, price, emoji FROM coffees ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	coffees := []domain.Coffee{}
	for rows.Next() {
		coffee := domain.Coffee{}
		err := rows.Scan(&coffee.ID, &coffee.Name, &coffee.Price, &coffee.Emoji)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		coffees = append(coffees, coffee)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return coffees, nil
}
