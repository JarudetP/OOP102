// cmd/api/main.go
// Wiring Layer - ต่อแต่ละ layer เข้าด้วยกัน
// ไม่มี business logic, ไม่มี SQL, ไม่มี HTTP-specific code
// มีแค่ Dependency Injection และ Initialization

package main

import (
	"coffee-shop/oop102/backend/delivery/http"
	"coffee-shop/oop102/backend/repository/sqlite"
	"coffee-shop/oop102/backend/usecase"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// ========== LAYER 0: Database Connection ==========
	// เปิด SQLite connection
	dbPath := getDBPath()
	db, err := sqlite.InitDB(dbPath)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer db.Close()
	fmt.Printf("✓ Database connected: %s\n", dbPath)

	// ========== LAYER 1: Repository (Data Access) ==========
	// สร้าง repositories ที่ implement domain interfaces
	coffeeRepo := sqlite.NewSQLiteCoffeeRepo(db)
	orderRepo := sqlite.NewSQLiteOrderRepo(db, coffeeRepo)
	fmt.Println("✓ SQLite Repositories initialized")

	// ========== LAYER 2: UseCase (Business Logic) ==========
	// inject repositories เข้า usecase
	// usecase ไม่รู้ว่า repository มาจาก SQLite หรือ In-Memory
	orderUseCase := usecase.NewOrderUseCase(coffeeRepo, orderRepo)
	fmt.Println("✓ UseCase initialized")

	// ========== LAYER 3: Delivery/HTTP (API Interface) ==========
	// inject usecase เข้า HTTP handler
	router := http.NewRouter(orderUseCase)
	fmt.Println("✓ HTTP Router initialized")

	// ========== LAYER 4: Start Server ==========
	// ดึง port จาก environment หรือใช้ default
	port := getPort()
	addr := fmt.Sprintf(":%s", port)

	fmt.Println("==========================================")
	fmt.Printf("🚀 Coffee Shop API running on http://localhost%s\n", addr)
	fmt.Println("==========================================")
	fmt.Printf("📝 Endpoints:\n")
	fmt.Printf("   GET  http://localhost%s/coffees      (ดึงเมนู)\n", addr)
	fmt.Printf("   POST http://localhost%s/orders       (สั่งกาแฟ)\n", addr)
	fmt.Println("==========================================")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}

// ============ Helper Functions ============

// getPort = ดึง port จาก environment หรือใช้ default
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// getDBPath = ดึง database path จาก environment หรือใช้ default
// Path: ./coffee.db หรือจาก DATABASE_PATH env var
func getDBPath() string {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./coffee.db"
	}
	return dbPath
}
