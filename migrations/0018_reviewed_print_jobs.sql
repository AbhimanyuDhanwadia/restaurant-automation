ALTER TABLE print_jobs
  DROP CONSTRAINT IF EXISTS print_jobs_status_check;

ALTER TABLE print_jobs
  ADD CONSTRAINT print_jobs_status_check
  CHECK (status IN ('queued', 'printing', 'printed', 'failed', 'reviewed'));
