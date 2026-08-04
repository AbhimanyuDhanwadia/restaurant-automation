import {
  Activity,
  Archive,
  Bell,
  CalendarDays,
  ChefHat,
  ClipboardCheck,
  Database,
  FileClock,
  LayoutDashboard,
  ListChecks,
  PackageSearch,
  PlugZap,
  Printer,
  ReceiptText,
  ScrollText,
  Settings,
  ShieldCheck,
  UsersRound,
  BarChart3,
} from "lucide-react";
import { Link } from "@tanstack/react-router";
import type { ComponentType } from "react";

type NavIcon = ComponentType<{ size?: number; "aria-hidden"?: "true" }>;

const NAV_GROUPS: Array<{
  label: string;
  items: Array<{ to: string; label: string; icon: NavIcon; disabled?: boolean }>;
}> = [
  {
    label: "Workspace",
    items: [
      { to: "/", label: "Dashboard", icon: LayoutDashboard },
    ],
  },
  {
    label: "Operations",
    items: [
      { to: "/operations", label: "Live Operations", icon: ClipboardCheck },
      { to: "/orders", label: "Orders", icon: ReceiptText },
      { to: "/kitchen", label: "Kitchen", icon: ChefHat },
      { to: "/tables", label: "Tables", icon: CalendarDays },
      { to: "/inventory", label: "Inventory", icon: PackageSearch },
      { to: "/staff", label: "Staff", icon: UsersRound },
      { to: "/alerts", label: "Alerts", icon: Bell },
    ],
  },
  {
    label: "Automation",
    items: [
      { to: "/automation", label: "Automation Overview", icon: Activity },
      { to: "/automation/order-engine", label: "Order Engine", icon: ListChecks },
      { to: "/automation/integrations", label: "Integrations", icon: PlugZap },
      { to: "/automation/print-queue", label: "Print Queue", icon: Printer },
      { to: "/automation/printers", label: "Printers", icon: Printer },
      { to: "/automation/system-health", label: "System Health", icon: Activity },
      { to: "/automation/event-logs", label: "Event Logs", icon: ScrollText },
      { to: "/automation/queue-monitor", label: "Queue Monitor", icon: ListChecks },
    ],
  },
  {
    label: "Analytics",
    items: [
      { to: "/analytics", label: "Reports", icon: BarChart3 },
      { to: "/analytics/sales", label: "Sales", icon: BarChart3 },
      { to: "/analytics/kitchen", label: "Kitchen Performance", icon: ChefHat },
      { to: "/analytics/delivery", label: "Delivery Performance", icon: ReceiptText, disabled: true },
      { to: "/analytics/inventory", label: "Inventory Analytics", icon: PackageSearch },
    ],
  },
  {
    label: "Administration",
    items: [
      { to: "/settings", label: "Settings", icon: Settings },
      { to: "/admin/users", label: "Users", icon: UsersRound, disabled: true },
      { to: "/admin/roles", label: "Roles", icon: ShieldCheck, disabled: true },
      { to: "/admin/database", label: "Database", icon: Database },
      { to: "/admin/audit-logs", label: "Audit Logs", icon: FileClock },
      { to: "/admin/backups", label: "Backups", icon: Archive, disabled: true },
    ],
  },
];

export function Sidebar() {
  return (
    <aside className="sidebar" aria-label="Primary navigation">
      <div className="brand">
        <ChefHat size={26} aria-hidden="true" />
        <span>Restaurant Automation</span>
      </div>
      <nav>
        {NAV_GROUPS.map((group) => (
          <div className="nav-group" key={group.label}>
            <p className="nav-group-label">{group.label}</p>
            {group.items.map(({ to, label, icon: Icon, disabled }) => disabled ? (
              <button className="disabled-nav" disabled key={to} title="Available in a later phase">
                <Icon size={17} aria-hidden="true" />
                {label}
              </button>
            ) : (
              <Link
                key={to}
                to={to}
                activeProps={{ className: "active" }}
                activeOptions={to === "/" ? { exact: true } : undefined}
              >
                <Icon size={17} aria-hidden="true" />
                {label}
              </Link>
            ))}
          </div>
        ))}
      </nav>
    </aside>
  );
}
