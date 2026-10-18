# Implementation plan

1. Extend the notification domain with category, urgency, typed icon metadata and a dedicated inline HTML sanitizer. Keep existing metadata representations compatible; normalize both builder and broadcast writes, and sanitize legacy rows on read.
2. Extend the notification service abstraction with status and monotonic mark-read methods. SQLC queries count records newer than the stored JSON timestamp, expose a database snapshot time, and atomically clamp/advance read timestamps in user_data. Keep the key outside the generic user-data whitelist.
3. Add authenticated status/read endpoints alongside history. Admin send accepts presentation options. Keep SSE eligibility unchanged.
4. Share a notification record renderer between the reusable menu and moderation history. Fetch authoritative status on mount, refresh on focus and stream delivery, and acknowledge the snapshot only after a successful initial history load. Handle stale responses, duplicate events, concurrent requests and failed acknowledgements. Expose admin fields and English catalog keys.
5. Generate SQLC, test service validation/sanitization/metadata/read transitions and HTTP recipient isolation, build the frontend, and check the diff. Use available local DB integration coverage for timestamp SQL; document unavailable checks.
