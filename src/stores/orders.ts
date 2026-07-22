/**
 * Orders store.
 *
 * Manages the list of active orders and exposes mutation actions.
 * State is persisted to localStorage so it survives page reloads.
 * In Milestone 4, the initial state will be fetched from the Go API
 * and mutations will send PATCH requests; the store shape stays the same.
 */

import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { Order, OrderStatus } from "@/types/domain";
import { SEED_ORDERS } from "@/data/seeds";

interface OrdersState {
  orders: Order[];
  /** Cycle an order's status through the defined progression. */
  cycleStatus: (orderId: string) => void;
  /** Append a new order to the list. */
  addOrder: (order: Order) => void;
}

const STATUS_CYCLE: Record<OrderStatus, OrderStatus> = {
  Preparing: "Ready",
  Ready: "Preparing",
  Delayed: "Preparing",
};

export const useOrdersStore = create<OrdersState>()(
  persist(
    (set) => ({
      orders: SEED_ORDERS,

      cycleStatus: (orderId) =>
        set((state) => ({
          orders: state.orders.map((o) =>
            o.id === orderId
              ? { ...o, status: STATUS_CYCLE[o.status] }
              : o,
          ),
        })),

      addOrder: (order) =>
        set((state) => ({ orders: [...state.orders, order] })),
    }),
    { name: "restaurant-automation:orders" },
  ),
);
