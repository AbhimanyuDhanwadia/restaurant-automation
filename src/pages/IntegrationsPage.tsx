import { CheckCircle2, CircleAlert, PlugZap, RefreshCw, Webhook } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { listIntegrations, type IntegrationStatus } from "@/features/integrations/api";

function formatName(value: string) {
  return value.replace(/[._-]/g, " ").replace(/\b\w/g, (letter) => letter.toUpperCase());
}

function statusCopy(status: IntegrationStatus) {
  if (status === "connected") return "Accepting orders";
  if (status === "error") return "Requires attention";
  return "Not accepting orders";
}

export function IntegrationsPage() {
  const integrationsQuery = useQuery({
    queryKey: ["integrations"],
    queryFn: listIntegrations,
    refetchInterval: 10_000,
  });
  const integrations = integrationsQuery.data ?? [];
  const connected = integrations.filter((integration) => integration.status === "connected").length;
  const apiBaseUrl = (import.meta.env.VITE_API_URL ?? "http://localhost:8080").replace(/\/$/, "");

  return <section className="integrations-workspace">
    <div className="automation-toolbar">
      <div><p className="eyebrow">Automation</p><h2>Integrations</h2><p className="automation-subtitle">Provider health and inbound order collection.</p></div>
      <button type="button" className="icon-button" onClick={() => integrationsQuery.refetch()} disabled={integrationsQuery.isFetching} title="Refresh integrations" aria-label="Refresh integrations"><RefreshCw size={18} aria-hidden="true" className={integrationsQuery.isFetching ? "spin" : undefined} /></button>
    </div>

    {integrationsQuery.isError && <div className="orders-api-error" role="alert"><span>{integrationsQuery.error.message}</span><button type="button" className="text-button" onClick={() => integrationsQuery.refetch()}>Retry</button></div>}

    <section className="automation-metric-grid integration-metrics" aria-label="Integration metrics">
      <article className="automation-metric"><PlugZap size={20} aria-hidden="true" /><span>Registered providers</span><strong>{integrationsQuery.data ? integrations.length : "Unavailable"}</strong><small>Configured at API startup</small></article>
      <article className="automation-metric"><CheckCircle2 size={20} aria-hidden="true" /><span>Connected</span><strong>{integrationsQuery.data ? connected : "Unavailable"}</strong><small>Currently accepting orders</small></article>
      <article className="automation-metric"><Webhook size={20} aria-hidden="true" /><span>Webhook route</span><strong>/webhooks</strong><small>Signed provider intake</small></article>
      <article className="automation-metric"><CircleAlert size={20} aria-hidden="true" /><span>Provider errors</span><strong>{integrationsQuery.data ? integrations.filter((integration) => integration.status === "error").length : "Unavailable"}</strong><small>Require operator review</small></article>
    </section>

    <Panel label="Provider connections"><PanelHeading eyebrow="Provider registry" title="Connection status" action={<PlugZap size={20} aria-hidden="true" />} /><div className="integration-list">{integrationsQuery.isPending && <EmptyState message="Loading registered providers..." />}{integrations.map((integration) => <article className="integration-row" key={integration.name}><div><strong>{formatName(integration.name)}</strong><span>{statusCopy(integration.status)}</span></div><strong className={`automation-status status-${integration.status}`}>{integration.status}</strong></article>)}{!integrationsQuery.isPending && !integrationsQuery.isError && integrations.length === 0 && <EmptyState message="No integration providers are registered." />}</div></Panel>

    <Panel label="Webhook intake"><PanelHeading eyebrow="Inbound orders" title="Signed webhook intake" action={<Webhook size={20} aria-hidden="true" />} /><div className="integration-webhook-list">{integrations.map((integration) => <div key={integration.name}><span>{formatName(integration.name)}</span><code>{apiBaseUrl}/api/v1/webhooks/{integration.name}</code></div>)}{!integrationsQuery.isPending && !integrationsQuery.isError && integrations.length === 0 && <EmptyState message="Webhook routes appear after a provider is configured." />}</div><p className="automation-footnote">Webhook credentials are configured only on the API server and are never exposed in this workspace.</p></Panel>
  </section>;
}
