create table if not exists import_jobs (
    id uuid primary key,
    import_id uuid not null,
    job_id uuid not null,
    status int not null,
    error text null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    foreign key (import_id) references imports (id)
);
