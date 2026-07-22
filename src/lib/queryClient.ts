/**
 * TanStack Query client configuration.
 *
 * Centralises retry policy, stale time, and error handling so every
 * query in the app inherits consistent behaviour. In Milestone 4 this
 * will be configured to talk to the Go REST API.
 */

import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // Data is considered fresh for 30 seconds before a background refetch.
      staleTime: 30_000,
      // Retry failed requests up to 2 times with exponential backoff.
      retry: 2,
      retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30_000),
    },
  },
});
