// src/components/OrderSummary.tsx
// Component แสดงสรุปการสั่งหลังสำเร็จ
// ✅ Success confirmation view

import React from 'react';
import { Order } from '../types';
import './OrderSummary.css';

interface OrderSummaryProps {
  order: Order | null;
  onNewOrder: () => void;
}

export const OrderSummary: React.FC<OrderSummaryProps> = ({ order, onNewOrder }) => {
  if (!order) {
    return null;
  }

  return (
    <div className="order-summary-modal">
      <div className="modal-content">
        <div className="success-icon">✅</div>

        <h2>Order Confirmed!</h2>

        <div className="order-details">
          <div className="order-id">
            <strong>Order ID:</strong>
            <span>{order.id}</span>
          </div>

          <div className="order-items">
            <h3>Items:</h3>
            {order.items.map((item, index) => (
              <div key={index} className="item">
                <span className="item-name">{item.coffee_name}</span>
                <span className="item-qty">× {item.quantity}</span>
                <span className="item-price">฿{item.price}</span>
              </div>
            ))}
          </div>

          <div className="order-total">
            <strong>Total:</strong>
            <span className="total-amount">฿{order.total}</span>
          </div>

          <div className="order-status">
            <strong>Status:</strong>
            <span className="status">{order.status}</span>
          </div>

          <div className="order-time">
            <strong>Time:</strong>
            <span>{order.created_at}</span>
          </div>
        </div>

        <button className="new-order-btn" onClick={onNewOrder}>
          Place Another Order
        </button>
      </div>
    </div>
  );
};
