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
import type { Session } from "@supabase/supabase-js";
import { Sidebar } from "@/components/layout/Sidebar";
import { Topbar } from "@/components/layout/Topbar";

const PAGE_META: Record<string, { eyebrow: string; title: string }> = {
  "/": { eyebrow: "Dinner service", title: "Live Operations" },
  "/orders": { eyebrow: "Order intake", title: "Orders" },
  "/tables": { eyebrow: "Floor plan", title: "Tables" },
  "/inventory": { eyebrow: "Stock control", title: "Inventory" },
  "/staff": { eyebrow: "Shift handoff", title: "Staff" },
  "/alerts": { eyebrow: "Attention queue", title: "Alerts" },
};

interface AppShellProps {
  session: Session;
}

export function AppShell({ session }: AppShellProps) {
  const location = useLocation();
  const meta = PAGE_META[location.pathname] ?? { eyebrow: "Restaurant Automation", title: "Dashboard" };

  return (
    <main className="app-shell">
      <Sidebar />
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
