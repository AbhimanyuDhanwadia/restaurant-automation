import { apiRequest } from "@/lib/api";

export type AlertSeverity = "high" | "medium" | "low";
export interface Alert { id: string; label: string; detail: string; severity: AlertSeverity; created_at: string }
export interface CreateAlertInput { label: string; detail: string; severity: AlertSeverity }
export const listAlerts = () => apiRequest<Alert[]>("/api/v1/alerts");
export const createAlert = (input: CreateAlertInput) => apiRequest<Alert>("/api/v1/alerts", { method: "POST", body: JSON.stringify(input) });
export const acknowledgeAlert = (id: string) => apiRequest<void>(`/api/v1/alerts/${encodeURIComponent(id)}/acknowledge`, { method: "PATCH" });
