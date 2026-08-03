import { apiRequest } from "@/lib/api";

export interface Migration {
  version: number;
  applied_at: string;
}

export interface DatabaseStatus {
  status: "available" | "not_configured";
  migrations: Migration[];
}

export const getDatabaseStatus = () => apiRequest<DatabaseStatus>("/api/v1/admin/database");
