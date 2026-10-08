# Tasks: Adult-content and warning preferences

Classification boundary and service abstraction agreed on 2026-10-08.
Broader filtering requirements remain provisional.

- [x] Record NC-17 plus adult tags behind a replaceable application service.
- [x] Implement and verify the classification service and custom-provider behavior.
- [x] Remove the independent editor switch and verify legacy rejection before saving.

- [ ] Resolve open questions and finalize observable requirements (AC-1–12).
- [ ] Inventory public surfaces, APIs, caches, and availability checks (AC-4–8, AC-11).
- [ ] Implement central classification/preferences policy and boundary tests (AC-3, AC-8, AC-11).
- [x] Fix editor classification controls and legacy API contract (AC-1; browser verification pending).
- [ ] Implement canonical descriptor storage/transport and migration fixtures (AC-2, AC-9–10).
- [ ] Implement settings selector, validation, translations, and save/error states (AC-2, AC-12).
- [ ] Add SQL/OpenSearch constraints, count parity, index backfill, and invalidation (AC-4, AC-10).
- [ ] Apply policy to random selection and remaining discovery/API surfaces (AC-4).
- [ ] Add server-rendered book/chapter/preview interstitials and safe back action (AC-5, AC-12).
- [ ] Implement scoped session exceptions and retire global cookie bypass (AC-6, AC-9, AC-11).
- [ ] Apply saved-item placeholders and verify hidden payload handling (AC-5, AC-7).
- [ ] Verify preference precedence, failure behavior, and cache isolation (AC-8).
- [ ] Verify tag lifecycle, ID precision, and moderation conflict handling (AC-3, AC-10).
- [ ] Validate management authorization and public availability restrictions (AC-11).
- [ ] Complete browser accessibility, theme, mobile, and no-JavaScript checks (AC-12).
- [ ] Document rollout/rollback and run proportional repository checks; record evidence for every criterion (AC-1–12).
