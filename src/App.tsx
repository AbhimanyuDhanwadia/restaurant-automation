import {
  AlertTriangle,
  Bell,
  ChefHat,
  Clock3,
  ClipboardCheck,
  Flame,
  PackageSearch,
  Plus,
  ReceiptText,
  Search,
  TableProperties,
  UsersRound,
  X,
} from "lucide-react";
import { useMemo, useState } from "react";

type OrderStatus = "Preparing" | "Ready" | "Delayed";

const stats = [
  { label: "Open orders", value: "28", detail: "+6 in 15 min", icon: ReceiptText },
  { label: "Kitchen load", value: "82%", detail: "Grill at capacity", icon: Flame },
  { label: "Tables seated", value: "19/24", detail: "5 turning soon", icon: TableProperties },
  { label: "Staff online", value: "14", detail: "2 tasks overdue", icon: UsersRound },
];

const orders: Array<{
  id: string;
  table: string;
  channel: string;
  items: string;
  eta: string;
  status: OrderStatus;
}> = [
  {
    id: "ORD-1842",
    table: "Table 12",
    channel: "Dine-in",
    items: "2 mains, 1 starter",
    eta: "7 min",
    status: "Preparing",
  },
  {
    id: "ORD-1843",
    table: "Delivery",
    channel: "Aggregator",
    items: "4 entrees",
    eta: "Ready now",
    status: "Ready",
  },
  {
    id: "ORD-1844",
    table: "Table 3",
    channel: "Dine-in",
    items: "1 tasting menu",
    eta: "18 min",
    status: "Delayed",
  },
  {
    id: "ORD-1845",
    table: "Pickup",
    channel: "Web",
    items: "3 bowls, 2 drinks",
    eta: "11 min",
    status: "Preparing",
  },
];

const initialAlerts = [
  { label: "Romaine lettuce", detail: "Below par by 4 cases", severity: "High" },
  { label: "Dish station", detail: "Needs support before dinner rush", severity: "Medium" },
  { label: "Table 8", detail: "Guest has waited 9 min for check", severity: "Low" },
];

const initialTasks = [
  { title: "Approve prep list", owner: "Sous chef", due: "4:30 PM" },
  { title: "Confirm delivery partner SLA", owner: "Manager", due: "5:00 PM" },
  { title: "Restock bar garnishes", owner: "Bar lead", due: "5:15 PM" },
];

function App() {
  const [orderState, setOrderState] = useState(orders);
  const [searchTerm, setSearchTerm] = useState("");
  const [visibleAlerts, setVisibleAlerts] = useState(initialAlerts);
  const [visibleTasks, setVisibleTasks] = useState(initialTasks);

  const filteredOrders = useMemo(() => {
    const query = searchTerm.trim().toLowerCase();

    if (!query) {
      return orderState;
    }

    return orderState.filter((order) =>
      [order.id, order.table, order.channel, order.items, order.status]
        .join(" ")
        .toLowerCase()
        .includes(query),
    );
  }, [orderState, searchTerm]);

  const cycleOrderStatus = (orderId: string) => {
    setOrderState((currentOrders) =>
      currentOrders.map((order) => {
        if (order.id !== orderId) {
          return order;
        }

        const nextStatus: Record<OrderStatus, OrderStatus> = {
          Preparing: "Ready",
          Ready: "Delayed",
          Delayed: "Preparing",
        };

        return { ...order, status: nextStatus[order.status] };
      }),
    );
  };

  return (
    <main className="app-shell">
      <aside className="sidebar" aria-label="Primary navigation">
        <div className="brand">
          <ChefHat size={26} aria-hidden="true" />
          <span>Restaurant Automation</span>
        </div>
        <nav>
          <a className="active" href="#operations">
            <ClipboardCheck size={18} aria-hidden="true" />
            Operations
          </a>
          <a href="#orders">
            <ReceiptText size={18} aria-hidden="true" />
            Orders
          </a>
          <a href="#inventory">
            <PackageSearch size={18} aria-hidden="true" />
            Inventory
          </a>
          <a href="#alerts">
            <Bell size={18} aria-hidden="true" />
            Alerts
          </a>
        </nav>
      </aside>

      <section className="workspace">
        <header className="topbar">
          <div>
            <p className="eyebrow">Dinner service</p>
            <h1>Live Operations</h1>
          </div>
          <div className="topbar-actions">
            <label className="search-box">
              <Search size={17} aria-hidden="true" />
              <input
                aria-label="Search operations"
                placeholder="Search orders, tables, items"
                value={searchTerm}
                onChange={(event) => setSearchTerm(event.target.value)}
              />
            </label>
            <button className="icon-button" aria-label="Create action">
              <Plus size={20} aria-hidden="true" />
            </button>
          </div>
        </header>

        <section className="stats-grid" aria-label="Operational summary">
          {stats.map((stat) => (
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

        <section className="content-grid">
          <section className="panel orders-panel" id="orders">
            <div className="panel-heading">
              <div>
                <p className="eyebrow">Kitchen queue</p>
                <h2>Active Orders</h2>
              </div>
              <button className="text-button">View all</button>
            </div>

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
                  <button
                    className={`status status-button ${order.status.toLowerCase()}`}
                    onClick={() => cycleOrderStatus(order.id)}
                    title={`Advance ${order.id} status`}
                  >
                    {order.status}
                  </button>
                </article>
              ))}
              {filteredOrders.length === 0 && (
                <p className="empty-state">No orders match “{searchTerm}”.</p>
              )}
            </div>
          </section>

          <aside className="side-stack">
            <section className="panel" id="alerts">
              <div className="panel-heading">
                <div>
                  <p className="eyebrow">Attention</p>
                  <h2>Alerts</h2>
                </div>
                <AlertTriangle size={20} aria-hidden="true" />
              </div>
              <div className="alert-list">
                {visibleAlerts.map((alert) => (
                  <article className="alert-row" key={alert.label}>
                    <div>
                      <strong>{alert.label}</strong>
                      <span>{alert.detail}</span>
                    </div>
                    <small>{alert.severity}</small>
                    <button
                      className="dismiss-button"
                      aria-label={`Dismiss ${alert.label} alert`}
                      onClick={() => setVisibleAlerts((current) => current.filter((item) => item.label !== alert.label))}
                    >
                      <X size={16} aria-hidden="true" />
                    </button>
                  </article>
                ))}
                {visibleAlerts.length === 0 && <p className="empty-state">All alerts are cleared.</p>}
              </div>
            </section>

            <section className="panel">
              <div className="panel-heading">
                <div>
                  <p className="eyebrow">Automation</p>
                  <h2>Next Tasks</h2>
                </div>
              </div>
              <div className="task-list">
                {visibleTasks.map((task) => (
                  <article className="task-row" key={task.title}>
                    <div>
                      <strong>{task.title}</strong>
                      <span>{task.owner}</span>
                    </div>
                    <time>{task.due}</time>
                    <button
                      className="task-complete"
                      onClick={() => setVisibleTasks((current) => current.filter((item) => item.title !== task.title))}
                    >
                      Done
                    </button>
                  </article>
                ))}
                {visibleTasks.length === 0 && <p className="empty-state">No pending tasks.</p>}
              </div>
            </section>
          </aside>
        </section>
      </section>
    </main>
  );
}

export { App };
