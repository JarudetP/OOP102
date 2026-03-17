// src/api/orderApi.ts
// API Layer สำหรับ Order - จัดการการสั่งกาแฟ

import { CreateOrderRequest, Order, APIResponse } from '../types';
import { API_CONFIG, API_ENDPOINTS } from '../config/api';

/**
 * สั่งกาแฟใหม่
 * @param request - ข้อมูลการสั่ง (coffeeId + quantity)
 * @returns Order ที่สร้างขึ้น
 */
export async function createOrder(request: CreateOrderRequest): Promise<Order> {
  try {
    const url = `${API_CONFIG.BASE_URL}${API_ENDPOINTS.ORDERS}`;

    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      throw new Error(`HTTP Error: ${response.status}`);
    }

    const data: APIResponse<Order> = await response.json();

    if (data.status !== 'success') {
      throw new Error(data.message);
    }

    return data.data;
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : 'Unknown error';
    console.error('Failed to create order:', errorMessage);
    throw error;
  }
}
