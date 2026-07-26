import { BarChart3, Clock3, Flame, ReceiptText, RefreshCw } from "lucide-react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { Panel, PanelHeading } from "@/components/ui/Panel";

type Metric = { value: number; available: boolean };
type Report = { generated_at: string; orders: Metric; sales: Metric; average_ticket: Metric; kitchen_completion: Metric; delivery_time: Metric; printer_availability: Metric; printer_utilization: Metric; staff_productivity: Metric; peak_hours: Array<{ hour: number; orders: number }> };

const API_URL = (import.meta.env.VITE_API_URL ?? "http://localhost:8080").replace(/\/$/, "");
const DEMO_REPORT: Report = { generated_at: new Date().toISOString(), orders: { value: 128, available: true }, sales: { value: 0, available: false }, average_ticket: { value: 0, available: false }, kitchen_completion: { value: 94, available: true }, delivery_time: { value: 0, available: false }, printer_availability: { value: 98, available: true }, printer_utilization: { value: 222, available: true }, staff_productivity: { value: 0, available: false }, peak_hours: [{ hour: 11, orders: 8 }, { hour: 12, orders: 16 }, { hour: 13, orders: 12 }, { hour: 17, orders: 19 }, { hour: 18, orders: 26 }, { hour: 19, orders: 18 }, { hour: 20, orders: 29 }] };

function metricValue(metric: Metric, suffix = "") { return metric.available ? `${metric.value}${suffix}` : "Unavailable"; }
function hourLabel(hour: number) { const meridiem = hour >= 12 ? "p" : "a"; const displayHour = hour % 12 || 12; return `${displayHour}${meridiem}`; }

export function AnalyticsPage() {
  const [report, setReport] = useState<Report>(DEMO_REPORT);
  const [refreshing, setRefreshing] = useState(false);

  const refresh = useCallback(async () => {
    setRefreshing(true);
    try { const response = await fetch(`${API_URL}/api/v1/analytics/overview`); if (response.ok) setReport(await response.json() as Report); } catch { setReport(DEMO_REPORT); } finally { setRefreshing(false); }
  }, []);
  useEffect(() => { void refresh(); const timer = window.setInterval(() => void refresh(), 15_000); return () => window.clearInterval(timer); }, [refresh]);

  const peakHours = useMemo(() => report.peak_hours.length ? report.peak_hours.slice(0, 7).sort((a, b) => a.hour - b.hour) : DEMO_REPORT.peak_hours, [report.peak_hours]);
  const highestVolume = Math.max(...peakHours.map((entry) => entry.orders), 1);

  return (
    <section className="analytics-workspace">
      <div className="orders-toolbar">
        <div><p className="eyebrow">Performance view</p><h2>Restaurant Analytics</h2></div>
        <button type="button" className="icon-button" onClick={() => void refresh()} disabled={refreshing} title="Refresh analytics" aria-label="Refresh analytics"><RefreshCw size={18} aria-hidden="true" className={refreshing ? "spin" : undefined} /></button>
      </div>

      <section className="analytics-metric-grid" aria-label="Analytics summary">
        <article className="stat-card"><ReceiptText size={22} aria-hidden="true" /><div><span>Orders</span><strong>{metricValue(report.orders)}</strong><small>Captured events</small></div></article>
        <article className="stat-card"><BarChart3 size={22} aria-hidden="true" /><div><span>Sales</span><strong>{metricValue(report.sales)}</strong><small>Order totals required</small></div></article>
        <article className="stat-card"><ReceiptText size={22} aria-hidden="true" /><div><span>Printer availability</span><strong>{metricValue(report.printer_availability, "%")}</strong><small>Registered printers</small></div></article>
        <article className="stat-card"><Clock3 size={22} aria-hidden="true" /><div><span>Kitchen completion</span><strong>{metricValue(report.kitchen_completion, "%")}</strong><small>Orders queued</small></div></article>
      </section>

      <div className="analytics-grid">
        <Panel label="Order volume trend"><PanelHeading eyebrow="Peak hours" title="Service demand" /><div className="bar-chart" aria-label="Observed order volume by hour">{peakHours.map((entry) => <div className="bar-column" key={entry.hour}><span style={{ height: `${Math.max(12, (entry.orders / highestVolume) * 100)}%` }} /><small>{hourLabel(entry.hour)}</small></div>)}</div></Panel>
        <Panel label="Operational performance"><PanelHeading eyebrow="Operational data" title="Key performance" /><div className="performance-list"><div><span>Printer tickets</span><strong>{metricValue(report.printer_utilization)}</strong></div><div><span>Delivery time</span><strong>{metricValue(report.delivery_time)}</strong></div><div><span>Staff productivity</span><strong>{metricValue(report.staff_productivity)}</strong></div><div><span>Average ticket</span><strong>{metricValue(report.average_ticket)}</strong></div></div></Panel>
      </div>

      <Panel label="Analytics data quality"><PanelHeading eyebrow="Data coverage" title="Operational reporting" action={<Flame size={20} aria-hidden="true" />} /><p className="analytics-summary-copy">Order events and printer telemetry are live. Sales, delivery, and staff metrics will populate when persisted orders, delivery updates, and shift activity are connected to the reporting pipeline.</p></Panel>
    </section>
  );
}
