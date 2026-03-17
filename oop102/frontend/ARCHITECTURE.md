# Frontend Architecture Overview

## 🎯 Architecture Philosophy: Separation of Concerns

The frontend follows a **3-layer architecture** to keep the codebase clean and maintainable:

```
┌─────────────────────────────────────────────────┐
│          UI Layer (Components)                  │
│  MenuList | OrderForm | OrderSummary            │
│  ↓ Function calls, Props, State management      │
└────────────────────┬────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────┐
│     Orchestration Layer (HomePage)              │
│  State management, Business logic coordination  │
│        ↓ API function calls                     │
└────────────────────┬────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────┐
│         API Layer (api/)                        │
│  coffeeApi.ts | orderApi.ts                     │
│  ↓ HTTP requests to backend                     │
└────────────────────┬────────────────────────────┘
                     │
                 Backend API
                (Go Backend)
```

---

## 📁 Directory Structure

```
frontend/
├── public/
│   └── index.html                    ← React app entrypoint
├── src/
│   ├── api/                          ← 🔌 API Service Layer
│   │   ├── coffeeApi.ts              (Functions that call backend)
│   │   ├── orderApi.ts
│   │   └── index.ts
│   │
│   ├── components/                   ← 🎨 UI Components Layer
│   │   ├── MenuList.tsx              (Display menu - no logic)
│   │   ├── MenuList.css
│   │   ├── OrderForm.tsx             (Show cart - no logic)
│   │   ├── OrderForm.css
│   │   ├── OrderSummary.tsx          (Show confirmation - no logic)
│   │   ├── OrderSummary.css
│   │   └── index.ts
│   │
│   ├── config/                       ← ⚙️ Configuration
│   │   └── api.ts                    (API base URL & endpoints)
│   │
│   ├── pages/                        ← 📄 Page Layer (Orchestration)
│   │   ├── HomePage.tsx              (Main logic & coordination)
│   │   ├── HomePage.css
│   │   └── index.ts
│   │
│   ├── types/                        ← 📋 TypeScript Types
│   │   └── index.ts                  (Shared interfaces)
│   │
│   ├── App.tsx
│   ├── App.css
│   ├── index.tsx                     (React root)
│   └── index.css
│
├── .env.example
├── package.json
├── tsconfig.json
└── README.md
```

---

## 🔄 Data Flow Example: Placing an Order

```
1. User selects coffee quantity and clicks "Add to Order"
   ↓
2. MenuList component calls onSelectCoffee() callback
   ↓
3. HomePage.handleSelectCoffee() updates cartItems state
   ↓
4. OrderForm component receives new items via props
   ↓
5. OrderForm re-renders showing updated cart
   ↓
6. User clicks "Place Order" button
   ↓
7. OrderForm calls onSubmit() callback
   ↓
8. HomePage.handlePlaceOrder():
   ├─ Creates CreateOrderRequest from cartItems
   ├─ Calls createOrder(request) from API layer
   ├─ API layer calls fetch() to backend
   ├─ Backend processes and returns Order
   └─ HomePage updates completedOrder state
   ↓
9. OrderSummary component receives order via props
   ↓
10. OrderSummary modal displays confirmation
```

---

## 📦 Layer Responsibilities

### API Layer (`api/`)
**Domain:** Backend communication only

```typescript
// ✅ Responsibility: Fetch and return data
export async function fetchCoffees(): Promise<Coffee[]> {
  const response = await fetch('/coffees');
  const data = await response.json();
  return data.data;
}

// ✅ Responsibility: Handle errors
// ❌ Not: Business logic
// ❌ Not: UI decisions
// ❌ Not: State management
```

### UI Components (`components/`)
**Domain:** Rendering and user interaction only

```typescript
// ✅ Responsibility: Render UI
function MenuList({ onSelectCoffee }) {
  return <button onClick={() => onSelectCoffee(coffee)}>Add</button>;
}

// ✅ Responsibility: Show loading/error states
// ❌ Not: Call API
// ❌ Not: Manage global state
// ❌ Not: Business decisions
```

### Page/Container (`pages/`)
**Domain:** Orchestration and business logic

```typescript
// ✅ Responsibility: Manage state
// ✅ Responsibility: Coordinate components
// ✅ Responsibility: Call API functions
// ✅ Responsibility: Make business decisions
export function HomePage() {
  const [cartItems, setCartItems] = useState([]);
  
  const handlePlaceOrder = async () => {
    const order = await createOrder(cartItems);  // API call
    setCompletedOrder(order);                    // Update state
  };
  
  return (
    <>
      <MenuList onSelectCoffee={...} />
      <OrderForm onSubmit={handlePlaceOrder} />
      <OrderSummary order={completedOrder} />
    </>
  );
}
```

---

## 🌐 API Contract

### Backend Endpoints

```
GET /coffees
└─ Returns: { status: "success", data: Coffee[] }

POST /orders
├─ Body: { items: [{coffee_id, quantity}] }
└─ Returns: { status: "success", data: Order }
```

### Type Definitions (src/types/index.ts)

```typescript
interface Coffee {
  id: string;
  name: string;
  price: number;
  emoji: string;
}

interface Order {
  id: string;
  items: OrderItemResponse[];
  total: number;
  status: string;
  created_at: string;
}
```

---

## 🎛️ State Management Strategy

### State Location: Top-Down Approach

```
HomePage (main state holder)
  ├─ cartItems: OrderItem[]         (managed here)
  ├─ completedOrder: Order | null   (managed here)
  ├─ isSubmitting: boolean          (managed here)
  ├─ submitError: string | null     (managed here)
  │
  └─ Pass to children:
      ├─ MenuList
      │  └─ Emits: onSelectCoffee()
      │
      ├─ OrderForm
      │  ├─ Receives: items, loading, error
      │  └─ Emits: onRemoveItem(), onQuantityChange(), onSubmit()
      │
      └─ OrderSummary
         ├─ Receives: order
         └─ Emits: onNewOrder()
```

**Why top-down?**
- ✅ Single source of truth
- ✅ Data flows one direction (easier to debug)
- ✅ Components stay pure (same props = same output)
- ✅ Easy to test in isolation

---

## 🔌 API Service Pattern

### Best Practice: Custom Hook (Future)

```typescript
// hooks/useCoffees.ts (Not yet implemented, for reference)
export function useCoffees() {
  const [coffees, setCoffees] = useState<Coffee[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    fetchCoffees()
      .then(setCoffees)
      .catch(setError)
      .finally(() => setLoading(false));
  }, []);

  return { coffees, loading, error };
}

// In component:
function MenuList() {
  const { coffees, loading, error } = useCoffees();
  return ...;
}
```

---

## 🧪 Testing Strategy

### Test API Layer (Unit)
```typescript
// api/coffeeApi.test.ts
test('fetchCoffees calls correct endpoint', async () => {
  global.fetch = jest.fn(() =>
    Promise.resolve({
      ok: true,
      json: () => Promise.resolve({
        status: 'success',
        data: [mockCoffee]
      })
    })
  );

  const result = await fetchCoffees();
  expect(fetch).toHaveBeenCalledWith('http://localhost:8080/coffees');
  expect(result).toEqual([mockCoffee]);
});
```

### Test Components (Isolated)
```typescript
// components/MenuList.test.tsx
test('calls onSelectCoffee when button clicked', async () => {
  const mockHandler = jest.fn();
  
  jest.mock('../api', () => ({
    fetchCoffees: jest.fn().mockResolvedValue([mockCoffee])
  }));

  render(<MenuList onSelectCoffee={mockHandler} />);
  const button = screen.getByText('Add to Order');
  fireEvent.click(button);

  expect(mockHandler).toHaveBeenCalledWith(mockCoffee, 1);
});
```

### Test Page (Integration)
```typescript
// pages/HomePage.test.tsx
test('full order flow works', async () => {
  render(<HomePage />);

  // Click add to order
  const addBtn = await screen.findByText('Add to Order');
  fireEvent.click(addBtn);

  // Check cart updated
  expect(screen.getByText('Latte')).toBeInTheDocument();

  // Submit order
  const submitBtn = screen.getByText('Place Order');
  fireEvent.click(submitBtn);

  // Check confirmation
  await waitFor(() => {
    expect(screen.getByText('Order Confirmed!')).toBeInTheDocument();
  });
});
```

---

## 🚀 Getting Started

### 1. Install Dependencies
```bash
cd frontend
npm install
```

### 2. Configure Backend URL
```bash
# Copy example to .env.local
cp .env.example .env.local

# Edit .env.local:
REACT_APP_API_URL=http://localhost:8080
```

### 3. Start Development Server
```bash
npm start
# Opens http://localhost:3000
```

### 4. Build for Production
```bash
npm run build
# Outputs to build/
```

---

## 🔐 Error Handling

### API Layer
```typescript
try {
  const response = await fetch(url);
  if (!response.ok) throw new Error(`${response.status}`);
  const data = await response.json();
  return data.data;
} catch (error) {
  console.error('API Error:', error);
  throw error;  // Re-throw for component to handle
}
```

### Component Layer
```typescript
const [error, setError] = useState<string | null>(null);

useEffect(() => {
  fetchCoffees()
    .then(setCoffees)
    .catch(err => {
      const msg = err instanceof Error ? err.message : 'Unknown error';
      setError(msg);
    });
}, []);

if (error) return <div>{error}</div>;
```

### Page Layer
```typescript
try {
  const order = await createOrder(cartItems);
  setCompletedOrder(order);
  setCartItems([]);  // Clear cart
} catch (error) {
  const msg = error instanceof Error ? error.message : 'Failed';
  setSubmitError(msg);  // Show to user
}
```

---

## 🎓 Key Takeaways

| Principle | Why | How |
|-----------|-----|-----|
| **Separation of Concerns** | Easier to maintain & test | API / Components / Pages |
| **Unidirectional Data** | Less bugs & easier to debug | Props down, callbacks up |
| **Pure Functions** | Testable & predictable | Same input = same output |
| **Single Responsibility** | Easy to change | Each layer has 1 job |
| **No Business Logic in Components** | Reusable components | Logic in pages/hooks |

---

## 📚 Further Reading

- [React Hooks](https://react.dev/reference/react)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Clean Architecture (UI Layer)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Fetch API](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API)
