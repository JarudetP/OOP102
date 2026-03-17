// repository/order_repo.go
// In-Memory implementation ของ domain.OrderRepository

package repository

import (
	"coffee-shop/oop102/backend/domain"
	"fmt"
	"sort"
)

// InMemoryOrderRepo = เก็บข้อมูลออเดอร์ในหน่วยความจำ
type InMemoryOrderRepo struct {
	orders map[string]*domain.Order
}

// NewInMemoryOrderRepo = สร้าง repo ใหม่
func NewInMemoryOrderRepo() domain.OrderRepository {
	return &InMemoryOrderRepo{
		orders: make(map[string]*domain.Order),
	}
}

// Save = บันทึกออเดอร์
func (r *InMemoryOrderRepo) Save(order *domain.Order) error {
	r.orders[order.ID] = order
	return nil
}

// FindByID = ค้นหาออเดอร์ด้วย ID
func (r *InMemoryOrderRepo) FindByID(id string) (*domain.Order, error) {
	order, ok := r.orders[id]
	if !ok {
		return nil, fmt.Errorf("order not found: %s", id)
	}
	return order, nil
}

// FindAll = ดึงออเดอร์ทั้งหมด เรียงตาม ID
func (r *InMemoryOrderRepo) FindAll() ([]domain.Order, error) {
	var result []domain.Order
	for _, o := range r.orders {
		result = append(result, *o)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result, nil
}
