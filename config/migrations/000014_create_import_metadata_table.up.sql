create table if not exists import_metadata (
    import_id uuid primary key,
    new_jobs int default 0,
    updated_jobs int default 0,
    no_change_jobs int default 0,
    missing_jobs int default 0,
    errors int default 0,
    published int default 0,
    late_published int default 0,
    missing_published int default 0,
    foreign key (import_id) references imports (id)
)
