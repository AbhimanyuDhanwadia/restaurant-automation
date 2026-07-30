import { ClipboardCheck, RefreshCw } from "lucide-react";
import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { DirectoryRow } from "@/components/ui/DirectoryRow";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import { acknowledgeAlert, createAlert, listAlerts, type AlertSeverity } from "@/features/alerts/api";

const severities: AlertSeverity[] = ["high", "medium", "low"];
const severityLabel = (severity: AlertSeverity) => `${severity.slice(0, 1).toUpperCase()}${severity.slice(1)}`;

export function AlertsPage() {
  const [filter, setFilter] = useState<AlertSeverity | "all">("all");
  const [selectedID, setSelectedID] = useState("");
  const [label, setLabel] = useState(""); const [detail, setDetail] = useState(""); const [severity, setSeverity] = useState<AlertSeverity>("medium");
  const queryClient = useQueryClient();
  const alertsQuery = useQuery({ queryKey: ["alerts"], queryFn: listAlerts });
  const refreshAlerts = () => queryClient.invalidateQueries({ queryKey: ["alerts"] });
  const createMutation = useMutation({ mutationFn: createAlert, onSuccess: (alert) => { setSelectedID(alert.id); setLabel(""); setDetail(""); setSeverity("medium"); refreshAlerts(); } });
  const acknowledgeMutation = useMutation({ mutationFn: acknowledgeAlert, onSuccess: refreshAlerts });
  const alerts = useMemo(() => alertsQuery.data ?? [], [alertsQuery.data]);
  const directory = useMemo(() => filter === "all" ? alerts : alerts.filter((alert) => alert.severity === filter), [alerts, filter]);
  const selected = directory.find((alert) => alert.id === selectedID) ?? directory[0];
  const error = alertsQuery.error ?? createMutation.error ?? acknowledgeMutation.error;
  function createSubmit(event: React.FormEvent<HTMLFormElement>) { event.preventDefault(); createMutation.mutate({ label: label.trim(), detail: detail.trim(), severity }); }

  return <section className="alerts-workspace">
    <div className="orders-toolbar"><div><p className="eyebrow">Attention queue</p><h2>Alert Center</h2></div><label className="filter-control"><span>Severity</span><select value={filter} onChange={(event) => setFilter(event.target.value as AlertSeverity | "all")}><option value="all">All alerts</option>{severities.map((value) => <option key={value} value={value}>{severityLabel(value)}</option>)}</select></label></div>
    <Panel label="Create manual alert"><form className="alert-create-form" onSubmit={createSubmit}><label><span>Alert</span><input value={label} onChange={(event) => setLabel(event.target.value)} required /></label><label><span>Detail</span><input value={detail} onChange={(event) => setDetail(event.target.value)} required /></label><label><span>Severity</span><select value={severity} onChange={(event) => setSeverity(event.target.value as AlertSeverity)}>{severities.map((value) => <option key={value} value={value}>{severityLabel(value)}</option>)}</select></label><button type="submit" className="primary-button" disabled={createMutation.isPending}>{createMutation.isPending ? "Creating..." : "Create alert"}</button></form></Panel>
    {alertsQuery.isError && <div className="orders-api-error" role="alert"><span>{alertsQuery.error.message}</span><button type="button" className="text-button" onClick={() => alertsQuery.refetch()}><RefreshCw size={15} aria-hidden="true" /> Retry</button></div>}{error && !alertsQuery.isError && <p className="orders-api-error" role="alert">{error.message}</p>}
    <div className="alerts-directory"><Panel className="alert-directory-list" label="Active alerts">{alertsQuery.isPending && <EmptyState message="Loading alerts..." />}{directory.map((alert) => <DirectoryRow key={alert.id} selected={selected?.id === alert.id} onClick={() => setSelectedID(alert.id)} primary={alert.label} secondary={alert.detail} badge={<StatusPill variant={`severity-${alert.severity}`}>{severityLabel(alert.severity)}</StatusPill>} label={`Select ${alert.label} alert`} />)}{!alertsQuery.isPending && !alertsQuery.isError && directory.length === 0 && <EmptyState message="No alerts match this severity." />}</Panel>{selected && <Panel className="order-detail" label={`Details for ${selected.label}`}><PanelHeading eyebrow="Selected alert" title={selected.label} action={<StatusPill variant={`severity-${selected.severity}`}>{severityLabel(selected.severity)}</StatusPill>} /><p className="alert-detail-copy">{selected.detail}</p><button type="button" className="back-button" disabled={acknowledgeMutation.isPending} onClick={() => acknowledgeMutation.mutate(selected.id)}><ClipboardCheck size={16} aria-hidden="true" />{acknowledgeMutation.isPending ? "Acknowledging..." : "Acknowledge alert"}</button></Panel>}</div>
  </section>;
}
