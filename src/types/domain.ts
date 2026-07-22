/**
 * Domain types for the Restaurant Automation platform.
 *
 * These are pure TypeScript interfaces / union types with no runtime
 * dependencies. They represent the canonical shape of every entity used
 * across the application. Changing a type here propagates type errors to
 * every consumer immediately, making refactors safe.
 */

// ---------------------------------------------------------------------------
// Primitive scalars
// ---------------------------------------------------------------------------

export type OrderStatus = "Preparing" | "Ready" | "Delayed";
export type TableStatus = "Available" | "Seated" | "Reserved" | "Needs Check";
export type InventoryStatus = "In Stock" | "Low Stock" | "On Order";
export type StaffStatus = "On shift" | "On break" | "Off shift";
export type AlertSeverity = "High" | "Medium" | "Low";

// ---------------------------------------------------------------------------
// Order
// ---------------------------------------------------------------------------

export interface Order {
  id: string;
  table: string;
  channel: string;
  items: string;
  eta: string;
  status: OrderStatus;
}

// ---------------------------------------------------------------------------
// Table
// ---------------------------------------------------------------------------

export interface RestaurantTable {
  id: string;
  seats: number;
  guest: string;
  reservation: string;
  status: TableStatus;
}

// ---------------------------------------------------------------------------
// Inventory
// ---------------------------------------------------------------------------

export interface InventoryItem {
  id: string;
  item: string;
  category: string;
  onHand: string;
  par: string;
  supplier: string;
  status: InventoryStatus;
}

// ---------------------------------------------------------------------------
// Staff
// ---------------------------------------------------------------------------

export interface StaffMember {
  id: string;
  name: string;
  role: string;
  station: string;
  status: StaffStatus;
  handoff: string;
}

// ---------------------------------------------------------------------------
// Alert
// ---------------------------------------------------------------------------

export interface Alert {
  label: string;
  detail: string;
  severity: AlertSeverity;
}

// ---------------------------------------------------------------------------
// Task
// ---------------------------------------------------------------------------

export interface ShiftTask {
  title: string;
  owner: string;
  due: string;
}

// ---------------------------------------------------------------------------
// Stats (operations dashboard summary cards)
// ---------------------------------------------------------------------------

export interface StatCard {
  label: string;
  value: string;
  detail: string;
  icon: React.ComponentType<{ size?: number; "aria-hidden"?: "true" }>;
}
