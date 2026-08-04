import { apiRequest } from "@/lib/api";

export interface CurrentUserAccess {
  subject: string;
  email: string;
  auth_role: string;
  role_id: string;
  role_name: string;
  permissions: string[];
}

export const getCurrentUserAccess = () => apiRequest<CurrentUserAccess>("/api/v1/me");
