create type currency_type as enum ('USD', 'EUR', 'MXN');
create type currency_rate_status as enum ('PENDING', 'SUCCESS', 'ERROR');

create table if not exists currency_rate (
    id         uuid primary key,
    created_at timestamptz           not null default now(),
    updated_at timestamptz           not null default now(),
    rate_time  timestamptz,
    status     currency_rate_status not null,
    base       currency_type        not null,
    target     currency_type        not null,
    rate       real
);

create index if not exists currency_rate_base_target_idx
    on currency_rate(base, target);