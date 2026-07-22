/**
 * Router configuration.
 *
 * Uses TanStack Router with a code-based route tree (not file-based) to keep
 * the router setup explicit and fully type-safe without requiring a Vite plugin
 * to generate route files. Every route is protected by the AuthGuard at the
 * root layout level.
 *
 * Route tree:
 *   /            → OperationsPage
 *   /orders      → OrdersPage
 *   /tables      → TablesPage
 *   /inventory   → InventoryPage
 *   /staff       → StaffPage
 *   /alerts      → AlertsPage
 *
 * In Milestone 5 we will add:
 *   /printers    → PrintersPage
 *   /analytics   → AnalyticsPage
 *   /settings    → SettingsPage
 */

import {
  createRootRoute,
  createRoute,
  createRouter,
} from "@tanstack/react-router";
import { AuthGuard } from "@/features/auth/AuthGuard";
import { AppShell } from "@/components/layout/AppShell";
import { OperationsPage } from "@/pages/OperationsPage";
import { OrdersPage } from "@/pages/OrdersPage";
import { TablesPage } from "@/pages/TablesPage";
import { InventoryPage } from "@/pages/InventoryPage";
import { StaffPage } from "@/pages/StaffPage";
import { AlertsPage } from "@/pages/AlertsPage";

// ---------------------------------------------------------------------------
// Root route — wraps everything in the AuthGuard.
// AppShell renders its own <Outlet /> internally — no children needed here.
// ---------------------------------------------------------------------------

const rootRoute = createRootRoute({
  component: () => (
    <AuthGuard>
      {(session) => <AppShell session={session} />}
    </AuthGuard>
  ),
});

// ---------------------------------------------------------------------------
// Page routes
// ---------------------------------------------------------------------------

// We use a shared context type so pages can receive searchTerm from AppShell.
// TanStack Router's context system is used in Milestone 5 for API context.

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: OperationsPage,
});

const ordersRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/orders",
  component: OrdersPage,
});

const tablesRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/tables",
  component: TablesPage,
});

const inventoryRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/inventory",
  component: InventoryPage,
});

const staffRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/staff",
  component: StaffPage,
});

const alertsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/alerts",
  component: AlertsPage,
});

// ---------------------------------------------------------------------------
// Router instance
// ---------------------------------------------------------------------------

const routeTree = rootRoute.addChildren([
  indexRoute,
  ordersRoute,
  tablesRoute,
  inventoryRoute,
  staffRoute,
  alertsRoute,
]);

export const router = createRouter({ routeTree });

// TanStack Router type registration — enables full TypeScript inference on
// useNavigate, Link, useParams, etc.
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
