DROP INDEX IF EXISTS idx_jobs_status_publish_status;
ALTER TABLE jobs DROP COLUMN status;
ALTER TABLE jobs DROP COLUMN publish_status;
