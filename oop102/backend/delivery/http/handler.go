// delivery/http/handler.go
// HTTP Handler Layer - รับ HTTP request → ส่งให้ UseCase → ส่ง JSON response
// ห้ามเข้า Database โดยตรง! ต้องผ่าน UseCase เท่านั้น

package http

import (
	"coffee-shop/oop102/backend/domain"
	"coffee-shop/oop102/backend/usecase"
	"encoding/json"
	"io"
	"net/http"
)

// HTTPHandler = ตัวจัดการ HTTP requests
type HTTPHandler struct {
	orderUC *usecase.OrderUseCase
}

// NewHTTPHandler = สร้าง handler ใหม่
// รับ UseCase เป็น dependency (Dependency Injection)
func NewHTTPHandler(orderUC *usecase.OrderUseCase) *HTTPHandler {
	return &HTTPHandler{
		orderUC: orderUC,
	}
}

// ============ Request DTOs ============
// (รับมาจาก HTTP client - ไม่ต้องเหมือน domain)

// CreateOrderRequest = client ส่งมา (ต่างจาก domain.Order!)
type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items"`
}

// OrderItemRequest = 1 รายการกาแฟในออเดอร์
type OrderItemRequest struct {
	CoffeeID string `json:"coffee_id"`
	Quantity int    `json:"quantity"`
}

// ============ Response DTOs ============
// (ส่งออกไปหา HTTP client - ไม่ต้องใช้ domain model โดยตรง)

type CoffeeResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Emoji string  `json:"emoji"`
}

// OrderResponse = ส่งกลับเมื่อสั่งกาแฟสำเร็จ
type OrderResponse struct {
	ID        string              `json:"id"`
	Items     []OrderItemResponse `json:"items"`
	Total     float64             `json:"total"`
	Status    string              `json:"status"`
	CreatedAt string              `json:"created_at"`
}

// OrderItemResponse = 1 รายการในออเดอร์ (response)
type OrderItemResponse struct {
	CoffeeID   string  `json:"coffee_id"`
	CoffeeName string  `json:"coffee_name"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

// ============ Handlers ============

// GetCoffees = GET /coffees - ดึงเมนูกาแฟทั้งหมด
// ✅ เรียก UseCase เท่านั้น
// ✅ ไม่เข้า Database โดยตรง
// ✅ ส่ง JSON response
func (h *HTTPHandler) GetCoffees(w http.ResponseWriter, r *http.Request) {
	// 1. เรียก UseCase เพื่อดึงเมนู
	coffees, err := h.orderUC.GetMenu()
	if err != nil {
		// ถ้าเกิด error → return HTTP 500
		respondError(w, 500, err.Error())
		return
	}

	// 2. แปลง domain.Coffee → CoffeeResponse (ให้ HTTP layer เป็นเจ้าของการแปลง)
	responses := make([]CoffeeResponse, len(coffees))
	for i, coffee := range coffees {
		responses[i] = CoffeeResponse{
			ID:    coffee.ID,
			Name:  coffee.Name,
			Price: coffee.Price,
			Emoji: coffee.Emoji,
		}
	}

	// 3. ส่ง JSON response
	respondSuccess(w, http.StatusOK, responses)
}

// CreateOrder = POST /orders - สั่งกาแฟใหม่
// ✅ Parse JSON request → แปลงเป็น UseCase request → เรียก UseCase
// ✅ Business rules อยู่ใน UseCase ไม่ใช่ที่นี่
// ✅ ส่ง JSON response
func (h *HTTPHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// 1. Parse JSON request body
	var req CreateOrderRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// 2. แปลง HTTP DTO → UseCase DTO
	// (นี่คือการแปลง format จากภายนอก → ภายใน)
	items := make([]usecase.OrderItemRequest, len(req.Items))
	for i, item := range req.Items {
		items[i] = usecase.OrderItemRequest{
			CoffeeID: item.CoffeeID,
			Quantity: item.Quantity,
		}
	}

	// 3. เรียก UseCase - ให้ UseCase ตัดสินใจทางธุรกิจ
	// ที่นี่ UseCase จะเช็ค:
	//   - ออเดอร์ต้องไม่ว่าง
	//   - จำนวนต้องมากกว่า 0
	//   - Coffee ID ต้องเจอ
	//   - คำนวณเงิน คิด status เป็นต้น
	order, err := h.orderUC.PlaceOrder(usecase.OrderRequest{Items: items})
	if err != nil {
		// Business error (validation failed)
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 4. แปลง domain.Order → OrderResponse (ให้ HTTP layer เป็นเจ้าของการแปลง)
	orderResp := h.mapOrderToResponse(order)

	// 5. ส่ง JSON response
	respondSuccess(w, http.StatusCreated, orderResp)
}

// ============ Helper Functions ============

// mapOrderToResponse = แปลง domain.Order → OrderResponse
// ทำให้ usecase ไม่รู้เรื่อง JSON format
func (h *HTTPHandler) mapOrderToResponse(order *domain.Order) *OrderResponse {
	itemResponses := make([]OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		itemResponses[i] = OrderItemResponse{
			CoffeeID:   item.Coffee.ID,
			CoffeeName: item.Coffee.Name,
			Quantity:   item.Quantity,
			Price:      item.Coffee.Price * float64(item.Quantity),
		}
	}

	return &OrderResponse{
		ID:        order.ID,
		Items:     itemResponses,
		Total:     order.Total,
		Status:    order.Status,
		CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// decodeJSON = Parse JSON from request body
func decodeJSON(r *http.Request, v interface{}) error {
	// Limit body size เพื่อป้องกัน DOS attack
	r.Body = io.NopCloser(io.LimitReader(r.Body, 1<<20)) // 1MB

	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
