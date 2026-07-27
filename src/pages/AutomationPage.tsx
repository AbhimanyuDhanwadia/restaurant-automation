import { Activity, CheckCircle2, CircleAlert, Database, PlugZap, Printer, RefreshCw, Sparkles, Wifi } from "lucide-react";
import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { getAutomationDashboard } from "@/features/automation/api";

function statusClass(status: string) {
  return `automation-status status-${status}`;
}

function displayName(value: string) {
  return value.replace(/[._-]/g, " ").replace(/\b\w/g, (letter) => letter.toUpperCase());
}

export function AutomationPage() {
  const dashboardQuery = useQuery({
    queryKey: ["automation", "dashboard"],
    queryFn: getAutomationDashboard,
    refetchInterval: 5_000,
  });
  const dashboard = dashboardQuery.data;
  const providers = dashboard?.providers ?? [];
  const printers = dashboard?.printers ?? [];
  const events = useMemo(() => (dashboard?.events ?? []).slice(-8).reverse(), [dashboard?.events]);
  const insights = dashboard?.insights ?? [];
  const queue = dashboard?.queue;
  const readyPrinters = printers.filter((printer) => printer.status === "ready").length;
  const connectedProviders = providers.filter((provider) => provider.status === "connected").length;
  const capturedOrders = new Set((dashboard?.events ?? []).filter((event) => event.type === "order.received").map((event) => event.order_id)).size;
  const retryRate = queue && queue.events ? `${((queue.retried / queue.events) * 100).toFixed(1)}%` : "Unavailable";
  const queuePercent = queue && queue.queue_capacity ? Math.min(100, (queue.queue_depth / queue.queue_capacity) * 100) : 0;
  const lastUpdated = dashboardQuery.dataUpdatedAt ? new Date(dashboardQuery.dataUpdatedAt) : null;

  return (
    <section className="automation-workspace">
      <div className="automation-toolbar">
        <div><p className="eyebrow">Developer operations</p><h2>Automation Overview</h2><p className="automation-subtitle">Runtime control surface for orders, queues, integrations, and printers.</p></div>
        <button type="button" className="icon-button" onClick={() => dashboardQuery.refetch()} disabled={dashboardQuery.isFetching} title="Refresh automation status" aria-label="Refresh automation status"><RefreshCw size={18} aria-hidden="true" className={dashboardQuery.isFetching ? "spin" : undefined} /></button>
      </div>

      {dashboardQuery.isError && <div className="orders-api-error" role="alert"><span>{dashboardQuery.error.message}</span><button type="button" className="text-button" onClick={() => dashboardQuery.refetch()}>Retry</button></div>}

      <section className="automation-metric-grid" aria-label="Automation metrics">
        <article className="automation-metric"><Activity size={20} aria-hidden="true" /><span>Captured orders</span><strong>{dashboard ? capturedOrders : "Unavailable"}</strong><small>Current event stream</small></article>
        <article className="automation-metric"><Printer size={20} aria-hidden="true" /><span>Average print time</span><strong>Unavailable</strong><small>Telemetry required</small></article>
        <article className="automation-metric"><Wifi size={20} aria-hidden="true" /><span>Queue depth</span><strong>{queue?.queue_depth ?? "Unavailable"}</strong><small>{queue ? `${queue.queue_capacity} capacity` : "Queue unavailable"}</small></article>
        <article className="automation-metric"><CircleAlert size={20} aria-hidden="true" /><span>Retry rate</span><strong>{retryRate}</strong><small>{queue ? `${queue.failed} failed jobs` : "Queue unavailable"}</small></article>
      </section>

      <div className="automation-grid">
        <Panel label="System health"><PanelHeading eyebrow="System health" title="Runtime status" action={<CheckCircle2 size={20} aria-hidden="true" />} /><div className="automation-health-list"><div><span><Database size={16} aria-hidden="true" />API data</span><strong className={dashboard ? "health-good" : ""}>{dashboard ? "Available" : "Unavailable"}</strong></div><div><span><Activity size={16} aria-hidden="true" />Automation workers</span><strong className={queue ? "health-good" : ""}>{queue ? `${queue.workers} online` : "Unavailable"}</strong></div><div><span><Wifi size={16} aria-hidden="true" />Printer fleet</span><strong>{printers.length ? `${readyPrinters}/${printers.length} ready` : "Unavailable"}</strong></div><div><span><PlugZap size={16} aria-hidden="true" />Integrations</span><strong>{providers.length ? `${connectedProviders}/${providers.length} connected` : "Unavailable"}</strong></div></div></Panel>
        <Panel label="Queue metrics"><PanelHeading eyebrow="Queue" title="Order processing" action={<Activity size={20} aria-hidden="true" />} /><div className="queue-gauge"><strong>{queue?.queue_depth ?? "-"}</strong><span>{queue ? `waiting of ${queue.queue_capacity}` : "Queue unavailable"}</span><div><i style={{ width: `${queuePercent}%` }} /></div></div><div className="automation-mini-stats"><span>Events<strong>{queue?.events ?? "-"}</strong></span><span>Retries<strong>{queue?.retried ?? "-"}</strong></span><span>Failures<strong>{queue?.failed ?? "-"}</strong></span></div></Panel>
      </div>

      <div className="automation-grid">
        <Panel label="Integration status"><PanelHeading eyebrow="Integrations" title="Provider connections" action={<PlugZap size={20} aria-hidden="true" />} /><div className="automation-directory">{providers.map((provider) => <div className="automation-row" key={provider.name}><span>{displayName(provider.name)}</span><strong className={statusClass(provider.status)}>{provider.status}</strong></div>)}{!dashboardQuery.isPending && providers.length === 0 && <EmptyState message="No integration providers are registered." />}</div></Panel>
        <Panel label="Printer status"><PanelHeading eyebrow="Printers" title="Print infrastructure" action={<Printer size={20} aria-hidden="true" />} /><div className="automation-directory">{printers.map((printer) => <div className="automation-row" key={printer.name}><span>{displayName(printer.name)}</span><div><strong className={statusClass(printer.status)}>{printer.status}</strong><small>{printer.queue_depth} queued · {printer.failed} failures</small></div></div>)}{!dashboardQuery.isPending && printers.length === 0 && <EmptyState message="No printers are registered." />}</div><p className="automation-footnote">{printers.length ? `${readyPrinters} of ${printers.length} printers ready` : "Printer availability unavailable"}</p></Panel>
      </div>

      <Panel label="Operational intelligence"><PanelHeading eyebrow="Operational intelligence" title="Attention signals" action={<Sparkles size={20} aria-hidden="true" />} /><div className="insight-list">{insights.map((insight) => <article className={`insight-row insight-${insight.severity}`} key={insight.id}><CircleAlert size={18} aria-hidden="true" /><div><strong>{insight.title}</strong><p>{insight.detail}</p><small>{displayName(insight.kind)} · {displayName(insight.resource)}</small></div></article>)}{!dashboardQuery.isPending && insights.length === 0 && <EmptyState message="No automation signals require attention." />}</div></Panel>

      <Panel label="Event stream"><PanelHeading eyebrow="Event stream" title="Latest automation events" action={events.length ? <span className="automation-live-dot">Live</span> : undefined} /><div className="event-stream">{events.map((event) => <div className="event-row" key={event.id}><time>{new Date(event.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}</time><span className="event-marker" /><div><strong>{displayName(event.type)}</strong><small>{event.order_id}</small></div></div>)}{dashboardQuery.isPending && <EmptyState message="Loading automation events..." />}{!dashboardQuery.isPending && !dashboardQuery.isError && events.length === 0 && <EmptyState message="No automation events have been recorded." />}</div></Panel>
      {lastUpdated && <p className="automation-updated">Last checked {lastUpdated.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", second: "2-digit" })}</p>}
    </section>
  );
}
