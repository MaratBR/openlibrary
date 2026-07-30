-- name: Analytics_InsertEvent :copyfrom
insert into ol_analytics.interaction_event 
    (user_key, book_id, event_type, value, created_at)
    values
    ($1, $2, $3, $4, $5);

-- name: Analytics_UpdateMetrics :exec
WITH 
    input_rows AS (
        SELECT
            (sqlc.arg('book_ids')::bigint[])[i] AS book_id,
            (sqlc.arg('metrics')::text[])[i] AS metric,
            (sqlc.arg('values')::double precision[])[i] AS increment_value,
            (sqlc.arg('samples')::bigint[])[i] AS increment_samples
        FROM generate_subscripts(
            sqlc.arg('book_ids')::bigint[],
            1
        ) AS indexes(i)
    ),
    bucket_definitions AS (
        SELECT
            'all'::ol_analytics.bucket_period_type AS bucket_type,
            TIMESTAMPTZ '1970-01-01 00:00:00+00' AS bucket_start

        UNION ALL

        SELECT
            'year'::ol_analytics.bucket_period_type,
            date_trunc('year', sqlc.arg(day_start)::timestamptz, 'UTC')

        UNION ALL

        SELECT
            'month'::ol_analytics.bucket_period_type,
            date_trunc('month', sqlc.arg(day_start)::timestamptz, 'UTC')

        UNION ALL

        SELECT
            'week'::ol_analytics.bucket_period_type,
            date_trunc('week', sqlc.arg(day_start)::timestamptz, 'UTC')

        UNION ALL

        SELECT
            'day'::ol_analytics.bucket_period_type,
            date_trunc('day', sqlc.arg(day_start)::timestamptz, 'UTC')
    ),
    aggregated_buckets AS (
        SELECT
            input.book_id,
            input.metric,
            buckets.bucket_type,
            buckets.bucket_start,
            sum(input.increment_value) AS increment_value,
            sum(input.increment_samples) AS increment_samples
        FROM input_rows AS input
        CROSS JOIN bucket_definitions AS buckets
        GROUP BY
            input.book_id,
            input.metric,
            buckets.bucket_type,
            buckets.bucket_start
    )
INSERT INTO ol_analytics.bucket (
    book_id,
    metric,
    bucket_type,
    bucket_start,
    value_sum,
    samples_count
)
SELECT
    book_id,
    metric,
    bucket_type,
    bucket_start,
    increment_value,
    increment_samples
FROM aggregated_buckets
ON CONFLICT (
    book_id,
    metric,
    bucket_type,
    bucket_start
)
DO UPDATE SET
    value_sum = ol_analytics.bucket.value_sum + EXCLUDED.value_sum,
    samples_count = ol_analytics.bucket.samples_count + EXCLUDED.samples_count,
    updated_at = now();


-- name: Analytics_RecalculateBookPopularity :exec
WITH metric_weights AS (
    SELECT
        mw.metric,
        mw.weight_text::double precision AS weight
    FROM jsonb_each_text(sqlc.arg(weights)::jsonb)
        AS mw(metric, weight_text)
)
INSERT INTO ol_analytics.book_popularity_bucket (
    book_id,
    bucket_type,
    bucket_start,
    value,
    updated_at
)
SELECT
    b.book_id,
    b.bucket_type,
    b.bucket_start,
    SUM(b.value_sum * w.weight),
    now()
FROM ol_analytics.bucket AS b
JOIN metric_weights AS w ON w.metric = b.metric
WHERE b.bucket_type = sqlc.arg(bucket_type)
AND b.bucket_start = sqlc.arg(bucket_start)
GROUP BY b.book_id, b.bucket_type, b.bucket_start
ON CONFLICT (book_id, bucket_type, bucket_start)
DO UPDATE SET
    value = EXCLUDED.value,
    updated_at = EXCLUDED.updated_at;

-- name: Analytics_GetMetricValue :one
select samples_count, value_sum
from ol_analytics.bucket
where bucket_type = $1 and book_id = $2 and metric = $3;

-- name: Analytics_GetMetricValues :many
select book_id, metric, bucket_type, samples_count, value_sum
from ol_analytics.bucket
where book_id = ANY(sqlc.arg('book_ids')::int8[]) and metric = $1;

-- name: Analytics_GetTopBooksBySamplesCount :many
select book_id, samples_count, value_sum
from ol_analytics.bucket
where bucket_type = $1 and metric = $2 and (bucket_start = $3 or $1 = 'all')
order by samples_count desc
offset sqlc.arg('skip') limit sqlc.arg('limit');

-- name: Analytics_GetTopBooksByValueSum :many
select book_id, samples_count, value_sum
from ol_analytics.bucket
where bucket_type = $1 and metric = $2 and bucket_start = $3
order by value_sum desc
offset sqlc.arg('skip') limit sqlc.arg('limit');

-- name: Analytics_GetEvents :many
select * 
from ol_analytics.interaction_event
where id > $1
order by id asc
limit $2;

-- name: Analytics_GetWorkerState :one
select *
from ol_analytics.worker_state
where worker_name = $1;

-- name: Analytics_SetWorkerState :exec
insert into ol_analytics.worker_state (worker_name, last_launch, last_cursor)
values ($1, $2, $3)
on conflict (worker_name) do update set
    last_launch = $2,
    last_cursor = $3;