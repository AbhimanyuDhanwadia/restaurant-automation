/**
 * Panel
 *
 * The primary card-like container used throughout the dashboard.
 * Wraps content with the standard border, background, shadow, and padding
 * defined in the design system.
 */

import type { ReactNode } from "react";

interface PanelProps {
  children: ReactNode;
  className?: string;
  /** Optional aria-label for screen readers. */
  label?: string;
  /** HTML id attribute, used for landmark navigation. */
  id?: string;
}

export function Panel({ children, className = "", label, id }: PanelProps) {
  return (
    <section
      className={`panel ${className}`.trim()}
      aria-label={label}
      id={id}
    >
      {children}
    </section>
  );
}

// ---------------------------------------------------------------------------
// PanelHeading — the two-column header (title + action) inside a Panel.
// ---------------------------------------------------------------------------

interface PanelHeadingProps {
  eyebrow: string;
  title: string;
  action?: ReactNode;
  id?: string;
}

export function PanelHeading({ eyebrow, title, action, id }: PanelHeadingProps) {
  return (
    <div className="panel-heading">
      <div>
        <p className="eyebrow">{eyebrow}</p>
        <h2 id={id}>{title}</h2>
      </div>
      {action}
    </div>
  );
}
