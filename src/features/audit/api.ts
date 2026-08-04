import { apiRequest } from "@/lib/api";
import type { AutomationEvent } from "@/features/automation/api";

export interface AuditLogsResponse {
  source: "persistent" | "runtime";
  events: AutomationEvent[];
}

export interface RoleAssignmentEvent {
  id: string;
  actor_subject: string;
  actor_email: string;
  target_subject: string;
  previous_role_id: string;
  previous_role_name: string;
  new_role_id: string;
  new_role_name: string;
  created_at: string;
}

export const listAuditLogs = () => apiRequest<AuditLogsResponse>("/api/v1/admin/audit-logs");
export const listRoleAssignmentEvents = () => apiRequest<RoleAssignmentEvent[]>("/api/v1/admin/access-audit");
