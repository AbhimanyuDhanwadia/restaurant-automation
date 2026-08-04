import { apiRequest } from "@/lib/api";

export interface AppUser {
  auth_subject: string;
  email: string;
  role_id: string;
  role_name: string;
  created_at: string;
  last_seen_at: string;
  updated_at: string;
}

export const listUsers = () => apiRequest<AppUser[]>("/api/v1/admin/users");

export const updateUserRole = ({ subject, roleID }: { subject: string; roleID: string }) => apiRequest<AppUser>(`/api/v1/admin/users/${encodeURIComponent(subject)}/role`, {
  method: "PATCH",
  body: JSON.stringify({ role_id: roleID }),
});
