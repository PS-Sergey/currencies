update currency_rate
set updated_at = now(),
    rate_time = :rate_time,
    status = :status,
    base = :base,
    target = :target,
    rate = :rate
where id = :id;