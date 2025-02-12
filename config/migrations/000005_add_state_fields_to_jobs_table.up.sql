ALTER TABLE jobs ADD COLUMN status int NOT NULL DEFAULT 0;
ALTER TABLE jobs ADD COLUMN publish_status int NOT NULL DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_jobs_status_publish_status ON jobs(status, publish_status);
