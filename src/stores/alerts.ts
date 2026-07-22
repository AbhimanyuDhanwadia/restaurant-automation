/**
 * Alerts store.
 *
 * Manages operational alerts and their acknowledgement lifecycle.
 * State is persisted to localStorage so unacknowledged alerts survive reloads.
 */

import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { Alert } from "@/types/domain";
import { SEED_ALERTS } from "@/data/seeds";

interface AlertsState {
  alerts: Alert[];
  /** Remove an alert from the active list (acknowledge it). */
  acknowledge: (alertLabel: string) => void;
}

export const useAlertsStore = create<AlertsState>()(
  persist(
    (set) => ({
      alerts: SEED_ALERTS,

      acknowledge: (alertLabel) =>
        set((state) => ({
          alerts: state.alerts.filter((a) => a.label !== alertLabel),
        })),
    }),
    { name: "restaurant-automation:alerts" },
  ),
);
