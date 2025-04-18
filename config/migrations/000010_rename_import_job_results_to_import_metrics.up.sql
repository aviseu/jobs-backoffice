alter table import_job_results rename to import_metrics;
alter table import_metrics rename column result to metric_type;
