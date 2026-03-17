// delivery/http/response.go
// Helper functions สำหรับส่ง JSON response
// ห้ามใช้ได้เฉพาะ HTTP layer เท่านั้น

package http

import (
	"encoding/json"
	"net/http"
)

// successResponse = โครงสร้าง JSON ตอบสำเร็จ
type successResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// errorResponse = โครงสร้าง JSON ตอบ error
type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// respondSuccess = ส่ง JSON response สำเร็จ
func respondSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := successResponse{
		Status:  "success",
		Message: "OK",
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// respondError = ส่ง JSON response error
func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := errorResponse{
		Status:  "error",
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}
