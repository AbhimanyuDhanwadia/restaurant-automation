import { apiRequest } from "@/lib/api";

export interface Metric {
  value: number;
  available: boolean;
}

export interface MonetaryMetric extends Metric {
  currency: string;
  included_orders: number;
  excluded_orders: number;
}

export interface DeliveryMetric extends Metric {
  delivered_orders: number;
  active_orders: number;
}

export interface StaffMetric extends Metric {
  completed_tasks: number;
  total_tasks: number;
}

export interface KitchenMetric extends Metric {
  completed_orders: number;
  eligible_orders: number;
}

export interface AnalyticsReport {
  generated_at: string;
  orders: Metric;
  order_volume_source: "durable" | "runtime";
  sales: MonetaryMetric;
  average_ticket: MonetaryMetric;
  kitchen_completion: KitchenMetric;
  kitchen_completion_source: "durable" | "runtime";
  delivery_time: DeliveryMetric;
  printer_availability: Metric;
  printer_utilization: Metric;
  staff_productivity: StaffMetric;
  peak_hours: Array<{ hour: number; orders: number }>;
}

export function getAnalyticsOverview() {
  return apiRequest<AnalyticsReport>("/api/v1/analytics/overview");
}
