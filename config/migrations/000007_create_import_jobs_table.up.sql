create table if not exists import_jobs (
    import_id uuid not null,
    job_id uuid not null,
    status int not null,
    primary key(import_id, job_id),
    foreign key (import_id) references imports (id)
);
