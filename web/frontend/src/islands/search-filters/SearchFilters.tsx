import { useState } from 'preact/hooks'
import {
  DetailedBookSearchQuery,
  detailedBookSearchQuerySchema,
  getDefaultDetailedBookSearchQuery,
  getQueryParams,
} from './api'
import TagsInput from './TagsInput'
import { _ } from '@/common/i18n'
import RangeInput from './RangeInput'

export default function SearchFilters() {
  const [filters, setFilters] = useState<DetailedBookSearchQuery>(getDetailedBookSearchQuery)

  function handleSubmit(event: SubmitEvent) {
    event.preventDefault()

    const queryParams = getQueryParams(filters)

    const url = new URL(window.location.href)
    url.search = queryParams.toString()
    window.location.href = url.href
  }

  return (
    <form onSubmit={handleSubmit}>
      <div class="mb-4">
        <label class="ol-label font-semibold mb-2 text-md">{_('search.words')}</label>
        <RangeInput
          disableNegative
          value={filters.words}
          onInput={(words) => setFilters({ ...filters, words })}
        />
      </div>
      <div class="mb-4">
        <label class="ol-label font-semibold mb-2 text-md">{_('search.chapters')}</label>
        <RangeInput
          disableNegative
          value={filters.chapters}
          onInput={(chapters) => setFilters({ ...filters, chapters })}
        />
      </div>

      <div class="mb-4">
        <label class="ol-label font-semibold mb-2 text-md">{_('search.chapters')}</label>
        <RangeInput
          disableNegative
          value={filters.wordsPerChapter}
          onInput={(wordsPerChapter) => setFilters({ ...filters, wordsPerChapter })}
        />
      </div>

      <div class="mb-4">
        <label class="ol-label font-semibold mb-2 text-md">{_('search.includeTags')}</label>
        <TagsInput
          tags={filters.includeTags}
          onInput={(tags) => setFilters({ ...filters, includeTags: tags })}
        />
      </div>

      <div class="mb-4">
        <label class="ol-label font-semibold mb-2 text-md">{_('search.excludeTags')}</label>
        <TagsInput
          tags={filters.excludeTags}
          onInput={(tags) => setFilters({ ...filters, excludeTags: tags })}
        />
      </div>

      <button type="submit" class="ol-btn ol-btn--lg ol-btn--primary rounded-full">
        {_('search.doSearch')}
      </button>
    </form>
  )
}

function getDetailedBookSearchQuery(): DetailedBookSearchQuery {
  const el = document.getElementById('data-search-explained-query')
  if (el instanceof HTMLTemplateElement) {
    try {
      const parsed = JSON.parse(el.content.textContent || '')
      return detailedBookSearchQuerySchema.parse(parsed)
    } catch {
      // no-op
    }
  }

  return getDefaultDetailedBookSearchQuery()
}
