/**
 * EmptyState
 *
 * Displayed when a filtered list has no results. Keeps the empty-state
 * copy consistent across all views.
 */

interface EmptyStateProps {
  message: string;
}

export function EmptyState({ message }: EmptyStateProps) {
  return <p className="empty-state">{message}</p>;
}
