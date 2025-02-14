create table if not exists imports (
    id uuid primary key,
    channel_id uuid not null,
    status int not null,
    new_jobs int not null,
    updated_jobs int not null,
    no_change_jobs int not null,
    missing_jobs int not null,
    failed_jobs int not null,
    error text null,
    started_at timestamptz not null,
    ended_at timestamptz null,
    foreign key (channel_id) references channels (id)
)
