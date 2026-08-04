import { RefreshCw, ShieldCheck, UserRoundCheck, UsersRound } from "lucide-react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { listRoles } from "@/features/roles/api";
import { listUsers, updateUserRole } from "@/features/users/api";

export function UsersPage() {
  const queryClient = useQueryClient();
  const usersQuery = useQuery({ queryKey: ["admin", "users"], queryFn: listUsers, refetchInterval: 15_000 });
  const rolesQuery = useQuery({ queryKey: ["admin", "roles"], queryFn: listRoles });
  const roleMutation = useMutation({
    mutationFn: updateUserRole,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "users"] });
      queryClient.invalidateQueries({ queryKey: ["auth", "current-user"] });
    },
  });
  const users = usersQuery.data ?? [];
  const roles = rolesQuery.data ?? [];

  return <section className="users-workspace">
    <div className="automation-toolbar">
      <div><p className="eyebrow">Administration</p><h2>Users</h2><p className="automation-subtitle">Observed restaurant identities and their assigned application roles.</p></div>
      <button type="button" className="icon-button" onClick={() => { usersQuery.refetch(); rolesQuery.refetch(); }} disabled={usersQuery.isFetching || rolesQuery.isFetching} title="Refresh users" aria-label="Refresh users"><RefreshCw size={18} aria-hidden="true" className={usersQuery.isFetching || rolesQuery.isFetching ? "spin" : undefined} /></button>
    </div>

    {usersQuery.isError && <div className="orders-api-error" role="alert"><span>{usersQuery.error.message}</span><button type="button" className="text-button" onClick={() => usersQuery.refetch()}>Retry</button></div>}
    {rolesQuery.isError && <div className="orders-api-error" role="alert"><span>{rolesQuery.error.message}</span><button type="button" className="text-button" onClick={() => rolesQuery.refetch()}>Retry</button></div>}
    {roleMutation.isError && <p className="orders-api-error" role="alert">{roleMutation.error.message}</p>}

    <section className="automation-metric-grid" aria-label="User summary"><article className="automation-metric"><UsersRound size={20} aria-hidden="true" /><span>Observed users</span><strong>{usersQuery.data ? users.length : "Unavailable"}</strong><small>Authenticated API identities</small></article><article className="automation-metric"><ShieldCheck size={20} aria-hidden="true" /><span>Role templates</span><strong>{rolesQuery.data ? roles.length : "Unavailable"}</strong><small>Available assignments</small></article><article className="automation-metric"><UserRoundCheck size={20} aria-hidden="true" /><span>Administrators</span><strong>{usersQuery.data ? users.filter((user) => user.role_name === "Administrator").length : "Unavailable"}</strong><small>Assigned in the local directory</small></article><article className="automation-metric"><RefreshCw size={20} aria-hidden="true" /><span>Directory source</span><strong>JWT</strong><small>Supabase identities are observed, not provisioned</small></article></section>

    <Panel label="Application users"><PanelHeading eyebrow="Identity directory" title="Observed users" action={<span className="task-count">{usersQuery.data ? `${users.length} users` : "Unavailable"}</span>} /><div className="user-directory-list">{usersQuery.isPending && <EmptyState message="Loading observed users..." />}{users.map((user) => { const saving = roleMutation.isPending && roleMutation.variables?.subject === user.auth_subject; return <article className="user-directory-row" key={user.auth_subject}><div><strong>{user.email}</strong><span>Last seen {new Date(user.last_seen_at).toLocaleString()}</span><small>Subject: {user.auth_subject}</small></div><label className="user-role-control"><span>Application role</span><select value={user.role_id} onChange={(event) => roleMutation.mutate({ subject: user.auth_subject, roleID: event.target.value })} disabled={rolesQuery.isPending || saving}>{roles.map((role) => <option key={role.id} value={role.id}>{role.name}</option>)}</select>{saving && <small>Saving...</small>}</label></article>; })}{!usersQuery.isPending && !usersQuery.isError && users.length === 0 && <EmptyState message="No identities have reached the authenticated API yet." />}</div></Panel>

    <Panel label="User administration boundary"><PanelHeading eyebrow="Security boundary" title="Local role assignments" action={<ShieldCheck size={20} aria-hidden="true" />} /><p className="analytics-summary-copy">Users appear automatically after authenticated API activity. Only email addresses in the API&apos;s `ADMIN_EMAILS` configuration can view this workspace or change assignments. This workspace never creates, deletes, or lists Supabase Auth users and requires no Supabase service-role credential.</p></Panel>
  </section>;
}
