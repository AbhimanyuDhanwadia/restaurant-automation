/**
 * OperationsPage — /
 *
 * The live operations dashboard: summary stat cards, active order queue,
 * alert centre, and shift task list.
 */

import { AlertTriangle, Clock3, X } from "lucide-react";
import { useMemo } from "react";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { EmptyState } from "@/components/ui/EmptyState";
import { StatusPill } from "@/components/ui/StatusPill";
import { useOrdersStore } from "@/stores/orders";
import { useAlertsStore } from "@/stores/alerts";
import { useStaffStore } from "@/stores/staff";
import { useUIStore } from "@/stores/ui";
import { SEED_STATS } from "@/data/seeds";
import { useNavigate } from "@tanstack/react-router";

export function OperationsPage() {
  const orders = useOrdersStore((s) => s.orders);
  const cycleOrderStatus = useOrdersStore((s) => s.cycleStatus);
  const alerts = useAlertsStore((s) => s.alerts);
  const acknowledgeAlert = useAlertsStore((s) => s.acknowledge);
  const tasks = useStaffStore((s) => s.tasks);
  const completeTask = useStaffStore((s) => s.completeTask);
  const searchTerm = useUIStore((s) => s.searchTerm);
  const navigate = useNavigate();

  const filteredOrders = useMemo(() => {
    const query = searchTerm.trim().toLowerCase();
    if (!query) return orders;
    return orders.filter((o) =>
      [o.id, o.table, o.channel, o.items, o.status].join(" ").toLowerCase().includes(query),
    );
  }, [orders, searchTerm]);

  return (
    <>
      {/* Stat cards */}
      <section className="stats-grid" aria-label="Operational summary">
        {SEED_STATS.map((stat) => (
          <article className="stat-card" key={stat.label}>
            <stat.icon size={22} aria-hidden="true" />
            <div>
              <span>{stat.label}</span>
              <strong>{stat.value}</strong>
              <small>{stat.detail}</small>
            </div>
          </article>
        ))}
      </section>

      {/* Main grid: orders + sidebar */}
      <section className="content-grid">
        <Panel label="Active Orders" id="orders">
          <PanelHeading
            eyebrow="Kitchen queue"
            title="Active Orders"
            action={
              <button
                type="button"
                className="text-button"
                onClick={() => navigate({ to: "/orders" })}
              >
                View all
              </button>
            }
          />
          <div className="order-list">
            {filteredOrders.map((order) => (
              <article className="order-row" key={order.id}>
                <div>
                  <strong>{order.id}</strong>
                  <span>{order.table} · {order.channel}</span>
                </div>
                <span>{order.items}</span>
                <span className="eta">
                  <Clock3 size={16} aria-hidden="true" />
                  {order.eta}
                </span>
                <StatusPill
                  variant={order.status.toLowerCase()}
                  onClick={() => cycleOrderStatus(order.id)}
                >
                  {order.status}
                </StatusPill>
              </article>
            ))}
            {filteredOrders.length === 0 && (
              <EmptyState message={`No orders match "${searchTerm}".`} />
            )}
          </div>
        </Panel>

        <aside className="side-stack">
          <Panel label="Alerts" id="alerts">
            <PanelHeading
              eyebrow="Attention"
              title="Alerts"
              action={<AlertTriangle size={20} aria-hidden="true" />}
            />
            <div className="alert-list">
              {alerts.map((alert) => (
                <article className="alert-row" key={alert.label}>
                  <div>
                    <strong>{alert.label}</strong>
                    <span>{alert.detail}</span>
                  </div>
                  <small>{alert.severity}</small>
                  <button
                    type="button"
                    className="dismiss-button"
                    aria-label={`Dismiss ${alert.label} alert`}
                    onClick={() => acknowledgeAlert(alert.label)}
                  >
                    <X size={16} aria-hidden="true" />
                  </button>
                </article>
              ))}
              {alerts.length === 0 && <EmptyState message="All alerts are cleared." />}
            </div>
          </Panel>

          <Panel>
            <PanelHeading eyebrow="Automation" title="Next Tasks" />
            <div className="task-list">
              {tasks.map((task) => (
                <article className="task-row" key={task.title}>
                  <div>
                    <strong>{task.title}</strong>
                    <span>{task.owner}</span>
                  </div>
                  <time>{task.due}</time>
                  <button
                    type="button"
                    className="task-complete"
                    onClick={() => completeTask(task.title)}
                  >
                    Done
                  </button>
                </article>
              ))}
              {tasks.length === 0 && <EmptyState message="No pending tasks." />}
            </div>
          </Panel>
        </aside>
      </section>
    </>
  );
}
