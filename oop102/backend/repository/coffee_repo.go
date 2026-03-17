// repository/coffee_repo.go
// In-Memory implementation ของ domain.CoffeeRepository

package repository

import (
	"coffee-shop/oop102/backend/domain"
	"fmt"
	"sort"
)

// InMemoryCoffeeRepo = เก็บข้อมูลกาแฟในหน่วยความจำ
type InMemoryCoffeeRepo struct {
	coffees map[string]domain.Coffee
}

// NewInMemoryCoffeeRepo = สร้าง repo ใหม่พร้อมเมนูเริ่มต้น
func NewInMemoryCoffeeRepo() domain.CoffeeRepository {
	return &InMemoryCoffeeRepo{
		coffees: map[string]domain.Coffee{
			"1": {ID: "1", Name: "Latte", Price: 65, Emoji: "☕"},
			"2": {ID: "2", Name: "Mocha", Price: 75, Emoji: "☕"},
			"3": {ID: "3", Name: "Americano", Price: 55, Emoji: "☕"},
			"4": {ID: "4", Name: "Matcha Latte", Price: 70, Emoji: "🍵"},
			"5": {ID: "5", Name: "Thai Tea", Price: 50, Emoji: "🧋"},
		},
	}
}

// FindByID = ค้นหากาแฟด้วย ID
func (r *InMemoryCoffeeRepo) FindByID(id string) (*domain.Coffee, error) {
	c, ok := r.coffees[id]
	if !ok {
		return nil, fmt.Errorf("coffee not found: %s", id)
	}
	return &c, nil
}

// FindAll = ดึงกาแฟทั้งหมด เรียงตาม ID
func (r *InMemoryCoffeeRepo) FindAll() ([]domain.Coffee, error) {
	var result []domain.Coffee
	for _, c := range r.coffees {
		result = append(result, c)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result, nil
}
