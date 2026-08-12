import { apiRequest } from "@/lib/api";

export type PrintJobStatus = "queued" | "printing" | "printed" | "failed" | "reviewed";

export interface PrintLine {
  text: string;
  quantity: number;
}

export interface PrintJob {
  id: string;
  order_id: string;
  destination: string;
  lines: PrintLine[];
  reprint: boolean;
  status: PrintJobStatus;
  attempts: number;
  last_error?: string;
  created_at: string;
  updated_at: string;
}

export const listPrintJobs = () => apiRequest<PrintJob[]>("/api/v1/printers/queue");
export const requeuePrintJob = (jobID: string) => apiRequest<PrintJob>(`/api/v1/printers/queue/${encodeURIComponent(jobID)}/requeue`, { method: "POST" });
export const retryPrintJob = (jobID: string) => apiRequest<PrintJob>(`/api/v1/printers/queue/${encodeURIComponent(jobID)}/retry`, { method: "POST" });
export const reprintOrder = (orderID: string) => apiRequest<PrintJob>(`/api/v1/printers/tickets/${encodeURIComponent(orderID)}/reprint`, { method: "POST" });
