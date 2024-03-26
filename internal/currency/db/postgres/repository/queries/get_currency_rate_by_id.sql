select id,
       rate_time,
       status,
       base,
       target,
       rate
from currency_rate
where id = :id;