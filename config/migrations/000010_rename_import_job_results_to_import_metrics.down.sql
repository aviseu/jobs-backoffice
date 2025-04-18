alter table import_metrics rename to import_job_results;
alter table import_job_results rename column metric_type to result;
