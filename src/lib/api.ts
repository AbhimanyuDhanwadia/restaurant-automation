import { supabase } from "@/lib/supabase";

const apiBaseUrl = (import.meta.env.VITE_API_URL ?? "http://localhost:8080").replace(/\/$/, "");

/** Performs an authenticated request against the Restaurant Automation API. */
export async function apiRequest<T>(path: string, init: RequestInit = {}): Promise<T> {
  const session = supabase ? (await supabase.auth.getSession()).data.session : null;
  const headers = new Headers(init.headers);

  headers.set("Accept", "application/json");
  if (init.body) headers.set("Content-Type", "application/json");
  if (session) headers.set("Authorization", `Bearer ${session.access_token}`);

  let response: Response;
  try {
    response = await fetch(`${apiBaseUrl}${path}`, { ...init, headers });
  } catch {
    throw new Error("The API service is unreachable. Check that the API is running.");
  }

  if (!response.ok) {
    const detail = await response.text();
    throw new Error(detail || `API request failed (${response.status}).`);
  }

  if (response.status === 204) return undefined as T;

  return response.json() as Promise<T>;
}
