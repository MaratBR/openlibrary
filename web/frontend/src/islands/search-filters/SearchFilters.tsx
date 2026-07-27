import { SubmitEvent, useMemo, useState } from 'react'
import TagsInput from '../../components/TagsInput'

import RangeInput from './RangeInput'
import { ReactIslandProps } from '../common/react-island'
import { z } from 'zod'
import {
  DetailedBookSearchQuery,
  getDefaultDetailedBookSearchQuery,
  getQueryParams,
} from '@/api/search'

const dataSchema = z
  .object({
    searchInputId: z.string().optional().nullable(),
  })
  .nullable()
  .optional()

export default function SearchFilters({ data }: ReactIslandProps) {
  const parsedData = useMemo(() => dataSchema.parse(data), [data])
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
        <select
          id="search-sort"
          name="sort"
          className="input w-full"
          value={filters.sort}
          onChange={(event) =>
            setFilters({ ...filters, sort: event.target.value as DetailedBookSearchQuery['sort'] })
          }
        >
          <option value="">{window._('search.sortRelevance')}</option>
          <option value="chapters">{window._('search.sortChapters')}</option>
          <option value="words">{window._('search.sortWords')}</option>
          <option value="words-per-chapter">{window._('search.sortWordsPerChapter')}</option>
          <option value="last-updated">{window._('search.sortLastUpdated')}</option>
          <option value="created-at">{window._('search.sortCreatedAt')}</option>
          <option value="reviews">{window._('search.sortReviews')}</option>
          <option value="readers">{window._('search.sortReaders')}</option>
          <option value="weighted-rating">{window._('search.sortWeightedRating')}</option>
          <option value="random">{window._('search.sortRandom')}</option>
        </select>
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

      <button type="submit" className="btn btn--lg btn--default">
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
