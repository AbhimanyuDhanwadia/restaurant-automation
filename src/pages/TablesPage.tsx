/**
 * TablesPage — /tables
 *
 * Grid overview of all restaurant tables with status filtering and an
 * inline detail panel for the selected table.
 */

import { ArrowLeft } from "lucide-react";
import { useMemo, useState } from "react";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import { useTablesStore } from "@/stores/tables";
import { toSlug } from "@/lib/utils";
import type { TableStatus } from "@/types/domain";
import { useNavigate } from "@tanstack/react-router";

export function TablesPage() {
  const tables = useTablesStore((s) => s.tables);
  const cycleStatus = useTablesStore((s) => s.cycleStatus);
  const [filter, setFilter] = useState<TableStatus | "All">("All");
  const [selectedId, setSelectedId] = useState(tables[0]?.id ?? "");
  const navigate = useNavigate();

  const directory = useMemo(
    () => (filter === "All" ? tables : tables.filter((t) => t.status === filter)),
    [tables, filter],
  );

  const selected = directory.find((t) => t.id === selectedId) ?? directory[0];

  return (
    <section className="tables-workspace">
      <div className="orders-toolbar">
        <div>
          <p className="eyebrow">Floor plan</p>
          <h2>Table Status</h2>
        </div>
        <label className="filter-control">
          <span>Status</span>
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value as TableStatus | "All")}
          >
            <option value="All">All tables</option>
            <option value="Available">Available</option>
            <option value="Seated">Seated</option>
            <option value="Reserved">Reserved</option>
            <option value="Needs Check">Needs check</option>
          </select>
        </label>
      </div>

      <div className="tables-directory">
        <Panel className="table-grid" label="Restaurant table grid">
          {directory.map((table) => (
            <button
              key={table.id}
              type="button"
              className={`table-tile ${toSlug(table.status)}${selected?.id === table.id ? " selected" : ""}`}
              onClick={() => setSelectedId(table.id)}
            >
              <strong>{table.id}</strong>
              <span>{table.seats} seats</span>
              <small>{table.status}</small>
            </button>
          ))}
          {directory.length === 0 && (
            <EmptyState message="No tables match this status." />
          )}
        </Panel>

        {selected && (
          <Panel className="order-detail" label={`Details for ${selected.id}`}>
            <PanelHeading
              eyebrow="Selected table"
              title={selected.id}
              action={
                <StatusPill
                  variant={`table-status-${toSlug(selected.status)}`}
                  onClick={() => cycleStatus(selected.id)}
                >
                  {selected.status}
                </StatusPill>
              }
            />
            <dl className="detail-list">
              <div><dt>Capacity</dt><dd>{selected.seats} guests</dd></div>
              <div><dt>Guest</dt><dd>{selected.guest}</dd></div>
              <div><dt>Reservation</dt><dd>{selected.reservation}</dd></div>
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
