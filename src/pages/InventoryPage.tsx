import { ClipboardCheck, RefreshCw } from "lucide-react";
import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { DirectoryRow } from "@/components/ui/DirectoryRow";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import {
  createInventoryItem,
  listInventory,
  updateInventoryStatus,
  updateInventoryStock,
  type InventoryItem,
  type InventoryStatus,
} from "@/features/inventory/api";

const statuses: InventoryStatus[] = ["in_stock", "low_stock", "on_order"];

function statusLabel(status: InventoryStatus) {
  return ({ in_stock: "In stock", low_stock: "Low stock", on_order: "On order" })[status];
}

function statusClass(status: InventoryStatus) {
  return status.replace("_", "-");
}

function quantity(item: InventoryItem, value = item.on_hand) {
  return `${value} ${item.unit}`;
}

export function InventoryPage() {
  const [filter, setFilter] = useState<InventoryStatus | "all">("all");
  const [selectedID, setSelectedID] = useState("");
  const [name, setName] = useState("");
  const [category, setCategory] = useState("");
  const [onHand, setOnHand] = useState(0);
  const [parLevel, setParLevel] = useState(0);
  const [unit, setUnit] = useState("units");
  const [supplier, setSupplier] = useState("");
  const queryClient = useQueryClient();
  const inventoryQuery = useQuery({ queryKey: ["inventory"], queryFn: listInventory });
  const refreshInventory = () => queryClient.invalidateQueries({ queryKey: ["inventory"] });
  const createMutation = useMutation({
    mutationFn: createInventoryItem,
    onSuccess: (item) => {
      setSelectedID(item.id);
      setName("");
      setCategory("");
      setOnHand(0);
      setParLevel(0);
      setUnit("units");
      setSupplier("");
      refreshInventory();
    },
  });
  const stockMutation = useMutation({ mutationFn: ({ id, amount }: { id: string; amount: number }) => updateInventoryStock(id, amount), onSuccess: refreshInventory });
  const statusMutation = useMutation({ mutationFn: ({ id, status }: { id: string; status: InventoryStatus }) => updateInventoryStatus(id, status), onSuccess: refreshInventory });
  const items = useMemo(() => inventoryQuery.data ?? [], [inventoryQuery.data]);
  const directory = useMemo(() => (filter === "all" ? items : items.filter((item) => item.status === filter)), [filter, items]);
  const selected = directory.find((item) => item.id === selectedID) ?? directory[0];
  const mutationError = createMutation.error ?? stockMutation.error ?? statusMutation.error;

  function handleCreate(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    createMutation.mutate({ name: name.trim(), category: category.trim(), on_hand: onHand, par_level: parLevel, unit: unit.trim(), supplier: supplier.trim() });
  }

  function handleStockUpdate(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!selected) return;
    const value = new FormData(event.currentTarget).get("on_hand");
    stockMutation.mutate({ id: selected.id, amount: Number(value) });
  }

  return (
    <section className="inventory-workspace">
      <div className="orders-toolbar"><div><p className="eyebrow">Stock control</p><h2>Inventory &amp; Purchasing</h2></div><label className="filter-control"><span>Status</span><select value={filter} onChange={(event) => setFilter(event.target.value as InventoryStatus | "all")}><option value="all">All items</option>{statuses.map((status) => <option key={status} value={status}>{statusLabel(status)}</option>)}</select></label></div>

      <Panel label="Create inventory item"><form className="inventory-create-form" onSubmit={handleCreate}><label><span>Item</span><input value={name} onChange={(event) => setName(event.target.value)} required /></label><label><span>Category</span><input value={category} onChange={(event) => setCategory(event.target.value)} required /></label><label><span>On hand</span><input type="number" min="0" step="any" value={onHand} onChange={(event) => setOnHand(Number(event.target.value))} required /></label><label><span>Par level</span><input type="number" min="0" step="any" value={parLevel} onChange={(event) => setParLevel(Number(event.target.value))} required /></label><label><span>Unit</span><input value={unit} onChange={(event) => setUnit(event.target.value)} required /></label><label><span>Supplier</span><input value={supplier} onChange={(event) => setSupplier(event.target.value)} /></label><button type="submit" className="primary-button" disabled={createMutation.isPending}>{createMutation.isPending ? "Creating..." : "Create item"}</button></form></Panel>

      {inventoryQuery.isError && <div className="orders-api-error" role="alert"><span>{inventoryQuery.error.message}</span><button type="button" className="text-button" onClick={() => inventoryQuery.refetch()}><RefreshCw size={15} aria-hidden="true" /> Retry</button></div>}
      {mutationError && <p className="orders-api-error" role="alert">{mutationError.message}</p>}

      <div className="inventory-directory">
        <Panel className="inventory-list" label="Inventory items">
          {inventoryQuery.isPending && <EmptyState message="Loading inventory..." />}
          {directory.map((item) => <DirectoryRow key={item.id} selected={selected?.id === item.id} onClick={() => setSelectedID(item.id)} primary={item.name} secondary={`${item.category} · ${quantity(item)} on hand`} badge={<StatusPill variant={`inventory-status-${statusClass(item.status)}`}>{statusLabel(item.status)}</StatusPill>} label={`Select ${item.name}`} />)}
          {!inventoryQuery.isPending && !inventoryQuery.isError && directory.length === 0 && <EmptyState message="No inventory matches this status." />}
        </Panel>

        {selected && <Panel className="order-detail" label={`Details for ${selected.name}`}><PanelHeading eyebrow="Selected item" title={selected.name} action={<label className="order-status-control"><span>Stock status</span><select value={selected.status} disabled={statusMutation.isPending} onChange={(event) => statusMutation.mutate({ id: selected.id, status: event.target.value as InventoryStatus })}>{statuses.map((status) => <option key={status} value={status}>{statusLabel(status)}</option>)}</select></label>} /><dl className="detail-list"><div><dt>On hand</dt><dd>{quantity(selected)}</dd></div><div><dt>Par level</dt><dd>{quantity(selected, selected.par_level)}</dd></div><div><dt>Supplier</dt><dd>{selected.supplier || "Unassigned"}</dd></div></dl><form className="inventory-stock-form" key={selected.id} onSubmit={handleStockUpdate}><label><span>Set on hand</span><input name="on_hand" type="number" min="0" step="any" defaultValue={selected.on_hand} required /></label><button type="submit" className="back-button" disabled={stockMutation.isPending}><ClipboardCheck size={16} aria-hidden="true" />{stockMutation.isPending ? "Updating..." : "Update stock"}</button></form></Panel>}
      </div>
    </section>
  );
}
