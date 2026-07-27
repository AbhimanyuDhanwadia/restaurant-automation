import { supabase } from "@/lib/supabase";

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
  items: OrderItem[];
  created_at: string;
  updated_at: string;
}

export interface CreateOrderInput {
  channel: string;
  notes?: string;
  items: OrderItem[];
}

const apiBaseUrl = (import.meta.env.VITE_API_URL ?? "http://localhost:8080").replace(/\/$/, "");

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const session = supabase ? (await supabase.auth.getSession()).data.session : null;
  const headers = new Headers(init.headers);

  headers.set("Accept", "application/json");
  if (init.body) headers.set("Content-Type", "application/json");
  if (session) headers.set("Authorization", `Bearer ${session.access_token}`);

  let response: Response;
  try {
    response = await fetch(`${apiBaseUrl}${path}`, { ...init, headers });
  } catch {
    throw new Error("The orders service is unreachable. Check that the API is running.");
  }

  if (!response.ok) {
    const detail = await response.text();
    throw new Error(detail || `Orders request failed (${response.status}).`);
  }

  return response.json() as Promise<T>;
}

export function listOrders() {
  return request<Order[]>("/api/v1/orders");
}

export function createOrder(input: CreateOrderInput) {
  return request<Order>("/api/v1/orders", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function updateOrderStatus(id: string, status: OrderStatus) {
  return request<Order>(`/api/v1/orders/${encodeURIComponent(id)}/status`, {
    method: "PATCH",
    body: JSON.stringify({ status }),
  });
}
