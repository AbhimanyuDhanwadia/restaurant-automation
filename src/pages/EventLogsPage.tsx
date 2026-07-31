import { Filter, RefreshCw, ScrollText, Search } from "lucide-react";
import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { getAutomationEvents } from "@/features/automation/api";

function displayName(value: string) {
  return value.replace(/[._-]/g, " ").replace(/\b\w/g, (letter) => letter.toUpperCase());
}

export function EventLogsPage() {
  const [eventType, setEventType] = useState("all");
  const [search, setSearch] = useState("");
  const eventsQuery = useQuery({ queryKey: ["automation", "events"], queryFn: getAutomationEvents, refetchInterval: 2_000 });
  const events = useMemo(() => eventsQuery.data ?? [], [eventsQuery.data]);
  const eventTypes = useMemo(() => [...new Set(events.map((event) => event.type))].sort(), [events]);
  const visibleEvents = useMemo(() => {
    const term = search.trim().toLowerCase();
    return [...events].reverse().filter((event) => (eventType === "all" || event.type === eventType) && (!term || event.order_id.toLowerCase().includes(term) || event.type.toLowerCase().includes(term)));
  }, [eventType, events, search]);

  return <section className="event-logs-workspace">
    <div className="automation-toolbar">
      <div><p className="eyebrow">Automation</p><h2>Event Logs</h2><p className="automation-subtitle">Live event stream from the active order-processing engine.</p></div>
      <button type="button" className="icon-button" onClick={() => eventsQuery.refetch()} disabled={eventsQuery.isFetching} title="Refresh event stream" aria-label="Refresh event stream"><RefreshCw size={18} aria-hidden="true" className={eventsQuery.isFetching ? "spin" : undefined} /></button>
    </div>

    {eventsQuery.isError && <div className="orders-api-error" role="alert"><span>{eventsQuery.error.message}</span><button type="button" className="text-button" onClick={() => eventsQuery.refetch()}>Retry</button></div>}

    <Panel label="Automation event stream"><PanelHeading eyebrow="Current process" title="Live events" action={<span className="task-count">{eventsQuery.data ? `${events.length} captured` : "Unavailable"}</span>} /><div className="event-log-controls"><label className="filter-control"><span><Filter size={14} aria-hidden="true" /> Event type</span><select value={eventType} onChange={(event) => setEventType(event.target.value)}><option value="all">All events</option>{eventTypes.map((type) => <option key={type} value={type}>{displayName(type)}</option>)}</select></label><label className="event-search"><Search size={16} aria-hidden="true" /><input aria-label="Search event logs" placeholder="Search order or event type" value={search} onChange={(event) => setSearch(event.target.value)} /></label></div><div className="event-log-list">{eventsQuery.isPending && <EmptyState message="Loading automation events..." />}{visibleEvents.map((event) => <article className="event-log-row" key={event.id}><time>{new Date(event.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", second: "2-digit" })}</time><div><strong>{displayName(event.type)}</strong><span>{event.order_id}</span></div><code>{event.id.slice(0, 8)}</code></article>)}{!eventsQuery.isPending && !eventsQuery.isError && visibleEvents.length === 0 && <EmptyState message={events.length ? "No events match the current filters." : "No automation events have been recorded."} />}</div></Panel>

    <Panel label="Event retention"><PanelHeading eyebrow="Retention" title="Runtime stream" action={<ScrollText size={20} aria-hidden="true" />} /><p className="event-retention-copy">This view shows the active engine’s in-memory stream. PostgreSQL deployments also persist operational events independently, but browsing durable audit history will be introduced with the Administration Audit Logs milestone.</p></Panel>
  </section>;
}
