import { Database, FileClock, Filter, RefreshCw, Search } from "lucide-react";
import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { listAuditLogs } from "@/features/audit/api";

function displayName(value: string) {
  return value.replace(/[._-]/g, " ").replace(/\b\w/g, (letter) => letter.toUpperCase());
}

export function AuditLogsPage() {
  const [eventType, setEventType] = useState("all");
  const [search, setSearch] = useState("");
  const logsQuery = useQuery({ queryKey: ["audit", "logs"], queryFn: listAuditLogs, refetchInterval: 10_000 });
  const events = useMemo(() => logsQuery.data?.events ?? [], [logsQuery.data?.events]);
  const eventTypes = useMemo(() => [...new Set(events.map((event) => event.type))].sort(), [events]);
  const visibleEvents = useMemo(() => events.filter((event) => {
    const term = search.trim().toLowerCase();
    return (eventType === "all" || event.type === eventType) && (!term || event.order_id.toLowerCase().includes(term) || event.type.toLowerCase().includes(term));
  }), [eventType, events, search]);
  const isPersistent = logsQuery.data?.source === "persistent";

  return <section className="audit-logs-workspace">
    <div className="automation-toolbar">
      <div><p className="eyebrow">Administration</p><h2>Audit Logs</h2><p className="automation-subtitle">Traceable operational events for review and incident investigation.</p></div>
      <button type="button" className="icon-button" onClick={() => logsQuery.refetch()} disabled={logsQuery.isFetching} title="Refresh audit logs" aria-label="Refresh audit logs"><RefreshCw size={18} aria-hidden="true" className={logsQuery.isFetching ? "spin" : undefined} /></button>
    </div>
    {logsQuery.isError && <div className="orders-api-error" role="alert"><span>{logsQuery.error.message}</span><button type="button" className="text-button" onClick={() => logsQuery.refetch()}>Retry</button></div>}
    <Panel label="Operational audit logs"><PanelHeading eyebrow="Event history" title="Operational records" action={logsQuery.data ? <span className={isPersistent ? "audit-source audit-persistent" : "audit-source"}>{isPersistent ? "PostgreSQL history" : "Runtime fallback"}</span> : undefined} /><div className="event-log-controls"><label className="filter-control"><span><Filter size={14} aria-hidden="true" /> Event type</span><select value={eventType} onChange={(event) => setEventType(event.target.value)}><option value="all">All events</option>{eventTypes.map((type) => <option key={type} value={type}>{displayName(type)}</option>)}</select></label><label className="event-search"><Search size={16} aria-hidden="true" /><input aria-label="Search audit logs" placeholder="Search order or event type" value={search} onChange={(event) => setSearch(event.target.value)} /></label></div><div className="audit-log-list">{logsQuery.isPending && <EmptyState message="Loading audit history..." />}{visibleEvents.map((event) => <article className="audit-log-row" key={event.id}><time>{new Date(event.created_at).toLocaleString([], { dateStyle: "medium", timeStyle: "medium" })}</time><div><strong>{displayName(event.type)}</strong><span>{event.order_id}</span></div><code>{event.id}</code></article>)}{!logsQuery.isPending && !logsQuery.isError && visibleEvents.length === 0 && <EmptyState message={events.length ? "No events match the current filters." : "No operational events are available."} />}</div></Panel>
    <Panel label="Audit retention"><PanelHeading eyebrow="Storage" title="Event source" action={<Database size={20} aria-hidden="true" />} /><p className="audit-retention-copy">{isPersistent ? "These events are read from the configured PostgreSQL operational-events store." : "PostgreSQL is not configured, so this view is showing the active engine stream and will reset when the API restarts."}</p><FileClock className="audit-retention-icon" size={18} aria-hidden="true" /></Panel>
  </section>;
}
