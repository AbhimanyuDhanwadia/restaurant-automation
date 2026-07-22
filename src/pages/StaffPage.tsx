/**
 * StaffPage — /staff
 *
 * Shift roster, staff detail, manager handoff note, and open shift task list.
 */

import { ClipboardCheck } from "lucide-react";
import { useMemo, useState } from "react";
import { DirectoryRow } from "@/components/ui/DirectoryRow";
import { EmptyState } from "@/components/ui/EmptyState";
import { Panel, PanelHeading } from "@/components/ui/Panel";
import { StatusPill } from "@/components/ui/StatusPill";
import { useStaffStore } from "@/stores/staff";
import { toSlug } from "@/lib/utils";
import type { StaffStatus } from "@/types/domain";

export function StaffPage() {
  const staff = useStaffStore((s) => s.staff);
  const tasks = useStaffStore((s) => s.tasks);
  const handoffNote = useStaffStore((s) => s.handoffNote);
  const completeTask = useStaffStore((s) => s.completeTask);
  const saveHandoffNote = useStaffStore((s) => s.saveHandoffNote);

  const [filter, setFilter] = useState<StaffStatus | "All">("All");
  const [selectedId, setSelectedId] = useState(staff[0]?.id ?? "");
  const [noteInput, setNoteInput] = useState(handoffNote);
  const [noteSaved, setNoteSaved] = useState(false);

  const directory = useMemo(
    () => (filter === "All" ? staff : staff.filter((s) => s.status === filter)),
    [staff, filter],
  );

  const selected = directory.find((s) => s.id === selectedId) ?? directory[0];

  const handleSaveNote = () => {
    saveHandoffNote(noteInput);
    setNoteSaved(true);
  };

  return (
    <section className="staff-workspace">
      <div className="orders-toolbar">
        <div>
          <p className="eyebrow">Shift handoff</p>
          <h2>Staff &amp; Tasks</h2>
        </div>
        <label className="filter-control">
          <span>Availability</span>
          <select
            value={filter}
            onChange={(e) => setFilter(e.target.value as StaffStatus | "All")}
          >
            <option value="All">Everyone</option>
            <option value="On shift">On shift</option>
            <option value="On break">On break</option>
            <option value="Off shift">Off shift</option>
          </select>
        </label>
      </div>

      <div className="staff-directory">
        <Panel className="staff-roster" label="Staff roster">
          {directory.map((member) => (
            <DirectoryRow
              key={member.id}
              selected={selected?.id === member.id}
              onClick={() => setSelectedId(member.id)}
              primary={member.name}
              secondary={`${member.role} · ${member.station}`}
              badge={
                <StatusPill variant={`staff-status-${toSlug(member.status)}`}>
                  {member.status}
                </StatusPill>
              }
              label={`Select ${member.name}`}
            />
          ))}
          {directory.length === 0 && (
            <EmptyState message="No staff match this availability." />
          )}
        </Panel>

        <div className="staff-side-stack">
          {selected && (
            <Panel className="staff-detail" label={`Handoff for ${selected.name}`}>
              <PanelHeading
                eyebrow="Selected staff member"
                title={selected.name}
                action={
                  <StatusPill variant={`staff-status-${toSlug(selected.status)}`}>
                    {selected.status}
                  </StatusPill>
                }
              />
              <dl className="detail-list">
                <div><dt>Role</dt><dd>{selected.role}</dd></div>
                <div><dt>Station</dt><dd>{selected.station}</dd></div>
              </dl>
              <p className="eyebrow handoff-label">Handoff note</p>
              <p className="handoff-copy">{selected.handoff}</p>
            </Panel>
          )}

          <Panel className="handoff-panel" label="Shift handoff note">
            <PanelHeading
              eyebrow="Manager note"
              title="Pass to next shift"
              action={noteSaved ? <span className="saved-label">Saved</span> : undefined}
            />
            <textarea
              aria-label="Shift handoff note"
              placeholder="Add a note for the next shift"
              value={noteInput}
              onChange={(e) => {
                setNoteInput(e.target.value);
                setNoteSaved(false);
              }}
            />
            <button
              type="button"
              className="back-button"
              onClick={handleSaveNote}
              disabled={!noteInput.trim()}
            >
              <ClipboardCheck size={16} aria-hidden="true" />
              Save handoff
            </button>
          </Panel>
        </div>
      </div>

      <Panel className="staff-task-panel" label="Open shift tasks">
        <PanelHeading
          eyebrow="Automation queue"
          title="Open shift tasks"
          action={<span className="task-count">{tasks.length} open</span>}
        />
        <div className="task-list">
          {tasks.map((task) => (
            <article className="task-row" key={task.title}>
              <div>
                <strong>{task.title}</strong>
                <span>{task.owner}</span>
              </div>
              <time>{task.due}</time>
              <button
                type="button"
                className="task-complete"
                onClick={() => completeTask(task.title)}
              >
                Done
              </button>
            </article>
          ))}
          {tasks.length === 0 && <EmptyState message="All shift tasks are complete." />}
        </div>
      </Panel>
    </section>
  );
}
