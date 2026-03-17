# Coffee Shop API - Backend

HTTP API Backend สำหรับร้านกาแฟ ใช้ Clean Architecture

## 🏗️ Architecture Layers

```
┌─────────────────────────────────────────┐
│  HTTP Handler (delivery/http/)          │
│  • Parse JSON + ส่ง JSON                 │
│  • ไม่รู้เรื่อง Business Logic           │
└────────────┬────────────────────────────┘
             │ (calls)
┌────────────▼────────────────────────────┐
│  UseCase (usecase/)                     │
│  • Business Rules & Logic               │
│  • ไม่รู้เรื่อง Database หรือ HTTP      │
└────────────┬────────────────────────────┘
             │ (calls interface)
┌────────────▼────────────────────────────┐
│  Repository (repository/sqlite/)        │
│  • SQL Queries                          │
│  • Implement domain interfaces          │
└────────────┬────────────────────────────┘
             │
┌────────────▼────────────────────────────┐
│  Domain (domain/)                       │
│  • Business Models & Rules              │
│  • ไม่import อะไรจากภายนอก              │
└─────────────────────────────────────────┘
```

## 🚀 How to Run

### 1. ติดตั้ง dependencies

```bash
go get github.com/mattn/go-sqlite3
```

### 2. เตรียม environment variables (optional)

```bash
# ถ้าไม่ตั้งค่า จะใช้ defaults
export PORT=8080
export DATABASE_PATH=./coffee.db
```

### 3. รัน server

```bash
cd cmd/api
go run main.go
```

### 4. ลอง API

```bash
# ดึงเมนูกาแฟ
curl http://localhost:8080/coffees

# สั่งกาแฟ
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"coffee_id": "1", "quantity": 2},
      {"coffee_id": "3", "quantity": 1}
    ]
  }'
```

## 📋 API Endpoints

### GET /coffees
ดึงเมนูกาแฟทั้งหมด

**Response:**
```json
{
  "status": "success",
  "message": "OK",
  "data": [
    {"id": "1", "name": "Latte", "price": 65, "emoji": "☕"},
    {"id": "2", "name": "Mocha", "price": 75, "emoji": "☕"}
  ]
}
```

### POST /orders
สั่งกาแฟใหม่

**Request:**
```json
{
  "items": [
    {"coffee_id": "1", "quantity": 2},
    {"coffee_id": "2", "quantity": 1}
  ]
}
```

**Response:**
```json
{
  "status": "success",
  "message": "OK",
  "data": {
    "id": "ORD-001",
    "items": [
      {"coffee_id": "1", "coffee_name": "Latte", "quantity": 2, "price": 130},
      {"coffee_id": "2", "coffee_name": "Mocha", "quantity": 1, "price": 75}
    ],
    "total": 205,
    "status": "กำลังเตรียม ☕",
    "created_at": "2026-03-17 10:30:00"
  }
}
```

## 🔌 main.go - Wiring Layer

`cmd/api/main.go` เป็น **Wiring Layer** ที่ต่อ dependencies ทั้งหมด:

### สิ่งที่ main.go ทำ:
```
1. เปิด SQLite database
   db := sqlite.InitDB("./coffee.db")

2. สร้าง Repositories (data access layer)
   coffeeRepo := sqlite.NewSQLiteCoffeeRepo(db)
   orderRepo := sqlite.NewSQLiteOrderRepo(db, coffeeRepo)

3. สร้าง UseCase (business logic layer)
   usecase := usecase.NewOrderUseCase(coffeeRepo, orderRepo)

4. สร้าง HTTP Router (delivery layer)
   router := http.NewRouter(usecase)

5. เริ่ม HTTP server
   http.ListenAndServe(":8080", router)
```

### สิ่งที่ main.go ไม่ควรทำ:
- ❌ ไม่มี SQL queries
- ❌ ไม่มี business logic
- ❌ ไม่มี HTTP-specific code
- ✅ มีแค่ Dependency Injection

## 📦 Project Structure

```
oop102/backend/
├── cmd/api/
│   └── main.go                    ◄─ Entry point (Wiring Layer)
├── domain/                        ◄─ Business Models & Rules
│   ├── coffee.go
│   ├── order.go
│   └── repository.go             (Interfaces only)
├── usecase/                       ◄─ Business Logic
│   └── order_usecase.go
├── repository/                    ◄─ Data Access
│   ├── coffee_repo.go            (In-Memory - old)
│   ├── order_repo.go             (In-Memory - old)
│   └── sqlite/                   (SQLite - new)
│       ├── db.go
│       ├── coffee_repo.go
│       └── order_repo.go
├── delivery/http/                 ◄─ HTTP Interface
│   ├── handler.go                (Parse JSON + Call UseCase)
│   ├── response.go               (Format JSON)
│   └── router.go                 (Setup routes)
└── go.mod
```

## 🔄 Request Flow

```
Client HTTP Request
    ↓
HTTP Handler (delivery/http)
    ├─ Parse JSON
    ├─ Validate format
    └─ Call UseCase
    ↓
UseCase (business logic)
    ├─ Validate business rules
    ├─ Call Repository
    ↓
Repository (SQLite)
    ├─ Execute SQL
    ├─ Return domain models
    ↓
UseCase (returns result)
    ↓
Handler (convert to JSON)
    ↓
HTTP Response
```

## 💡 Key Principles

### 1️⃣ Dependency Injection
- ✅ UseCase ไม่รู้ว่า Repository มาจาก SQLite หรือ In-Memory
- ✅ Handler ไม่รู้เรื่อง Database
- ✅ Domain ไม่ import อะไรจากภายนอก

### 2️⃣ Separation of Concerns
- ✅ Domain = Business Models only
- ✅ UseCase = Business Rules only
- ✅ Repository = Data Access only
- ✅ Handler = HTTP Interface only
- ✅ main.go = Wiring only

### 3️⃣ Easy to Switch Implementation
```go
// เปลี่ยนจาก SQLite → PostgreSQL แค่เปลี่ยน main.go
// coffeeRepo := sqlite.NewSQLiteCoffeeRepo(db)
coffeeRepo := postgres.NewPostgresCoffeeRepo(db)  // ✅ same interface!
```

### 4️⃣ Easy to Test
```go
// ใช้ Mock repositories สำหรับ Unit Test
mockCoffeeRepo := &MockCoffeeRepo{...}
usecase := usecase.NewOrderUseCase(mockCoffeeRepo, ...)
```

## 🛠️ Development Tips

### ใช้ In-Memory Repository (สำหรับ development เร็ว)
```go
// ใน main.go แทน SQLite
coffeeRepo := repository.NewInMemoryCoffeeRepo()
orderRepo := repository.NewInMemoryOrderRepo()
```

### ดู Database ที่สร้างขึ้น
```bash
# ใช้ sqlite3 CLI
sqlite3 coffee.db

# ใน SQLite prompt
.tables
SELECT * FROM coffees;
SELECT * FROM orders;
```

## 📝 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DATABASE_PATH` | `./coffee.db` | SQLite database file path |

## ✨ Features

- ✅ Clean Architecture (Domain, UseCase, Repository, Delivery)
- ✅ Dependency Injection
- ✅ SQLite Database
- ✅ RESTful API
- ✅ JSON Request/Response
- ✅ Transaction Support (for Orders)
- ✅ Error Handling
- ✅ Easy to Test
- ✅ Easy to Switch Database

## 🐛 Troubleshooting

### Error: "Failed to implement interface"
- ตรวจสอบว่า Repository implement ครบทุก method ใน domain interface

### Error: "Foreign key constraint failed"
- ตรวจสอบว่า Coffee ID ที่ส่ง exists ในตาราง coffees

### Database file grows large
- SQLite stores data in single file, ลบ `coffee.db` เพื่อยืนใหม่

## 📚 References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go by Example](https://gobyexample.com/)
- [SQLite with Go](https://github.com/mattn/go-sqlite3)
