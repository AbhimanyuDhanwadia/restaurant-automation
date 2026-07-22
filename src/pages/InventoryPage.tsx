/**
 * InventoryPage — /inventory
 *
 * Master-detail view of inventory items with status filtering and a receive-stock action.
 */

import { ClipboardCheck } from "lucide-react";
import { useMemo, useState } from "react";
import { DirectoryRow } from "@/components/ui/DirectoryRow";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import { useInventoryStore } from "@/stores/inventory";
import { toSlug } from "@/lib/utils";
import type { InventoryStatus } from "@/types/domain";

export function InventoryPage() {
  const items = useInventoryStore((s) => s.items);
  const receiveItem = useInventoryStore((s) => s.receiveItem);
  const [filter, setFilter] = useState<InventoryStatus | "All">("All");
  const [selectedId, setSelectedId] = useState(items[0]?.id ?? "");

  const directory = useMemo(
    () => (filter === "All" ? items : items.filter((i) => i.status === filter)),
    [items, filter],
  );

  const selected = directory.find((i) => i.id === selectedId) ?? directory[0];

  return (
    <section className="inventory-workspace">
      <div className="orders-toolbar">
        <div>
          <p className="eyebrow">Stock control</p>
          <h2>Inventory &amp; Purchasing</h2>
        </div>
        <label className="filter-control">
          <span>Status</span>
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value as InventoryStatus | "All")}
          >
            <option value="All">All items</option>
            <option value="In Stock">In stock</option>
            <option value="Low Stock">Low stock</option>
            <option value="On Order">On order</option>
          </select>
        </label>
      </div>

      <div className="inventory-directory">
        <Panel className="inventory-list" label="Inventory items">
          {directory.map((item) => (
            <DirectoryRow
              key={item.id}
              selected={selected?.id === item.id}
              onClick={() => setSelectedId(item.id)}
              primary={item.item}
              secondary={`${item.category} · ${item.onHand} on hand`}
              badge={
                <StatusPill variant={`inventory-status-${toSlug(item.status)}`}>
                  {item.status}
                </StatusPill>
              }
              label={`Select ${item.item}`}
            />
          ))}
          {directory.length === 0 && (
            <EmptyState message="No inventory matches this status." />
          )}
        </Panel>

        {selected && (
          <Panel className="order-detail" label={`Details for ${selected.item}`}>
            <PanelHeading
              eyebrow="Selected item"
              title={selected.item}
              action={
                <StatusPill variant={`inventory-status-${toSlug(selected.status)}`}>
                  {selected.status}
                </StatusPill>
              }
            />
            <dl className="detail-list">
              <div><dt>On hand</dt><dd>{selected.onHand}</dd></div>
              <div><dt>Par level</dt><dd>{selected.par}</dd></div>
              <div><dt>Supplier</dt><dd>{selected.supplier}</dd></div>
            </dl>
            <button
              type="button"
              className="back-button"
              onClick={() => receiveItem(selected.id)}
              disabled={selected.status === "In Stock"}
            >
              <ClipboardCheck size={16} aria-hidden="true" />
              Mark received
            </button>
          </Panel>
        )}
      </div>
    </section>
  );
}
