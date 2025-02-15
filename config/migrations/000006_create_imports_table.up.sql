create table if not exists imports (
    id uuid primary key,
    channel_id uuid not null,
    status int not null,
    new_jobs int default 0,
    updated_jobs int default 0,
    no_change_jobs int default 0,
    missing_jobs int default 0,
    failed_jobs int default 0,
    error text null,
    started_at timestamptz not null,
    ended_at timestamptz null,
    foreign key (channel_id) references channels (id)
)
