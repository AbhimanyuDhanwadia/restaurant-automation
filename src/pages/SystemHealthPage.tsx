import { Activity, CheckCircle2, CircleAlert, Database, PlugZap, Printer, RefreshCw } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { getSystemHealth, type SystemHealthComponent } from "@/features/automation/api";

const componentIcons = { database: Database, queue: Activity, integrations: PlugZap, printers: Printer };

function componentLabel(name: SystemHealthComponent["name"]) {
  return name === "database" ? "Database" : name === "queue" ? "Automation queue" : name === "integrations" ? "Integrations" : "Printers";
}

export function SystemHealthPage() {
  const healthQuery = useQuery({ queryKey: ["system", "health"], queryFn: getSystemHealth, refetchInterval: 5_000 });
  const health = healthQuery.data;
  const healthy = health?.components.filter((component) => component.status === "healthy").length ?? 0;

  return <section className="system-health-workspace">
    <div className="automation-toolbar">
      <div><p className="eyebrow">Automation</p><h2>System Health</h2><p className="automation-subtitle">Availability of the operational services behind order automation.</p></div>
      <button type="button" className="icon-button" onClick={() => healthQuery.refetch()} disabled={healthQuery.isFetching} title="Refresh system health" aria-label="Refresh system health"><RefreshCw size={18} aria-hidden="true" className={healthQuery.isFetching ? "spin" : undefined} /></button>
    </div>
    {healthQuery.isError && <div className="orders-api-error" role="alert"><span>{healthQuery.error.message}</span><button type="button" className="text-button" onClick={() => healthQuery.refetch()}>Retry</button></div>}
    <section className="automation-metric-grid system-health-metrics" aria-label="System health summary"><article className="automation-metric"><CheckCircle2 size={20} aria-hidden="true" /><span>Overall status</span><strong>{health?.status ?? "Unavailable"}</strong><small>Current API health snapshot</small></article><article className="automation-metric"><Activity size={20} aria-hidden="true" /><span>Healthy services</span><strong>{health ? `${healthy}/${health.components.length}` : "Unavailable"}</strong><small>Configured services only</small></article><article className="automation-metric"><CircleAlert size={20} aria-hidden="true" /><span>Unavailable</span><strong>{health ? health.components.filter((component) => component.status === "unavailable").length : "Unavailable"}</strong><small>Require operator attention</small></article><article className="automation-metric"><Database size={20} aria-hidden="true" /><span>Not configured</span><strong>{health ? health.components.filter((component) => component.status === "not_configured").length : "Unavailable"}</strong><small>Optional services not enabled</small></article></section>
    <Panel label="Service health"><PanelHeading eyebrow="Components" title="Operational dependencies" action={<Activity size={20} aria-hidden="true" />} /><div className="system-health-list">{healthQuery.isPending && <EmptyState message="Loading system health..." />}{health?.components.map((component) => { const Icon = componentIcons[component.name]; return <article className="system-health-row" key={component.name}><Icon size={20} aria-hidden="true" /><div><strong>{componentLabel(component.name)}</strong><span>{component.detail}</span></div><strong className={`automation-status status-${component.status}`}>{component.status.replace(/_/g, " ")}</strong></article>; })}{!healthQuery.isPending && !healthQuery.isError && !health?.components.length && <EmptyState message="No service health details are available." />}</div></Panel>
  </section>;
}
