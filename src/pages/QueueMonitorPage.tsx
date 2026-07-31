import { Activity, CircleAlert, Gauge, RefreshCw, RotateCcw, UsersRound } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { getQueue } from "@/features/automation/api";

export function QueueMonitorPage() {
  const queueQuery = useQuery({ queryKey: ["automation", "queue"], queryFn: getQueue, refetchInterval: 2_000 });
  const queue = queueQuery.data;
  const percent = queue && queue.queue_capacity > 0 ? Math.min(100, (queue.queue_depth / queue.queue_capacity) * 100) : 0;
  const retryRate = queue && queue.events > 0 ? `${((queue.retried / queue.events) * 100).toFixed(1)}%` : "0.0%";
  const capacityState = percent >= 75 ? "queue-warning" : "queue-healthy";

  return <section className="queue-monitor-workspace">
    <div className="automation-toolbar">
      <div><p className="eyebrow">Automation</p><h2>Queue Monitor</h2><p className="automation-subtitle">Live order-processing capacity, worker activity, and retry health.</p></div>
      <button type="button" className="icon-button" onClick={() => queueQuery.refetch()} disabled={queueQuery.isFetching} title="Refresh queue metrics" aria-label="Refresh queue metrics"><RefreshCw size={18} aria-hidden="true" className={queueQuery.isFetching ? "spin" : undefined} /></button>
    </div>

    {queueQuery.isError && <div className="orders-api-error" role="alert"><span>{queueQuery.error.message}</span><button type="button" className="text-button" onClick={() => queueQuery.refetch()}>Retry</button></div>}

    <section className="automation-metric-grid queue-monitor-metrics" aria-label="Queue metrics">
      <article className="automation-metric"><Gauge size={20} aria-hidden="true" /><span>Queue depth</span><strong>{queue?.queue_depth ?? "Unavailable"}</strong><small>{queue ? `${queue.queue_capacity} total capacity` : "Queue unavailable"}</small></article>
      <article className="automation-metric"><UsersRound size={20} aria-hidden="true" /><span>Workers</span><strong>{queue?.workers ?? "Unavailable"}</strong><small>Configured at API startup</small></article>
      <article className="automation-metric"><RotateCcw size={20} aria-hidden="true" /><span>Retry rate</span><strong>{queue ? retryRate : "Unavailable"}</strong><small>{queue ? `${queue.retried} retries` : "Queue unavailable"}</small></article>
      <article className="automation-metric"><CircleAlert size={20} aria-hidden="true" /><span>Failed jobs</span><strong>{queue?.failed ?? "Unavailable"}</strong><small>Since API startup</small></article>
    </section>

    <div className="queue-monitor-grid">
      <Panel label="Queue capacity"><PanelHeading eyebrow="Capacity" title="Active workload" action={<Activity size={20} aria-hidden="true" />} /><div className="queue-monitor-gauge"><div><strong>{queue ? `${percent.toFixed(0)}%` : "-"}</strong><span>{queue ? `${queue.queue_depth} of ${queue.queue_capacity} waiting` : "Queue unavailable"}</span></div><div className="queue-monitor-track"><i className={capacityState} style={{ width: `${percent}%` }} /></div></div><p className="automation-footnote">The engine rejects submissions when the bounded queue reaches its configured capacity.</p></Panel>
      <Panel label="Queue counters"><PanelHeading eyebrow="Processing" title="Current process counters" action={<Activity size={20} aria-hidden="true" />} /><div className="queue-counter-list"><div><span>Automation events</span><strong>{queue?.events ?? "Unavailable"}</strong></div><div><span>Retry attempts</span><strong>{queue?.retried ?? "Unavailable"}</strong></div><div><span>Failed jobs</span><strong>{queue?.failed ?? "Unavailable"}</strong></div></div></Panel>
    </div>

    <Panel label="Queue behavior"><PanelHeading eyebrow="Runtime behavior" title="Worker lifecycle" action={<UsersRound size={20} aria-hidden="true" />} /><p className="queue-behavior-copy">Workers, retry limits, and queue capacity are established when the API starts. This monitor reports the active configuration and does not alter processing behavior from the browser.</p></Panel>
  </section>;
}
