import { ClipboardCheck, RefreshCw } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { DirectoryRow } from "@/components/ui/DirectoryRow";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import { completeShiftTask, createShiftTask, createStaffMember, getShiftHandoff, listShiftTasks, listStaff, saveShiftHandoff, updateStaffHandoff, updateStaffStatus, type StaffStatus } from "@/features/staff/api";

const statuses: StaffStatus[] = ["on_shift", "on_break", "off_shift"];
const statusLabel = (status: StaffStatus) => ({ on_shift: "On shift", on_break: "On break", off_shift: "Off shift" })[status];
const statusClass = (status: StaffStatus) => status.replace(/_/g, "-");

export function StaffPage() {
  const [filter, setFilter] = useState<StaffStatus | "all">("all");
  const [selectedID, setSelectedID] = useState("");
  const [name, setName] = useState(""); const [role, setRole] = useState(""); const [station, setStation] = useState("");
  const [noteInput, setNoteInput] = useState(""); const [memberHandoff, setMemberHandoff] = useState("");
  const [taskTitle, setTaskTitle] = useState(""); const [taskOwner, setTaskOwner] = useState(""); const [taskDue, setTaskDue] = useState("");
  const queryClient = useQueryClient();
  const staffQuery = useQuery({ queryKey: ["staff"], queryFn: listStaff });
  const tasksQuery = useQuery({ queryKey: ["staff", "tasks"], queryFn: listShiftTasks });
  const handoffQuery = useQuery({ queryKey: ["staff", "handoff"], queryFn: getShiftHandoff });
  const refreshStaff = () => queryClient.invalidateQueries({ queryKey: ["staff"] });
  const createMember = useMutation({ mutationFn: createStaffMember, onSuccess: (member) => { setSelectedID(member.id); setName(""); setRole(""); setStation(""); refreshStaff(); } });
  const statusMutation = useMutation({ mutationFn: ({ id, status }: { id: string; status: StaffStatus }) => updateStaffStatus(id, status), onSuccess: refreshStaff });
  const memberHandoffMutation = useMutation({ mutationFn: ({ id, handoff }: { id: string; handoff: string }) => updateStaffHandoff(id, handoff), onSuccess: refreshStaff });
  const createTask = useMutation({ mutationFn: createShiftTask, onSuccess: () => { setTaskTitle(""); setTaskOwner(""); setTaskDue(""); refreshStaff(); } });
  const completeTask = useMutation({ mutationFn: completeShiftTask, onSuccess: refreshStaff });
  const saveHandoff = useMutation({ mutationFn: saveShiftHandoff, onSuccess: refreshStaff });
  const members = useMemo(() => staffQuery.data ?? [], [staffQuery.data]);
  const tasks = tasksQuery.data ?? [];
  const directory = useMemo(() => filter === "all" ? members : members.filter((member) => member.status === filter), [filter, members]);
  const selected = directory.find((member) => member.id === selectedID) ?? directory[0];
  const error = staffQuery.error ?? tasksQuery.error ?? handoffQuery.error ?? createMember.error ?? statusMutation.error ?? memberHandoffMutation.error ?? createTask.error ?? completeTask.error ?? saveHandoff.error;
  useEffect(() => { setNoteInput(handoffQuery.data?.note ?? ""); }, [handoffQuery.data?.note]);
  useEffect(() => { setMemberHandoff(selected?.handoff ?? ""); }, [selected?.id, selected?.handoff]);

  function createMemberSubmit(event: React.FormEvent<HTMLFormElement>) { event.preventDefault(); createMember.mutate({ name: name.trim(), role: role.trim(), station: station.trim() }); }
  function createTaskSubmit(event: React.FormEvent<HTMLFormElement>) { event.preventDefault(); createTask.mutate({ title: taskTitle.trim(), owner: taskOwner.trim(), due_label: taskDue.trim() }); }
  function saveMemberHandoff(event: React.FormEvent<HTMLFormElement>) { event.preventDefault(); if (selected) memberHandoffMutation.mutate({ id: selected.id, handoff: memberHandoff }); }

  return <section className="staff-workspace">
    <div className="orders-toolbar"><div><p className="eyebrow">Shift handoff</p><h2>Staff &amp; Tasks</h2></div><label className="filter-control"><span>Availability</span><select value={filter} onChange={(event) => setFilter(event.target.value as StaffStatus | "all")}><option value="all">Everyone</option>{statuses.map((status) => <option key={status} value={status}>{statusLabel(status)}</option>)}</select></label></div>
    <Panel label="Create staff member"><form className="staff-create-form" onSubmit={createMemberSubmit}><label><span>Name</span><input value={name} onChange={(event) => setName(event.target.value)} required /></label><label><span>Role</span><input value={role} onChange={(event) => setRole(event.target.value)} required /></label><label><span>Station</span><input value={station} onChange={(event) => setStation(event.target.value)} required /></label><button type="submit" className="primary-button" disabled={createMember.isPending}>{createMember.isPending ? "Creating..." : "Create staff"}</button></form></Panel>
    {(staffQuery.isError || tasksQuery.isError || handoffQuery.isError) && <div className="orders-api-error" role="alert"><span>{error?.message}</span><button type="button" className="text-button" onClick={refreshStaff}><RefreshCw size={15} aria-hidden="true" /> Retry</button></div>}{error && !staffQuery.isError && !tasksQuery.isError && !handoffQuery.isError && <p className="orders-api-error" role="alert">{error.message}</p>}
    <div className="staff-directory"><Panel className="staff-roster" label="Staff roster">{staffQuery.isPending && <EmptyState message="Loading staff..." />}{directory.map((member) => <DirectoryRow key={member.id} selected={selected?.id === member.id} onClick={() => setSelectedID(member.id)} primary={member.name} secondary={`${member.role} · ${member.station}`} badge={<StatusPill variant={`staff-status-${statusClass(member.status)}`}>{statusLabel(member.status)}</StatusPill>} label={`Select ${member.name}`} />)}{!staffQuery.isPending && !staffQuery.isError && directory.length === 0 && <EmptyState message="No staff match this availability." />}</Panel><div className="staff-side-stack">{selected && <Panel className="staff-detail" label={`Handoff for ${selected.name}`}><PanelHeading eyebrow="Selected staff member" title={selected.name} action={<label className="order-status-control"><span>Availability</span><select value={selected.status} disabled={statusMutation.isPending} onChange={(event) => statusMutation.mutate({ id: selected.id, status: event.target.value as StaffStatus })}>{statuses.map((status) => <option key={status} value={status}>{statusLabel(status)}</option>)}</select></label>} /><dl className="detail-list"><div><dt>Role</dt><dd>{selected.role}</dd></div><div><dt>Station</dt><dd>{selected.station}</dd></div></dl><form className="member-handoff-form" onSubmit={saveMemberHandoff}><label><span>Handoff note</span><textarea value={memberHandoff} onChange={(event) => setMemberHandoff(event.target.value)} /></label><button type="submit" className="back-button" disabled={memberHandoffMutation.isPending}><ClipboardCheck size={16} aria-hidden="true" /> Save handoff</button></form></Panel>}<Panel className="handoff-panel" label="Shift handoff note"><PanelHeading eyebrow="Manager note" title="Pass to next shift" action={saveHandoff.isSuccess ? <span className="saved-label">Saved</span> : undefined} /><textarea aria-label="Shift handoff note" placeholder="Add a note for the next shift" value={noteInput} onChange={(event) => setNoteInput(event.target.value)} /><button type="button" className="back-button" onClick={() => saveHandoff.mutate(noteInput)} disabled={saveHandoff.isPending}><ClipboardCheck size={16} aria-hidden="true" /> Save handoff</button></Panel></div></div>
    <Panel className="staff-task-panel" label="Open shift tasks"><PanelHeading eyebrow="Automation queue" title="Open shift tasks" action={<span className="task-count">{tasks.length} open</span>} /><form className="task-create-form" onSubmit={createTaskSubmit}><label><span>Task</span><input value={taskTitle} onChange={(event) => setTaskTitle(event.target.value)} required /></label><label><span>Owner</span><input value={taskOwner} onChange={(event) => setTaskOwner(event.target.value)} required /></label><label><span>Due</span><input value={taskDue} onChange={(event) => setTaskDue(event.target.value)} /></label><button type="submit" className="text-button" disabled={createTask.isPending}>Add task</button></form><div className="task-list">{tasks.map((task) => <article className="task-row" key={task.id}><div><strong>{task.title}</strong><span>{task.owner}</span></div><time>{task.due_label || "No due time"}</time><button type="button" className="task-complete" disabled={completeTask.isPending} onClick={() => completeTask.mutate(task.id)}>Done</button></article>)}{!tasksQuery.isPending && !tasksQuery.isError && tasks.length === 0 && <EmptyState message="All shift tasks are complete." />}</div></Panel>
  </section>;
}
