// src/components/MenuList.tsx
// Component แสดงเมนูกาแฟทั้งหมด
// ☕ Presentation Layer - ไม่มี business logic

import React, { useEffect, useState } from 'react';
import { fetchCoffees } from '../api';
import { Coffee } from '../types';
import './MenuList.css';

interface MenuListProps {
  onSelectCoffee: (coffee: Coffee, quantity: number) => void;
}

export const MenuList: React.FC<MenuListProps> = ({ onSelectCoffee }) => {
  const [coffees, setCoffees] = useState<Coffee[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [quantities, setQuantities] = useState<{ [key: string]: number }>({});

  useEffect(() => {
    loadCoffees();
  }, []);

  const loadCoffees = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await fetchCoffees();
      setCoffees(data);
      // Initialize quantities to 1
      const initialQuantities: { [key: string]: number } = {};
      data.forEach((coffee) => {
        initialQuantities[coffee.id] = 1;
      });
      setQuantities(initialQuantities);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load coffees');
    } finally {
      setLoading(false);
    }
  };

  const handleQuantityChange = (coffeeId: string, value: number) => {
    if (value >= 1) {
      setQuantities((prev) => ({
        ...prev,
        [coffeeId]: value,
      }));
    }
  };

  const handleAddToCart = (coffee: Coffee) => {
    const quantity = quantities[coffee.id] || 1;
    onSelectCoffee(coffee, quantity);
    // Reset quantity after adding
    handleQuantityChange(coffee.id, 1);
  };

  if (loading) {
    return <div className="menu-list loading">Loading menu...</div>;
  }

  if (error) {
    return <div className="menu-list error">Error: {error}</div>;
  }

  return (
    <div className="menu-list">
      <h2>☕ Coffee Menu</h2>
      <div className="coffee-grid">
        {coffees.map((coffee) => (
          <div key={coffee.id} className="coffee-card">
            <div className="coffee-emoji">{coffee.emoji}</div>
            <h3>{coffee.name}</h3>
            <p className="price">฿{coffee.price}</p>

            <div className="quantity-selector">
              <button
                onClick={() =>
                  handleQuantityChange(
                    coffee.id,
                    (quantities[coffee.id] || 1) - 1
                  )
                }
              >
                −
              </button>
              <input
                type="number"
                min="1"
                value={quantities[coffee.id] || 1}
                onChange={(e) =>
                  handleQuantityChange(coffee.id, Math.max(1, parseInt(e.target.value)))
                }
              />
              <button
                onClick={() =>
                  handleQuantityChange(
                    coffee.id,
                    (quantities[coffee.id] || 1) + 1
                  )
                }
              >
                +
              </button>
            </div>

            <button
              className="add-to-cart-btn"
              onClick={() => handleAddToCart(coffee)}
            >
              Add to Order
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};
