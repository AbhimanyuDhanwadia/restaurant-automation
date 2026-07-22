/**
 * OrdersPage — /orders
 *
 * Master-detail view of all orders with status filtering and an inline
 * detail panel for the selected order.
 */

import { ArrowLeft } from "lucide-react";
import { useMemo, useState } from "react";
import { DirectoryRow } from "@/components/ui/DirectoryRow";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import { useOrdersStore } from "@/stores/orders";
import type { OrderStatus } from "@/types/domain";
import { useNavigate } from "@tanstack/react-router";

export function OrdersPage() {
  const orders = useOrdersStore((s) => s.orders);
  const cycleStatus = useOrdersStore((s) => s.cycleStatus);
  const [statusFilter, setStatusFilter] = useState<OrderStatus | "All">("All");
  const [selectedId, setSelectedId] = useState(orders[0]?.id ?? "");
  const navigate = useNavigate();

  const directory = useMemo(
    () => (statusFilter === "All" ? orders : orders.filter((o) => o.status === statusFilter)),
    [orders, statusFilter],
  );

  const selected = directory.find((o) => o.id === selectedId) ?? directory[0];

  return (
    <section className="orders-workspace">
      <div className="orders-toolbar">
        <div>
          <p className="eyebrow">Order intake</p>
          <h2>All Orders</h2>
        </div>
        <label className="filter-control">
          <span>Status</span>
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value as OrderStatus | "All")}
          >
            <option value="All">All statuses</option>
            <option value="Preparing">Preparing</option>
            <option value="Ready">Ready</option>
            <option value="Delayed">Delayed</option>
          </select>
        </label>
      </div>

      <div className="orders-directory">
        <Panel className="order-directory-list" label="Order list">
          {directory.map((order) => (
            <DirectoryRow
              key={order.id}
              selected={selected?.id === order.id}
              onClick={() => setSelectedId(order.id)}
              primary={order.id}
              secondary={`${order.table} · ${order.channel}`}
              badge={
                <StatusPill variant={order.status.toLowerCase()}>
                  {order.status}
                </StatusPill>
              }
              label={`Select ${order.id}`}
            />
          ))}
          {directory.length === 0 && (
            <EmptyState message="No orders match this status." />
          )}
        </Panel>

        {selected && (
          <Panel className="order-detail" label={`Details for ${selected.id}`}>
            <PanelHeading
              eyebrow="Selected order"
              title={selected.id}
              action={
                <StatusPill
                  variant={selected.status.toLowerCase()}
                  onClick={() => cycleStatus(selected.id)}
                >
                  {selected.status}
                </StatusPill>
              }
            />
            <dl className="detail-list">
              <div><dt>Service</dt><dd>{selected.table}</dd></div>
              <div><dt>Channel</dt><dd>{selected.channel}</dd></div>
              <div><dt>Items</dt><dd>{selected.items}</dd></div>
              <div><dt>Expected</dt><dd>{selected.eta}</dd></div>
            </dl>
            <button
              type="button"
              className="back-button"
              onClick={() => navigate({ to: "/" })}
            >
              <ArrowLeft size={16} aria-hidden="true" />
              Back to operations
            </button>
          </Panel>
        )}
      </div>
    </section>
  );
}
