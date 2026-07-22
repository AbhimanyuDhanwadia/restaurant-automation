/**
 * Inventory store.
 *
 * Manages inventory items and exposes actions for stock updates.
 * State is persisted to localStorage.
 */

import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { InventoryItem } from "@/types/domain";
import { SEED_INVENTORY } from "@/data/seeds";

interface InventoryState {
  items: InventoryItem[];
  /** Mark an item as received (sets status to In Stock). */
  receiveItem: (itemId: string) => void;
  /** Add a new inventory item. */
  addItem: (item: InventoryItem) => void;
}

export const useInventoryStore = create<InventoryState>()(
  persist(
    (set) => ({
      items: SEED_INVENTORY,

      receiveItem: (itemId) =>
        set((state) => ({
          items: state.items.map((i) =>
            i.id === itemId ? { ...i, status: "In Stock" as const } : i,
          ),
        })),

      addItem: (item) =>
        set((state) => ({ items: [...state.items, item] })),
    }),
    { name: "restaurant-automation:inventory" },
  ),
);
