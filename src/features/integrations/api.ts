import { apiRequest } from "@/lib/api";

export type IntegrationStatus = "connected" | "disconnected" | "error";

export interface Integration {
  name: string;
  status: IntegrationStatus;
}

export const listIntegrations = () => apiRequest<Integration[]>("/api/v1/integrations");
