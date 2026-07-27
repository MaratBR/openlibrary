• The repository has enough data to support useful book discovery, but discovery is currently centered on basic search plus all-time views. The biggest opportunity is converting reading behavior into trustworthy, explainable signals
  such as “trending,” “often completed,” and “readers with similar tastes liked this.”

  ## What exists today

  The site already captures or models:

  - Genres, tags, warnings, adult/spoiler flags, and tag hierarchies.
  - Book length, chapter count, and average chapter length.
  - Ratings, rating count, reviews, and review distribution.
  - Book and chapter views over hour/day/week/month/year/all-time periods.
  - Reader states: want to read, reading, paused, completed, and DNF.
  - Chapter progress and finished-chapter history.
  - User-created collections.
  - Authors and author follows.
  - Book publication and chapter timestamps.

  Relevant implementation areas include internal/app/search.go, internal/app/analytics/analytics_views.go, internal/app/reading_list.go, and internal/store/query.reviews.sql.

  However:

  - Search has no explicit sorting or discovery ranking.
  - Search cards only show title, author, tags, and summary.
  - The homepage only has an all-time “most viewed” carousel.
  - Book pages show rating, review count, and total views, but no momentum, completion, update cadence, or reader-fit information.
  - Book lifecycle status such as ongoing/completed/hiatus is not modeled.
  - Personalization and book-to-book similarity do not exist.

  ## Recommended reader-facing analytics

  ### 1. Popular now

  Use recent qualified reading activity instead of raw all-time book-page views.

  Show:

  - Trending today
  - Trending this week
  - Popular this month
  - All-time favorites

  A reasonable first ranking formula:

  trend_score =
      log(1 + qualified_chapter_readers_7d)
    + 0.5 × log(1 + want_to_read_adds_7d)
    + 0.7 × log(1 + readers_started_7d)
    + rating_confidence
    - age_decay

  “Qualified reader” should mean a unique reader who opened a chapter and spent enough time or made enough progress to distinguish reading from an accidental page load.

  This avoids letting old books permanently dominate the homepage.

  ### 2. Rising stories

  This should identify acceleration rather than absolute popularity:

  growth = qualified_readers_last_7d / max(qualified_readers_previous_7d, baseline)

  Require a minimum sample size and combine growth with absolute reader count. Otherwise a book growing from one reader to three would outrank meaningful successes.

  Useful shelf labels:

  - Rising this week
  - Hidden gems
  - New and gaining readers
  - Recently updated and trending

  ### 3. Reader commitment funnel

  For every book, aggregate:

  impression → book page → first chapter → third chapter
             → 25% → 50% → latest/final chapter

  The useful public signals are:

  - Readers who continued past chapter 1
  - Readers who reached chapter 3
  - Estimated completion rate
  - DNF rate
  - Typical stopping point

  Do not show every number prominently. A reader-friendly presentation would be:

  - “82% continue after chapter 1”
  - “Frequently completed”
  - “Readers usually know whether it suits them within 3 chapters”

  These signals are much more informative than views because they describe whether the story holds attention.

  For ongoing serials, calculate “caught up to the latest chapter” separately from “completed.”

  ### 4. Rating quality with confidence

  A raw average rating is unreliable for books with few votes. Rank using a Bayesian or Wilson-style confidence score:

  weighted_rating =
      (votes / (votes + minimum_votes)) × book_average
    + (minimum_votes / (votes + minimum_votes)) × site_average

  Expose both average and sample size:

  - 8.7 · 1,240 ratings
  - Rating distribution
  - “Highly rated by progression-fantasy readers”
  - “Readers who completed it rated it 9.1”

  Keep the ordinary average visible, but use the confidence-adjusted value for rankings.

  ### 5. Book similarity

  Create “More like this” using progressively better signals:

  1. Shared genres and tags.
  2. Shared appearances in public collections.
  3. Readers who meaningfully read both books.
  4. Similar completion and rating behavior.
  5. Text embeddings of summaries later, if needed.

  A simple initial score could combine:

  similarity =
      tag_jaccard
    + weighted_shared_collections
    + normalized_reader_overlap

  Explain recommendations:

  - “Similar tags: academy, progression, weak-to-strong”
  - “Frequently collected with this book”
  - “Readers of this also finished…”

  The explanation is important because tags may include warnings or incidental story elements.

  ### 6. Personalized recommendations

  For signed-in users, build a taste profile from:

  - Completed and highly rated books: strong positive signal.
  - Meaningful chapter progress: positive signal.
  - Want-to-read: weak positive signal.
  - Paused: weak negative or neutral signal.
  - DNF: negative signal.
  - Explicit excluded/censored tags: hard exclusion.
  - Followed authors: positive signal.

  Recommended homepage shelves:

  - Because you finished…
  - More from tags you read
  - Continue with similar stories
  - New from authors you follow
  - Outside your usual genres
  - Short reads you may like

  Do not treat a click or book-page view as strong preference. Reading depth and explicit ratings are much better signals.

  ### 7. Taste-match explanation

  Each recommended card should answer “Why am I seeing this?”

  Examples:

  - “Matches 4 of your favorite tags”
  - “Popular with readers who finished X”
  - “Short chapters, completed story”
  - “Does not contain your excluded tags”
  - “New release from an author you follow”

  This makes personalization feel helpful rather than opaque.

  ### 8. Update health and commitment

  Serial-fiction readers care heavily about abandonment risk. Add explicit book status:

  - Ongoing
  - Completed
  - Hiatus
  - Stub
  - Dropped

  Then derive:

  - Last updated
  - Chapters published in the last 30/90 days
  - Median days between chapters
  - Update consistency
  - Estimated reading time
  - Estimated time to catch up

  User-facing examples:

  - “Updates about twice a week”
  - “38 chapters · ~7 hours”
  - “Completed”
  - “No update for 5 months”

  This requires adding book lifecycle status; it cannot be inferred reliably from timestamps alone.

  ### 9. Tag and preference analytics

  Tags should become discovery facets, not just labels.

  Useful signals:

  - Popular books within a tag
  - Rising books within a tag
  - Highly completed books within a tag
  - Commonly paired tags
  - Related tags
  - Books matching all selected tags
  - “Exclude” suggestions based on warnings

  The existing include/exclude tag system is a strong foundation. Add result counts beside filter choices and allow sorting within the filtered set.

  ### 10. Community curation

  Public collections can power:

  - Most saved to collections
  - Curated lists containing this book
  - Frequently collected together
  - Popular collections for a genre
  - Lists from readers with similar tastes

  Collection signals should be weighted to resist spam:

  - Count unique collection owners.
  - Give less weight to enormous collections.
  - Ignore private collections.
  - Require minimum account/activity quality where abuse becomes a problem.

  ## Search improvements

  The most urgent UX issue is that search lacks sorting. Add:

  - Best match
  - Trending
  - Most read
  - Highest rated
  - Most completed
  - Recently updated
  - Newest
  - Shortest/longest
  - Most chapters

  Search cards should display the information needed to choose without opening every result:

  - Rating and rating count
  - Book status
  - Total length and chapter count
  - Last update or update cadence
  - One contextual signal based on the active sort
  - Matched and warning tags

  For example, a “Trending” result should show weekly readers or growth, while “Highest rated” should show its rating and vote count.

  ## Privacy and integrity rules

  Reader-facing analytics should be aggregate and anonymous:

  - Never reveal an individual’s reading or DNF behavior without explicit consent.
  - Suppress percentages below a minimum reader count.
  - Separate anonymous sessions from authenticated readers carefully.
  - Deduplicate bots, repeated refreshes, and author self-traffic.
  - Do not count a chapter request as a completed read.
  - Keep adult/censored-tag preferences as hard recommendation constraints.
  - Allow users to disable personalization or reset their taste profile.

  ## Repository issues to resolve first

  A few implementation details currently limit reliable discovery analytics:

  - IncrChapterView writes to counterBooks instead of counterChapters; its composite bookID,chapterID key will later fail integer parsing during the book-counter commit. Chapter analytics therefore appear incomplete or broken in
    internal/app/analytics/analytics_views_impl.go:105.

  - Search has no sort field and the SQL search query has no ORDER BY, so stable or meaningful discovery order is not guaranteed in internal/store/book_search.go:79.
  - Favorites remains in search-extreme types and SQL, but no corresponding migrated books.favorites column exists. This appears to be stale unfinished functionality.
  - “Most viewed books in 2026” is hard-coded while the controller requests the all-time period in web/public/templates/home.templ:67 and web/public/routes_home.go:43.
  - Chapter history stores progress and finished_reading, but the current service mainly updates the last chapter; the behavioral funnel needs consistent progress events.
  - Reading-list history lacks timestamps, which makes retention and time-windowed completion analysis difficult.

  ## Suggested implementation order

  1. Fix chapter analytics and define qualified reading events.
  2. Add book lifecycle status and update-cadence fields.
  3. Add search sorting and richer result-card statistics.
  4. Launch Trending, Rising, Recently Updated, and Confidence-Rated shelves.
  5. Add start/continue/catch-up/completion/DNF aggregates.
  6. Build tag- and collection-based “More like this.”
  7. Add personalized recommendations after enough behavior exists.
  8. Add controlled exploration so new books can receive exposure.

  The best first release would be: weekly trending, rising stories, recently updated, confidence-adjusted top rated, explicit story status, and richer search sorting. Those provide substantial discovery value before requiring a
  sophisticated recommendation system.