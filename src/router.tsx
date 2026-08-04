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
 *   /operations   → OperationsPage
 *   /orders      → OrdersPage
 *   /kitchen     → KitchenPage
 *   /tables      → TablesPage
 *   /inventory   → InventoryPage
 *   /staff       → StaffPage
 *   /alerts      → AlertsPage
 *   /analytics   → AnalyticsPage
 *   /settings    → SettingsPage
 *
 * Future phases will add:
 *   /printers    → PrintersPage
 *   /analytics   → AnalyticsPage
 *   /settings    → SettingsPage
 */

import {
  createRootRoute,
  createRoute,
  createRouter,
  lazyRouteComponent,
} from "@tanstack/react-router";
import { AuthGuard } from "@/features/auth/AuthGuard";
import { AppShell } from "@/components/layout/AppShell";

const OperationsPage = lazyRouteComponent(() => import("@/pages/OperationsPage"), "OperationsPage");
const OrdersPage = lazyRouteComponent(() => import("@/pages/OrdersPage"), "OrdersPage");
const TablesPage = lazyRouteComponent(() => import("@/pages/TablesPage"), "TablesPage");
const InventoryPage = lazyRouteComponent(() => import("@/pages/InventoryPage"), "InventoryPage");
const StaffPage = lazyRouteComponent(() => import("@/pages/StaffPage"), "StaffPage");
const AlertsPage = lazyRouteComponent(() => import("@/pages/AlertsPage"), "AlertsPage");
const KitchenPage = lazyRouteComponent(() => import("@/pages/KitchenPage"), "KitchenPage");
const AnalyticsPage = lazyRouteComponent(() => import("@/pages/AnalyticsPage"), "AnalyticsPage");
const SettingsPage = lazyRouteComponent(() => import("@/pages/SettingsPage"), "SettingsPage");
const AutomationPage = lazyRouteComponent(() => import("@/pages/AutomationPage"), "AutomationPage");
const IntegrationsPage = lazyRouteComponent(() => import("@/pages/IntegrationsPage"), "IntegrationsPage");
const PrintersPage = lazyRouteComponent(() => import("@/pages/PrintersPage"), "PrintersPage");
const QueueMonitorPage = lazyRouteComponent(() => import("@/pages/QueueMonitorPage"), "QueueMonitorPage");
const EventLogsPage = lazyRouteComponent(() => import("@/pages/EventLogsPage"), "EventLogsPage");
const SystemHealthPage = lazyRouteComponent(() => import("@/pages/SystemHealthPage"), "SystemHealthPage");
const AuditLogsPage = lazyRouteComponent(() => import("@/pages/AuditLogsPage"), "AuditLogsPage");
const DatabasePage = lazyRouteComponent(() => import("@/pages/DatabasePage"), "DatabasePage");
const PrintQueuePage = lazyRouteComponent(() => import("@/pages/PrintQueuePage"), "PrintQueuePage");
const OrderEnginePage = lazyRouteComponent(() => import("@/pages/OrderEnginePage"), "OrderEnginePage");
const KitchenPerformancePage = lazyRouteComponent(() => import("@/pages/KitchenPerformancePage"), "KitchenPerformancePage");
const InventoryAnalyticsPage = lazyRouteComponent(() => import("@/pages/InventoryAnalyticsPage"), "InventoryAnalyticsPage");
const SalesPage = lazyRouteComponent(() => import("@/pages/SalesPage"), "SalesPage");

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

const operationsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/operations",
  component: OperationsPage,
});

const kitchenRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/kitchen",
  component: KitchenPage,
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

const analyticsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/analytics",
  component: AnalyticsPage,
});

const settingsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/settings",
  component: SettingsPage,
});

const automationRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/automation",
  component: AutomationPage,
});

const integrationsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/automation/integrations",
  component: IntegrationsPage,
});

const printersRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/automation/printers",
  component: PrintersPage,
});

const queueMonitorRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/automation/queue-monitor",
  component: QueueMonitorPage,
});

const eventLogsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/automation/event-logs",
  component: EventLogsPage,
});

const systemHealthRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/automation/system-health",
  component: SystemHealthPage,
});

const auditLogsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/admin/audit-logs",
  component: AuditLogsPage,
});

const databaseRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/admin/database",
  component: DatabasePage,
});

const printQueueRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/automation/print-queue",
  component: PrintQueuePage,
});

const orderEngineRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/automation/order-engine",
  component: OrderEnginePage,
});

const kitchenPerformanceRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/analytics/kitchen",
  component: KitchenPerformancePage,
});

const inventoryAnalyticsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/analytics/inventory",
  component: InventoryAnalyticsPage,
});

const salesRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/analytics/sales",
  component: SalesPage,
});

// ---------------------------------------------------------------------------
// Router instance
// ---------------------------------------------------------------------------

const routeTree = rootRoute.addChildren([
  indexRoute,
  operationsRoute,
  ordersRoute,
  kitchenRoute,
  tablesRoute,
  inventoryRoute,
  staffRoute,
  alertsRoute,
  analyticsRoute,
  settingsRoute,
  automationRoute,
  integrationsRoute,
  printersRoute,
  queueMonitorRoute,
  eventLogsRoute,
  systemHealthRoute,
  auditLogsRoute,
  databaseRoute,
  printQueueRoute,
  orderEngineRoute,
  kitchenPerformanceRoute,
  inventoryAnalyticsRoute,
  salesRoute,
]);

const basepath = import.meta.env.BASE_URL.replace(/\/$/, "") || "/";

export const router = createRouter({ routeTree, basepath });

// TanStack Router type registration — enables full TypeScript inference on
// useNavigate, Link, useParams, etc.
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
