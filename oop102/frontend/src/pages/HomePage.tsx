// src/pages/HomePage.tsx
// 🏠 Main page - ประสานงาน UI components กับ API layer
// Logic orchestration layer

import React, { useState, useCallback } from 'react';
import { MenuList, OrderForm, OrderSummary } from '../components';
import { createOrder } from '../api';
import { Coffee, Order, OrderItem, CreateOrderRequest } from '../types';
import './HomePage.css';

export const HomePage: React.FC = () => {
  // State management: เก็บรายการสั่งและผลลัพธ์
  const [cartItems, setCartItems] = useState<
    Array<OrderItem & { coffeeName: string; coffeeEmoji: string; price: number }>
  >([]);
  const [completedOrder, setCompletedOrder] = useState<Order | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);

  // ==========================================
  // Event Handlers - ตัดสินใจธุรกิจระดับหน้า
  // ==========================================

  /**
   * เมื่อเลือกกาแฟจากเมนู
   * การตัดสินใจ: เพิ่มหรือ update quantity ถ้า coffee มี item แล้ว
   */
  const handleSelectCoffee = useCallback((coffee: Coffee, quantity: number) => {
    setCartItems((prev) => {
      // ตรวจสอบว่ากาแฟตัวนี้มี item แล้วหรือไม่
      const existingItem = prev.find((item) => item.coffee_id === coffee.id);

      if (existingItem) {
        // ถ้ามี: เพิ่ม quantity
        return prev.map((item) =>
          item.coffee_id === coffee.id
            ? { ...item, quantity: item.quantity + quantity }
            : item
        );
      } else {
        // ถ้าไม่มี: เพิ่ม item ใหม่
        return [
          ...prev,
          {
            coffee_id: coffee.id,
            quantity,
            coffeeName: coffee.name,
            coffeeEmoji: coffee.emoji,
            price: coffee.price,
          },
        ];
      }
    });

    // Clear error when adding items
    setSubmitError(null);
  }, []);

  /**
   * ลบรายการจากรถเข็น
   */
  const handleRemoveItem = useCallback((coffeeId: string) => {
    setCartItems((prev) => prev.filter((item) => item.coffee_id !== coffeeId));
  }, []);

  /**
   * เปลี่ยนจำนวนสินค้า
   */
  const handleQuantityChange = useCallback(
    (coffeeId: string, quantity: number) => {
      if (quantity <= 0) {
        handleRemoveItem(coffeeId);
      } else {
        setCartItems((prev) =>
          prev.map((item) =>
            item.coffee_id === coffeeId ? { ...item, quantity } : item
          )
        );
      }
    },
    [handleRemoveItem]
  );

  /**
   * ส่งการสั่งไปที่ backend
   * Main business logic: แปลง cart items → API request → จัดการ response
   */
  const handlePlaceOrder = useCallback(async () => {
    if (cartItems.length === 0) {
      setSubmitError('Please select at least one item');
      return;
    }

    try {
      setIsSubmitting(true);
      setSubmitError(null);

      // แปลง cart items → CreateOrderRequest format
      const orderRequest: CreateOrderRequest = {
        items: cartItems.map((item) => ({
          coffee_id: item.coffee_id,
          quantity: item.quantity,
        })),
      };

      // เรียก API layer
      const order = await createOrder(orderRequest);

      // สำเร็จ: บันทึก order และ clear cart
      setCompletedOrder(order);
      setCartItems([]);
    } catch (error) {
      // Error handling
      const errorMessage =
        error instanceof Error ? error.message : 'Failed to place order';
      setSubmitError(errorMessage);
    } finally {
      setIsSubmitting(false);
    }
  }, [cartItems]);

  /**
   * สั่งกาแฟใหม่ (หลังสำเร็จ)
   * Reset state สำหรับการสั่งครั้งต่อไป
   */
  const handleNewOrder = useCallback(() => {
    setCompletedOrder(null);
    setCartItems([]);
    setSubmitError(null);
  }, []);

  // ==========================================
  // Render
  // ==========================================

  return (
    <div className="home-page">
      <header className="header">
        <h1>☕ Coffee Shop</h1>
        <p>Place your order easily</p>
      </header>

      <main className="main-content">
        <div className="two-column">
          {/* Left Column: Menu */}
          <section className="menu-section">
            <MenuList onSelectCoffee={handleSelectCoffee} />
          </section>

          {/* Right Column: Order */}
          <section className="order-section">
            <OrderForm
              items={cartItems}
              onRemoveItem={handleRemoveItem}
              onQuantityChange={handleQuantityChange}
              onSubmit={handlePlaceOrder}
              loading={isSubmitting}
              error={submitError}
            />
          </section>
        </div>
      </main>

      {/* Modal: Order Confirmation */}
      <OrderSummary order={completedOrder} onNewOrder={handleNewOrder} />
    </div>
  );
};
