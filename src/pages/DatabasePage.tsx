import { CheckCircle2, Database, FileClock, RefreshCw, Server } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { getDatabaseStatus } from "@/features/database/api";

export function DatabasePage() {
  const databaseQuery = useQuery({ queryKey: ["admin", "database"], queryFn: getDatabaseStatus, refetchInterval: 15_000 });
  const database = databaseQuery.data;
  const migrations = database?.migrations ?? [];
  const latest = migrations[0];
  const available = database?.status === "available";

  return <section className="database-workspace">
    <div className="automation-toolbar">
      <div><p className="eyebrow">Administration</p><h2>Database</h2><p className="automation-subtitle">Connection state and versioned schema migration history.</p></div>
      <button type="button" className="icon-button" onClick={() => databaseQuery.refetch()} disabled={databaseQuery.isFetching} title="Refresh database status" aria-label="Refresh database status"><RefreshCw size={18} aria-hidden="true" className={databaseQuery.isFetching ? "spin" : undefined} /></button>
    </div>
    {databaseQuery.isError && <div className="orders-api-error" role="alert"><span>{databaseQuery.error.message}</span><button type="button" className="text-button" onClick={() => databaseQuery.refetch()}>Retry</button></div>}
    <section className="automation-metric-grid database-metrics" aria-label="Database summary"><article className="automation-metric"><Database size={20} aria-hidden="true" /><span>Database status</span><strong>{database?.status.replace(/_/g, " ") ?? "Unavailable"}</strong><small>Read-only metadata access</small></article><article className="automation-metric"><FileClock size={20} aria-hidden="true" /><span>Applied migrations</span><strong>{database ? migrations.length : "Unavailable"}</strong><small>Schema migration ledger</small></article><article className="automation-metric"><CheckCircle2 size={20} aria-hidden="true" /><span>Latest version</span><strong>{latest ? latest.version : "-"}</strong><small>{latest ? "Most recent applied migration" : "No ledger available"}</small></article><article className="automation-metric"><Server size={20} aria-hidden="true" /><span>Persistence</span><strong>{available ? "Enabled" : database ? "Disabled" : "Unavailable"}</strong><small>Configured at API startup</small></article></section>
    <Panel label="Schema migrations"><PanelHeading eyebrow="Migration ledger" title="Applied schema versions" action={<FileClock size={20} aria-hidden="true" />} /><div className="migration-list">{databaseQuery.isPending && <EmptyState message="Loading database metadata..." />}{migrations.map((migration) => <article className="migration-row" key={migration.version}><strong>Version {migration.version}</strong><time>{new Date(migration.applied_at).toLocaleString([], { dateStyle: "medium", timeStyle: "short" })}</time></article>)}{!databaseQuery.isPending && !databaseQuery.isError && !available && <EmptyState message="PostgreSQL is not configured for this API instance." />}{!databaseQuery.isPending && !databaseQuery.isError && available && migrations.length === 0 && <EmptyState message="No applied schema migrations were returned." />}</div></Panel>
    <Panel label="Database access"><PanelHeading eyebrow="Access policy" title="Read-only metadata" action={<Database size={20} aria-hidden="true" />} /><p className="database-access-copy">This workspace only reports schema migration metadata. Database connection strings, credentials, table contents, and arbitrary query execution are never exposed through the application.</p></Panel>
  </section>;
}
