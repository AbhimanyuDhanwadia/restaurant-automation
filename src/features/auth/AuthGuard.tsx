/**
 * AuthGuard
 *
 * Wraps protected routes. Renders a loading state while the Supabase
 * session is being restored, redirects to the login page when there is
 * no session, and renders children when authenticated.
 *
 * This is intentionally a component guard (not a route-level `beforeLoad`)
 * for Milestone 1. In Milestone 2 we will migrate to a proper TanStack
 * Router `beforeLoad` + `redirect` pattern once the Go API validates JWTs.
 */

import { useEffect, useState } from "react";
import type { Session } from "@supabase/supabase-js";
import { supabase, isSupabaseConfigured } from "@/lib/supabase";
import { LoginPage } from "./LoginPage";

interface AuthGuardProps {
  children: (session: Session) => React.ReactNode;
}

export function AuthGuard({ children }: AuthGuardProps) {
  const [session, setSession] = useState<Session | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!supabase) {
      setLoading(false);
      return;
    }

    supabase.auth.getSession().then(({ data }) => {
      setSession(data.session);
      setLoading(false);
    });

    const { data: listener } = supabase.auth.onAuthStateChange((_event, next) => {
      setSession(next);
      setLoading(false);
    });

    return () => listener.subscription.unsubscribe();
  }, []);

  if (!isSupabaseConfigured) {
    return <LoginPage />;
  }

  if (loading) {
    return (
      <main className="error-shell">
        <section className="error-panel" aria-live="polite">
          <p className="eyebrow">Restaurant Automation</p>
          <h1>Restoring your session...</h1>
        </section>
      </main>
    );
  }

  if (!session) {
    return <LoginPage />;
  }

  return <>{children(session)}</>;
}
