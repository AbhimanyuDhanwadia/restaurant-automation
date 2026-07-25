import { ChefHat, Clock3, Flame } from "lucide-react";
import { useMemo, useState } from "react";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";

type TicketStatus = "Queued" | "Cooking" | "Ready";
type Station = "Grill" | "Cold line" | "Pastry";

const STATIONS = [
  { name: "Grill", load: "At capacity", active: 8, capacity: 8, tone: "critical" },
  { name: "Cold line", load: "Steady", active: 4, capacity: 6, tone: "healthy" },
  { name: "Pastry", load: "Available", active: 2, capacity: 5, tone: "healthy" },
];

const INITIAL_TICKETS: Array<{
  id: string;
  order: string;
  table: string;
  items: string;
  station: Station;
  elapsed: string;
  status: TicketStatus;
}> = [
  { id: "K-1842", order: "ORD-1842", table: "Table 12", items: "2 mains, 1 starter", station: "Grill", elapsed: "7 min", status: "Cooking" },
  { id: "K-1843", order: "ORD-1843", table: "Delivery", items: "4 entrees", station: "Cold line", elapsed: "3 min", status: "Queued" },
  { id: "K-1844", order: "ORD-1844", table: "Table 3", items: "1 tasting menu", station: "Pastry", elapsed: "18 min", status: "Cooking" },
  { id: "K-1845", order: "ORD-1845", table: "Pickup", items: "3 bowls, 2 drinks", station: "Cold line", elapsed: "11 min", status: "Ready" },
];

const STATUS_CYCLE: Record<TicketStatus, TicketStatus> = {
  Queued: "Cooking",
  Cooking: "Ready",
  Ready: "Queued",
};

export function KitchenPage() {
  const [tickets, setTickets] = useState(INITIAL_TICKETS);
  const [stationFilter, setStationFilter] = useState<Station | "All">("All");

  const directory = useMemo(
    () => stationFilter === "All" ? tickets : tickets.filter((ticket) => ticket.station === stationFilter),
    [stationFilter, tickets],
  );

  const cycleTicket = (ticketId: string) => {
    setTickets((current) => current.map((ticket) => ticket.id === ticketId
      ? { ...ticket, status: STATUS_CYCLE[ticket.status] }
      : ticket));
  };

  return (
    <section className="kitchen-workspace">
      <div className="orders-toolbar">
        <div>
          <p className="eyebrow">Production floor</p>
          <h2>Kitchen Queue</h2>
        </div>
        <label className="filter-control">
          <span>Station</span>
          <select value={stationFilter} onChange={(event) => setStationFilter(event.target.value as Station | "All")}>
            <option value="All">All stations</option>
            <option value="Grill">Grill</option>
            <option value="Cold line">Cold line</option>
            <option value="Pastry">Pastry</option>
          </select>
        </label>
      </div>

      <section className="kitchen-station-grid" aria-label="Kitchen station load">
        {STATIONS.map((station) => (
          <article className="stat-card kitchen-station-card" key={station.name}>
            {station.name === "Grill" ? <Flame size={22} aria-hidden="true" /> : <ChefHat size={22} aria-hidden="true" />}
            <div>
              <span>{station.name}</span>
              <strong>{station.active}/{station.capacity}</strong>
              <small className={`station-load ${station.tone}`}>{station.load}</small>
            </div>
          </article>
        ))}
      </section>

      <Panel label="Kitchen tickets">
        <PanelHeading eyebrow="Active tickets" title={`${directory.length} tickets in queue`} />
        <div className="kitchen-ticket-list">
          {directory.map((ticket) => (
            <article className="kitchen-ticket" key={ticket.id}>
              <div className="kitchen-ticket-main">
                <strong>{ticket.order}</strong>
                <span>{ticket.table} · {ticket.items}</span>
              </div>
              <span className="kitchen-station-label">{ticket.station}</span>
              <span className="eta"><Clock3 size={16} aria-hidden="true" />{ticket.elapsed}</span>
              <StatusPill variant={`ticket-${ticket.status.toLowerCase()}`} onClick={() => cycleTicket(ticket.id)}>
                {ticket.status}
              </StatusPill>
            </article>
          ))}
          {directory.length === 0 && <EmptyState message="No tickets match this station." />}
        </div>
      </Panel>
    </section>
  );
}
