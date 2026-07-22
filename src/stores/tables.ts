/**
 * Tables store.
 *
 * Manages restaurant table status and allows status cycling through the
 * seating lifecycle. State is persisted to localStorage.
 */

import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { RestaurantTable, TableStatus } from "@/types/domain";
import { SEED_TABLES } from "@/data/seeds";

interface TablesState {
  tables: RestaurantTable[];
  /** Advance a table through its status lifecycle. */
  cycleStatus: (tableId: string) => void;
  /** Add a new table to the floor plan. */
  addTable: (table: RestaurantTable) => void;
}

const STATUS_CYCLE: Record<TableStatus, TableStatus> = {
  Available: "Reserved",
  Reserved: "Seated",
  Seated: "Needs Check",
  "Needs Check": "Available",
};

export const useTablesStore = create<TablesState>()(
  persist(
    (set) => ({
      tables: SEED_TABLES,

      cycleStatus: (tableId) =>
        set((state) => ({
          tables: state.tables.map((t) =>
            t.id === tableId
              ? { ...t, status: STATUS_CYCLE[t.status] }
              : t,
          ),
        })),

      addTable: (table) =>
        set((state) => ({ tables: [...state.tables, table] })),
    }),
    { name: "restaurant-automation:tables" },
  ),
);
