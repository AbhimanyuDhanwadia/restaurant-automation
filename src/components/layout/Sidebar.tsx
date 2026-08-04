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
  items: Array<{ to: string; label: string; icon: NavIcon; requiredPermission?: string }>;
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
      { to: "/automation", label: "Automation Overview", icon: Activity, requiredPermission: "automation.view" },
      { to: "/automation/order-engine", label: "Order Engine", icon: ListChecks, requiredPermission: "automation.view" },
      { to: "/automation/integrations", label: "Integrations", icon: PlugZap, requiredPermission: "integrations.manage" },
      { to: "/automation/print-queue", label: "Print Queue", icon: Printer, requiredPermission: "printers.manage" },
      { to: "/automation/printers", label: "Printers", icon: Printer, requiredPermission: "printers.manage" },
      { to: "/automation/system-health", label: "System Health", icon: Activity, requiredPermission: "automation.view" },
      { to: "/automation/event-logs", label: "Event Logs", icon: ScrollText, requiredPermission: "automation.view" },
      { to: "/automation/queue-monitor", label: "Queue Monitor", icon: ListChecks, requiredPermission: "automation.view" },
    ],
  },
  {
    label: "Analytics",
    items: [
      { to: "/analytics", label: "Reports", icon: BarChart3, requiredPermission: "analytics.view" },
      { to: "/analytics/sales", label: "Sales", icon: BarChart3, requiredPermission: "analytics.view" },
      { to: "/analytics/kitchen", label: "Kitchen Performance", icon: ChefHat, requiredPermission: "analytics.view" },
      { to: "/analytics/delivery", label: "Delivery Performance", icon: ReceiptText, requiredPermission: "analytics.view" },
      { to: "/analytics/inventory", label: "Inventory Analytics", icon: PackageSearch, requiredPermission: "analytics.view" },
    ],
  },
  {
    label: "Administration",
    items: [
      { to: "/settings", label: "Settings", icon: Settings, requiredPermission: "settings.manage" },
      { to: "/admin/users", label: "Users", icon: UsersRound, requiredPermission: "users.manage" },
      { to: "/admin/roles", label: "Roles", icon: ShieldCheck, requiredPermission: "roles.manage" },
      { to: "/admin/database", label: "Database", icon: Database, requiredPermission: "database.view" },
      { to: "/admin/audit-logs", label: "Audit Logs", icon: FileClock, requiredPermission: "audit.view" },
      { to: "/admin/backups", label: "Backups", icon: Archive, requiredPermission: "backups.manage" },
    ],
  },
];

interface SidebarProps {
  permissions?: string[];
}

export function Sidebar({ permissions }: SidebarProps) {
  return (
    <aside className="sidebar" aria-label="Primary navigation">
      <div className="brand">
        <ChefHat size={26} aria-hidden="true" />
        <span>Restaurant Automation</span>
      </div>
      <nav>
        {NAV_GROUPS.map((group) => {
          const items = permissions ? group.items.filter((item) => !item.requiredPermission || permissions.includes(item.requiredPermission)) : group.items;
          if (items.length === 0) return null;
          return (
          <div className="nav-group" key={group.label}>
            <p className="nav-group-label">{group.label}</p>
            {items.map(({ to, label, icon: Icon }) => (
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
          );
        })}
      </nav>
    </aside>
  );
}
