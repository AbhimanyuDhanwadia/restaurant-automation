import { Activity, CheckCircle2, CircleAlert, ListChecks, RefreshCw, UsersRound } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { useMemo } from "react";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { getAutomationEvents, getQueue, type AutomationEvent } from "@/features/automation/api";

const PIPELINE_STAGES = [
  { type: "order.received", label: "Received" },
  { type: "order.normalized", label: "Normalized" },
  { type: "order.validated", label: "Validated" },
  { type: "order.stored", label: "Stored" },
  { type: "order.queued", label: "Queued" },
] as const;

function displayName(value: string) {
  return value.replace(/[._-]/g, " ").replace(/\b\w/g, (letter) => letter.toUpperCase());
}

interface ObservedOrder {
  orderID: string;
  completedStages: Set<string>;
  latestEvent: AutomationEvent;
}

export function OrderEnginePage() {
  const engineQuery = useQuery({
    queryKey: ["automation", "order-engine"],
    queryFn: async () => {
      const [queue, events] = await Promise.all([getQueue(), getAutomationEvents()]);
      return { queue, events };
    },
    refetchInterval: 2_000,
  });
  const queue = engineQuery.data?.queue;
  const events = useMemo(() => engineQuery.data?.events ?? [], [engineQuery.data?.events]);
  const stageCounts = useMemo(() => new Map<string, number>(PIPELINE_STAGES.map(({ type }) => [type, 0])), []);
  const orders = useMemo(() => {
    const grouped = new Map<string, ObservedOrder>();

    for (const event of events) {
      if (!event.order_id) continue;
      const current = grouped.get(event.order_id);
      if (current) {
        current.completedStages.add(event.type);
        if (new Date(event.created_at) > new Date(current.latestEvent.created_at)) current.latestEvent = event;
      } else {
        grouped.set(event.order_id, { orderID: event.order_id, completedStages: new Set([event.type]), latestEvent: event });
      }
    }

    return [...grouped.values()].sort((left, right) => new Date(right.latestEvent.created_at).getTime() - new Date(left.latestEvent.created_at).getTime());
  }, [events]);
  const counts = useMemo(() => {
    const nextCounts = new Map(stageCounts);
    for (const event of events) {
      if (nextCounts.has(event.type)) nextCounts.set(event.type, (nextCounts.get(event.type) ?? 0) + 1);
    }
    return nextCounts;
  }, [events, stageCounts]);
  const completedOrders = orders.filter((order) => PIPELINE_STAGES.every(({ type }) => order.completedStages.has(type))).length;
  const retryRate = queue && queue.events > 0 ? `${((queue.retried / queue.events) * 100).toFixed(1)}%` : "0.0%";

  return <section className="order-engine-workspace">
    <div className="automation-toolbar">
      <div><p className="eyebrow">Automation</p><h2>Order Engine</h2><p className="automation-subtitle">Live pipeline progress from received orders through queue submission.</p></div>
      <button type="button" className="icon-button" onClick={() => engineQuery.refetch()} disabled={engineQuery.isFetching} title="Refresh order engine" aria-label="Refresh order engine"><RefreshCw size={18} aria-hidden="true" className={engineQuery.isFetching ? "spin" : undefined} /></button>
    </div>

    {engineQuery.isError && <div className="orders-api-error" role="alert"><span>{engineQuery.error.message}</span><button type="button" className="text-button" onClick={() => engineQuery.refetch()}>Retry</button></div>}

    <section className="automation-metric-grid" aria-label="Order engine metrics">
      <article className="automation-metric"><ListChecks size={20} aria-hidden="true" /><span>Orders observed</span><strong>{engineQuery.data ? orders.length : "Unavailable"}</strong><small>In the active event stream</small></article>
      <article className="automation-metric"><CheckCircle2 size={20} aria-hidden="true" /><span>Completed pipeline</span><strong>{engineQuery.data ? completedOrders : "Unavailable"}</strong><small>Reached queue submission</small></article>
      <article className="automation-metric"><Activity size={20} aria-hidden="true" /><span>Queue depth</span><strong>{queue?.queue_depth ?? "Unavailable"}</strong><small>{queue ? `${queue.queue_capacity} total capacity` : "Queue unavailable"}</small></article>
      <article className="automation-metric"><CircleAlert size={20} aria-hidden="true" /><span>Failed jobs</span><strong>{queue?.failed ?? "Unavailable"}</strong><small>Since API startup</small></article>
    </section>

    <Panel label="Order pipeline"><PanelHeading eyebrow="Pipeline" title="Current event stages" action={<span className="task-count">{events.length} events</span>} /><div className="engine-pipeline">{PIPELINE_STAGES.map(({ type, label }, index) => <div className="engine-stage" key={type}><span>{index + 1}</span><div><strong>{label}</strong><small>{counts.get(type) ?? 0} recorded</small></div></div>)}</div><p className="automation-footnote">These are the stages currently emitted by the order engine. Printing and kitchen handling are managed by their own automation services.</p></Panel>

    <div className="order-engine-grid">
      <Panel label="Recent order progress"><PanelHeading eyebrow="Observed orders" title="Pipeline progress" action={<ListChecks size={20} aria-hidden="true" />} /><div className="engine-order-list">{engineQuery.isPending && <EmptyState message="Loading order-engine events..." />}{orders.map((order) => {
        const completedStages = PIPELINE_STAGES.filter(({ type }) => order.completedStages.has(type)).length;
        const percent = (completedStages / PIPELINE_STAGES.length) * 100;
        return <article className="engine-order-row" key={order.orderID}><div><strong>{order.orderID}</strong><span>{displayName(order.latestEvent.type)}</span></div><div className="engine-order-progress"><span>{completedStages} of {PIPELINE_STAGES.length} stages</span><div role="progressbar" aria-label={`Pipeline progress for ${order.orderID}`} aria-valuemin={0} aria-valuemax={PIPELINE_STAGES.length} aria-valuenow={completedStages}><i style={{ width: `${percent}%` }} /></div></div><time>{new Date(order.latestEvent.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", second: "2-digit" })}</time></article>;
      })}{!engineQuery.isPending && !engineQuery.isError && orders.length === 0 && <EmptyState message="No order-engine events have been recorded." />}</div></Panel>

      <Panel label="Engine configuration"><PanelHeading eyebrow="Runtime configuration" title="Active worker settings" action={<UsersRound size={20} aria-hidden="true" />} /><div className="queue-counter-list"><div><span>Workers</span><strong>{queue?.workers ?? "Unavailable"}</strong></div><div><span>Queue capacity</span><strong>{queue?.queue_capacity ?? "Unavailable"}</strong></div><div><span>Retry rate</span><strong>{queue ? retryRate : "Unavailable"}</strong></div><div><span>Automation events</span><strong>{queue?.events ?? "Unavailable"}</strong></div></div><p className="automation-footnote">Runtime settings are established when the API starts and cannot be changed from the browser.</p></Panel>
    </div>
  </section>;
}
