ALTER TABLE import_metrics ADD COLUMN error text;
ALTER TABLE import_metrics ADD COLUMN created_at timestamptz NOT NULL DEFAULT now();
CREATE INDEX IF NOT EXISTS idx_import_metrics_import_id_created_at ON import_metrics(import_id, created_at);
CREATE INDEX IF NOT EXISTS idx_import_metrics_job_id_created_at ON import_metrics(job_id, created_at);
