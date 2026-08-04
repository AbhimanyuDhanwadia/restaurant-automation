# Backups API

Migration `0013_backup_records.sql` creates durable `backup_records`. Each record stores an external artifact location, size, status, completion time, and creation time.

`GET /api/v1/admin/backups` lists backup records. `POST /api/v1/admin/backups` registers an external artifact as `recorded`; it cannot mark it verified or execute a dump. Verification and failure reporting are reserved for a future trusted backup-worker integration.
