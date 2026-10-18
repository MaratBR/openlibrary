# Notification inbox

Extend the existing notification service and reusable menus; retain the existing SSE transport and role eligibility.

## Acceptance criteria

- A1: The bell displays the server-derived unread count across reloads and dashboards. The notification service owns last-read state in `user_data` under `notifications:last-read-at`; clients cannot edit it through generic user-data APIs. Missing state means all stored notifications are unread.
- A2: Opening an inbox marks all notifications through its server snapshot read only after history loads successfully. A failed load or acknowledgement preserves unread state. Newer arrivals remain unread; concurrent or stale acknowledgements cannot move the timestamp backwards or into the future. The UI offers recovery if acknowledging fails.
- A3: Bodies allow bold, italic, underline, strike, sub/sup, line breaks, code and safe links. Other tags, styles, event handlers and dangerous URLs are removed. Titles remain plain text. Stored legacy bodies are sanitized on output too.
- A4: Notifications have an optional urgent flag, category (`system`, `moderation`, `admin`, default `system`), and freeform type (default `info`, bounded to 100 characters). Urgent records have a distinct theme-aware color and slightly larger type, plus an accessible text label. Type has no behavioral effect.
- A5: Moderation and admin use shield and dashboard icons by default; system has none. Explicit emoji, Font Awesome or safe image icons override defaults. Existing metadata remains compatible.
- A6: Both the builder/sink and admin send paths persist the new fields. Admin send controls expose these fields, with translated labels, validation and error states. Menu and moderation history use the same safe record renderer.
- A7: Notification history/status/read APIs require authentication and scope every operation to that account. Regular users can access stored notifications; live SSE remains restricted to its existing roles. Invalid input and storage failures return appropriate errors, without clearing the badge.

No schema migration is required: notification metadata and user_data already exist. Automated notification-producing workflows and per-notification read flags are outside scope.
