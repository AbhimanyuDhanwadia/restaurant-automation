# Staff API

All endpoints are authenticated under `/api/v1`.

- `GET/POST /staff` lists and creates staff members.
- `PATCH /staff/{staffID}/status` updates `on_shift`, `on_break`, or `off_shift`.
- `PATCH /staff/{staffID}/handoff` saves a member handoff note.
- `GET/POST /staff/tasks` lists open tasks and creates a task.
- `PATCH /staff/tasks/{taskID}/complete` marks a task complete.
- `GET/PUT /staff/handoff` reads and writes the current manager shift note.

The `0006_staff.sql` migration persists members, tasks, and the current handoff note in PostgreSQL. The in-memory repository remains available for local development without a database.

Completed tasks remain stored for analytics even though `GET /staff/tasks` returns only open work. The analytics overview uses the staff service's aggregate task summary and does not expose completed task records through the operational task list.
