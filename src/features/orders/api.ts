import { apiRequest } from "@/lib/api";

export type OrderStatus = "received" | "preparing" | "ready" | "delivered" | "cancelled";

export interface OrderItem {
  name: string;
  quantity: number;
}

export interface Order {
  id: string;
  channel: string;
  status: OrderStatus;
  notes?: string;
  total_minor: number | null;
  currency?: string;
  items: OrderItem[];
  created_at: string;
  updated_at: string;
}

export interface CreateOrderInput {
  channel: string;
  notes?: string;
  total_minor?: number;
  currency?: string;
  items: OrderItem[];
}

export function listOrders() {
  return apiRequest<Order[]>("/api/v1/orders");
}

export function createOrder(input: CreateOrderInput) {
  return apiRequest<Order>("/api/v1/orders", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function updateOrderStatus(id: string, status: OrderStatus) {
  return apiRequest<Order>(`/api/v1/orders/${encodeURIComponent(id)}/status`, {
    method: "PATCH",
    body: JSON.stringify({ status }),
  });
}
