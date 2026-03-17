// src/config/api.ts
// API Configuration - สำหรับเก็บ base URL และการตั้งค่าอื่น ๆ

export const API_CONFIG = {
  BASE_URL: process.env.REACT_APP_API_URL || 'http://localhost:8080',
  TIMEOUT: 10000,
};

export const API_ENDPOINTS = {
  COFFEES: '/coffees',
  ORDERS: '/orders',
};
