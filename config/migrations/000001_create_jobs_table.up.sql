create table jobs (
    uuid serial primary key,
    url text not null,
    title text not null,
    description jsonb not null,
    source text not null,
    location text not null,
    remote bool not null default false,
    posted_at timestamptz not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
)
