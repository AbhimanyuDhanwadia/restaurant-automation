import { AlertTriangle, BarChart3, PackageSearch, RefreshCw } from "lucide-react";
import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import { listInventory, type InventoryItem } from "@/features/inventory/api";

function statusLabel(status: InventoryItem["status"]) {
  return ({ in_stock: "In stock", low_stock: "Low stock", on_order: "On order" })[status];
}

function statusClass(status: InventoryItem["status"]) {
  return status.replace("_", "-");
}

function quantity(value: number, unit: string) {
  return `${new Intl.NumberFormat(undefined, { maximumFractionDigits: 2 }).format(value)} ${unit}`;
}

function stockRatio(item: InventoryItem) {
  return item.par_level > 0 ? item.on_hand / item.par_level : 1;
}

export function InventoryAnalyticsPage() {
  const inventoryQuery = useQuery({ queryKey: ["analytics", "inventory"], queryFn: listInventory, refetchInterval: 15_000 });
  const items = useMemo(() => inventoryQuery.data ?? [], [inventoryQuery.data]);
  const lowStockItems = useMemo(() => items.filter((item) => item.status === "low_stock" || (item.par_level > 0 && item.on_hand < item.par_level)).sort((left, right) => stockRatio(left) - stockRatio(right)), [items]);
  const inStockItems = items.filter((item) => item.status === "in_stock");
  const onOrderItems = items.filter((item) => item.status === "on_order");
  const coverage = items.length > 0 ? `${((inStockItems.length / items.length) * 100).toFixed(0)}%` : "0%";
  const categories = useMemo(() => {
    const grouped = new Map<string, { name: string; items: number; atRisk: number }>();
    for (const item of items) {
      const current = grouped.get(item.category) ?? { name: item.category, items: 0, atRisk: 0 };
      current.items += 1;
      if (item.status === "low_stock" || (item.par_level > 0 && item.on_hand < item.par_level)) current.atRisk += 1;
      grouped.set(item.category, current);
    }
    return [...grouped.values()].sort((left, right) => right.atRisk - left.atRisk || left.name.localeCompare(right.name));
  }, [items]);

  return <section className="inventory-analytics-workspace">
    <div className="automation-toolbar">
      <div><p className="eyebrow">Analytics</p><h2>Inventory Analytics</h2><p className="automation-subtitle">Stock coverage, replenishment risk, and category health from durable inventory records.</p></div>
      <button type="button" className="icon-button" onClick={() => inventoryQuery.refetch()} disabled={inventoryQuery.isFetching} title="Refresh inventory analytics" aria-label="Refresh inventory analytics"><RefreshCw size={18} aria-hidden="true" className={inventoryQuery.isFetching ? "spin" : undefined} /></button>
    </div>

    {inventoryQuery.isError && <div className="orders-api-error" role="alert"><span>{inventoryQuery.error.message}</span><button type="button" className="text-button" onClick={() => inventoryQuery.refetch()}>Retry</button></div>}

    <section className="analytics-metric-grid" aria-label="Inventory analytics metrics">
      <article className="stat-card"><PackageSearch size={22} aria-hidden="true" /><div><span>Tracked items</span><strong>{inventoryQuery.data ? items.length : "Unavailable"}</strong><small>Across all categories</small></div></article>
      <article className="stat-card"><BarChart3 size={22} aria-hidden="true" /><div><span>Stock coverage</span><strong>{inventoryQuery.data ? coverage : "Unavailable"}</strong><small>Items marked in stock</small></div></article>
      <article className="stat-card"><AlertTriangle size={22} aria-hidden="true" /><div><span>Below par</span><strong>{inventoryQuery.data ? lowStockItems.length : "Unavailable"}</strong><small>Need replenishment review</small></div></article>
      <article className="stat-card"><PackageSearch size={22} aria-hidden="true" /><div><span>On order</span><strong>{inventoryQuery.data ? onOrderItems.length : "Unavailable"}</strong><small>Pending replenishment</small></div></article>
    </section>

    <div className="inventory-analytics-grid">
      <Panel label="Inventory by category"><PanelHeading eyebrow="Category health" title="Tracked stock" action={<BarChart3 size={20} aria-hidden="true" />} /><div className="inventory-category-list">{inventoryQuery.isPending && <EmptyState message="Loading inventory analytics..." />}{categories.map((category) => <div key={category.name}><span>{category.name}</span><strong>{category.items} items</strong><small>{category.atRisk ? `${category.atRisk} below par` : "Coverage healthy"}</small></div>)}{!inventoryQuery.isPending && !inventoryQuery.isError && categories.length === 0 && <EmptyState message="No inventory categories are available yet." />}</div></Panel>
      <Panel label="Inventory data coverage"><PanelHeading eyebrow="Data coverage" title="Available reporting" action={<PackageSearch size={20} aria-hidden="true" />} /><p className="analytics-summary-copy">This workspace intentionally compares each item only against its own par level. It does not add quantities across units such as kilograms, litres, and pieces. Purchase cost, consumption velocity, and supplier lead-time reporting require future transaction and vendor records.</p></Panel>
    </div>

    <Panel label="Inventory replenishment risks"><PanelHeading eyebrow="Replenishment" title="Items below par" action={<span className="task-count">{inventoryQuery.data ? `${lowStockItems.length} at risk` : "Unavailable"}</span>} /><div className="inventory-risk-list">{inventoryQuery.isPending && <EmptyState message="Loading stock risks..." />}{lowStockItems.map((item) => { const ratio = stockRatio(item); const shortfall = Math.max(0, item.par_level - item.on_hand); return <article className="inventory-risk-row" key={item.id}><div><strong>{item.name}</strong><span>{item.category} · {quantity(item.on_hand, item.unit)} on hand of {quantity(item.par_level, item.unit)} par</span></div><strong>{ratio === 0 ? "Empty" : `${(ratio * 100).toFixed(0)}% of par`}</strong><div><span>{shortfall > 0 ? `${quantity(shortfall, item.unit)} short` : "Review status"}</span><StatusPill variant={`inventory-status-${statusClass(item.status)}`}>{statusLabel(item.status)}</StatusPill></div></article>; })}{!inventoryQuery.isPending && !inventoryQuery.isError && lowStockItems.length === 0 && <EmptyState message="All tracked inventory is at or above par." />}</div></Panel>
  </section>;
}
