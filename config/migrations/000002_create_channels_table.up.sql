create table channels (
    id uuid primary key,
    name text not null,
    integration text not null,
    status int not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
)