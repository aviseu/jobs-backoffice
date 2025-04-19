CREATE INDEX IF NOT EXISTS idx_jobs_channel_id_posted_at ON jobs(channel_id, posted_at desc);

