# Search filters

Users can select 20, 50 or 100 results per page (default 20). HTTP page sizes are clamped to 1–100; missing or invalid values default to 20.

Search supports optional age rating, inclusive creation/update date bounds, and nonnegative ranges for views, readers, reviews, popularity and weighted rating from the existing index. Blank or invalid bounds are ignored. Public visibility and trash restrictions remain enforced. Filters round-trip through URLs and pagination; applying filters starts at page one. Existing filters remain available.
