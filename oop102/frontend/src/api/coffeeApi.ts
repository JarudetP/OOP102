// src/api/coffeeApi.ts
// API Layer สำหรับ Coffee - ส่วนนี้รับผิดชอบการเรียก backend
// UI Components ไม่ควรเรียก fetch() โดยตรง แต่ให้เรียก function ที่นี่

import { Coffee, APIResponse } from '../types';
import { API_CONFIG, API_ENDPOINTS } from '../config/api';

/**
 * ดึงเมนูกาแฟทั้งหมด
 * UI Component ต้องเรียก function นี้เท่านั้น
 */
export async function fetchCoffees(): Promise<Coffee[]> {
  try {
    const url = `${API_CONFIG.BASE_URL}${API_ENDPOINTS.COFFEES}`;
    const response = await fetch(url);

    if (!response.ok) {
      throw new Error(`HTTP Error: ${response.status}`);
    }

    const data: APIResponse<Coffee[]> = await response.json();

    if (data.status !== 'success') {
      throw new Error(data.message);
    }

    return data.data;
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : 'Unknown error';
    console.error('Failed to fetch coffees:', errorMessage);
    throw error;
  }
}
