import { apiRequest } from "@/lib/api";

export type PrinterStatus = "ready" | "offline" | "error";

export interface PrinterHealth {
  name: string;
  status: PrinterStatus;
  queue_depth: number;
  printed: number;
  failed: number;
  busy_ms: number;
  utilization: number;
  utilization_available: boolean;
}

export const listPrinters = () => apiRequest<PrinterHealth[]>("/api/v1/printers");
