# Backups UI

The Backups workspace refreshes durable backup records every 15 seconds. Operators can register the location and size of an externally-created artifact, which is displayed as recorded rather than verified.

The browser cannot run `pg_dump`, access database credentials, or download backup files. A scheduled infrastructure worker must perform backup creation and report verification separately.
