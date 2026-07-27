/** Durable orders workspace backed by the Go REST API. */

import { ArrowLeft, RefreshCw } from "lucide-react";
import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { DirectoryRow } from "@/components/ui/DirectoryRow";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import {
  createOrder,
  listOrders,
  updateOrderStatus,
  type Order,
  type OrderStatus,
} from "@/features/orders/api";

const statuses: OrderStatus[] = ["received", "preparing", "ready", "delivered", "cancelled"];

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

export function OrdersPage() {
  const [statusFilter, setStatusFilter] = useState<OrderStatus | "all">("all");
  const [selectedId, setSelectedId] = useState("");
  const [channel, setChannel] = useState("counter");
  const [itemName, setItemName] = useState("");
  const [quantity, setQuantity] = useState(1);
  const [notes, setNotes] = useState("");
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const ordersQuery = useQuery({ queryKey: ["orders"], queryFn: listOrders });
  const refreshOrders = () => queryClient.invalidateQueries({ queryKey: ["orders"] });
  const createMutation = useMutation({
    mutationFn: createOrder,
    onSuccess: (order) => {
      setSelectedId(order.id);
      setItemName("");
      setQuantity(1);
      setNotes("");
      refreshOrders();
    },
  });
  const statusMutation = useMutation({
    mutationFn: ({ id, status }: { id: string; status: OrderStatus }) => updateOrderStatus(id, status),
    onSuccess: refreshOrders,
  });

  const orders = useMemo(() => ordersQuery.data ?? [], [ordersQuery.data]);
  const directory = useMemo(
    () => (statusFilter === "all" ? orders : orders.filter((order) => order.status === statusFilter)),
    [orders, statusFilter],
  );
  const selected = directory.find((order) => order.id === selectedId) ?? directory[0];
  const mutationError = createMutation.error ?? statusMutation.error;

  function handleCreate(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    createMutation.mutate({
      channel: channel.trim(),
      notes: notes.trim(),
      items: [{ name: itemName.trim(), quantity }],
    });
  }

  return (
    <section className="orders-workspace">
      <div className="orders-toolbar">
        <div>
          <p className="eyebrow">Order intake</p>
          <h2>All Orders</h2>
        </div>
        <label className="filter-control">
          <span>Status</span>
          <select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value as OrderStatus | "all")}>
            <option value="all">All statuses</option>
            {statuses.map((status) => <option key={status} value={status}>{statusLabel(status)}</option>)}
          </select>
        </label>
      </div>

      <Panel label="Create durable order">
        <form className="order-create-form" onSubmit={handleCreate}>
          <label><span>Channel</span><input value={channel} onChange={(event) => setChannel(event.target.value)} required /></label>
          <label><span>Item</span><input value={itemName} onChange={(event) => setItemName(event.target.value)} required /></label>
          <label><span>Quantity</span><input type="number" min="1" value={quantity} onChange={(event) => setQuantity(Number(event.target.value))} required /></label>
          <label><span>Notes</span><input value={notes} onChange={(event) => setNotes(event.target.value)} /></label>
          <button type="submit" className="primary-button" disabled={createMutation.isPending}>
            {createMutation.isPending ? "Creating..." : "Create order"}
          </button>
        </form>
      </Panel>

      {ordersQuery.isError && (
        <div className="orders-api-error" role="alert">
          <span>{ordersQuery.error.message}</span>
          <button type="button" className="text-button" onClick={() => ordersQuery.refetch()}>
            <RefreshCw size={15} aria-hidden="true" /> Retry
          </button>
        </div>
      )}
      {mutationError && <p className="orders-api-error" role="alert">{mutationError.message}</p>}

      <div className="orders-directory">
        <Panel className="order-directory-list" label="Order list">
          {ordersQuery.isPending && <EmptyState message="Loading orders..." />}
          {directory.map((order) => (
            <DirectoryRow
              key={order.id}
              selected={selected?.id === order.id}
              onClick={() => setSelectedId(order.id)}
              primary={order.id.slice(0, 8)}
              secondary={`${order.channel} · ${orderSummary(order)}`}
              badge={<StatusPill variant={statusVariant(order.status)}>{statusLabel(order.status)}</StatusPill>}
              label={`Select order ${order.id}`}
            />
          ))}
          {!ordersQuery.isPending && !ordersQuery.isError && directory.length === 0 && <EmptyState message="No orders match this status." />}
        </Panel>

        {selected && (
          <Panel className="order-detail" label={`Details for ${selected.id}`}>
            <PanelHeading
              eyebrow="Selected order"
              title={selected.id.slice(0, 8)}
              action={
                <label className="order-status-control">
                  <span>Order status</span>
                  <select
                    value={selected.status}
                    disabled={statusMutation.isPending}
                    onChange={(event) => statusMutation.mutate({
                      id: selected.id,
                      status: event.target.value as OrderStatus,
                    })}
                  >
                    {statuses.map((status) => <option key={status} value={status}>{statusLabel(status)}</option>)}
                  </select>
                </label>
              }
            />
            <dl className="detail-list">
              <div><dt>Channel</dt><dd>{selected.channel}</dd></div>
              <div><dt>Items</dt><dd>{orderSummary(selected)}</dd></div>
              <div><dt>Created</dt><dd>{new Date(selected.created_at).toLocaleString()}</dd></div>
              <div><dt>Notes</dt><dd>{selected.notes || "None"}</dd></div>
            </dl>
            <button type="button" className="back-button" onClick={() => navigate({ to: "/" })}>
              <ArrowLeft size={16} aria-hidden="true" /> Back to operations
            </button>
          </Panel>
        )}
      </div>
    </section>
  );
}
