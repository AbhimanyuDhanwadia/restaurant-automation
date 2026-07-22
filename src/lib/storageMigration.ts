const STORAGE_KEYS = {
  orders: "restaurant-automation:orders",
  alerts: "restaurant-automation:alerts",
  tables: "restaurant-automation:tables",
  inventory: "restaurant-automation:inventory",
  staff: "restaurant-automation:staff",
  tasks: "restaurant-automation:tasks",
  handoffNote: "restaurant-automation:handoff-note",
} as const;

function readStorage(key: string): unknown {
  try {
    const value = window.localStorage.getItem(key);
    return value ? JSON.parse(value) : undefined;
  } catch {
    return undefined;
  }
}

function isPersistedState(value: unknown): value is { state: Record<string, unknown> } {
  return typeof value === "object" && value !== null && !Array.isArray(value) && "state" in value;
}

function writePersistedState(key: string, state: Record<string, unknown>) {
  window.localStorage.setItem(key, JSON.stringify({ state, version: 0 }));
}

export function migrateLegacyStorage() {
  if (typeof window === "undefined") {
    return;
  }

  try {
    const collectionMigrations = [
      [STORAGE_KEYS.orders, "orders"],
      [STORAGE_KEYS.alerts, "alerts"],
      [STORAGE_KEYS.tables, "tables"],
      [STORAGE_KEYS.inventory, "items"],
    ] as const;

    for (const [key, stateKey] of collectionMigrations) {
      const stored = readStorage(key);
      if (Array.isArray(stored)) {
        writePersistedState(key, { [stateKey]: stored });
      }
    }

    const storedStaff = readStorage(STORAGE_KEYS.staff);
    const storedTasks = readStorage(STORAGE_KEYS.tasks);
    const storedHandoffNote = readStorage(STORAGE_KEYS.handoffNote);
    const staffState = isPersistedState(storedStaff) ? { ...storedStaff.state } : {};
    const staffNeedsMigration = Array.isArray(storedStaff) || Array.isArray(storedTasks) || typeof storedHandoffNote === "string";

    if (Array.isArray(storedStaff)) {
      staffState.staff = storedStaff;
    }
    if (Array.isArray(storedTasks)) {
      staffState.tasks = storedTasks;
    }
    if (typeof storedHandoffNote === "string") {
      staffState.handoffNote = storedHandoffNote;
    }

    if (staffNeedsMigration) {
      writePersistedState(STORAGE_KEYS.staff, staffState);
      window.localStorage.removeItem(STORAGE_KEYS.tasks);
      window.localStorage.removeItem(STORAGE_KEYS.handoffNote);
    }
  } catch {
    // Keep the seed state when storage is unavailable or malformed.
  }
}
