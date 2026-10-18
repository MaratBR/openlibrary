# Book title highlighting

Search results emphasize words in book names matched by OpenSearch, including fuzzy matches. Full titles and links remain intact. Browsing without a query and results without highlights display ordinary titles. Book names containing HTML or highlight-like text must remain safe literal text.

Acceptance criteria:
- AC1: Query searches request full-name highlights and carry them through cached results.
- AC2: Search result titles render matched portions as semantic marks; other book listings remain unchanged.
- AC3: Missing, malformed, or inconsistent highlights fall back to the original title; title content cannot introduce HTML.
