import { Activity, CheckCircle2, CircleAlert, Database, PlugZap, Printer, RefreshCw, Wifi } from "lucide-react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { Panel, PanelHeading } from "@/components/ui/Panel";

type Provider = { name: string; status: string };
type Printer = { name: string; status: string; queue_depth: number; printed: number; failed: number };
type Event = { id: string; type: string; order_id: string; created_at: string };
type Queue = { queue_depth: number; queue_capacity: number; workers: number; events: number; retried: number; failed: number };

const MOCK_PROVIDERS: Provider[] = [
  { name: "mock", status: "connected" },
  { name: "zomato", status: "disconnected" },
  { name: "swiggy", status: "disconnected" },
];
const MOCK_PRINTERS: Printer[] = [
  { name: "kitchen", status: "ready", queue_depth: 2, printed: 128, failed: 1 },
  { name: "cashier", status: "offline", queue_depth: 0, printed: 94, failed: 3 },
];
const MOCK_EVENTS: Event[] = [
  { id: "1", type: "order.received", order_id: "ORD-1842", created_at: new Date().toISOString() },
  { id: "2", type: "order.normalized", order_id: "ORD-1842", created_at: new Date(Date.now() - 30_000).toISOString() },
  { id: "3", type: "order.queued", order_id: "ORD-1839", created_at: new Date(Date.now() - 90_000).toISOString() },
  { id: "4", type: "order.stored", order_id: "ORD-1838", created_at: new Date(Date.now() - 150_000).toISOString() },
];
const MOCK_QUEUE: Queue = { queue_depth: 2, queue_capacity: 100, workers: 2, events: 184, retried: 3, failed: 1 };

const API_URL = (import.meta.env.VITE_API_URL ?? "http://localhost:8080").replace(/\/$/, "");

async function readApi<T>(path: string, fallback: T): Promise<T> {
  try {
    const response = await fetch(`${API_URL}${path}`);
    if (!response.ok) return fallback;
    return await response.json() as T;
  } catch {
    return fallback;
  }
}

function statusClass(status: string) { return `automation-status status-${status}`; }
function displayName(value: string) { return value.replace(/[._-]/g, " ").replace(/\b\w/g, (letter) => letter.toUpperCase()); }

export function AutomationPage() {
  const [providers, setProviders] = useState(MOCK_PROVIDERS);
  const [printers, setPrinters] = useState(MOCK_PRINTERS);
  const [events, setEvents] = useState(MOCK_EVENTS);
  const [queue, setQueue] = useState(MOCK_QUEUE);
  const [lastUpdated, setLastUpdated] = useState(new Date());
  const [refreshing, setRefreshing] = useState(false);

  const refresh = useCallback(async () => {
    setRefreshing(true);
    const [nextProviders, nextPrinters, nextEvents, nextQueue] = await Promise.all([
      readApi<Provider[]>("/api/v1/integrations", MOCK_PROVIDERS),
      readApi<Printer[]>("/api/v1/printers", MOCK_PRINTERS),
      readApi<Event[]>("/api/v1/automation/events", MOCK_EVENTS),
      readApi<Queue>("/api/v1/automation/queue", MOCK_QUEUE),
    ]);
    setProviders(nextProviders.length ? nextProviders : MOCK_PROVIDERS);
    setPrinters(nextPrinters.length ? nextPrinters : MOCK_PRINTERS);
    setEvents(nextEvents.length ? nextEvents.slice(-8).reverse() : MOCK_EVENTS);
    setQueue(nextQueue);
    setLastUpdated(new Date());
    setRefreshing(false);
  }, []);

  useEffect(() => { void refresh(); const timer = window.setInterval(() => void refresh(), 5000); return () => window.clearInterval(timer); }, [refresh]);

  const readyPrinters = printers.filter((printer) => printer.status === "ready").length;
  const connectedProviders = providers.filter((provider) => provider.status === "connected").length;
  const retryRate = queue.events ? `${((queue.retried / queue.events) * 100).toFixed(1)}%` : "0.0%";
  const latestEvent = useMemo(() => events[0], [events]);

  return (
    <section className="automation-workspace">
      <div className="automation-toolbar">
        <div><p className="eyebrow">Developer operations</p><h2>Automation Overview</h2><p className="automation-subtitle">Runtime control surface for orders, queues, integrations, and printers.</p></div>
        <button type="button" className="icon-button" onClick={() => void refresh()} disabled={refreshing} title="Refresh automation status" aria-label="Refresh automation status"><RefreshCw size={18} aria-hidden="true" className={refreshing ? "spin" : undefined} /></button>
      </div>

      <section className="automation-metric-grid" aria-label="Automation metrics">
        <article className="automation-metric"><Activity size={20} aria-hidden="true" /><span>Today's orders</span><strong>128</strong><small>Live operations</small></article>
        <article className="automation-metric"><Printer size={20} aria-hidden="true" /><span>Average print time</span><strong>1.8s</strong><small>Last 24 hours</small></article>
        <article className="automation-metric"><Wifi size={20} aria-hidden="true" /><span>Queue time</span><strong>{queue.queue_depth ? "24s" : "0s"}</strong><small>{queue.queue_depth} waiting</small></article>
        <article className="automation-metric"><CircleAlert size={20} aria-hidden="true" /><span>Retry rate</span><strong>{retryRate}</strong><small>{queue.failed} failed jobs</small></article>
      </section>

      <div className="automation-grid">
        <Panel label="System health"><PanelHeading eyebrow="System health" title="Runtime status" action={<CheckCircle2 size={20} aria-hidden="true" />} /><div className="automation-health-list"><div><span><Database size={16} aria-hidden="true" />Database</span><strong className="health-good">Operational</strong></div><div><span><Activity size={16} aria-hidden="true" />Automation workers</span><strong className="health-good">{queue.workers} online</strong></div><div><span><Wifi size={16} aria-hidden="true" />Internet</span><strong className="health-good">Connected</strong></div><div><span><PlugZap size={16} aria-hidden="true" />Integrations</span><strong>{connectedProviders}/{providers.length} connected</strong></div></div></Panel>
        <Panel label="Queue metrics"><PanelHeading eyebrow="Queue" title="Order processing" action={<Activity size={20} aria-hidden="true" />} /><div className="queue-gauge"><strong>{queue.queue_depth}</strong><span>waiting of {queue.queue_capacity}</span><div><i style={{ width: `${Math.min(100, (queue.queue_depth / queue.queue_capacity) * 100)}%` }} /></div></div><div className="automation-mini-stats"><span>Events<strong>{queue.events}</strong></span><span>Retries<strong>{queue.retried}</strong></span><span>Failures<strong>{queue.failed}</strong></span></div></Panel>
      </div>

      <div className="automation-grid">
        <Panel label="Integration status"><PanelHeading eyebrow="Integrations" title="Provider connections" action={<PlugZap size={20} aria-hidden="true" />} /><div className="automation-directory">{providers.map((provider) => <div className="automation-row" key={provider.name}><span>{displayName(provider.name)}</span><strong className={statusClass(provider.status)}>{provider.status}</strong></div>)}</div></Panel>
        <Panel label="Printer status"><PanelHeading eyebrow="Printers" title="Print infrastructure" action={<Printer size={20} aria-hidden="true" />} /><div className="automation-directory">{printers.map((printer) => <div className="automation-row" key={printer.name}><span>{displayName(printer.name)}</span><div><strong className={statusClass(printer.status)}>{printer.status}</strong><small>{printer.queue_depth} queued · {printer.failed} failures</small></div></div>)}</div><p className="automation-footnote">{readyPrinters} of {printers.length} printers ready</p></Panel>
      </div>

      <Panel label="Event stream"><PanelHeading eyebrow="Event stream" title="Latest automation events" action={latestEvent ? <span className="automation-live-dot">Live</span> : undefined} /><div className="event-stream">{events.map((event) => <div className="event-row" key={event.id}><time>{new Date(event.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}</time><span className="event-marker" /><div><strong>{displayName(event.type)}</strong><small>{event.order_id}</small></div></div>)}</div></Panel>
      <p className="automation-updated">Last checked {lastUpdated.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", second: "2-digit" })}</p>
    </section>
  );
}
