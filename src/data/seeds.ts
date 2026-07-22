/**
 * Seed / mock data used to populate stores before the Go API exists.
 *
 * In Milestone 4, these will be replaced by real API responses fetched via
 * TanStack Query. Keeping seeds here (not inside stores) keeps stores free
 * of data-fetching concerns.
 */

import type {
  Order,
  RestaurantTable,
  InventoryItem,
  StaffMember,
  Alert,
  ShiftTask,
} from "@/types/domain";
import { Flame, ReceiptText, TableProperties, UsersRound } from "lucide-react";

export const SEED_STATS = [
  { label: "Open orders", value: "28", detail: "+6 in 15 min", icon: ReceiptText },
  { label: "Kitchen load", value: "82%", detail: "Grill at capacity", icon: Flame },
  { label: "Tables seated", value: "19/24", detail: "5 turning soon", icon: TableProperties },
  { label: "Staff online", value: "14", detail: "2 tasks overdue", icon: UsersRound },
];

export const SEED_ORDERS: Order[] = [
  { id: "ORD-1842", table: "Table 12", channel: "Dine-in", items: "2 mains, 1 starter", eta: "7 min", status: "Preparing" },
  { id: "ORD-1843", table: "Delivery", channel: "Aggregator", items: "4 entrees", eta: "Ready now", status: "Ready" },
  { id: "ORD-1844", table: "Table 3", channel: "Dine-in", items: "1 tasting menu", eta: "18 min", status: "Delayed" },
  { id: "ORD-1845", table: "Pickup", channel: "Web", items: "3 bowls, 2 drinks", eta: "11 min", status: "Preparing" },
];

export const SEED_ALERTS: Alert[] = [
  { label: "Romaine lettuce", detail: "Below par by 4 cases", severity: "High" },
  { label: "Dish station", detail: "Needs support before dinner rush", severity: "Medium" },
  { label: "Table 8", detail: "Guest has waited 9 min for check", severity: "Low" },
];

export const SEED_TASKS: ShiftTask[] = [
  { title: "Approve prep list", owner: "Sous chef", due: "4:30 PM" },
  { title: "Confirm delivery partner SLA", owner: "Manager", due: "5:00 PM" },
  { title: "Restock bar garnishes", owner: "Bar lead", due: "5:15 PM" },
];

export const SEED_TABLES: RestaurantTable[] = [
  { id: "T1", seats: 2, guest: "Available", reservation: "Walk-in", status: "Available" },
  { id: "T2", seats: 4, guest: "The Mehta party", reservation: "6:00 PM", status: "Seated" },
  { id: "T3", seats: 2, guest: "A. Kapoor", reservation: "6:15 PM", status: "Needs Check" },
  { id: "T4", seats: 6, guest: "The Shah party", reservation: "6:30 PM", status: "Reserved" },
  { id: "T5", seats: 4, guest: "Available", reservation: "Walk-in", status: "Available" },
  { id: "T6", seats: 2, guest: "R. Iyer", reservation: "5:45 PM", status: "Seated" },
  { id: "T7", seats: 8, guest: "Available", reservation: "Walk-in", status: "Available" },
  { id: "T8", seats: 4, guest: "N. Rao", reservation: "5:50 PM", status: "Needs Check" },
];

export const SEED_INVENTORY: InventoryItem[] = [
  { id: "INV-01", item: "Romaine lettuce", category: "Produce", onHand: "2 cases", par: "6 cases", supplier: "Greenline Farms", status: "Low Stock" },
  { id: "INV-02", item: "Chicken breast", category: "Protein", onHand: "18 kg", par: "24 kg", supplier: "Metro Provisions", status: "Low Stock" },
  { id: "INV-03", item: "Sparkling water", category: "Beverage", onHand: "9 cases", par: "8 cases", supplier: "Beverage Co.", status: "In Stock" },
  { id: "INV-04", item: "Wild-caught salmon", category: "Protein", onHand: "12 kg", par: "18 kg", supplier: "Ocean Table", status: "On Order" },
  { id: "INV-05", item: "Sourdough loaves", category: "Bakery", onHand: "14 loaves", par: "12 loaves", supplier: "Daily Crumb", status: "In Stock" },
];

export const SEED_STAFF: StaffMember[] = [
  { id: "STAFF-01", name: "Maya Patel", role: "Manager", station: "Front of house", status: "On shift", handoff: "Confirm the delivery partner SLA before 5:00 PM." },
  { id: "STAFF-02", name: "Arjun Shah", role: "Sous chef", station: "Hot line", status: "On shift", handoff: "Approve the prep list and watch grill capacity." },
  { id: "STAFF-03", name: "Nisha Rao", role: "Bar lead", station: "Bar", status: "On break", handoff: "Restock bar garnishes before the dinner rush." },
  { id: "STAFF-04", name: "Kabir Mehta", role: "Server", station: "Section B", status: "On shift", handoff: "Check in on Table 8 and close the open check." },
];
