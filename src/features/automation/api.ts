import { apiRequest } from "@/lib/api";

export interface Provider {
  name: string;
  status: "connected" | "disconnected" | "error";
}

export interface Printer {
  name: string;
  status: "ready" | "offline" | "error";
  queue_depth: number;
  printed: number;
  failed: number;
}

export interface AutomationEvent {
  id: string;
  type: string;
  order_id: string;
  created_at: string;
}

export interface Queue {
  queue_depth: number;
  queue_capacity: number;
  workers: number;
  events: number;
  retried: number;
  failed: number;
}

export interface Insight {
  id: string;
  kind: string;
  severity: "info" | "warning" | "critical";
  resource: string;
  title: string;
  detail: string;
  updated_at: string;
}

export interface AutomationDashboard {
  providers: Provider[];
  printers: Printer[];
  events: AutomationEvent[];
  queue: Queue;
  insights: Insight[];
}

export interface SystemHealthComponent {
  name: "database" | "queue" | "integrations" | "printers";
  status: "healthy" | "unavailable" | "not_configured";
  detail: string;
}

export interface SystemHealth {
  status: "healthy" | "degraded";
  components: SystemHealthComponent[];
}

export const getQueue = () => apiRequest<Queue>("/api/v1/automation/queue");
export const getAutomationEvents = () => apiRequest<AutomationEvent[]>("/api/v1/automation/events");
export const getSystemHealth = () => apiRequest<SystemHealth>("/api/v1/system/health");

export async function getAutomationDashboard(): Promise<AutomationDashboard> {
  const [providers, printers, events, queue, insights] = await Promise.all([
    apiRequest<Provider[]>("/api/v1/integrations"),
    apiRequest<Printer[]>("/api/v1/printers"),
    apiRequest<AutomationEvent[]>("/api/v1/automation/events"),
    apiRequest<Queue>("/api/v1/automation/queue"),
    apiRequest<Insight[]>("/api/v1/insights"),
  ]);

  return { providers, printers, events, queue, insights };
}
