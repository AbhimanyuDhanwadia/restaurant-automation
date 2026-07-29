/**
 * Topbar
 *
 * The horizontal page header containing the page title, search box,
 * create-action trigger, and sign-out button.
 * The search term is managed via UIStore so OperationsPage can read it directly.
 */

import { LogOut, Plus, Search, X } from "lucide-react";
import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import type { Session } from "@supabase/supabase-js";
import { supabase } from "@/lib/supabase";
import { useUIStore } from "@/stores/ui";

interface TopbarProps {
  eyebrow: string;
  title: string;
  session: Session;
}

export function Topbar({ eyebrow, title, session }: TopbarProps) {
  const [createOpen, setCreateOpen] = useState(false);
  const searchTerm = useUIStore((s) => s.searchTerm);
  const setSearchTerm = useUIStore((s) => s.setSearchTerm);
  const navigate = useNavigate();

  const signOut = () => supabase?.auth.signOut();

  const createQuickAction = (type: "order" | "table" | "inventory") => {
    if (type === "order") {
      navigate({ to: "/orders" });
      setCreateOpen(false);
      return;
    }
    if (type === "table") {
      navigate({ to: "/tables" });
      setCreateOpen(false);
      return;
    }
    if (type === "inventory") {
      navigate({ to: "/inventory" });
      setCreateOpen(false);
      return;
    }
    setCreateOpen(false);
  };

  return (
    <>
      <header className="topbar">
        <div>
          <p className="eyebrow">{eyebrow}</p>
          <h1>{title}</h1>
        </div>
        <div className="topbar-actions">
          <label className="search-box">
            <Search size={17} aria-hidden="true" />
            <input
              aria-label="Search operations"
              placeholder="Search orders, tables, items"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </label>
          <button
            id="create-action-trigger"
            type="button"
            className="icon-button"
            aria-label="Create action"
            onClick={() => setCreateOpen(true)}
          >
            <Plus size={20} aria-hidden="true" />
          </button>
          <button
            id="sign-out-button"
            type="button"
            className="icon-button"
            aria-label="Sign out"
            title={`Sign out ${session.user.email ?? "user"}`}
            onClick={signOut}
          >
            <LogOut size={20} aria-hidden="true" />
          </button>
        </div>
      </header>

      {createOpen && (
        <div
          className="modal-backdrop"
          role="presentation"
          onClick={() => setCreateOpen(false)}
        >
          <section
            className="modal-panel"
            role="dialog"
            aria-modal="true"
            aria-labelledby="create-action-title"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="panel-heading">
              <div>
                <p className="eyebrow">Quick action</p>
                <h2 id="create-action-title">Create new</h2>
              </div>
              <button
                type="button"
                className="dismiss-button"
                aria-label="Close create action"
                onClick={() => setCreateOpen(false)}
              >
                <X size={18} aria-hidden="true" />
              </button>
            </div>
            <div className="create-action-list">
              <button type="button" className="create-action" onClick={() => createQuickAction("order")}>
                New order
              </button>
              <button type="button" className="create-action" onClick={() => createQuickAction("table")}>
                New table
              </button>
              <button type="button" className="create-action" onClick={() => createQuickAction("inventory")}>
                New inventory item
              </button>
            </div>
          </section>
        </div>
      )}
    </>
  );
}
