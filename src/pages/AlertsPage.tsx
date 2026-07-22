/**
 * AlertsPage — /alerts
 *
 * Alert centre with severity filtering and a detail panel for the selected alert.
 */

import { ClipboardCheck } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { DirectoryRow } from "@/components/ui/DirectoryRow";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import { useAlertsStore } from "@/stores/alerts";
import type { AlertSeverity } from "@/types/domain";

export function AlertsPage() {
  const alerts = useAlertsStore((s) => s.alerts);
  const acknowledge = useAlertsStore((s) => s.acknowledge);
  const [filter, setFilter] = useState<AlertSeverity | "All">("All");
  const [selectedLabel, setSelectedLabel] = useState(alerts[0]?.label ?? "");

  const directory = useMemo(
    () => (filter === "All" ? alerts : alerts.filter((a) => a.severity === filter)),
    [alerts, filter],
  );

  const selected = directory.find((a) => a.label === selectedLabel) ?? directory[0];

  // Keep selectedLabel in sync when the selected alert is acknowledged.
  useEffect(() => {
    if (!alerts.some((a) => a.label === selectedLabel)) {
      setSelectedLabel(alerts[0]?.label ?? "");
    }
  }, [alerts, selectedLabel]);

  return (
    <section className="alerts-workspace">
      <div className="orders-toolbar">
        <div>
          <p className="eyebrow">Attention queue</p>
          <h2>Alert Center</h2>
        </div>
        <label className="filter-control">
          <span>Severity</span>
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value as AlertSeverity | "All")}
          >
            <option value="All">All alerts</option>
            <option value="High">High</option>
            <option value="Medium">Medium</option>
            <option value="Low">Low</option>
          </select>
        </label>
      </div>

      <div className="alerts-directory">
        <Panel className="alert-directory-list" label="Active alerts">
          {directory.map((alert) => (
            <DirectoryRow
              key={alert.label}
              selected={selected?.label === alert.label}
              onClick={() => setSelectedLabel(alert.label)}
              primary={alert.label}
              secondary={alert.detail}
              badge={
                <StatusPill variant={`severity-${alert.severity.toLowerCase()}`}>
                  {alert.severity}
                </StatusPill>
              }
              label={`Select ${alert.label} alert`}
            />
          ))}
          {directory.length === 0 && (
            <EmptyState message="No alerts match this severity." />
          )}
        </Panel>

        {selected && (
          <Panel className="order-detail" label={`Details for ${selected.label}`}>
            <PanelHeading
              eyebrow="Selected alert"
              title={selected.label}
              action={
                <StatusPill variant={`severity-${selected.severity.toLowerCase()}`}>
                  {selected.severity}
                </StatusPill>
              }
            />
            <p className="alert-detail-copy">{selected.detail}</p>
            <button
              type="button"
              className="back-button"
              onClick={() => acknowledge(selected.label)}
            >
              <ClipboardCheck size={16} aria-hidden="true" />
              Acknowledge alert
            </button>
          </Panel>
        )}
      </div>
    </section>
  );
}
