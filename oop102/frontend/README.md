# Coffee Shop Frontend - React Architecture Guide

## 🏗️ Project Structure

```
frontend/src/
├── api/                          # 🔌 API Service Layer
│   ├── coffeeApi.ts              # Functions to call coffee endpoints
│   ├── orderApi.ts               # Functions to call order endpoints
│   └── index.ts                  # Export all API functions
├── components/                   # 🎨 Presentation Components
│   ├── MenuList.tsx              # Display coffee menu
│   ├── MenuList.css
│   ├── OrderForm.tsx             # Shopping cart UI
│   ├── OrderForm.css
│   ├── OrderSummary.tsx          # Confirmation modal
│   ├── OrderSummary.css
│   └── index.ts                  # Export all components
├── config/                       # ⚙️ Configuration
│   └── api.ts                    # API endpoints & base URL
├── types/                        # 📋 TypeScript Interfaces
│   └── index.ts                  # Shared types for API
├── pages/                        # 📄 Page Components
│   ├── HomePage.tsx              # Main orchestration page
│   └── HomePage.css              # Page styling
└── App.tsx                       # Main React app
```

---

## 🔄 Architecture Layers

### Layer 1: API Service (`api/`)
**Purpose:** Handle all backend communication

```typescript
// ✅ Good: Centralized API logic
import { fetchCoffees, createOrder } from '../api';

const coffees = await fetchCoffees();  // Component calls this
```

**Why separate?**
- UI components stay dumb (no fetch() calls)
- Easy to change API endpoints (one place)
- Easy to add debugging/logging
- Easy to mock for testing

### Layer 2: Components (`components/`)
**Purpose:** Pure UI rendering (no API calls)

```typescript
// ✅ Good: Components receive props, emit events
<MenuList onSelectCoffee={handleSelectCoffee} />

// ❌ Bad: Don't fetch here
const handleClick = () => {
  fetch('/api/coffees')  // ❌ NO!
}
```

**Component Responsibilities:**
- Render UI
- Handle user input (clicks, typing)
- Call callback functions (props)
- Manage local UI state (loading, error)

### Layer 3: Page/Container (`pages/`)
**Purpose:** Orchestrate components + API layer

```typescript
// HomePage = "Brain" of the app
// - Manages cart state
// - Calls API functions
// - Passes data to components
// - Handles business logic (add to cart, etc)
```

**Page Responsibilities:**
- State management (cart items)
- API orchestration (when to call API)
- Event handling (linking components together)
- Business logic decisions

### Layer 4: Types (`types/`)
**Purpose:** Shared TypeScript interfaces

```typescript
// One source of truth for all data shapes
interface Coffee {
  id: string;
  name: string;
  price: number;
  emoji: string;
}
```

---

## 📊 Data Flow

### ✅ Correct Data Flow (Unidirectional)

```
User clicks "Add to Order"
    ↓
MenuList component emits onSelectCoffee()
    ↓
HomePage.handleSelectCoffee() 
    ├─ Updates cartItems state
    └─ Passes to OrderForm via props
    ↓
OrderForm component re-renders
    ├─ Shows updated cart
    └─ Emits onSubmit()
    ↓
HomePage.handlePlaceOrder()
    ├─ Calls API: createOrder()
    ├─ API returns Order
    ├─ Updates completedOrder state
    └─ OrderSummary modal shows
```

### ❌ Wrong Data Flow (Avoid)

```
Component directly fetches:
Component
  └─ fetch('/api/coffees')  ❌ NO!

Multiple components manage same state:
MenuList → state
OrderForm → state  ❌ Causes bugs!
HomePage → state

Component calls API + handles business logic:
Component
  ├─ fetch() ❌
  ├─ filter() ❌
  └─ calculate() ❌
```

---

## 🎯 Component Responsibilities

### MenuList Component
```typescript
Receives: onSelectCoffee(coffee, quantity) callback
Does:
  ✅ Fetch coffees on mount (via API layer)
  ✅ Display coffee grid
  ✅ Handle quantity input
  ✅ Call onSelectCoffee() when user clicks
Does NOT:
  ❌ Know about global cart state
  ❌ Know about order submission
  ❌ Save to database
```

### OrderForm Component
```typescript
Receives:
  ✅ items (cart items)
  ✅ onRemoveItem() callback
  ✅ onQuantityChange() callback
  ✅ onSubmit() callback
Does:
  ✅ Display cart items
  ✅ Show total price
  ✅ Handle quantity changes
  ✅ Show loading state
  ✅ Show error messages
Does NOT:
  ❌ Call API
  ❌ Manage global state
```

### OrderSummary Component
```typescript
Receives:
  ✅ order (completed order data)
  ✅ onNewOrder() callback
Does:
  ✅ Display confirmation
  ✅ Show order details
Does NOT:
  ❌ Call API
  ❌ Modify state
```

### HomePage
```typescript
Manages:
  ✅ cartItems state
  ✅ completedOrder state
  ✅ isSubmitting flag
  ✅ submitError message

Handles:
  ✅ handleSelectCoffee() - Add item to cart
  ✅ handleRemoveItem() - Remove from cart
  ✅ handleQuantityChange() - Update quantity
  ✅ handlePlaceOrder() - Call API + manage response
  ✅ handleNewOrder() - Reset for new order

Orchestrates:
  ✅ MenuList → API
  ✅ OrderForm + API → OrderSummary
```

---

## 🔌 API Layer Pattern

### Before (❌ Bad)
```typescript
// Component mixing concerns
function MenuList() {
  const [coffees, setCoffees] = useState([]);

  useEffect(() => {
    // ❌ API logic inside component
    fetch('/coffees')
      .then(r => r.json())
      .then(setCoffees);
  }, []);

  return <div>{...}</div>;
}
```

### After (✅ Good)
```typescript
// api/coffeeApi.ts - Pure API logic
export async function fetchCoffees(): Promise<Coffee[]> {
  const response = await fetch(`${API_CONFIG.BASE_URL}/coffees`);
  const data: APIResponse<Coffee[]> = await response.json();
  return data.data;
}

// Component - Pure UI
function MenuList({ onSelectCoffee }: Props) {
  const [coffees, setCoffees] = useState([]);

  useEffect(() => {
    // ✅ Call API function
    fetchCoffees().then(setCoffees);
  }, []);

  return <div>{...}</div>;
}
```

---

## 🎨 Component Props Pattern

### ✅ Good Component Props
```typescript
interface MenuListProps {
  // Action callbacks (functions)
  onSelectCoffee: (coffee: Coffee, quantity: number) => void;
}

interface OrderFormProps {
  // Data passed down
  items: OrderItem[];
  loading: boolean;
  error: string | null;

  // Actions passed down
  onRemoveItem: (coffeeId: string) => void;
  onQuantityChange: (coffeeId: string, quantity: number) => void;
  onSubmit: () => void;
}
```

This pattern:
- Components are **pure** (same props = same output)
- Components are **testable** (easy to mock props)
- Components are **reusable** (not tied to specific data source)

---

## 🧪 Testing Patterns

### Test Component in Isolation
```typescript
// MenuList doesn't need backend
test('MenuList renders coffee items', () => {
  const mockCoffees = [
    { id: '1', name: 'Latte', price: 65, emoji: '☕' }
  ];

  // Mock the API
  jest.mock('../api', () => ({
    fetchCoffees: jest.fn(() => Promise.resolve(mockCoffees))
  }));

  render(<MenuList onSelectCoffee={jest.fn()} />);
  // assertions...
});
```

### Test Page Orchestration
```typescript
test('HomePage adds coffee to cart', async () => {
  render(<HomePage />);

  // User clicks "Add to Order"
  const addBtn = screen.getByText('Add to Order');
  fireEvent.click(addBtn);

  // Check if item appears in order
  await waitFor(() => {
    expect(screen.getByText('Latte')).toBeInTheDocument();
  });
});
```

---

## 🚀 Setup & Run

### Install Dependencies
```bash
npm install
# or
yarn install
```

### Configure API Base URL
Create `.env` file:
```env
REACT_APP_API_URL=http://localhost:8080
```

### Run Frontend
```bash
npm start
# or
yarn start
```

Frontend runs on: `http://localhost:3000`

---

## 🔗 Integration with Backend

### Backend API Expected
```
GET /coffees
  → { status: "success", data: [Coffee] }

POST /orders
  Body: { items: [{coffee_id, quantity}] }
  → { status: "success", data: Order }
```

### If API Changes
Change **one place** only:
- `src/config/api.ts` - Update endpoints
- `src/types/index.ts` - Update interfaces

All components will work with new API! ✅

---

## 💡 Key Principles

1. **Separation of Concerns**
   - API = fetch only
   - Component = render only
   - Page = orchestrate only

2. **Unidirectional Data Flow**
   - Parent → Child (props)
   - Child → Parent (callbacks)
   - Never side-ways or backwards

3. **Pure Components**
   - Same props = Same output
   - No side effects in render
   - Easy to test & debug

4. **Single Responsibility**
   - MenuList: Show menu
   - OrderForm: Show cart
   - OrderSummary: Show confirmation
   - HomePage: Coordinate all

5. **Future-Proof**
   - Change API? Update 1-2 files
   - Change layout? Update components
   - Change business logic? Update page
   - Easy to move to Next.js, Vue, etc

---

## 📝 Example: Adding New Feature

### Want to add "Favorites"?

1. **Add to API layer** (`api/`)
   ```typescript
   export async function saveFavorite(coffeeId: string) { ... }
   ```

2. **Add to types** (`types/`)
   ```typescript
   interface MenuItem extends Coffee {
     isFavorite: boolean;
   }
   ```

3. **Update component** (`components/`)
   ```typescript
   <favorite-button onClick={() => saveFavorite(coffee.id)} />
   ```

4. **Update page** (`pages/`)
   ```typescript
   const [favorites, setFavorites] = useState([]);
   const handleAddFavorite = (coffeeId) => { ... }
   ```

No other component breaks! ✅

---

## ✨ Best Practices

- ✅ One API file per feature (coffeeApi.ts, orderApi.ts)
- ✅ Pass callbacks from parent to child
- ✅ Use TypeScript for all data structures
- ✅ Keep components under 200 lines
- ✅ Import API functions near top of component
- ✅ Handle loading/error states
- ✅ Show user feedback (loading spinner, error message)
- ✅ Test components in isolation
- ✅ Mock API calls in tests
