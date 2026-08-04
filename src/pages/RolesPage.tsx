import { KeyRound, Plus, RefreshCw } from "lucide-react";
import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { createRole, listRoles } from "@/features/roles/api";

const PERMISSIONS = [
  { value: "operations.view", label: "Operations" }, { value: "orders.manage", label: "Orders" }, { value: "kitchen.manage", label: "Kitchen" }, { value: "inventory.manage", label: "Inventory" }, { value: "staff.manage", label: "Staff" }, { value: "analytics.view", label: "Analytics" }, { value: "automation.view", label: "Automation" }, { value: "integrations.manage", label: "Integrations" }, { value: "printers.manage", label: "Printers" }, { value: "settings.manage", label: "Settings" }, { value: "audit.view", label: "Audit logs" }, { value: "database.view", label: "Database" }, { value: "roles.manage", label: "Roles" },
];

export function RolesPage() {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [permissions, setPermissions] = useState<string[]>(["operations.view"]);
  const queryClient = useQueryClient();
  const rolesQuery = useQuery({ queryKey: ["admin", "roles"], queryFn: listRoles });
  const createMutation = useMutation({
    mutationFn: createRole,
    onSuccess: () => {
      setName("");
      setDescription("");
      setPermissions(["operations.view"]);
      queryClient.invalidateQueries({ queryKey: ["admin", "roles"] });
    },
  });
  const roles = useMemo(() => rolesQuery.data ?? [], [rolesQuery.data]);

  function togglePermission(value: string) {
    setPermissions((current) => current.includes(value) ? current.filter((permission) => permission !== value) : [...current, value]);
  }

  return <section className="roles-workspace">
    <div className="automation-toolbar">
      <div><p className="eyebrow">Administration</p><h2>Roles</h2><p className="automation-subtitle">Durable permission templates for authenticated user assignments.</p></div>
      <button type="button" className="icon-button" onClick={() => rolesQuery.refetch()} disabled={rolesQuery.isFetching} title="Refresh roles" aria-label="Refresh roles"><RefreshCw size={18} aria-hidden="true" className={rolesQuery.isFetching ? "spin" : undefined} /></button>
    </div>

    {rolesQuery.isError && <div className="orders-api-error" role="alert"><span>{rolesQuery.error.message}</span><button type="button" className="text-button" onClick={() => rolesQuery.refetch()}>Retry</button></div>}
    {createMutation.isError && <p className="orders-api-error" role="alert">{createMutation.error.message}</p>}

    <div className="roles-grid">
      <Panel label="Role catalog"><PanelHeading eyebrow="Permission templates" title="Available roles" action={<span className="task-count">{rolesQuery.data ? `${roles.length} roles` : "Unavailable"}</span>} /><div className="role-list">{rolesQuery.isPending && <EmptyState message="Loading roles..." />}{roles.map((role) => <article className="role-row" key={role.id}><div><strong>{role.name}</strong><span>{role.description || "No description"}</span><small>{role.permissions.join(" · ")}</small></div><span className={role.system ? "role-system" : "role-custom"}>{role.system ? "System" : "Custom"}</span></article>)}{!rolesQuery.isPending && !rolesQuery.isError && roles.length === 0 && <EmptyState message="No roles have been configured." />}</div></Panel>

      <Panel label="Create role"><PanelHeading eyebrow="New template" title="Create custom role" action={<Plus size={20} aria-hidden="true" />} /><form className="role-create-form" onSubmit={(event) => { event.preventDefault(); createMutation.mutate({ name: name.trim(), description: description.trim(), permissions }); }}><label><span>Role name</span><input value={name} onChange={(event) => setName(event.target.value)} required maxLength={80} /></label><label><span>Description</span><input value={description} onChange={(event) => setDescription(event.target.value)} /></label><fieldset><legend>Permissions</legend><div className="role-permission-grid">{PERMISSIONS.map((permission) => <label key={permission.value}><input type="checkbox" checked={permissions.includes(permission.value)} onChange={() => togglePermission(permission.value)} /><span>{permission.label}</span></label>)}</div></fieldset><button type="submit" className="primary-button" disabled={createMutation.isPending || permissions.length === 0}>{createMutation.isPending ? "Creating..." : "Create role"}</button></form></Panel>
    </div>

    <Panel label="Role assignment boundary"><PanelHeading eyebrow="Access boundary" title="User assignments" action={<KeyRound size={20} aria-hidden="true" />} /><p className="analytics-summary-copy">Roles are durable permission templates. The Users workspace assigns them to identities already observed from validated Supabase JWTs; it does not provision or administer Supabase Auth users.</p></Panel>
  </section>;
}
