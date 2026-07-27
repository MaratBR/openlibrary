-- name: Analytics_InsertEvent :copyfrom
insert into ol_analytics.interaction_event 
    (user_key, book_id, event_type, value)
    values
    ($1, $2, $3, $4);

-- name: Analytics_IncrementBucketCounters :exec
WITH 
    input_rows AS (
        SELECT
            (sqlc.arg('book_ids')::bigint[])[i] AS book_id,
            (sqlc.arg('metrics')::ol_analytics.counter_type[])[i] AS metric,
            (sqlc.arg('values')::bigint[])[i] AS increment_value
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
            date_trunc('year', now(), 'UTC')

        UNION ALL

        SELECT
            'month'::ol_analytics.bucket_period_type,
            date_trunc('month', now(), 'UTC')

        UNION ALL

        SELECT
            'week'::ol_analytics.bucket_period_type,
            date_trunc('week', now(), 'UTC')

        UNION ALL

        SELECT
            'day'::ol_analytics.bucket_period_type,
            date_trunc('day', now(), 'UTC')
    ),
    aggregated_buckets AS (
        SELECT
            input.book_id,
            input.metric,
            buckets.bucket_type,
            buckets.bucket_start,
            sum(input.increment_value) AS increment_value
        FROM input_rows AS input
        CROSS JOIN bucket_definitions AS buckets
        GROUP BY
            input.book_id,
            input.metric,
            buckets.bucket_type,
            buckets.bucket_start
    )
INSERT INTO ol_analytics.bucket_counter (
    book_id,
    metric,
    bucket_type,
    bucket_start,
    value
)
SELECT
    book_id,
    metric,
    bucket_type,
    bucket_start,
    increment_value
FROM aggregated_buckets
ON CONFLICT (
    book_id,
    metric,
    bucket_type,
    bucket_start
)
DO UPDATE SET
    value = ol_analytics.bucket_counter.value + EXCLUDED.value,
    updated_at = now();


-- name: Analytics_UpdatePopularity :exec
WITH 
    params AS (
        SELECT statement_timestamp() AS run_at
    ),
    expanded_events AS (
        SELECT
            event.book_id,
            event.event_type,
            event.created_at,
            event_bucket.bucket_type,
            event_bucket.bucket_start
        FROM ol_analytics.interaction_event AS event
        CROSS JOIN LATERAL (
            SELECT
                'all'::ol_analytics.bucket_period_type AS bucket_type,
                TIMESTAMPTZ '1970-01-01 00:00:00+00' AS bucket_start

            UNION ALL

            SELECT
                'year'::ol_analytics.bucket_period_type,
                date_trunc('year', event.created_at AT TIME ZONE 'UTC') AT TIME ZONE 'UTC'

            UNION ALL

            SELECT
                'month'::ol_analytics.bucket_period_type,
                date_trunc('month', event.created_at AT TIME ZONE 'UTC') AT TIME ZONE 'UTC'

            UNION ALL

            SELECT
                'week'::ol_analytics.bucket_period_type,
                date_trunc('week', event.created_at AT TIME ZONE 'UTC') AT TIME ZONE 'UTC'

            UNION ALL

            SELECT
                'day'::ol_analytics.bucket_period_type,
                date_trunc('day', event.created_at AT TIME ZONE 'UTC') AT TIME ZONE 'UTC'
        ) AS event_bucket
    ),
    unprocessed_events AS (
        SELECT
            event.*,
            existing.updated_at AS previous_updated_at
        FROM expanded_events AS event
        LEFT JOIN ol_analytics.bucket_popularity AS existing
            ON existing.book_id = event.book_id
        AND existing.bucket_type = event.bucket_type
        AND existing.bucket_start = event.bucket_start
        WHERE event.created_at > COALESCE(
            existing.updated_at,
            TIMESTAMPTZ '-infinity'
        )
    ),
    new_scores AS (
        SELECT
            event.book_id,
            event.bucket_type,
            event.bucket_start,
            SUM(
                CASE event.event_type
                    WHEN 'book_view'
                        THEN sqlc.arg('book_view_score')::double precision
                    WHEN 'search_click'
                        THEN sqlc.arg('search_click_score')::double precision
                    WHEN 'chapter_view'
                        THEN sqlc.arg('chapter_view_score')::double precision
                    WHEN 'started_reading'
                        THEN sqlc.arg('started_reading_score')::double precision
                    WHEN 'completed'
                        THEN sqlc.arg('completed_score')::double precision
                    WHEN 'dropped'
                        THEN sqlc.arg('dropped_score')::double precision
                    WHEN 'finished_chapter'
                        THEN sqlc.arg('finished_chapter_score')::double precision
                    ELSE 0
                END
                *
                POWER(
                    0.5,
                    EXTRACT(
                        EPOCH FROM params.run_at - event.created_at
                    )
                    / sqlc.arg('half_life_seconds')::double precision
                )
            ) AS value
        FROM unprocessed_events AS event
        CROSS JOIN params
        GROUP BY
            event.book_id,
            event.bucket_type,
            event.bucket_start
    )
INSERT INTO ol_analytics.bucket_popularity (
    book_id,
    bucket_type,
    bucket_start,
    value,
    updated_at
)
SELECT
    score.book_id,
    score.bucket_type,
    score.bucket_start,
    score.value,
    params.run_at
FROM new_scores AS score
CROSS JOIN params
ON CONFLICT (
    book_id,
    bucket_type,
    bucket_start
)
DO UPDATE SET
    value =
        ol_analytics.bucket_popularity.value
        * POWER(
            0.5,
            EXTRACT(
                EPOCH FROM (
                    EXCLUDED.updated_at
                    - ol_analytics.bucket_popularity.updated_at
                )
            )
            / sqlc.arg('half_life_seconds')::double precision
        )
        + EXCLUDED.value,
    updated_at = EXCLUDED.updated_at;


