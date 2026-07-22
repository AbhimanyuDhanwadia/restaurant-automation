/**
 * Utility helpers shared across the application.
 */

/**
 * Convert a string value (e.g. "Needs Check") into a kebab-case CSS class
 * segment (e.g. "needs-check"). Used to derive status-variant class names.
 */
export const toSlug = (value: string): string =>
  value.toLowerCase().replace(/\s+/g, "-");
