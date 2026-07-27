import { ChefHat, Clock3, Flame, RefreshCw } from "lucide-react";
import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import {
  listOrders,
  updateOrderStatus,
  type Order,
  type OrderStatus,
} from "@/features/orders/api";

type KitchenStatus = "all" | "received" | "preparing" | "ready";

const kitchenStatuses: Exclude<KitchenStatus, "all">[] = ["received", "preparing", "ready"];

function kitchenLabel(status: Exclude<KitchenStatus, "all">) {
  return ({ received: "Queued", preparing: "Cooking", ready: "Ready" })[status];
}

function ticketVariant(status: Exclude<KitchenStatus, "all">) {
  return `ticket-${({ received: "queued", preparing: "cooking", ready: "ready" })[status]}`;
}

function orderSummary(order: Order) {
  return order.items.map((item) => `${item.quantity}x ${item.name}`).join(", ");
}

function elapsedTime(createdAt: string) {
  const minutes = Math.max(0, Math.floor((Date.now() - new Date(createdAt).getTime()) / 60_000));
  return minutes < 1 ? "Just now" : `${minutes} min`;
}

function nextKitchenStatus(status: Exclude<KitchenStatus, "all">): OrderStatus | null {
  if (status === "received") return "preparing";
  if (status === "preparing") return "ready";
  return null;
}

export function KitchenPage() {
  const [statusFilter, setStatusFilter] = useState<KitchenStatus>("all");
  const queryClient = useQueryClient();
  const ordersQuery = useQuery({ queryKey: ["orders"], queryFn: listOrders });
  const advanceMutation = useMutation({
    mutationFn: ({ id, status }: { id: string; status: OrderStatus }) => updateOrderStatus(id, status),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["orders"] }),
  });

  const activeOrders = useMemo(
    () => (ordersQuery.data ?? []).filter((order) => kitchenStatuses.includes(order.status as Exclude<KitchenStatus, "all">)),
    [ordersQuery.data],
  );
  const tickets = useMemo(
    () => (statusFilter === "all" ? activeOrders : activeOrders.filter((order) => order.status === statusFilter)),
    [activeOrders, statusFilter],
  );
  const counts = useMemo(
    () => Object.fromEntries(kitchenStatuses.map((status) => [status, activeOrders.filter((order) => order.status === status).length])) as Record<Exclude<KitchenStatus, "all">, number>,
    [activeOrders],
  );

  return (
    <section className="kitchen-workspace">
      <div className="orders-toolbar">
        <div>
          <p className="eyebrow">Production floor</p>
          <h2>Kitchen Queue</h2>
        </div>
        <label className="filter-control">
          <span>Status</span>
          <select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value as KitchenStatus)}>
            <option value="all">All active tickets</option>
            {kitchenStatuses.map((status) => <option key={status} value={status}>{kitchenLabel(status)}</option>)}
          </select>
        </label>
      </div>

      <section className="kitchen-station-grid" aria-label="Kitchen queue summary">
        {kitchenStatuses.map((status) => {
          const count = counts[status];
          return (
            <article className="stat-card kitchen-station-card" key={status}>
              {status === "preparing" ? <Flame size={22} aria-hidden="true" /> : <ChefHat size={22} aria-hidden="true" />}
              <div>
                <span>{kitchenLabel(status)}</span>
                <strong>{count}</strong>
                <small className={`station-load ${count > 5 ? "critical" : "healthy"}`}>{count === 1 ? "1 order" : `${count} orders`}</small>
              </div>
            </article>
          );
        })}
      </section>

      {ordersQuery.isError && (
        <div className="orders-api-error" role="alert">
          <span>{ordersQuery.error.message}</span>
          <button type="button" className="text-button" onClick={() => ordersQuery.refetch()}>
            <RefreshCw size={15} aria-hidden="true" /> Retry
          </button>
        </div>
      )}
      {advanceMutation.error && <p className="orders-api-error" role="alert">{advanceMutation.error.message}</p>}

      <Panel label="Kitchen tickets">
        <PanelHeading eyebrow="Active tickets" title={`${tickets.length} tickets in queue`} />
        <div className="kitchen-ticket-list">
          {ordersQuery.isPending && <EmptyState message="Loading kitchen tickets..." />}
          {tickets.map((order) => {
            const nextStatus = nextKitchenStatus(order.status as Exclude<KitchenStatus, "all">);
            return (
              <article className="kitchen-ticket" key={order.id}>
                <div className="kitchen-ticket-main">
                  <strong>{order.id.slice(0, 8)}</strong>
                  <span>{order.channel} · {orderSummary(order)}</span>
                </div>
                <span className="kitchen-station-label">{order.notes || "No preparation notes"}</span>
                <span className="eta"><Clock3 size={16} aria-hidden="true" />{elapsedTime(order.created_at)}</span>
                <StatusPill variant={ticketVariant(order.status as Exclude<KitchenStatus, "all">)}>{kitchenLabel(order.status as Exclude<KitchenStatus, "all">)}</StatusPill>
                {nextStatus && (
                  <button
                    type="button"
                    className="text-button kitchen-ticket-action"
                    disabled={advanceMutation.isPending}
                    onClick={() => advanceMutation.mutate({ id: order.id, status: nextStatus })}
                  >
                    {nextStatus === "preparing" ? "Start preparing" : "Mark ready"}
                  </button>
                )}
              </article>
            );
          })}
          {!ordersQuery.isPending && !ordersQuery.isError && tickets.length === 0 && <EmptyState message="No tickets match this status." />}
        </div>
      </Panel>
    </section>
  );
}
