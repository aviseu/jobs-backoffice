alter table import_metrics drop constraint import_job_results_pkey;
alter table import_metrics add column id uuid not null default gen_random_uuid() primary key;
