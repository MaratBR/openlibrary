import { SubmitEvent, useMemo, useState } from 'react'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components'
import TagsInput from '@/components/TagsInput'

import RangeInput from './RangeInput'
import { ReactIslandProps } from '../common/react-island'
import { Schema } from 'effect'
import {
  DetailedBookSearchQuery,
  getDefaultDetailedBookSearchQuery,
  getQueryParams,
} from '@/features/search'

const dataSchema = Schema.UndefinedOr(
  Schema.NullOr(
    Schema.Struct({
      searchInputId: Schema.optional(Schema.NullOr(Schema.String)),
    }),
  ),
)

const RELEVANCE_SORT_VALUE = 'relevance'

export default function SearchFilters({ data }: ReactIslandProps) {
  const parsedData = useMemo(() => Schema.decodeUnknownSync(dataSchema)(data), [data])
  const [filters, setFilters] = useState<DetailedBookSearchQuery>(getDetailedBookSearchQuery)

  function handleSubmit(event: SubmitEvent) {
    event.preventDefault()

    const queryParams = getQueryParams(filters)

    const url = new URL(window.location.href)

    if (parsedData?.searchInputId) {
      const input = document.getElementById(parsedData.searchInputId)
      if (input instanceof HTMLInputElement) {
        const value = input.value.trim()
        if (value) queryParams.set('q', input.value)
      }
    }

    url.search = queryParams.toString()
    window.location.href = url.href
  }

  return (
    <form onSubmit={handleSubmit}>
      <div className="mb-4">
        <label htmlFor="search-sort" className="label font-semibold mb-2 text-md">
          {window._('search.sortBy')}
        </label>
        <Select
          value={filters.sort || RELEVANCE_SORT_VALUE}
          onValueChange={(value) =>
            setFilters({
              ...filters,
              sort:
                value === RELEVANCE_SORT_VALUE ? '' : (value as DetailedBookSearchQuery['sort']),
            })
          }
        >
          <SelectTrigger id="search-sort">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value={RELEVANCE_SORT_VALUE}>{window._('search.sortRelevance')}</SelectItem>
            <SelectItem value="chapters">{window._('search.sortChapters')}</SelectItem>
            <SelectItem value="words">{window._('search.sortWords')}</SelectItem>
            <SelectItem value="words-per-chapter">
              {window._('search.sortWordsPerChapter')}
            </SelectItem>
            <SelectItem value="last-updated">{window._('search.sortLastUpdated')}</SelectItem>
            <SelectItem value="created-at">{window._('search.sortCreatedAt')}</SelectItem>
            <SelectItem value="reviews">{window._('search.sortReviews')}</SelectItem>
            <SelectItem value="readers">{window._('search.sortReaders')}</SelectItem>
            <SelectItem value="weighted-rating">{window._('search.sortWeightedRating')}</SelectItem>
            <SelectItem value="random">{window._('search.sortRandom')}</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="mb-4">
        <label className="label font-semibold mb-2 text-md">{window._('search.words')}</label>
        <RangeInput
          disableNegative
          value={filters.words}
          onInput={(words) => setFilters({ ...filters, words })}
        />
      </div>
      <div className="mb-4">
        <label className="label font-semibold mb-2 text-md">{window._('search.chapters')}</label>
        <RangeInput
          disableNegative
          value={filters.chapters}
          onInput={(chapters) => setFilters({ ...filters, chapters })}
        />
      </div>

      <div className="mb-4">
        <label className="label font-semibold mb-2 text-md">{window._('search.chapters')}</label>
        <RangeInput
          disableNegative
          value={filters.wordsPerChapter}
          onInput={(wordsPerChapter) => setFilters({ ...filters, wordsPerChapter })}
        />
      </div>

      <div className="mb-4">
        <label className="label font-semibold mb-2 text-md">{window._('search.includeTags')}</label>
        <TagsInput
          tags={filters.includeTags}
          onInput={(tags) => setFilters({ ...filters, includeTags: tags })}
        />
      </div>

      <div className="mb-4">
        <label className="label font-semibold mb-2 text-md">{window._('search.excludeTags')}</label>
        <TagsInput
          tags={filters.excludeTags}
          onInput={(tags) => setFilters({ ...filters, excludeTags: tags })}
        />
      </div>

      <button type="submit" className="btn btn--lg btn--primary">
        {window._('search.doSearch')}
      </button>
    </form>
  )
}

function getDetailedBookSearchQuery(): DetailedBookSearchQuery {
  const el = document.getElementById('data-search-explained-query')
  if (el instanceof HTMLTemplateElement) {
    try {
      return JSON.parse(el.content.textContent || '') as DetailedBookSearchQuery
    } catch {
      // no-op
    }
  }

  return getDefaultDetailedBookSearchQuery()
}
