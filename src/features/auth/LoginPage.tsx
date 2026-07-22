/**
 * LoginPage
 *
 * Renders the Supabase email/password sign-in form. Shown only when there
 * is no active session. It is rendered by AuthGuard until a dedicated auth
 * route is introduced.
 */

import { ChefHat } from "lucide-react";
import { useState, type FormEvent } from "react";
import { isSupabaseConfigured, supabase } from "@/lib/supabase";

export function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  if (!isSupabaseConfigured) {
    return (
      <main className="error-shell">
        <section className="error-panel" role="alert">
          <p className="eyebrow">Authentication setup</p>
          <h1>Connect Supabase before opening the workspace.</h1>
          <p>Add the Vite Supabase environment variables to this deployment.</p>
        </section>
      </main>
    );
  }

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSubmitting(true);
    setError("");

    const { error: authError } = await supabase!.auth.signInWithPassword({
      email,
      password,
    });

    if (authError) {
      setError(authError.message);
    }

    setSubmitting(false);
  };

  return (
    <main className="auth-shell">
      <section className="auth-panel">
        <div className="brand auth-brand">
          <ChefHat size={26} aria-hidden="true" />
          <span>Restaurant Automation</span>
        </div>
        <p className="eyebrow">Team access</p>
        <h1>Sign in to operations</h1>
        <p className="auth-copy">Use your restaurant team account to continue.</p>
        <form className="auth-form" onSubmit={handleSubmit} id="login-form">
          <label>
            Email
            <input
              id="login-email"
              type="email"
              autoComplete="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </label>
          <label>
            Password
            <input
              id="login-password"
              type="password"
              autoComplete="current-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </label>
          {error && (
            <p className="auth-error" role="alert">
              {error}
            </p>
          )}
          <button
            id="login-submit"
            className="auth-submit"
            type="submit"
            disabled={submitting}
          >
            {submitting ? "Signing in..." : "Sign in"}
          </button>
        </form>
      </section>
    </main>
  );
}
