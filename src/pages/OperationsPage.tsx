/** Live operations dashboard with a durable active-orders summary. */

import { AlertTriangle, CheckCircle2, Clock3, Flame, ReceiptText, RefreshCw, X } from "lucide-react";
import { useMemo } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import { listOrders, type Order, type OrderStatus } from "@/features/orders/api";
import { completeShiftTask, listShiftTasks } from "@/features/staff/api";
import { useAlertsStore } from "@/stores/alerts";
import { useUIStore } from "@/stores/ui";

const activeStatuses: OrderStatus[] = ["received", "preparing", "ready"];

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

export function OperationsPage() {
  const alerts = useAlertsStore((state) => state.alerts);
  const acknowledgeAlert = useAlertsStore((state) => state.acknowledge);
  const searchTerm = useUIStore((state) => state.searchTerm);
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const ordersQuery = useQuery({ queryKey: ["orders"], queryFn: listOrders });
  const tasksQuery = useQuery({ queryKey: ["staff", "tasks"], queryFn: listShiftTasks });
  const completeTask = useMutation({
    mutationFn: completeShiftTask,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["staff", "tasks"] }),
  });
  const orders = useMemo(() => ordersQuery.data ?? [], [ordersQuery.data]);
  const tasks = tasksQuery.data ?? [];
  const activeOrders = useMemo(
    () => orders.filter((order) => activeStatuses.includes(order.status)),
    [orders],
  );
  const filteredOrders = useMemo(() => {
    const query = searchTerm.trim().toLowerCase();
    if (!query) return activeOrders;
    return activeOrders.filter((order) => [
      order.id,
      order.channel,
      orderSummary(order),
      order.notes,
      order.status,
    ].join(" ").toLowerCase().includes(query));
  }, [activeOrders, searchTerm]);
  const count = (status: OrderStatus) => orders.filter((order) => order.status === status).length;
  const summary = [
    { label: "Active orders", value: activeOrders.length, detail: `${count("received")} queued`, icon: ReceiptText },
    { label: "Preparing", value: count("preparing"), detail: "In the kitchen", icon: Flame },
    { label: "Ready", value: count("ready"), detail: "Awaiting handoff", icon: Clock3 },
    { label: "Delivered", value: count("delivered"), detail: "Recorded orders", icon: CheckCircle2 },
  ];

  return (
    <>
      <section className="stats-grid" aria-label="Operational summary">
        {summary.map((stat) => (
          <article className="stat-card" key={stat.label}>
            <stat.icon size={22} aria-hidden="true" />
            <div><span>{stat.label}</span><strong>{stat.value}</strong><small>{stat.detail}</small></div>
          </article>
        ))}
      </section>

      {ordersQuery.isError && (
        <div className="orders-api-error" role="alert">
          <span>{ordersQuery.error.message}</span>
          <button type="button" className="text-button" onClick={() => ordersQuery.refetch()}>
            <RefreshCw size={15} aria-hidden="true" /> Retry
          </button>
        </div>
      )}
      {tasksQuery.isError && (
        <div className="orders-api-error" role="alert">
          <span>{tasksQuery.error.message}</span>
          <button type="button" className="text-button" onClick={() => tasksQuery.refetch()}>
            <RefreshCw size={15} aria-hidden="true" /> Retry
          </button>
        </div>
      )}
      {completeTask.error && <p className="orders-api-error" role="alert">{completeTask.error.message}</p>}

      <section className="content-grid">
        <Panel label="Active Orders" id="orders">
          <PanelHeading
            eyebrow="Kitchen queue"
            title="Active Orders"
            action={<button type="button" className="text-button" onClick={() => navigate({ to: "/orders" })}>View all</button>}
          />
          <div className="order-list">
            {ordersQuery.isPending && <EmptyState message="Loading active orders..." />}
            {filteredOrders.map((order) => (
              <article className="order-row" key={order.id}>
                <div><strong>{order.id.slice(0, 8)}</strong><span>{order.channel}</span></div>
                <span>{orderSummary(order)}</span>
                <span className="eta"><Clock3 size={16} aria-hidden="true" />{new Date(order.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}</span>
                <StatusPill variant={statusVariant(order.status)}>{statusLabel(order.status)}</StatusPill>
              </article>
            ))}
            {!ordersQuery.isPending && !ordersQuery.isError && filteredOrders.length === 0 && (
              <EmptyState message={searchTerm ? `No orders match "${searchTerm}".` : "No active orders."} />
            )}
          </div>
        </Panel>

        <aside className="side-stack">
          <Panel label="Alerts" id="alerts">
            <PanelHeading eyebrow="Attention" title="Alerts" action={<AlertTriangle size={20} aria-hidden="true" />} />
            <div className="alert-list">
              {alerts.map((alert) => (
                <article className="alert-row" key={alert.label}>
                  <div><strong>{alert.label}</strong><span>{alert.detail}</span></div>
                  <small>{alert.severity}</small>
                  <button type="button" className="dismiss-button" aria-label={`Dismiss ${alert.label} alert`} onClick={() => acknowledgeAlert(alert.label)}>
                    <X size={16} aria-hidden="true" />
                  </button>
                </article>
              ))}
              {alerts.length === 0 && <EmptyState message="All alerts are cleared." />}
            </div>
          </Panel>

          <Panel label="Next tasks">
            <PanelHeading eyebrow="Automation" title="Next Tasks" />
            <div className="task-list">
              {tasksQuery.isPending && <EmptyState message="Loading shift tasks..." />}
              {tasks.map((task) => (
                <article className="task-row" key={task.id}>
                  <div><strong>{task.title}</strong><span>{task.owner}</span></div>
                  <time>{task.due_label || "No due time"}</time>
                  <button type="button" className="task-complete" disabled={completeTask.isPending} onClick={() => completeTask.mutate(task.id)}>Done</button>
                </article>
              ))}
              {!tasksQuery.isPending && !tasksQuery.isError && tasks.length === 0 && <EmptyState message="No pending tasks." />}
            </div>
          </Panel>
        </aside>
      </section>
    </>
  );
}
