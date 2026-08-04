import { BarChart3, CircleAlert, ReceiptText, RefreshCw } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { getAnalyticsOverview, type MonetaryMetric } from "@/features/analytics/api";

function moneyValue(metric: MonetaryMetric | undefined) {
  if (!metric?.available) return "Unavailable";
  return new Intl.NumberFormat(undefined, { style: "currency", currency: metric.currency }).format(metric.value);
}

export function SalesPage() {
  const salesQuery = useQuery({ queryKey: ["analytics", "sales"], queryFn: getAnalyticsOverview, refetchInterval: 15_000 });
  const sales = salesQuery.data?.sales;
  const averageTicket = salesQuery.data?.average_ticket;
  const includedOrders = sales?.included_orders ?? 0;
  const excludedOrders = sales?.excluded_orders ?? 0;
  const coverage = includedOrders + excludedOrders > 0 ? `${((includedOrders / (includedOrders + excludedOrders)) * 100).toFixed(0)}%` : "0%";
  const mixedCurrencies = includedOrders > 0 && !sales?.available;

  return <section className="sales-workspace">
    <div className="automation-toolbar">
      <div><p className="eyebrow">Analytics</p><h2>Sales</h2><p className="automation-subtitle">Recorded order totals with explicit coverage and currency safeguards.</p></div>
      <button type="button" className="icon-button" onClick={() => salesQuery.refetch()} disabled={salesQuery.isFetching} title="Refresh sales analytics" aria-label="Refresh sales analytics"><RefreshCw size={18} aria-hidden="true" className={salesQuery.isFetching ? "spin" : undefined} /></button>
    </div>

    {salesQuery.isError && <div className="orders-api-error" role="alert"><span>{salesQuery.error.message}</span><button type="button" className="text-button" onClick={() => salesQuery.refetch()}>Retry</button></div>}

    <section className="analytics-metric-grid" aria-label="Sales analytics metrics">
      <article className="stat-card"><BarChart3 size={22} aria-hidden="true" /><div><span>Recorded sales</span><strong>{moneyValue(sales)}</strong><small>{sales?.available ? `${sales.currency} totals` : "Currency-safe reporting"}</small></div></article>
      <article className="stat-card"><ReceiptText size={22} aria-hidden="true" /><div><span>Average ticket</span><strong>{moneyValue(averageTicket)}</strong><small>Orders with recorded totals</small></div></article>
      <article className="stat-card"><ReceiptText size={22} aria-hidden="true" /><div><span>Orders with totals</span><strong>{salesQuery.data ? includedOrders : "Unavailable"}</strong><small>Eligible non-cancelled orders</small></div></article>
      <article className="stat-card"><CircleAlert size={22} aria-hidden="true" /><div><span>Revenue coverage</span><strong>{salesQuery.data ? coverage : "Unavailable"}</strong><small>{salesQuery.data ? `${excludedOrders} eligible orders excluded` : "Order totals unavailable"}</small></div></article>
    </section>

    <div className="sales-grid">
      <Panel label="Sales reporting status"><PanelHeading eyebrow="Reporting quality" title="Recorded total coverage" action={<BarChart3 size={20} aria-hidden="true" />} /><div className="performance-list"><div><span>Included orders</span><strong>{salesQuery.data ? includedOrders : "Unavailable"}</strong></div><div><span>Orders without totals</span><strong>{salesQuery.data ? excludedOrders : "Unavailable"}</strong></div><div><span>Reporting currency</span><strong>{sales?.available ? sales.currency : mixedCurrencies ? "Multiple currencies" : "Unavailable"}</strong></div></div></Panel>
      <Panel label="Sales scope"><PanelHeading eyebrow="Scope" title="What is counted" action={<ReceiptText size={20} aria-hidden="true" />} /><p className="analytics-summary-copy">Sales include non-cancelled durable orders with a recorded total in one currency. Orders without a total are excluded and shown in coverage. Mixed currencies are never combined into a single amount.</p></Panel>
    </div>
  </section>;
}
