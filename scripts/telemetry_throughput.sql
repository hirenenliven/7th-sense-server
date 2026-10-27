-- machine_telementry frequency and throughput. Change the window in each query
-- (now() - interval '5 minutes') to the period you want; rows are only changed values.

-- 1. Rows per second (newest first): how much arrives each second.
SELECT time_bucket('1 second', recorded_at) AS second,
       count(*)                    AS rows,
       count(DISTINCT machine_id)  AS machines,
       count(DISTINCT variable_id) AS variables
FROM machine_telementry
WHERE recorded_at > now() - interval '5 minutes'
GROUP BY 1
ORDER BY 1 DESC;

-- 2. Throughput summary: average and peak rows per second over the window.
SELECT sum(rows_per_second)               AS total_rows,
       count(*)                           AS seconds_with_data,
       min(rows_per_second)               AS min_per_second,
       round(avg(rows_per_second), 2)     AS avg_per_second,
       max(rows_per_second)               AS peak_per_second
FROM (
    SELECT time_bucket('1 second', recorded_at) AS second, count(*) AS rows_per_second
    FROM machine_telementry
    WHERE recorded_at > now() - interval '5 minutes'
    GROUP BY 1
) per_second;

-- 3. Per machine: rows, average / peak per second, first and last row.
SELECT m.name AS machine, t.machine_id,
       sum(t.rows_per_second)              AS rows,
       round(avg(t.rows_per_second), 2)    AS avg_per_second,
       max(t.rows_per_second)              AS peak_per_second,
       min(t.second) AS first_second, max(t.second) AS last_second
FROM (
    SELECT machine_id, time_bucket('1 second', recorded_at) AS second, count(*) AS rows_per_second
    FROM machine_telementry
    WHERE recorded_at > now() - interval '5 minutes'
    GROUP BY 1, 2
) t
JOIN machines m ON m.id = t.machine_id
GROUP BY m.name, t.machine_id
ORDER BY rows DESC;

-- 4. Per variable frequency: how often each variable changes (avg gap between rows).
SELECT machine_id, t.variable_id, v.code AS variable,
       count(*)                                                        AS rows,
       round(count(*) / greatest(extract(epoch FROM max(recorded_at) - min(recorded_at)), 1)::numeric, 2)
                                                                       AS rows_per_second,
       round(avg(gap_ms)::numeric, 1)                                  AS avg_gap_ms,
       round(min(gap_ms)::numeric, 1)                                  AS min_gap_ms
FROM (
    SELECT machine_id, variable_id, recorded_at,
           extract(epoch FROM recorded_at - lag(recorded_at) OVER (PARTITION BY machine_id, variable_id ORDER BY recorded_at)) * 1000 AS gap_ms
    FROM machine_telementry
    WHERE recorded_at > now() - interval '5 minutes'
) t
JOIN cnc_variables v ON v.id = t.variable_id
GROUP BY machine_id, t.variable_id, v.code
ORDER BY rows DESC;
