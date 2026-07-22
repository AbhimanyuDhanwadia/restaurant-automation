/**
 * StatusPill
 *
 * A pill-shaped badge that visually communicates the status of an entity.
 * The `variant` prop maps directly to CSS class suffixes defined in globals.css,
 * keeping the colour logic out of JavaScript and in the design system.
 *
 * Usage:
 *   <StatusPill variant="preparing">Preparing</StatusPill>
 *   <StatusPill variant="severity-high">High</StatusPill>
 */

import type { ReactNode } from "react";

interface StatusPillProps {
  /** CSS class suffix (e.g. "preparing", "severity-high", "staff-status-on-shift"). */
  variant: string;
  children: ReactNode;
  /** Optional click handler — turns the pill into a button for in-place status cycling. */
  onClick?: () => void;
  className?: string;
}

export function StatusPill({ variant, children, onClick, className = "" }: StatusPillProps) {
  const base = `status ${variant} ${className}`.trim();

  if (onClick) {
    return (
      <button type="button" className={`${base} status-button`} onClick={onClick}>
        {children}
      </button>
    );
  }

  return <span className={base}>{children}</span>;
}
