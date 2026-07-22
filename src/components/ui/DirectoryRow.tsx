/**
 * DirectoryRow
 *
 * A list item used in the master-detail panels (Orders, Tables, Inventory,
 * Staff, Alerts). Renders an icon chevron and highlights when selected.
 */

import { ChevronRight } from "lucide-react";
import type { ReactNode } from "react";

interface DirectoryRowProps {
  /** Whether this row is the currently selected item. */
  selected: boolean;
  onClick: () => void;
  /** Primary content (name, ID, etc.) */
  primary: ReactNode;
  /** Secondary / sub-label content. */
  secondary?: ReactNode;
  /** Status pill or badge rendered on the right. */
  badge?: ReactNode;
  /** aria-label for the button. */
  label?: string;
}

export function DirectoryRow({
  selected,
  onClick,
  primary,
  secondary,
  badge,
  label,
}: DirectoryRowProps) {
  return (
    <button
      type="button"
      className={`directory-row${selected ? " selected" : ""}`}
      onClick={onClick}
      aria-label={label}
      aria-pressed={selected}
    >
      <span>
        <strong>{primary}</strong>
        {secondary && <small>{secondary}</small>}
      </span>
      {badge}
      <ChevronRight size={18} aria-hidden="true" />
    </button>
  );
}
