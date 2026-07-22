/**
 * Application entry point.
 *
 * Wires together:
 *  - React 19 root
 *  - TanStack Router (replaces the old view-switch in App.tsx)
 *  - TanStack Query (ready for API calls in Milestone 4)
 *  - AppErrorBoundary (preserved from the original main.tsx)
 *
 * The CSS import pulls in the existing design system (styles.css).
 * Tailwind v4 directives will be added to styles.css alongside the
 * existing class definitions as we migrate component-by-component.
 */

import { Component, StrictMode, type ErrorInfo, type ReactNode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider } from "@tanstack/react-router";
import { QueryClientProvider } from "@tanstack/react-query";
import { router } from "./router";
import { queryClient } from "./lib/queryClient";
import "./styles.css";

// ---------------------------------------------------------------------------
// Error boundary — catches unhandled render errors and shows a recovery UI.
// ---------------------------------------------------------------------------

type ErrorBoundaryState = { hasError: boolean };

class AppErrorBoundary extends Component<{ children: ReactNode }, ErrorBoundaryState> {
  state: ErrorBoundaryState = { hasError: false };

  static getDerivedStateFromError(): ErrorBoundaryState {
    return { hasError: true };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error("Restaurant workspace failed to render", error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <main className="error-shell">
          <section className="error-panel" role="alert">
            <p className="eyebrow">Workspace unavailable</p>
            <h1>We could not load the operations console.</h1>
            <p>Reload the workspace to try again.</p>
            <button
              type="button"
              className="text-button"
              onClick={() => window.location.reload()}
            >
              Reload workspace
            </button>
          </section>
        </main>
      );
    }

    return this.props.children;
  }
}

// ---------------------------------------------------------------------------
// Mount
// ---------------------------------------------------------------------------

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <AppErrorBoundary>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </AppErrorBoundary>
  </StrictMode>,
);
