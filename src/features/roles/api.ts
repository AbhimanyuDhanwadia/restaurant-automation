import { apiRequest } from "@/lib/api";

export interface Role {
  id: string;
  name: string;
  description: string;
  permissions: string[];
  system: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateRoleInput {
  name: string;
  description: string;
  permissions: string[];
}

export const listRoles = () => apiRequest<Role[]>("/api/v1/admin/roles");

export const createRole = (input: CreateRoleInput) => apiRequest<Role>("/api/v1/admin/roles", {
  method: "POST",
  body: JSON.stringify(input),
});
