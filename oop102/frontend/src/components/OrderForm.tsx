// src/components/OrderForm.tsx
// Component แสดงรายการสั่งที่เลือก
// 🛒 ไม่มี business logic - แค่ UI

import React from 'react';
import { OrderItem, Coffee } from '../types';
import './OrderForm.css';

interface OrderFormProps {
  items: Array<OrderItem & { coffeeName: string; coffeeEmoji: string; price: number }>;
  onRemoveItem: (coffeeId: string) => void;
  onQuantityChange: (coffeeId: string, quantity: number) => void;
  onSubmit: () => void;
  loading: boolean;
  error: string | null;
}

export const OrderForm: React.FC<OrderFormProps> = ({
  items,
  onRemoveItem,
  onQuantityChange,
  onSubmit,
  loading,
  error,
}) => {
  const total = items.reduce((sum, item) => sum + item.price * item.quantity, 0);

  return (
    <div className="order-form">
      <h2>🛒 Your Order</h2>

      {error && <div className="error-message">{error}</div>}

      {items.length === 0 ? (
        <p className="empty-cart">No items selected yet</p>
      ) : (
        <>
          <div className="order-items">
            {items.map((item) => (
              <div key={item.coffee_id} className="order-item">
                <div className="item-info">
                  <span className="emoji">{item.coffeeEmoji}</span>
                  <div>
                    <h4>{item.coffeeName}</h4>
                    <p className="item-price">฿{item.price}</p>
                  </div>
                </div>

                <div className="item-actions">
                  <div className="quantity-control">
                    <button
                      onClick={() =>
                        item.quantity > 1 &&
                        onQuantityChange(item.coffee_id, item.quantity - 1)
                      }
                    >
                      −
                    </button>
                    <span>{item.quantity}</span>
                    <button
                      onClick={() =>
                        onQuantityChange(item.coffee_id, item.quantity + 1)
                      }
                    >
                      +
                    </button>
                  </div>
                  <span className="subtotal">
                    ฿{item.price * item.quantity}
                  </span>
                  <button
                    className="remove-btn"
                    onClick={() => onRemoveItem(item.coffee_id)}
                  >
                    Remove
                  </button>
                </div>
              </div>
            ))}
          </div>

          <div className="order-summary">
            <div className="total">
              <strong>Total: </strong>
              <span className="total-price">฿{total}</span>
            </div>

            <button
              className="checkout-btn"
              onClick={onSubmit}
              disabled={loading || items.length === 0}
            >
              {loading ? 'Placing Order...' : 'Place Order'}
            </button>
          </div>
        </>
      )}
    </div>
  );
};
