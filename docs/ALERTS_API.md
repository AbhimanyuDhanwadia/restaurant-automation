# Alerts API

Authenticated endpoints under `/api/v1`:

- `GET /alerts` lists active alerts.
- `POST /alerts` creates a manual alert with label, detail, and `high`, `medium`, or `low` severity.
- `PATCH /alerts/{alertID}/acknowledge` persists acknowledgement and removes the alert from active results.

The `0007_alerts.sql` migration stores alert acknowledgement history. Future automated alert producers should create alerts through this service rather than bypassing it.
