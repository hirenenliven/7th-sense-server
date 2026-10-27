-- Statistics per variable over the whole machine_telementry table (all machines together).
WITH gaps AS (
  SELECT machine_id, variable_id, value, recorded_at,
    EXTRACT(EPOCH FROM recorded_at - LAG(recorded_at)
      OVER (PARTITION BY machine_id, variable_id ORDER BY recorded_at)) * 1000 AS gap_ms
  FROM machine_telementry
),
stats AS (
  SELECT variable_id,
    COUNT(*) AS total_rows,
    COUNT(DISTINCT machine_id) AS machine_count,
    COUNT(DISTINCT date_trunc('second', recorded_at)) AS seconds_with_data,
    MIN(recorded_at) AS first_at,
    MAX(recorded_at) AS last_at,
    AVG(gap_ms) AS avg_gap_ms,
    MIN(gap_ms) AS min_gap_ms,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY gap_ms) AS median_gap_ms,
    MAX(gap_ms) AS max_gap_ms,
    MIN(value) AS min_value,
    MAX(value) AS max_value,
    AVG(value) AS avg_value
  FROM gaps
  GROUP BY variable_id
)
SELECT
  stats.variable_id,
  cnc_variables.code AS variable,
  stats.total_rows,
  ROUND(100.0 * stats.total_rows / SUM(stats.total_rows) OVER (), 2) AS percent_of_all_rows,
  stats.machine_count,
  stats.seconds_with_data,
  ROUND(stats.total_rows::numeric / NULLIF(stats.seconds_with_data, 0), 2) AS rows_per_active_second,
  ROUND(stats.avg_gap_ms::numeric, 1) AS avg_gap_ms,
  ROUND(stats.median_gap_ms::numeric, 1) AS median_gap_ms,
  ROUND(stats.min_gap_ms::numeric, 1) AS min_gap_ms,
  ROUND(stats.max_gap_ms::numeric / 1000, 1) AS max_gap_seconds,
  stats.min_value,
  stats.max_value,
  ROUND(stats.avg_value::numeric, 3) AS avg_value,
  stats.first_at,
  stats.last_at
FROM stats
JOIN cnc_variables ON cnc_variables.id = stats.variable_id
ORDER BY stats.total_rows DESC;
