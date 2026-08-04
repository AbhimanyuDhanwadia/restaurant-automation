import { Archive, CircleAlert, HardDriveDownload, RefreshCw } from "lucide-react";
import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { listBackups, recordBackup, type BackupStatus } from "@/features/backups/api";

function formatSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 ** 3) return `${(bytes / (1024 ** 2)).toFixed(1)} MB`;
  return `${(bytes / (1024 ** 3)).toFixed(1)} GB`;
}

function statusLabel(status: BackupStatus) {
  return `${status.slice(0, 1).toUpperCase()}${status.slice(1)}`;
}

export function BackupsPage() {
  const [target, setTarget] = useState("");
  const [sizeMegabytes, setSizeMegabytes] = useState("");
  const queryClient = useQueryClient();
  const backupsQuery = useQuery({ queryKey: ["admin", "backups"], queryFn: listBackups, refetchInterval: 15_000 });
  const recordMutation = useMutation({
    mutationFn: recordBackup,
    onSuccess: () => {
      setTarget("");
      setSizeMegabytes("");
      queryClient.invalidateQueries({ queryKey: ["admin", "backups"] });
    },
  });
  const records = useMemo(() => backupsQuery.data ?? [], [backupsQuery.data]);
  const latest = records[0];
  const verified = records.filter((record) => record.status === "verified").length;
  const failed = records.filter((record) => record.status === "failed").length;

  return <section className="backups-workspace">
    <div className="automation-toolbar"><div><p className="eyebrow">Administration</p><h2>Backups</h2><p className="automation-subtitle">Record external backup artifacts without exposing database credentials or dump controls.</p></div><button type="button" className="icon-button" onClick={() => backupsQuery.refetch()} disabled={backupsQuery.isFetching} title="Refresh backup records" aria-label="Refresh backup records"><RefreshCw size={18} aria-hidden="true" className={backupsQuery.isFetching ? "spin" : undefined} /></button></div>
    {backupsQuery.isError && <div className="orders-api-error" role="alert"><span>{backupsQuery.error.message}</span><button type="button" className="text-button" onClick={() => backupsQuery.refetch()}>Retry</button></div>}
    {recordMutation.isError && <p className="orders-api-error" role="alert">{recordMutation.error.message}</p>}

    <section className="automation-metric-grid" aria-label="Backup metrics"><article className="automation-metric"><Archive size={20} aria-hidden="true" /><span>Backup records</span><strong>{backupsQuery.data ? records.length : "Unavailable"}</strong><small>External artifacts registered</small></article><article className="automation-metric"><HardDriveDownload size={20} aria-hidden="true" /><span>Latest backup</span><strong>{latest ? new Date(latest.completed_at).toLocaleDateString() : "None"}</strong><small>{latest ? formatSize(latest.size_bytes) : "No records available"}</small></article><article className="automation-metric"><Archive size={20} aria-hidden="true" /><span>Verified</span><strong>{backupsQuery.data ? verified : "Unavailable"}</strong><small>Set by trusted automation</small></article><article className="automation-metric"><CircleAlert size={20} aria-hidden="true" /><span>Failures</span><strong>{backupsQuery.data ? failed : "Unavailable"}</strong><small>Reported by trusted automation</small></article></section>

    <div className="backups-grid"><Panel label="Backup records"><PanelHeading eyebrow="Artifact catalog" title="Recorded backups" action={<span className="task-count">{backupsQuery.data ? `${records.length} records` : "Unavailable"}</span>} /><div className="backup-record-list">{backupsQuery.isPending && <EmptyState message="Loading backup records..." />}{records.map((record) => <article className="backup-record-row" key={record.id}><div><strong>{record.target}</strong><span>{formatSize(record.size_bytes)} · {new Date(record.completed_at).toLocaleString()}</span></div><span className={`backup-status-${record.status}`}>{statusLabel(record.status)}</span></article>)}{!backupsQuery.isPending && !backupsQuery.isError && records.length === 0 && <EmptyState message="No backup artifacts have been recorded." />}</div></Panel>
    <Panel label="Register backup artifact"><PanelHeading eyebrow="External backup" title="Record artifact" action={<HardDriveDownload size={20} aria-hidden="true" />} /><form className="backup-record-form" onSubmit={(event) => { event.preventDefault(); recordMutation.mutate({ target: target.trim(), size_bytes: Math.round(Number(sizeMegabytes) * 1024 * 1024) }); }}><label><span>Artifact location</span><input value={target} onChange={(event) => setTarget(event.target.value)} placeholder="s3://bucket/backup.sql.gz" required /></label><label><span>Size (MB)</span><input type="number" min="0" step="0.01" inputMode="decimal" value={sizeMegabytes} onChange={(event) => setSizeMegabytes(event.target.value)} required /></label><button type="submit" className="primary-button" disabled={recordMutation.isPending}>{recordMutation.isPending ? "Recording..." : "Record backup"}</button></form><p className="automation-footnote">New records are marked recorded. Verification and failure states must come from a trusted backup worker, not this workspace.</p></Panel></div>

    <Panel label="Backup boundary"><PanelHeading eyebrow="Safety boundary" title="Infrastructure-owned execution" action={<Archive size={20} aria-hidden="true" />} /><p className="analytics-summary-copy">This application does not run `pg_dump`, expose database credentials, or download backup files through the browser. A scheduled infrastructure job should create artifacts and later report their verification result through a privileged integration.</p></Panel>
  </section>;
}
