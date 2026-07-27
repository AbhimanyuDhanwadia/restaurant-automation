import { ArrowLeft, RefreshCw } from "lucide-react";
import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import {
  createTable,
  listTables,
  updateTableStatus,
  type TableStatus,
} from "@/features/tables/api";

const statuses: TableStatus[] = ["available", "seated", "reserved", "needs_check"];

function statusLabel(status: TableStatus) {
  return ({ available: "Available", seated: "Seated", reserved: "Reserved", needs_check: "Needs check" })[status];
}

function statusClass(status: TableStatus) {
  return status.replace("_", "-");
}

export function TablesPage() {
  const [filter, setFilter] = useState<TableStatus | "all">("all");
  const [selectedID, setSelectedID] = useState("");
  const [id, setID] = useState("");
  const [seats, setSeats] = useState(2);
  const [guest, setGuest] = useState("");
  const [reservation, setReservation] = useState("");
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const tablesQuery = useQuery({ queryKey: ["tables"], queryFn: listTables });
  const createMutation = useMutation({
    mutationFn: createTable,
    onSuccess: (table) => {
      setSelectedID(table.id);
      setID("");
      setSeats(2);
      setGuest("");
      setReservation("");
      queryClient.invalidateQueries({ queryKey: ["tables"] });
    },
  });
  const statusMutation = useMutation({
    mutationFn: ({ tableID, status }: { tableID: string; status: TableStatus }) => updateTableStatus(tableID, status),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["tables"] }),
  });
  const tables = useMemo(() => tablesQuery.data ?? [], [tablesQuery.data]);
  const directory = useMemo(() => (filter === "all" ? tables : tables.filter((table) => table.status === filter)), [filter, tables]);
  const selected = directory.find((table) => table.id === selectedID) ?? directory[0];
  const mutationError = createMutation.error ?? statusMutation.error;

  function handleCreate(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    createMutation.mutate({ id: id.trim(), seats, guest: guest.trim(), reservation: reservation.trim() });
  }

  return (
    <section className="tables-workspace">
      <div className="orders-toolbar">
        <div><p className="eyebrow">Floor plan</p><h2>Table Status</h2></div>
        <label className="filter-control"><span>Status</span><select value={filter} onChange={(event) => setFilter(event.target.value as TableStatus | "all")}><option value="all">All tables</option>{statuses.map((status) => <option key={status} value={status}>{statusLabel(status)}</option>)}</select></label>
      </div>

      <Panel label="Create restaurant table">
        <form className="table-create-form" onSubmit={handleCreate}>
          <label><span>Table ID</span><input value={id} onChange={(event) => setID(event.target.value)} placeholder="T1" required /></label>
          <label><span>Seats</span><input type="number" min="1" value={seats} onChange={(event) => setSeats(Number(event.target.value))} required /></label>
          <label><span>Guest</span><input value={guest} onChange={(event) => setGuest(event.target.value)} /></label>
          <label><span>Reservation</span><input value={reservation} onChange={(event) => setReservation(event.target.value)} /></label>
          <button type="submit" className="primary-button" disabled={createMutation.isPending}>{createMutation.isPending ? "Creating..." : "Create table"}</button>
        </form>
      </Panel>

      {tablesQuery.isError && <div className="orders-api-error" role="alert"><span>{tablesQuery.error.message}</span><button type="button" className="text-button" onClick={() => tablesQuery.refetch()}><RefreshCw size={15} aria-hidden="true" /> Retry</button></div>}
      {mutationError && <p className="orders-api-error" role="alert">{mutationError.message}</p>}

      <div className="tables-directory">
        <Panel className="table-grid" label="Restaurant table grid">
          {tablesQuery.isPending && <EmptyState message="Loading tables..." />}
          {directory.map((table) => <button key={table.id} type="button" className={`table-tile ${statusClass(table.status)}${selected?.id === table.id ? " selected" : ""}`} onClick={() => setSelectedID(table.id)}><strong>{table.id}</strong><span>{table.seats} seats</span><small>{statusLabel(table.status)}</small></button>)}
          {!tablesQuery.isPending && !tablesQuery.isError && directory.length === 0 && <EmptyState message="No tables match this status." />}
        </Panel>

        {selected && <Panel className="order-detail" label={`Details for ${selected.id}`}><PanelHeading eyebrow="Selected table" title={selected.id} action={<label className="order-status-control"><span>Table status</span><select value={selected.status} disabled={statusMutation.isPending} onChange={(event) => statusMutation.mutate({ tableID: selected.id, status: event.target.value as TableStatus })}>{statuses.map((status) => <option key={status} value={status}>{statusLabel(status)}</option>)}</select></label>} /><dl className="detail-list"><div><dt>Capacity</dt><dd>{selected.seats} guests</dd></div><div><dt>Guest</dt><dd>{selected.guest || "None"}</dd></div><div><dt>Reservation</dt><dd>{selected.reservation || "None"}</dd></div><div><dt>Status</dt><dd><StatusPill variant={`table-status-${statusClass(selected.status)}`}>{statusLabel(selected.status)}</StatusPill></dd></div></dl><button type="button" className="back-button" onClick={() => navigate({ to: "/" })}><ArrowLeft size={16} aria-hidden="true" /> Back to operations</button></Panel>}
      </div>
    </section>
  );
}
