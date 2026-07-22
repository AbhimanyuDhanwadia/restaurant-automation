import { Component, StrictMode, type ErrorInfo, type ReactNode } from "react";
import { createRoot } from "react-dom/client";
import { App } from "./App";
import "./styles.css";

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
            <button className="text-button" onClick={() => window.location.reload()}>Reload workspace</button>
          </section>
        </main>
      );
    }

    return this.props.children;
  }
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <AppErrorBoundary>
      <App />
    </AppErrorBoundary>
  </StrictMode>,
);
