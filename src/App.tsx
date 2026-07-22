import {
  AlertTriangle,
  ArrowLeft,
  Bell,
  CalendarDays,
  ChefHat,
  Clock3,
  ClipboardCheck,
  ChevronRight,
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
type TableStatus = "Available" | "Seated" | "Reserved" | "Needs Check";
type InventoryStatus = "In Stock" | "Low Stock" | "On Order";
type StaffStatus = "On shift" | "On break" | "Off shift";
type View = "operations" | "orders" | "tables" | "inventory" | "staff";

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

const initialTables: Array<{
  id: string;
  seats: number;
  guest: string;
  reservation: string;
  status: TableStatus;
}> = [
  { id: "T1", seats: 2, guest: "Available", reservation: "Walk-in", status: "Available" },
  { id: "T2", seats: 4, guest: "The Mehta party", reservation: "6:00 PM", status: "Seated" },
  { id: "T3", seats: 2, guest: "A. Kapoor", reservation: "6:15 PM", status: "Needs Check" },
  { id: "T4", seats: 6, guest: "The Shah party", reservation: "6:30 PM", status: "Reserved" },
  { id: "T5", seats: 4, guest: "Available", reservation: "Walk-in", status: "Available" },
  { id: "T6", seats: 2, guest: "R. Iyer", reservation: "5:45 PM", status: "Seated" },
  { id: "T7", seats: 8, guest: "Available", reservation: "Walk-in", status: "Available" },
  { id: "T8", seats: 4, guest: "N. Rao", reservation: "5:50 PM", status: "Needs Check" },
];

const initialInventory: Array<{
  id: string;
  item: string;
  category: string;
  onHand: string;
  par: string;
  supplier: string;
  status: InventoryStatus;
}> = [
  { id: "INV-01", item: "Romaine lettuce", category: "Produce", onHand: "2 cases", par: "6 cases", supplier: "Greenline Farms", status: "Low Stock" },
  { id: "INV-02", item: "Chicken breast", category: "Protein", onHand: "18 kg", par: "24 kg", supplier: "Metro Provisions", status: "Low Stock" },
  { id: "INV-03", item: "Sparkling water", category: "Beverage", onHand: "9 cases", par: "8 cases", supplier: "Beverage Co.", status: "In Stock" },
  { id: "INV-04", item: "Wild-caught salmon", category: "Protein", onHand: "12 kg", par: "18 kg", supplier: "Ocean Table", status: "On Order" },
  { id: "INV-05", item: "Sourdough loaves", category: "Bakery", onHand: "14 loaves", par: "12 loaves", supplier: "Daily Crumb", status: "In Stock" },
];

const initialStaff: Array<{
  id: string;
  name: string;
  role: string;
  station: string;
  status: StaffStatus;
  handoff: string;
}> = [
  { id: "STAFF-01", name: "Maya Patel", role: "Manager", station: "Front of house", status: "On shift", handoff: "Confirm the delivery partner SLA before 5:00 PM." },
  { id: "STAFF-02", name: "Arjun Shah", role: "Sous chef", station: "Hot line", status: "On shift", handoff: "Approve the prep list and watch grill capacity." },
  { id: "STAFF-03", name: "Nisha Rao", role: "Bar lead", station: "Bar", status: "On break", handoff: "Restock bar garnishes before the dinner rush." },
  { id: "STAFF-04", name: "Kabir Mehta", role: "Server", station: "Section B", status: "On shift", handoff: "Check in on Table 8 and close the open check." },
];

function App() {
  const [orderState, setOrderState] = useState(orders);
  const [searchTerm, setSearchTerm] = useState("");
  const [visibleAlerts, setVisibleAlerts] = useState(initialAlerts);
  const [visibleTasks, setVisibleTasks] = useState(initialTasks);
  const [activeView, setActiveView] = useState<View>("operations");
  const [statusFilter, setStatusFilter] = useState<OrderStatus | "All">("All");
  const [selectedOrderId, setSelectedOrderId] = useState(orders[0].id);
  const [tableState, setTableState] = useState(initialTables);
  const [tableFilter, setTableFilter] = useState<TableStatus | "All">("All");
  const [selectedTableId, setSelectedTableId] = useState(initialTables[0].id);
  const [inventoryState, setInventoryState] = useState(initialInventory);
  const [inventoryFilter, setInventoryFilter] = useState<InventoryStatus | "All">("All");
  const [selectedInventoryId, setSelectedInventoryId] = useState(initialInventory[0].id);
  const [staffState] = useState(initialStaff);
  const [staffFilter, setStaffFilter] = useState<StaffStatus | "All">("All");
  const [selectedStaffId, setSelectedStaffId] = useState(initialStaff[0].id);
  const [handoffNote, setHandoffNote] = useState("");
  const [handoffSaved, setHandoffSaved] = useState(false);

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

  const orderDirectory = useMemo(() => {
    if (statusFilter === "All") {
      return orderState;
    }

    return orderState.filter((order) => order.status === statusFilter);
  }, [orderState, statusFilter]);

  const selectedOrder = orderDirectory.find((order) => order.id === selectedOrderId) ?? orderDirectory[0];

  const tableDirectory = useMemo(() => {
    if (tableFilter === "All") {
      return tableState;
    }

    return tableState.filter((table) => table.status === tableFilter);
  }, [tableFilter, tableState]);

  const selectedTable = tableState.find((table) => table.id === selectedTableId) ?? tableState[0];

  const inventoryDirectory = useMemo(() => {
    if (inventoryFilter === "All") {
      return inventoryState;
    }

    return inventoryState.filter((inventory) => inventory.status === inventoryFilter);
  }, [inventoryFilter, inventoryState]);

  const selectedInventory = inventoryState.find((inventory) => inventory.id === selectedInventoryId) ?? inventoryState[0];

  const staffDirectory = useMemo(() => {
    if (staffFilter === "All") {
      return staffState;
    }

    return staffState.filter((staff) => staff.status === staffFilter);
  }, [staffFilter, staffState]);

  const selectedStaff = staffState.find((staff) => staff.id === selectedStaffId) ?? staffState[0];

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

  const cycleTableStatus = (tableId: string) => {
    setTableState((currentTables) =>
      currentTables.map((table) => {
        if (table.id !== tableId) {
          return table;
        }

        const nextStatus: Record<TableStatus, TableStatus> = {
          Available: "Reserved",
          Reserved: "Seated",
          Seated: "Needs Check",
          "Needs Check": "Available",
        };

        return { ...table, status: nextStatus[table.status] };
      }),
    );
  };

  const receiveInventory = (inventoryId: string) => {
    setInventoryState((currentInventory) =>
      currentInventory.map((inventory) =>
        inventory.id === inventoryId ? { ...inventory, status: "In Stock" } : inventory,
      ),
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
          <button
            className={activeView === "operations" ? "active" : ""}
            onClick={() => setActiveView("operations")}
          >
            <ClipboardCheck size={18} aria-hidden="true" />
            Operations
          </button>
          <button
            className={activeView === "orders" ? "active" : ""}
            onClick={() => setActiveView("orders")}
          >
            <ReceiptText size={18} aria-hidden="true" />
            Orders
          </button>
          <button
            className={activeView === "tables" ? "active" : ""}
            onClick={() => setActiveView("tables")}
          >
            <CalendarDays size={18} aria-hidden="true" />
            Tables
          </button>
          <button
            className={activeView === "inventory" ? "active" : ""}
            onClick={() => setActiveView("inventory")}
          >
            <PackageSearch size={18} aria-hidden="true" />
            Inventory
          </button>
          <button
            className={activeView === "staff" ? "active" : ""}
            onClick={() => setActiveView("staff")}
          >
            <UsersRound size={18} aria-hidden="true" />
            Staff
          </button>
          <button className="disabled-nav" disabled>
            <Bell size={18} aria-hidden="true" />
            Alerts
          </button>
        </nav>
      </aside>

      <section className="workspace">
        <header className="topbar">
          <div>
            <p className="eyebrow">Dinner service</p>
            <h1>{activeView === "operations" ? "Live Operations" : activeView === "orders" ? "Orders" : activeView === "tables" ? "Tables" : activeView === "inventory" ? "Inventory" : "Staff"}</h1>
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

        {activeView === "operations" && <section className="stats-grid" aria-label="Operational summary">
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
        </section>}

        {activeView === "staff" ? (
          <section className="staff-workspace">
            <div className="orders-toolbar">
              <div>
                <p className="eyebrow">Shift handoff</p>
                <h2>Staff & Tasks</h2>
              </div>
              <label className="filter-control">
                <span>Availability</span>
                <select value={staffFilter} onChange={(event) => setStaffFilter(event.target.value as StaffStatus | "All")}>
                  <option value="All">Everyone</option>
                  <option value="On shift">On shift</option>
                  <option value="On break">On break</option>
                  <option value="Off shift">Off shift</option>
                </select>
              </label>
            </div>
            <div className="staff-directory">
              <section className="panel staff-roster" aria-label="Staff roster">
                {staffDirectory.map((staff) => (
                  <button
                    className={`directory-row ${selectedStaff?.id === staff.id ? "selected" : ""}`}
                    key={staff.id}
                    onClick={() => setSelectedStaffId(staff.id)}
                  >
                    <span>
                      <strong>{staff.name}</strong>
                      <small>{staff.role} · {staff.station}</small>
                    </span>
                    <span className={`status staff-status-${staff.status.toLowerCase().replace(" ", "-")}`}>{staff.status}</span>
                    <ChevronRight size={18} aria-hidden="true" />
                  </button>
                ))}
                {staffDirectory.length === 0 && <p className="empty-state">No staff match this availability.</p>}
              </section>
              <div className="staff-side-stack">
                {selectedStaff && (
                  <section className="panel staff-detail" aria-label={`Handoff for ${selectedStaff.name}`}>
                    <div className="panel-heading">
                      <div>
                        <p className="eyebrow">Selected staff member</p>
                        <h2>{selectedStaff.name}</h2>
                      </div>
                      <span className={`status staff-status-${selectedStaff.status.toLowerCase().replace(" ", "-")}`}>{selectedStaff.status}</span>
                    </div>
                    <dl className="detail-list">
                      <div><dt>Role</dt><dd>{selectedStaff.role}</dd></div>
                      <div><dt>Station</dt><dd>{selectedStaff.station}</dd></div>
                    </dl>
                    <p className="eyebrow handoff-label">Handoff note</p>
                    <p className="handoff-copy">{selectedStaff.handoff}</p>
                  </section>
                )}
                <section className="panel handoff-panel" aria-label="Shift handoff note">
                  <div className="panel-heading">
                    <div>
                      <p className="eyebrow">Manager note</p>
                      <h2>Pass to next shift</h2>
                    </div>
                    {handoffSaved && <span className="saved-label">Saved</span>}
                  </div>
                  <textarea aria-label="Shift handoff note" placeholder="Add a note for the next shift" value={handoffNote} onChange={(event) => { setHandoffNote(event.target.value); setHandoffSaved(false); }} />
                  <button className="back-button" onClick={() => setHandoffSaved(true)} disabled={!handoffNote.trim()}>
                    <ClipboardCheck size={16} aria-hidden="true" />
                    Save handoff
                  </button>
                </section>
              </div>
            </div>
            <section className="panel staff-task-panel" aria-label="Open shift tasks">
              <div className="panel-heading">
                <div>
                  <p className="eyebrow">Automation queue</p>
                  <h2>Open shift tasks</h2>
                </div>
                <span className="task-count">{visibleTasks.length} open</span>
              </div>
              <div className="task-list">
                {visibleTasks.map((task) => (
                  <article className="task-row" key={task.title}>
                    <div>
                      <strong>{task.title}</strong>
                      <span>{task.owner}</span>
                    </div>
                    <time>{task.due}</time>
                    <button className="task-complete" onClick={() => setVisibleTasks((current) => current.filter((item) => item.title !== task.title))}>Done</button>
                  </article>
                ))}
                {visibleTasks.length === 0 && <p className="empty-state">All shift tasks are complete.</p>}
              </div>
            </section>
          </section>
        ) : activeView === "inventory" ? (
          <section className="inventory-workspace">
            <div className="orders-toolbar">
              <div>
                <p className="eyebrow">Stock control</p>
                <h2>Inventory & Purchasing</h2>
              </div>
              <label className="filter-control">
                <span>Status</span>
                <select value={inventoryFilter} onChange={(event) => setInventoryFilter(event.target.value as InventoryStatus | "All")}>
                  <option value="All">All items</option>
                  <option value="In Stock">In stock</option>
                  <option value="Low Stock">Low stock</option>
                  <option value="On Order">On order</option>
                </select>
              </label>
            </div>
            <div className="inventory-directory">
              <section className="panel inventory-list" aria-label="Inventory items">
                {inventoryDirectory.map((inventory) => (
                  <button
                    className={`directory-row ${selectedInventory?.id === inventory.id ? "selected" : ""}`}
                    key={inventory.id}
                    onClick={() => setSelectedInventoryId(inventory.id)}
                  >
                    <span>
                      <strong>{inventory.item}</strong>
                      <small>{inventory.category} · {inventory.onHand} on hand</small>
                    </span>
                    <span className={`status inventory-status-${inventory.status.toLowerCase().replace(" ", "-")}`}>{inventory.status}</span>
                    <ChevronRight size={18} aria-hidden="true" />
                  </button>
                ))}
                {inventoryDirectory.length === 0 && <p className="empty-state">No inventory matches this status.</p>}
              </section>
              {selectedInventory && (
                <section className="panel order-detail" aria-label={`Details for ${selectedInventory.item}`}>
                  <div className="panel-heading">
                    <div>
                      <p className="eyebrow">Selected item</p>
                      <h2>{selectedInventory.item}</h2>
                    </div>
                    <span className={`status inventory-status-${selectedInventory.status.toLowerCase().replace(" ", "-")}`}>{selectedInventory.status}</span>
                  </div>
                  <dl className="detail-list">
                    <div><dt>On hand</dt><dd>{selectedInventory.onHand}</dd></div>
                    <div><dt>Par level</dt><dd>{selectedInventory.par}</dd></div>
                    <div><dt>Supplier</dt><dd>{selectedInventory.supplier}</dd></div>
                  </dl>
                  <button className="back-button" onClick={() => receiveInventory(selectedInventory.id)} disabled={selectedInventory.status === "In Stock"}>
                    <ClipboardCheck size={16} aria-hidden="true" />
                    Mark received
                  </button>
                </section>
              )}
            </div>
          </section>
        ) : activeView === "tables" ? (
          <section className="tables-workspace">
            <div className="orders-toolbar">
              <div>
                <p className="eyebrow">Floor plan</p>
                <h2>Table Status</h2>
              </div>
              <label className="filter-control">
                <span>Status</span>
                <select value={tableFilter} onChange={(event) => setTableFilter(event.target.value as TableStatus | "All")}>
                  <option value="All">All tables</option>
                  <option value="Available">Available</option>
                  <option value="Seated">Seated</option>
                  <option value="Reserved">Reserved</option>
                  <option value="Needs Check">Needs check</option>
                </select>
              </label>
            </div>
            <div className="tables-directory">
              <section className="panel table-grid" aria-label="Restaurant table grid">
                {tableDirectory.map((table) => (
                  <button
                    className={`table-tile ${table.status.toLowerCase().replace(" ", "-")} ${selectedTable?.id === table.id ? "selected" : ""}`}
                    key={table.id}
                    onClick={() => setSelectedTableId(table.id)}
                  >
                    <strong>{table.id}</strong>
                    <span>{table.seats} seats</span>
                    <small>{table.status}</small>
                  </button>
                ))}
                {tableDirectory.length === 0 && <p className="empty-state">No tables match this status.</p>}
              </section>
              {selectedTable && (
                <section className="panel order-detail" aria-label={`Details for ${selectedTable.id}`}>
                  <div className="panel-heading">
                    <div>
                      <p className="eyebrow">Selected table</p>
                      <h2>{selectedTable.id}</h2>
                    </div>
                    <button className={`status status-button table-status-${selectedTable.status.toLowerCase().replace(" ", "-")}`} onClick={() => cycleTableStatus(selectedTable.id)}>
                      {selectedTable.status}
                    </button>
                  </div>
                  <dl className="detail-list">
                    <div><dt>Capacity</dt><dd>{selectedTable.seats} guests</dd></div>
                    <div><dt>Guest</dt><dd>{selectedTable.guest}</dd></div>
                    <div><dt>Reservation</dt><dd>{selectedTable.reservation}</dd></div>
                  </dl>
                  <button className="back-button" onClick={() => setActiveView("operations")}>
                    <ArrowLeft size={16} aria-hidden="true" />
                    Back to operations
                  </button>
                </section>
              )}
            </div>
          </section>
        ) : activeView === "orders" ? (
          <section className="orders-workspace">
            <div className="orders-toolbar">
              <div>
                <p className="eyebrow">Order intake</p>
                <h2>All Orders</h2>
              </div>
              <label className="filter-control">
                <span>Status</span>
                <select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value as OrderStatus | "All")}>
                  <option value="All">All statuses</option>
                  <option value="Preparing">Preparing</option>
                  <option value="Ready">Ready</option>
                  <option value="Delayed">Delayed</option>
                </select>
              </label>
            </div>
            <div className="orders-directory">
              <section className="panel order-directory-list">
                {orderDirectory.map((order) => (
                  <button
                    className={`directory-row ${selectedOrder?.id === order.id ? "selected" : ""}`}
                    key={order.id}
                    onClick={() => setSelectedOrderId(order.id)}
                  >
                    <span>
                      <strong>{order.id}</strong>
                      <small>{order.table} · {order.channel}</small>
                    </span>
                    <span className={`status ${order.status.toLowerCase()}`}>{order.status}</span>
                    <ChevronRight size={18} aria-hidden="true" />
                  </button>
                ))}
                {orderDirectory.length === 0 && <p className="empty-state">No orders match this status.</p>}
              </section>
              {selectedOrder && (
                <section className="panel order-detail" aria-label={`Details for ${selectedOrder.id}`}>
                  <div className="panel-heading">
                    <div>
                      <p className="eyebrow">Selected order</p>
                      <h2>{selectedOrder.id}</h2>
                    </div>
                    <button className={`status status-button ${selectedOrder.status.toLowerCase()}`} onClick={() => cycleOrderStatus(selectedOrder.id)}>
                      {selectedOrder.status}
                    </button>
                  </div>
                  <dl className="detail-list">
                    <div><dt>Service</dt><dd>{selectedOrder.table}</dd></div>
                    <div><dt>Channel</dt><dd>{selectedOrder.channel}</dd></div>
                    <div><dt>Items</dt><dd>{selectedOrder.items}</dd></div>
                    <div><dt>Expected</dt><dd>{selectedOrder.eta}</dd></div>
                  </dl>
                  <button className="back-button" onClick={() => setActiveView("operations")}>
                    <ArrowLeft size={16} aria-hidden="true" />
                    Back to operations
                  </button>
                </section>
              )}
            </div>
          </section>
        ) : (
        <section className="content-grid">
          <section className="panel orders-panel" id="orders">
            <div className="panel-heading">
              <div>
                <p className="eyebrow">Kitchen queue</p>
                <h2>Active Orders</h2>
              </div>
              <button className="text-button" onClick={() => setActiveView("orders")}>View all</button>
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
        )}
      </section>
    </main>
  );
}

export { App };
