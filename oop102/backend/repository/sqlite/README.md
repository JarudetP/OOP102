// repository/sqlite/README.md
# SQLite Repository Implementation

## Usage

```go
package main

import (
	"coffee-shop/oop102/backend/repository/sqlite"
	"coffee-shop/oop102/backend/usecase"
	"coffee-shop/oop102/backend/delivery/http"
)

func main() {
	// 1. Initialize SQLite database
	db, err := sqlite.InitDB("./coffee.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 2. Create repositories with SQLite
	coffeeRepo := sqlite.NewSQLiteCoffeeRepo(db)
	orderRepo := sqlite.NewSQLiteOrderRepo(db)  // ✅ No coffeeRepo coupling!

	// 3. Create use case (same as before - no changes!)
	orderUseCase := usecase.NewOrderUseCase(coffeeRepo, orderRepo)

	// 4. Create HTTP handler
	router := http.NewRouter(orderUseCase)

	// 5. Start server
	http.ListenAndServe(":8080", router)
}
```

## Database Schema

The SQLite database will be initialized with the following tables:

- **coffees**: Menu items (id, name, price, emoji)
- **orders**: Order headers (id, total, status, created_at)
- **order_items**: Order line items (id, order_id, coffee_id, quantity)

## Key Points

✅ Same interface as in-memory repository
✅ All SQL queries inside this layer
✅ UseCase doesn't change at all
✅ Handler doesn't change at all
✅ Only `main.go` needs to change which repository to use
