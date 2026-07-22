/**
 * Staff store.
 *
 * Manages the shift roster, shift tasks, and manager handoff note.
 * State is persisted to localStorage.
 */

import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { StaffMember, ShiftTask } from "@/types/domain";
import { SEED_STAFF, SEED_TASKS } from "@/data/seeds";

interface StaffState {
  staff: StaffMember[];
  tasks: ShiftTask[];
  handoffNote: string;
  /** Mark a shift task as complete and remove it from the queue. */
  completeTask: (taskTitle: string) => void;
  /** Persist the manager's handoff note for the next shift. */
  saveHandoffNote: (note: string) => void;
}

export const useStaffStore = create<StaffState>()(
  persist(
    (set) => ({
      staff: SEED_STAFF,
      tasks: SEED_TASKS,
      handoffNote: "",

      completeTask: (taskTitle) =>
        set((state) => ({
          tasks: state.tasks.filter((t) => t.title !== taskTitle),
        })),

      saveHandoffNote: (note) => set({ handoffNote: note }),
    }),
    { name: "restaurant-automation:staff" },
  ),
);
