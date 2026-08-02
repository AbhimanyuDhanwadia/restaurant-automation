import { apiRequest } from "@/lib/api";
import type { AutomationEvent } from "@/features/automation/api";

export interface AuditLogsResponse {
  source: "persistent" | "runtime";
  events: AutomationEvent[];
}

export const listAuditLogs = () => apiRequest<AuditLogsResponse>("/api/v1/admin/audit-logs");
