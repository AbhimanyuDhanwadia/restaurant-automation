import { apiRequest } from "@/lib/api";

export type BackupStatus = "recorded" | "verified" | "failed";

export interface BackupRecord {
  id: string;
  target: string;
  size_bytes: number;
  status: BackupStatus;
  completed_at: string;
  created_at: string;
}

export interface RecordBackupInput {
  target: string;
  size_bytes: number;
}

export const listBackups = () => apiRequest<BackupRecord[]>("/api/v1/admin/backups");

export const recordBackup = (input: RecordBackupInput) => apiRequest<BackupRecord>("/api/v1/admin/backups", {
  method: "POST",
  body: JSON.stringify(input),
});
