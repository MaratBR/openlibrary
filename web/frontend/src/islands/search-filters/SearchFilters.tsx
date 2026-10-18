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

    const queryParams = getQueryParams({ ...filters, page: 1 })

    const url = new URL(window.location.href)
    if (url.searchParams.get('admin.link') === '1') queryParams.set('admin.link', '1')

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
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-x-6">
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
              <SelectItem value={RELEVANCE_SORT_VALUE}>
                {window._('search.sortRelevance')}
              </SelectItem>
              <SelectItem value="chapters">{window._('search.sortChapters')}</SelectItem>
              <SelectItem value="words">{window._('search.sortWords')}</SelectItem>
              <SelectItem value="words-per-chapter">
                {window._('search.sortWordsPerChapter')}
              </SelectItem>
              <SelectItem value="last-updated">{window._('search.sortLastUpdated')}</SelectItem>
              <SelectItem value="created-at">{window._('search.sortCreatedAt')}</SelectItem>
              <SelectItem value="reviews">{window._('search.sortReviews')}</SelectItem>
              <SelectItem value="readers">{window._('search.sortReaders')}</SelectItem>
              <SelectItem value="weighted-rating">
                {window._('search.sortWeightedRating')}
              </SelectItem>
              <SelectItem value="random">{window._('search.sortRandom')}</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="mb-4">
          <label className="label font-semibold mb-2 text-md">{window._('search.words')}</label>
          <RangeInput
            disableNegative
            label={window._('search.words')}
            value={filters.words}
            onInput={(words) => setFilters({ ...filters, words })}
          />
        </div>
        <div className="mb-4">
          <label className="label font-semibold mb-2 text-md">{window._('search.chapters')}</label>
          <RangeInput
            disableNegative
            label={window._('search.chapters')}
            value={filters.chapters}
            onInput={(chapters) => setFilters({ ...filters, chapters })}
          />
        </div>

        <div className="mb-4">
          <label className="label font-semibold mb-2 text-md">
            {window._('search.wordsPerChapter')}
          </label>
          <RangeInput
            disableNegative
            label={window._('search.wordsPerChapter')}
            value={filters.wordsPerChapter}
            onInput={(wordsPerChapter) => setFilters({ ...filters, wordsPerChapter })}
          />
        </div>

        <div className="mb-4">
          <label className="label font-semibold mb-2 text-md">
            {window._('search.includeTags')}
          </label>
          <TagsInput
            tags={filters.includeTags}
            onInput={(tags) => setFilters({ ...filters, includeTags: tags })}
          />
        </div>

        <div className="mb-4">
          <label className="label font-semibold mb-2 text-md">
            {window._('search.excludeTags')}
          </label>
          <TagsInput
            tags={filters.excludeTags}
            onInput={(tags) => setFilters({ ...filters, excludeTags: tags })}
          />
        </div>

        <div className="mb-4">
          <label htmlFor="search-page-size" className="label font-semibold mb-2 text-md">
            {window._('search.pageSize')}
          </label>
          <select
            id="search-page-size"
            className="input"
            value={[20, 50, 100].includes(filters.pageSize) ? filters.pageSize : 20}
            onChange={(event) => setFilters({ ...filters, pageSize: Number(event.target.value) })}
          >
            {[20, 50, 100].map((size) => (
              <option key={size} value={size}>
                {size}
              </option>
            ))}
          </select>
        </div>
        <div className="mb-4">
          <label htmlFor="search-rating" className="label font-semibold mb-2 text-md">
            {window._('search.ageRating')}
          </label>
          <select
            id="search-rating"
            className="input"
            value={filters.rating}
            onChange={(event) => setFilters({ ...filters, rating: event.target.value })}
          >
            <option value="">{window._('search.anyRating')}</option>
            {['?', 'G', 'PG', 'PG-13', 'R', 'NC-17'].map((rating) => (
              <option key={rating} value={rating}>
                {rating}
              </option>
            ))}
          </select>
        </div>
        {(['views', 'readers', 'reviews', 'popularity', 'weightedRating'] as const).map((key) => (
          <div className="mb-4" key={key}>
            <label className="label font-semibold mb-2 text-md">{window._(`search.${key}`)}</label>
            <RangeInput
              disableNegative
              allowDecimal={key === 'popularity' || key === 'weightedRating'}
              label={window._(`search.${key}`)}
              value={filters[key]}
              onInput={(value) => setFilters({ ...filters, [key]: value })}
            />
          </div>
        ))}
        {(['createdAt', 'updatedAt'] as const).map((key) => (
          <fieldset className="mb-4 min-w-0" key={key}>
            <legend className="label font-semibold mb-2 text-md">
              {window._(`search.${key}`)}
            </legend>
            <div className="grid grid-cols-2 gap-2">
              {(['min', 'max'] as const).map((bound) => (
                <label key={bound} className="min-w-0">
                  <span className="label">
                    {window._(bound === 'min' ? 'search.from' : 'search.to')}
                  </span>
                  <input
                    className="input w-full min-w-0"
                    type="date"
                    value={filters[key][bound]}
                    onChange={(event) =>
                      setFilters({
                        ...filters,
                        [key]: { ...filters[key], [bound]: event.target.value },
                      })
                    }
                  />
                </label>
              ))}
            </div>
          </fieldset>
        ))}
      </div>
      <button type="submit" className="Btn Btn--lg Btn--primary">
        {window._('search.doSearch')}
      </button>
    </form>
  )
}

function getDetailedBookSearchQuery(): DetailedBookSearchQuery {
  const el = document.getElementById('data-search-explained-query')
  if (el instanceof HTMLTemplateElement) {
    try {
      const query = JSON.parse(el.content.textContent || '') as DetailedBookSearchQuery
      return { ...query, pageSize: [20, 50, 100].includes(query.pageSize) ? query.pageSize : 20 }
    } catch {
      // no-op
    }
  }

  return getDefaultDetailedBookSearchQuery()
}
