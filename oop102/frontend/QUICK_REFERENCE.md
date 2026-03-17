# 📂 Complete Frontend File Structure

```
frontend/
│
├── 📄 README.md                    ← Architecture & best practices
├── 📄 ARCHITECTURE.md              ← Deep dive into design
├── 📄 PROJECT_SUMMARY.md           ← Quick reference (this file)
├── 📄 .env.example                 ← Environment config template
├── 📄 package.json                 ← Dependencies & scripts
├── 📄 tsconfig.json                ← TypeScript config
│
├── 📁 public/
│   └── 📄 index.html               ← HTML entrypoint
│
└── 📁 src/
    │
    ├── 📄 index.tsx                ← React root (ReactDOM.render)
    ├── 📄 index.css
    │
    ├── 📄 App.tsx                  ← Main React component
    ├── 📄 App.css
    │
    ├── 📁 api/                     ← 🔌 API SERVICE LAYER
    │   ├── 📄 coffeeApi.ts         ← Fetch coffee menu
    │   ├── 📄 orderApi.ts          ← Create orders
    │   └── 📄 index.ts             ← Export all API functions
    │
    ├── 📁 config/                  ← ⚙️ CONFIGURATION
    │   └── 📄 api.ts               ← Base URL & endpoints
    │
    ├── 📁 types/                   ← 📋 TYPESCRIPT TYPES
    │   └── 📄 index.ts             ← All interfaces (Coffee, Order, etc)
    │
    ├── 📁 components/              ← 🎨 PRESENTATIONAL COMPONENTS
    │   ├── 📄 MenuList.tsx         ← Display coffee grid
    │   ├── 📄 MenuList.css
    │   │
    │   ├── 📄 OrderForm.tsx        ← Display shopping cart
    │   ├── 📄 OrderForm.css
    │   │
    │   ├── 📄 OrderSummary.tsx     ← Display confirmation modal
    │   ├── 📄 OrderSummary.css
    │   │
    │   └── 📄 index.ts             ← Export all components
    │
    └── 📁 pages/                   ← 📄 PAGE LAYER (ORCHESTRATION)
        ├── 📄 HomePage.tsx         ← Main orchestration page
        ├── 📄 HomePage.css
        └── 📄 index.ts             ← Export all pages
```

---

## 📋 File-by-File Overview

### Root Configuration Files

| File | Purpose | Content |
|------|---------|---------|
| `package.json` | Node dependencies & scripts | React, TypeScript, dev tools |
| `tsconfig.json` | TypeScript config | Compiler options |
| `.env.example` | Environment template | `REACT_APP_API_URL` |

### Public

| File | Purpose |
|------|---------|
| `index.html` | HTML wrapper for React app |

### Source Code (`src/`)

#### Entry Points
- `index.tsx` — Renders React app into DOM
- `App.tsx` — Main React component (wraps HomePage)

#### API Layer (`src/api/`)
```
coffeeApi.ts    → fetchCoffees()
orderApi.ts     → createOrder()
index.ts        → export * from './coffeeApi'
                  export * from './orderApi'
```

#### Configuration (`src/config/`)
```
api.ts → API_CONFIG.BASE_URL
         API_ENDPOINTS.COFFEES, ORDERS
```

#### Types (`src/types/`)
```
index.ts → Coffee, Order, OrderItem, APIResponse, etc
           (All TypeScript interfaces used across app)
```

#### Components (`src/components/`)
```
MenuList.tsx      → {coffees} → UI   (no API calls)
OrderForm.tsx     → {items}   → UI   (no API calls)
OrderSummary.tsx  → {order}   → UI   (no API calls)
index.ts          → export all components
```

Each component has matching `.css` file for styling.

#### Pages (`src/pages/`)
```
HomePage.tsx → Orchestrates everything
              - State: cartItems, completedOrder
              - Functions: handleSelectCoffee, handlePlaceOrder
              - Renders: MenuList, OrderForm, OrderSummary
index.ts     → export { HomePage }
```

---

## 🔄 Component Dependency Tree

```
App
└── HomePage (Orchestrator)
    ├── MenuList (Pure UI)
    ├── OrderForm (Pure UI)
    └── OrderSummary (Pure UI)
```

---

## 📊 Data Flow Visualization

```
User Interaction (click buttons)
    ↓
Component callback fired (onSelectCoffee, onSubmit)
    ↓
HomePage handler (handleSelectCoffee, handlePlaceOrder)
    ├─ Maybe: Call API (api/coffeeApi.ts or orderApi.ts)
    └─ Maybe: Update state (cartItems, completedOrder)
    ↓
Component receives new props
    ↓
Component re-renders
    ↓
User sees updated UI
```

---

## 🎯 Architecture Layers

### Layer 1: API (`api/`)
```typescript
import { API_CONFIG, API_ENDPOINTS } from '../config';
import { Coffee, Order, APIResponse } from '../types';

// ONE JOB: Fetch data from backend
export async function fetchCoffees(): Promise<Coffee[]> {
  const response = await fetch(url);
  const data: APIResponse<Coffee[]> = await response.json();
  return data.data;
}
```

### Layer 2: Components (`components/`)
```typescript
interface Props {
  // Receive data & callbacks from parent
  onSelectCoffee: (coffee: Coffee, qty: number) => void;
}

// ONE JOB: Render UI + call callbacks
export function MenuList({ onSelectCoffee }: Props) {
  return <button onClick={() => onSelectCoffee(coffee, qty)}>
    Add to Order
  </button>;
}
```

### Layer 3: Pages (`pages/`)
```typescript
// ONE JOB: Orchestrate components + API + state
export function HomePage() {
  const [cartItems, setCartItems] = useState([]);
  
  // When component calls onSelectCoffee
  const handleSelectCoffee = (coffee, qty) => {
    setCartItems(prev => [...prev, { coffee, qty }]);
  };
  
  // When component calls onSubmit
  const handlePlaceOrder = async () => {
    const order = await createOrder(cartItems);
    setCompletedOrder(order);
  };
  
  return (
    <>
      <MenuList onSelectCoffee={handleSelectCoffee} />
      <OrderForm onSubmit={handlePlaceOrder} ... />
      <OrderSummary order={completedOrder} ... />
    </>
  );
}
```

---

## 🚀 Running the App

```bash
# Install
npm install

# Environment
cp .env.example .env.local
# Edit: REACT_APP_API_URL=http://localhost:8080

# Start
npm start              # Frontend on :3000

# In another terminal:
cd ../../backend/cmd/api
go run main.go         # Backend on :8080
```

---

## ✅ Code Quality Checklist

- ✅ All API calls in `api/` folder
- ✅ Components have no fetch() calls
- ✅ All state in HomePage
- ✅ One-way data flow (top-down)
- ✅ Props for data, callbacks for events
- ✅ All types in `types/index.ts`
- ✅ Configuration centralized in `config/`
- ✅ Each component in own file
- ✅ CSS co-located with components
- ✅ Error handling in API layer

---

## 🔍 Quick File Reference

**Need to...**

| Task | File | Line Numbers |
|------|------|--------------|
| Add API endpoint | `api/coffeeApi.ts` | ~30-40 |
| Change API URL | `config/api.ts` | Line 3 |
| Add new type | `types/index.ts` | Add interface |
| Modify component | `components/*.tsx` | Whole file |
| Add business logic | `pages/HomePage.tsx` | Handler functions |
| Style component | `components/*.css` | CSS file |
| Configure app | `.env.local` | `REACT_APP_API_URL` |

---

## 🎓 Patterns Used

### API Service Pattern
```typescript
// One function = One API endpoint
export async function fetchCoffees(): Promise<Coffee[]>
export async function createOrder(req): Promise<Order>
```

### Component Pattern
```typescript
// Receive data & callbacks as props
interface Props {
  items: OrderItem[];
  onRemoveItem: (id: string) => void;
  onSubmit: () => void;
}

// Emit events to parent, don't manage parent's state
<button onClick={() => onRemoveItem(id)}>Remove</button>
```

### Page/Container Pattern
```typescript
// Manage all state here
const [cartItems, setCartItems] = useState([]);

// Write all handlers here
const handleSelectCoffee = (coffee) => { ... };

// Connect everything here
<MenuList onSelectCoffee={handleSelectCoffee} />
```

---

## 🧠 State Management Strategy

All state lives in **HomePage** (component tree root):

```typescript
const [cartItems, setCartItems] = useState([]);
const [completedOrder, setCompletedOrder] = useState(null);
const [isSubmitting, setIsSubmitting] = useState(false);
const [submitError, setSubmitError] = useState(null);
```

Pass down via **props**:
```typescript
<MenuList onSelectCoffee={...} />
<OrderForm items={cartItems} onRemove={...} />
<OrderSummary order={completedOrder} />
```

Child components communicate back via **callbacks**:
```typescript
// Child calls this
onSelectCoffee(coffee, quantity)

// Parent updates state
setCartItems(prev => [...prev, ...])
```

---

## 🎯 Decision Points

| Scenario | Decision | Why |
|----------|----------|-----|
| Where to put fetch? | In `api/` folder | Central place, reusable, testable |
| Where to put logic? | In HomePage | Close to state management |
| Where to put types? | In `types/index.ts` | One source of truth |
| How to pass data? | Props down | Clear flow, easy to debug |
| How to handle events? | Callbacks up | Clean component interface |

---

## 🚨 Common Mistakes (Don't Do)

```typescript
// ❌ Don't fetch inside component
function MenuList() {
  useEffect(() => {
    fetch('/coffees')  // NO!
  }, []);
}

// ❌ Don't manage same state in multiple components
// In MenuList:
const [select, setSelected] = useState([]);  // NO!
// In OrderForm:
const [order, setOrder] = useState([]);      // NO!

// ❌ Don't send prop to prop to prop (prop drilling)
<App cartItems={items}>
  <Section cartItems={items}>
    <Component cartItems={items} />  // 3 levels deep - bad!

// ✅ Do: Use context or state management library for deep nesting
```

---

## 📚 File Purposes Quick Reference

- **API files**: "How to get data from backend"
- **Component files**: "How to show UI for this feature"
- **Page file**: "How all components & APIs work together"
- **Type files**: "What shape is the data"
- **Config files**: "Where is the backend"
- **CSS files**: "How does this look"

---

## 🎉 Summary

```
🏗️  Clean Architecture + React
🔌  Separated API layer
🎨  Pure UI components
📄  Orchestration page
📋  TypeScript types
✨  Easy to maintain & extend
```

You now have a frontend that:
- ✅ Calls backend API
- ✅ Shows menu & order form
- ✅ Submits orders
- ✅ Displays confirmation
- ✅ Follows best practices
- ✅ Is testable & maintainable
