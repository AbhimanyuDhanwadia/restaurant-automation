/**
 * AppShell
 *
 * The authenticated application layout: sidebar navigation on the left,
 * scrollable workspace on the right. Passes the active session down to
 * the Topbar for sign-out and user info display.
 *
 * The `<Outlet />` from TanStack Router renders the active page component
 * inside the workspace column, replacing the old if/else view switch in App.tsx.
 * Search state is kept in the UIStore so both Topbar and OperationsPage can
 * read it without prop-drilling.
 */

import { Outlet, useLocation } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import type { Session } from "@supabase/supabase-js";
import { Sidebar } from "@/components/layout/Sidebar";
import { Topbar } from "@/components/layout/Topbar";
import { getCurrentUserAccess } from "@/features/auth/api";

const PAGE_META: Record<string, { eyebrow: string; title: string }> = {
  "/": { eyebrow: "Dinner service", title: "Live Operations" },
  "/operations": { eyebrow: "Operations", title: "Live Operations" },
  "/orders": { eyebrow: "Order intake", title: "Orders" },
  "/kitchen": { eyebrow: "Production floor", title: "Kitchen" },
  "/tables": { eyebrow: "Floor plan", title: "Tables" },
  "/inventory": { eyebrow: "Stock control", title: "Inventory" },
  "/staff": { eyebrow: "Shift handoff", title: "Staff" },
  "/alerts": { eyebrow: "Attention queue", title: "Alerts" },
  "/analytics": { eyebrow: "Performance view", title: "Analytics" },
  "/analytics/sales": { eyebrow: "Performance view", title: "Sales" },
  "/analytics/kitchen": { eyebrow: "Performance view", title: "Kitchen Performance" },
  "/analytics/delivery": { eyebrow: "Performance view", title: "Delivery Performance" },
  "/analytics/inventory": { eyebrow: "Performance view", title: "Inventory Analytics" },
  "/settings": { eyebrow: "Administration", title: "Settings" },
  "/admin/users": { eyebrow: "Administration", title: "Users" },
  "/admin/roles": { eyebrow: "Administration", title: "Roles" },
  "/admin/backups": { eyebrow: "Administration", title: "Backups" },
  "/automation": { eyebrow: "Developer operations", title: "Automation Overview" },
  "/automation/order-engine": { eyebrow: "Developer operations", title: "Order Engine" },
};

interface AppShellProps {
  session: Session;
}

export function AppShell({ session }: AppShellProps) {
  const location = useLocation();
  const meta = PAGE_META[location.pathname] ?? { eyebrow: "Restaurant Automation", title: "Dashboard" };
  const accessQuery = useQuery({ queryKey: ["auth", "current-user"], queryFn: getCurrentUserAccess, retry: false, staleTime: 60_000 });

  return (
    <main className="app-shell">
      <Sidebar permissions={accessQuery.data?.permissions} />
      <section className="workspace">
        <Topbar
          eyebrow={meta.eyebrow}
          title={meta.title}
          session={session}
        />
        <Outlet />
      </section>
    </main>
  );
}
