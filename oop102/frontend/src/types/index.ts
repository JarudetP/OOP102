// src/types/index.ts
// Type definitions สำหรับ Coffee Shop API

export interface Coffee {
  id: string;
  name: string;
  price: number;
  emoji: string;
}

export interface OrderItem {
  coffee_id: string;
  quantity: number;
}

export interface CreateOrderRequest {
  items: OrderItem[];
}

export interface OrderItemResponse {
  coffee_id: string;
  coffee_name: string;
  quantity: number;
  price: number;
}

export interface Order {
  id: string;
  items: OrderItemResponse[];
  total: number;
  status: string;
  created_at: string;
}

export interface APIResponse<T> {
  status: string;
  message: string;
  data: T;
}
