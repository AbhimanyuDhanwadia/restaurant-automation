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

export interface AnalyticsReport {
  generated_at: string;
  orders: Metric;
  sales: MonetaryMetric;
  average_ticket: MonetaryMetric;
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
