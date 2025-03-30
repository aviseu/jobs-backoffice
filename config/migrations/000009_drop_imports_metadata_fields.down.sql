ALTER TABLE imports ADD COLUMN new_jobs int default 0;
ALTER TABLE imports ADD COLUMN updated_jobs int default 0;
ALTER TABLE imports ADD COLUMN no_change_jobs int default 0;
ALTER TABLE imports ADD COLUMN missing_jobs int default 0;
ALTER TABLE imports ADD COLUMN failed_jobs int default 0;
