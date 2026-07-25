import { BarChart3, Clock3, Flame, ReceiptText } from "lucide-react";
import { useState } from "react";
import { Panel, PanelHeading } from "@/components/ui/Panel";

const RANGE_DATA = {
  Today: { orders: "128", revenue: "$4,860", averageTicket: "$38", prepTime: "14 min", bars: [52, 66, 48, 74, 82, 61, 88] },
  "7 days": { orders: "842", revenue: "$31,420", averageTicket: "$37", prepTime: "13 min", bars: [64, 72, 58, 82, 76, 68, 91] },
  "30 days": { orders: "3,610", revenue: "$136,800", averageTicket: "$38", prepTime: "12 min", bars: [71, 80, 68, 86, 79, 75, 94] },
} as const;

type Range = keyof typeof RANGE_DATA;

export function AnalyticsPage() {
  const [range, setRange] = useState<Range>("Today");
  const data = RANGE_DATA[range];

  return (
    <section className="analytics-workspace">
      <div className="orders-toolbar">
        <div>
          <p className="eyebrow">Performance view</p>
          <h2>Restaurant Analytics</h2>
        </div>
        <label className="filter-control">
          <span>Range</span>
          <select value={range} onChange={(event) => setRange(event.target.value as Range)}>
            <option value="Today">Today</option>
            <option value="7 days">Last 7 days</option>
            <option value="30 days">Last 30 days</option>
          </select>
        </label>
      </div>

      <section className="analytics-metric-grid" aria-label="Analytics summary">
        <article className="stat-card"><ReceiptText size={22} aria-hidden="true" /><div><span>Orders</span><strong>{data.orders}</strong><small>{range}</small></div></article>
        <article className="stat-card"><BarChart3 size={22} aria-hidden="true" /><div><span>Revenue</span><strong>{data.revenue}</strong><small>Gross sales</small></div></article>
        <article className="stat-card"><ReceiptText size={22} aria-hidden="true" /><div><span>Average ticket</span><strong>{data.averageTicket}</strong><small>Per order</small></div></article>
        <article className="stat-card"><Clock3 size={22} aria-hidden="true" /><div><span>Prep time</span><strong>{data.prepTime}</strong><small>Kitchen average</small></div></article>
      </section>

      <div className="analytics-grid">
        <Panel label="Order volume trend">
          <PanelHeading eyebrow="Order volume" title="Service demand" />
          <div className="bar-chart" aria-label={`Mock order volume for ${range}`}>
            {data.bars.map((height, index) => <div className="bar-column" key={`${range}-${index}`}><span style={{ height: `${height}%` }} /><small>{["11a", "12p", "1p", "5p", "6p", "7p", "8p"][index]}</small></div>)}
          </div>
        </Panel>
        <Panel label="Operational performance">
          <PanelHeading eyebrow="Service health" title="Key performance" />
          <div className="performance-list">
            <div><span>Orders ready on time</span><strong>94%</strong></div>
            <div><span>Kitchen throughput</span><strong>87%</strong></div>
            <div><span>Table turn time</span><strong>68 min</strong></div>
            <div><span>Inventory variance</span><strong>2.4%</strong></div>
          </div>
        </Panel>
      </div>

      <Panel label="Operations summary">
        <PanelHeading eyebrow="Manager view" title="What changed" action={<Flame size={20} aria-hidden="true" />} />
        <p className="analytics-summary-copy">Grill demand is the current constraint. Friday service is trending above the weekly average, while prep time remains within the target window.</p>
      </Panel>
    </section>
  );
}
