create table if not exists import_job_results (
    import_id uuid not null,
    job_id uuid not null,
    result int not null,
    primary key(import_id, job_id),
    foreign key (import_id) references imports (id)
);
