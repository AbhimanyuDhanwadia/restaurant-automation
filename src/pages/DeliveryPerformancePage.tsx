import { CircleAlert, RefreshCw, Truck } from "lucide-react";
import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import { getAnalyticsOverview } from "@/features/analytics/api";
import { listOrders, type OrderStatus } from "@/features/orders/api";

function statusLabel(status: OrderStatus) {
  return `${status.slice(0, 1).toUpperCase()}${status.slice(1)}`;
}

function statusVariant(status: OrderStatus) {
  if (status === "received") return "preparing";
  if (status === "cancelled") return "delayed";
  return status;
}

function orderSummary(items: Array<{ name: string; quantity: number }>) {
  return items.map((item) => `${item.quantity}x ${item.name}`).join(", ");
}

export function DeliveryPerformancePage() {
  const deliveryQuery = useQuery({
    queryKey: ["analytics", "delivery-performance"],
    queryFn: async () => {
      const [report, orders] = await Promise.all([getAnalyticsOverview(), listOrders()]);
      return { report, orders };
    },
    refetchInterval: 15_000,
  });
  const report = deliveryQuery.data?.report;
  const deliveryOrders = useMemo(() => (deliveryQuery.data?.orders ?? []).filter((order) => Boolean(order.delivery_partner)), [deliveryQuery.data?.orders]);
  const activeOrders = deliveryOrders.filter((order) => order.status !== "delivered" && order.status !== "cancelled");
  const providers = useMemo(() => {
    const grouped = new Map<string, { name: string; active: number; delivered: number }>();
    for (const order of deliveryOrders) {
      const name = order.delivery_partner || "Unassigned";
      const current = grouped.get(name) ?? { name, active: 0, delivered: 0 };
      if (order.status === "delivered") current.delivered += 1;
      else if (order.status !== "cancelled") current.active += 1;
      grouped.set(name, current);
    }
    return [...grouped.values()].sort((left, right) => right.active - left.active || right.delivered - left.delivered || left.name.localeCompare(right.name));
  }, [deliveryOrders]);
  const deliveryMetric = report?.delivery_time;

  return <section className="delivery-performance-workspace">
    <div className="automation-toolbar">
      <div><p className="eyebrow">Analytics</p><h2>Delivery Performance</h2><p className="automation-subtitle">Durable delivery throughput and intake-to-delivery time by provider.</p></div>
      <button type="button" className="icon-button" onClick={() => deliveryQuery.refetch()} disabled={deliveryQuery.isFetching} title="Refresh delivery performance" aria-label="Refresh delivery performance"><RefreshCw size={18} aria-hidden="true" className={deliveryQuery.isFetching ? "spin" : undefined} /></button>
    </div>

    {deliveryQuery.isError && <div className="orders-api-error" role="alert"><span>{deliveryQuery.error.message}</span><button type="button" className="text-button" onClick={() => deliveryQuery.refetch()}>Retry</button></div>}

    <section className="analytics-metric-grid" aria-label="Delivery performance metrics">
      <article className="stat-card"><Truck size={22} aria-hidden="true" /><div><span>Active deliveries</span><strong>{deliveryQuery.data ? deliveryMetric?.active_orders ?? 0 : "Unavailable"}</strong><small>Non-cancelled orders in progress</small></div></article>
      <article className="stat-card"><Truck size={22} aria-hidden="true" /><div><span>Delivered orders</span><strong>{deliveryQuery.data ? deliveryMetric?.delivered_orders ?? 0 : "Unavailable"}</strong><small>With a completion timestamp</small></div></article>
      <article className="stat-card"><CircleAlert size={22} aria-hidden="true" /><div><span>Average delivery time</span><strong>{deliveryMetric?.available ? `${deliveryMetric.value.toFixed(0)} min` : "Unavailable"}</strong><small>Intake to delivered</small></div></article>
      <article className="stat-card"><Truck size={22} aria-hidden="true" /><div><span>Delivery partners</span><strong>{deliveryQuery.data ? providers.length : "Unavailable"}</strong><small>Active provider records</small></div></article>
    </section>

    <div className="delivery-performance-grid">
      <Panel label="Delivery providers"><PanelHeading eyebrow="Provider workload" title="Current delivery partners" action={<Truck size={20} aria-hidden="true" />} /><div className="delivery-provider-list">{deliveryQuery.isPending && <EmptyState message="Loading delivery performance..." />}{providers.map((provider) => <div key={provider.name}><span>{provider.name}</span><strong>{provider.active} active</strong><small>{provider.delivered} delivered</small></div>)}{!deliveryQuery.isPending && !deliveryQuery.isError && providers.length === 0 && <EmptyState message="No delivery-provider orders are available yet." />}</div></Panel>
      <Panel label="Delivery reporting scope"><PanelHeading eyebrow="Data coverage" title="Measurement policy" action={<CircleAlert size={20} aria-hidden="true" />} /><p className="analytics-summary-copy">Only orders with a delivery partner are eligible. Delivery duration begins when the order is created and ends at the first delivered status. Counter and dine-in orders are excluded from this metric.</p></Panel>
    </div>

    <Panel label="Active delivery orders"><PanelHeading eyebrow="Delivery queue" title="Orders in progress" action={<span className="task-count">{deliveryQuery.data ? `${activeOrders.length} active` : "Unavailable"}</span>} /><div className="delivery-order-list">{deliveryQuery.isPending && <EmptyState message="Loading active deliveries..." />}{activeOrders.map((order) => <article className="delivery-order-row" key={order.id}><div><strong>{order.id.slice(0, 8)}</strong><span>{order.delivery_partner} · {orderSummary(order.items)}</span></div><time>{new Date(order.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}</time><StatusPill variant={statusVariant(order.status)}>{statusLabel(order.status)}</StatusPill></article>)}{!deliveryQuery.isPending && !deliveryQuery.isError && activeOrders.length === 0 && <EmptyState message="No delivery orders currently need attention." />}</div></Panel>
  </section>;
}
