import { CheckCircle2, ChefHat, Clock3, Printer, RefreshCw } from "lucide-react";
import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import { getAnalyticsOverview } from "@/features/analytics/api";
import { listOrders, type Order, type OrderStatus } from "@/features/orders/api";

const KITCHEN_STATUSES: OrderStatus[] = ["received", "preparing", "ready", "delivered", "cancelled"];

function statusLabel(status: OrderStatus) {
  return `${status.slice(0, 1).toUpperCase()}${status.slice(1)}`;
}

function statusVariant(status: OrderStatus) {
  if (status === "received") return "preparing";
  if (status === "cancelled") return "delayed";
  return status;
}

function orderSummary(order: Order) {
  return order.items.map((item) => `${item.quantity}x ${item.name}`).join(", ");
}

export function KitchenPerformancePage() {
  const performanceQuery = useQuery({
    queryKey: ["analytics", "kitchen-performance"],
    queryFn: async () => {
      const [report, orders] = await Promise.all([getAnalyticsOverview(), listOrders()]);
      return { report, orders };
    },
    refetchInterval: 15_000,
  });
  const report = performanceQuery.data?.report;
  const orders = useMemo(() => performanceQuery.data?.orders ?? [], [performanceQuery.data?.orders]);
  const statusCounts = useMemo(() => {
    const counts = new Map<OrderStatus, number>(KITCHEN_STATUSES.map((status) => [status, 0]));
    for (const order of orders) counts.set(order.status, (counts.get(order.status) ?? 0) + 1);
    return counts;
  }, [orders]);
  const actionableOrders = orders.filter((order) => order.status === "received" || order.status === "preparing");
  const readyOrders = statusCounts.get("ready") ?? 0;
  const completedOrders = (statusCounts.get("ready") ?? 0) + (statusCounts.get("delivered") ?? 0);
  const totalNonCancelled = orders.length - (statusCounts.get("cancelled") ?? 0);
  const completionRate = totalNonCancelled > 0 ? `${((completedOrders / totalNonCancelled) * 100).toFixed(0)}%` : "0%";

  return <section className="kitchen-performance-workspace">
    <div className="automation-toolbar">
      <div><p className="eyebrow">Analytics</p><h2>Kitchen Performance</h2><p className="automation-subtitle">Order-state throughput and printer telemetry for the active restaurant.</p></div>
      <button type="button" className="icon-button" onClick={() => performanceQuery.refetch()} disabled={performanceQuery.isFetching} title="Refresh kitchen performance" aria-label="Refresh kitchen performance"><RefreshCw size={18} aria-hidden="true" className={performanceQuery.isFetching ? "spin" : undefined} /></button>
    </div>

    {performanceQuery.isError && <div className="orders-api-error" role="alert"><span>{performanceQuery.error.message}</span><button type="button" className="text-button" onClick={() => performanceQuery.refetch()}>Retry</button></div>}

    <section className="analytics-metric-grid" aria-label="Kitchen performance metrics">
      <article className="stat-card"><ChefHat size={22} aria-hidden="true" /><div><span>Kitchen backlog</span><strong>{performanceQuery.data ? actionableOrders.length : "Unavailable"}</strong><small>Received or preparing</small></div></article>
      <article className="stat-card"><CheckCircle2 size={22} aria-hidden="true" /><div><span>Ready now</span><strong>{performanceQuery.data ? readyOrders : "Unavailable"}</strong><small>Awaiting pickup or delivery</small></div></article>
      <article className="stat-card"><Clock3 size={22} aria-hidden="true" /><div><span>Order completion</span><strong>{performanceQuery.data ? completionRate : "Unavailable"}</strong><small>Ready or delivered orders</small></div></article>
      <article className="stat-card"><Printer size={22} aria-hidden="true" /><div><span>Printer tickets</span><strong>{report?.printer_utilization.available ? report.printer_utilization.value : "Unavailable"}</strong><small>Since API startup</small></div></article>
    </section>

    <div className="kitchen-performance-grid">
      <Panel label="Kitchen order states"><PanelHeading eyebrow="Current workload" title="Order-state distribution" action={<ChefHat size={20} aria-hidden="true" />} /><div className="kitchen-performance-states">{KITCHEN_STATUSES.map((status) => <div key={status}><span>{statusLabel(status)}</span><strong>{performanceQuery.data ? statusCounts.get(status) ?? 0 : "Unavailable"}</strong></div>)}</div></Panel>
      <Panel label="Kitchen telemetry"><PanelHeading eyebrow="Live telemetry" title="Automation coverage" action={<Printer size={20} aria-hidden="true" />} /><div className="performance-list"><div><span>Orders observed by engine</span><strong>{report?.orders.available ? report.orders.value : "Unavailable"}</strong></div><div><span>Orders queued by engine</span><strong>{report?.kitchen_completion.available ? `${report.kitchen_completion.value.toFixed(0)}%` : "Unavailable"}</strong></div><div><span>Printer availability</span><strong>{report?.printer_availability.available ? `${report.printer_availability.value.toFixed(0)}%` : "Unavailable"}</strong></div></div><p className="automation-footnote">Cycle-time reporting will become available once kitchen stage timestamps are persisted.</p></Panel>
    </div>

    <Panel label="Kitchen backlog"><PanelHeading eyebrow="Attention queue" title="Orders needing kitchen action" action={<span className="task-count">{performanceQuery.data ? `${actionableOrders.length} active` : "Unavailable"}</span>} /><div className="kitchen-performance-orders">{performanceQuery.isPending && <EmptyState message="Loading kitchen performance..." />}{actionableOrders.map((order) => <article className="kitchen-performance-order" key={order.id}><div><strong>{order.id.slice(0, 8)}</strong><span>{order.channel} · {orderSummary(order)}</span></div><time>{new Date(order.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}</time><StatusPill variant={statusVariant(order.status)}>{statusLabel(order.status)}</StatusPill></article>)}{!performanceQuery.isPending && !performanceQuery.isError && actionableOrders.length === 0 && <EmptyState message="No orders currently need kitchen action." />}</div></Panel>
  </section>;
}
