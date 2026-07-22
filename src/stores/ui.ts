/**
 * UI store.
 *
 * Stores ephemeral UI state that needs to be shared across layout components
 * and pages (e.g. the global search term entered in the Topbar). This is
 * intentionally not persisted — UI state resets on page load.
 */

import { create } from "zustand";

interface UIState {
  searchTerm: string;
  setSearchTerm: (term: string) => void;
}

export const useUIStore = create<UIState>()((set) => ({
  searchTerm: "",
  setSearchTerm: (term) => set({ searchTerm: term }),
}));
