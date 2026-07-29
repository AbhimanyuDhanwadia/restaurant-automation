import { apiRequest } from "@/lib/api";

export type InventoryStatus = "in_stock" | "low_stock" | "on_order";

export interface InventoryItem {
  id: string;
  name: string;
  category: string;
  on_hand: number;
  par_level: number;
  unit: string;
  supplier: string;
  status: InventoryStatus;
  created_at: string;
  updated_at: string;
}

export interface CreateInventoryInput {
  name: string;
  category: string;
  on_hand: number;
  par_level: number;
  unit: string;
  supplier?: string;
  status?: InventoryStatus;
}

export function listInventory() {
  return apiRequest<InventoryItem[]>("/api/v1/inventory");
}

export function createInventoryItem(input: CreateInventoryInput) {
  return apiRequest<InventoryItem>("/api/v1/inventory", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function updateInventoryStock(id: string, onHand: number) {
  return apiRequest<InventoryItem>(`/api/v1/inventory/${encodeURIComponent(id)}/stock`, {
    method: "PATCH",
    body: JSON.stringify({ on_hand: onHand }),
  });
}

export function updateInventoryStatus(id: string, status: InventoryStatus) {
  return apiRequest<InventoryItem>(`/api/v1/inventory/${encodeURIComponent(id)}/status`, {
    method: "PATCH",
    body: JSON.stringify({ status }),
  });
}
