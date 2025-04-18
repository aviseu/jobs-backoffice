DROP INDEX IF EXISTS idx_import_metrics_import_id_created_at;
DROP INDEX IF EXISTS idx_import_metrics_job_id_created_at;
ALTER TABLE import_metrics DROP COLUMN error;
ALTER TABLE import_metrics DROP COLUMN created_at;
