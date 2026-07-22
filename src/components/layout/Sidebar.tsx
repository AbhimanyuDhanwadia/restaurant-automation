/**
 * Sidebar
 *
 * The primary navigation sidebar. Uses TanStack Router's `<Link>` for
 * navigation so the active route is reflected correctly in the URL and
 * browser history. Active state is determined by the current pathname.
 */

import {
  Bell,
  CalendarDays,
  ChefHat,
  ClipboardCheck,
  PackageSearch,
  ReceiptText,
  UsersRound,
} from "lucide-react";
import { Link } from "@tanstack/react-router";

const NAV_ITEMS = [
  { to: "/", label: "Operations", icon: ClipboardCheck },
  { to: "/orders", label: "Orders", icon: ReceiptText },
  { to: "/tables", label: "Tables", icon: CalendarDays },
  { to: "/inventory", label: "Inventory", icon: PackageSearch },
  { to: "/staff", label: "Staff", icon: UsersRound },
  { to: "/alerts", label: "Alerts", icon: Bell },
] as const;

export function Sidebar() {
  return (
    <aside className="sidebar" aria-label="Primary navigation">
      <div className="brand">
        <ChefHat size={26} aria-hidden="true" />
        <span>Restaurant Automation</span>
      </div>
      <nav>
        {NAV_ITEMS.map(({ to, label, icon: Icon }) => (
          <Link
            key={to}
            to={to}
            activeProps={{ className: "active" }}
            activeOptions={to === "/" ? { exact: true } : undefined}
          >
            <Icon size={18} aria-hidden="true" />
            {label}
          </Link>
        ))}
      </nav>
    </aside>
  );
}
