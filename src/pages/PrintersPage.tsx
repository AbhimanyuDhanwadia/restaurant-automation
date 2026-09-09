import { CheckCircle2, CircleAlert, Printer, RefreshCw, SendHorizontal } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { listPrinters } from "@/features/printers/api";

function formatName(value: string) {
  return value.replace(/[._-]/g, " ").replace(/\b\w/g, (letter) => letter.toUpperCase());
}

export function PrintersPage() {
  const printersQuery = useQuery({ queryKey: ["printers"], queryFn: listPrinters, refetchInterval: 5_000 });
  const printers = printersQuery.data ?? [];
  const ready = printers.filter((printer) => printer.status === "ready").length;
  const queued = printers.reduce((total, printer) => total + printer.queue_depth, 0);
  const printed = printers.reduce((total, printer) => total + printer.printed, 0);
  const failed = printers.reduce((total, printer) => total + printer.failed, 0);
  const utilizationSamples = printers.filter((printer) => printer.utilization_available);
  const utilization = utilizationSamples.length ? utilizationSamples.reduce((total, printer) => total + printer.utilization, 0) / utilizationSamples.length : null;

  return <section className="printers-workspace">
    <div className="automation-toolbar">
      <div><p className="eyebrow">Automation</p><h2>Printers</h2><p className="automation-subtitle">Fleet health, print throughput, and queued tickets.</p></div>
      <button type="button" className="icon-button" onClick={() => printersQuery.refetch()} disabled={printersQuery.isFetching} title="Refresh printer health" aria-label="Refresh printer health"><RefreshCw size={18} aria-hidden="true" className={printersQuery.isFetching ? "spin" : undefined} /></button>
    </div>

    {printersQuery.isError && <div className="orders-api-error" role="alert"><span>{printersQuery.error.message}</span><button type="button" className="text-button" onClick={() => printersQuery.refetch()}>Retry</button></div>}

    <section className="automation-metric-grid printer-metrics" aria-label="Printer metrics">
      <article className="automation-metric"><Printer size={20} aria-hidden="true" /><span>Registered printers</span><strong>{printersQuery.data ? printers.length : "Unavailable"}</strong><small>Configured at API startup</small></article>
      <article className="automation-metric"><CheckCircle2 size={20} aria-hidden="true" /><span>Ready</span><strong>{printersQuery.data ? ready : "Unavailable"}</strong><small>Available for new tickets</small></article>
      <article className="automation-metric"><SendHorizontal size={20} aria-hidden="true" /><span>Queued tickets</span><strong>{printersQuery.data ? queued : "Unavailable"}</strong><small>Across all printer queues</small></article>
      <article className="automation-metric"><CircleAlert size={20} aria-hidden="true" /><span>Print failures</span><strong>{printersQuery.data ? failed : "Unavailable"}</strong><small>Since API startup</small></article>
    </section>

    <Panel label="Printer fleet"><PanelHeading eyebrow="Print infrastructure" title="Fleet status" action={<Printer size={20} aria-hidden="true" />} /><div className="printer-list">{printersQuery.isPending && <EmptyState message="Loading printer fleet..." />}{printers.map((printer) => <article className="printer-row" key={printer.name}><div><strong>{formatName(printer.name)}</strong><span>{printer.queue_depth} queued · {printer.printed} printed · {printer.failed} failed · {printer.utilization_available ? `${printer.utilization.toFixed(1)}% utilized` : "utilization unavailable"}</span></div><strong className={`automation-status status-${printer.status}`}>{printer.status}</strong></article>)}{!printersQuery.isPending && !printersQuery.isError && printers.length === 0 && <EmptyState message="No printers are registered." />}</div><p className="automation-footnote">{printers.length ? `${ready} of ${printers.length} printers ready · ${printed} tickets printed · ${utilization === null ? "utilization unavailable" : `${utilization.toFixed(1)}% fleet utilization`}` : "Printer availability is unavailable until a driver is configured."}</p></Panel>

    <Panel label="Printer configuration"><PanelHeading eyebrow="Configuration" title="Driver management" action={<CircleAlert size={20} aria-hidden="true" />} /><p className="printer-configuration-copy">Printer drivers and routing are configured on the API server. This workspace observes the active fleet without exposing addresses or driver credentials in the browser.</p></Panel>
  </section>;
}
