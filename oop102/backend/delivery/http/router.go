// delivery/http/router.go
// ตั้ง HTTP routes
// ใช้ Handler สำหรับตอบสนอง requests

package http

import (
	"coffee-shop/oop102/backend/usecase"
	"net/http"
)

// NewRouter = สร้าง HTTP router พร้อมทุก routes
func NewRouter(orderUC *usecase.OrderUseCase) http.Handler {
	handler := NewHTTPHandler(orderUC)

	mux := http.NewServeMux()

	// Coffee endpoints
	mux.HandleFunc("GET /coffees", handler.GetCoffees)

	// Order endpoints
	mux.HandleFunc("POST /orders", handler.CreateOrder)

	return mux
}
