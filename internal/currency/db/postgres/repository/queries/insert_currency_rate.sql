insert into currency_rate (id,
                           rate_time,
                           status,
                           base,
                           target,
                           rate)
values(:id, :rate_time, :status, :base, :target, :rate);