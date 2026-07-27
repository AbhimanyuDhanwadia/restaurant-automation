import { apiRequest } from "@/lib/api";

export interface Metric {
  value: number;
  available: boolean;
}

export interface AnalyticsReport {
  generated_at: string;
  orders: Metric;
  sales: Metric;
  average_ticket: Metric;
  kitchen_completion: Metric;
  delivery_time: Metric;
  printer_availability: Metric;
  printer_utilization: Metric;
  staff_productivity: Metric;
  peak_hours: Array<{ hour: number; orders: number }>;
}

export function getAnalyticsOverview() {
  return apiRequest<AnalyticsReport>("/api/v1/analytics/overview");
}
