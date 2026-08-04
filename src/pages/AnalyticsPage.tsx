import { BarChart3, Clock3, Flame, ReceiptText, RefreshCw } from "lucide-react";
import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { getAnalyticsOverview, type Metric, type MonetaryMetric } from "@/features/analytics/api";

function metricValue(metric: Metric | undefined, suffix = "") {
  return metric?.available ? `${metric.value}${suffix}` : "Unavailable";
}

function moneyValue(metric: MonetaryMetric | undefined) {
  return metric?.available ? new Intl.NumberFormat(undefined, { style: "currency", currency: metric.currency }).format(metric.value) : "Unavailable";
}

function hourLabel(hour: number) {
  const meridiem = hour >= 12 ? "p" : "a";
  const displayHour = hour % 12 || 12;
  return `${displayHour}${meridiem}`;
}

export function AnalyticsPage() {
  const overviewQuery = useQuery({
    queryKey: ["analytics", "overview"],
    queryFn: getAnalyticsOverview,
    refetchInterval: 15_000,
  });
  const report = overviewQuery.data;
  const peakHours = useMemo(
    () => (report?.peak_hours ?? []).slice(0, 7).sort((left, right) => left.hour - right.hour),
    [report?.peak_hours],
  );
  const highestVolume = Math.max(...peakHours.map((entry) => entry.orders), 1);

  return (
    <section className="analytics-workspace">
      <div className="orders-toolbar">
        <div><p className="eyebrow">Performance view</p><h2>Restaurant Analytics</h2></div>
        <button type="button" className="icon-button" onClick={() => overviewQuery.refetch()} disabled={overviewQuery.isFetching} title="Refresh analytics" aria-label="Refresh analytics">
          <RefreshCw size={18} aria-hidden="true" className={overviewQuery.isFetching ? "spin" : undefined} />
        </button>
      </div>

      {overviewQuery.isError && (
        <div className="orders-api-error" role="alert">
          <span>{overviewQuery.error.message}</span>
          <button type="button" className="text-button" onClick={() => overviewQuery.refetch()}>Retry</button>
        </div>
      )}

      <section className="analytics-metric-grid" aria-label="Analytics summary">
        <article className="stat-card"><ReceiptText size={22} aria-hidden="true" /><div><span>Orders</span><strong>{metricValue(report?.orders)}</strong><small>Captured events</small></div></article>
        <article className="stat-card"><BarChart3 size={22} aria-hidden="true" /><div><span>Sales</span><strong>{moneyValue(report?.sales)}</strong><small>{report?.sales.available ? `${report.sales.included_orders} totals recorded` : "Order totals required"}</small></div></article>
        <article className="stat-card"><ReceiptText size={22} aria-hidden="true" /><div><span>Printer availability</span><strong>{metricValue(report?.printer_availability, "%")}</strong><small>Registered printers</small></div></article>
        <article className="stat-card"><Clock3 size={22} aria-hidden="true" /><div><span>Kitchen completion</span><strong>{metricValue(report?.kitchen_completion, "%")}</strong><small>Orders queued</small></div></article>
      </section>

      <div className="analytics-grid">
        <Panel label="Order volume trend">
          <PanelHeading eyebrow="Peak hours" title="Service demand" />
          {overviewQuery.isPending && <EmptyState message="Loading analytics..." />}
          {!overviewQuery.isPending && !overviewQuery.isError && peakHours.length === 0 && <EmptyState message="No order-event volume is available yet." />}
          {peakHours.length > 0 && <div className="bar-chart" aria-label="Observed order volume by hour">{peakHours.map((entry) => <div className="bar-column" key={entry.hour}><span style={{ height: `${Math.max(12, (entry.orders / highestVolume) * 100)}%` }} /><small>{hourLabel(entry.hour)}</small></div>)}</div>}
        </Panel>
        <Panel label="Operational performance"><PanelHeading eyebrow="Operational data" title="Key performance" /><div className="performance-list"><div><span>Printer tickets</span><strong>{metricValue(report?.printer_utilization)}</strong></div><div><span>Delivery time</span><strong>{metricValue(report?.delivery_time)}</strong></div><div><span>Staff productivity</span><strong>{metricValue(report?.staff_productivity)}</strong></div><div><span>Average ticket</span><strong>{moneyValue(report?.average_ticket)}</strong></div></div></Panel>
      </div>

      <Panel label="Analytics data quality"><PanelHeading eyebrow="Data coverage" title="Operational reporting" action={<Flame size={20} aria-hidden="true" />} /><p className="analytics-summary-copy">Order events, printer telemetry, and recorded durable order totals are live. Delivery and staff metrics will populate when delivery updates and shift activity are connected to the reporting pipeline.</p></Panel>
    </section>
  );
}
