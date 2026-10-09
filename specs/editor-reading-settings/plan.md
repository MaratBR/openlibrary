# Implementation plan

Reuse ChapterReaderSettings and the BookReader Alpine controller inside the editor iframe. Load preferences using the existing ReaderPreferencesService and cookie fallback. Scope a settings button to the editor canvas. No new API or storage schema.

## Verification
- AC-1/AC-2: Template regression test for controls and serialized preferences; public package tests.
- AC-3/AC-4: Browser fixture checks at desktop and mobile widths.
- Frontend build and diff whitespace check.

Browser verification used a rendered template fixture with a mocked preference-save response. Confirmed inherited preferences, font/size updates, unchanged explicitly formatted spans, cookies and API payload, page color/theme, mobile fit, and close/Escape/backdrop dismissal.
