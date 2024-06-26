select id,
       rate_time,
       status,
       base,
       target,
       rate
from currency_rate
where base = :base
and target = :target
and status = 'SUCCESS'
order by rate_time desc
limit 1;