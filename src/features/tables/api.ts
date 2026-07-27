import { apiRequest } from "@/lib/api";

export type TableStatus = "available" | "seated" | "reserved" | "needs_check";

export interface RestaurantTable {
  id: string;
  seats: number;
  guest: string;
  reservation: string;
  status: TableStatus;
  created_at: string;
  updated_at: string;
}

export interface CreateTableInput {
  id: string;
  seats: number;
  guest?: string;
  reservation?: string;
  status?: TableStatus;
}

export function listTables() {
  return apiRequest<RestaurantTable[]>("/api/v1/tables");
}

export function createTable(input: CreateTableInput) {
  return apiRequest<RestaurantTable>("/api/v1/tables", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function updateTableStatus(id: string, status: TableStatus) {
  return apiRequest<RestaurantTable>(`/api/v1/tables/${encodeURIComponent(id)}/status`, {
    method: "PATCH",
    body: JSON.stringify({ status }),
  });
}
