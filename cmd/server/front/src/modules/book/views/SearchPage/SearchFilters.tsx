import './SearchFilters.css'
import { useBookSearchParams, useSearchState } from './search-params'
import { Button } from '@/components/ui/button'
import { ButtonSpinner } from '@/components/spinner'
import { httpGetBookExtremes } from '../../api/api'
import { useQuery } from '@tanstack/react-query'
import React from 'react'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import { ChevronDown, Filter } from 'lucide-react'
import TagsField from '../../components/tags-field'
import { RangeInput } from './RangeInput'
import { useTranslation } from 'react-i18next'

export default function SearchFilters() {
  const { t } = useTranslation()
  const [isMobileOpen, setMobileOpen] = React.useState(false)

  const extremes = useSearchState((s) => s.extremes)
  const {
    params,
    setChapters,
    setFavorites,
    setWords,
    setWordsPerChapter,
    setExcludeTags,
    setIncludeTags,
    applyChanges,
  } = useBookSearchParams()

  useQuery({
    queryKey: ['search', 'extremes'],
    queryFn: async () => {
      const result = await httpGetBookExtremes()
      useSearchState.getState().setExtremes({
        chapters: result.chapters,
        words: result.words,
        wordsPerChapter: result.wordsPerChapter,
        favorites: result.favorites,
      })
      return result
    },
  })

  // if (!initialized) return

  return (
    <>
      <Button
        variant="secondary"
        className="h-14 w-full md:hidden"
        onClick={() => setMobileOpen(true)}
      >
        <Filter /> Filters
      </Button>
      <div
        className="
        inset-0 fixed bg-background z-50 overflow-auto data-[mobile-open=false]:invisible 
        md:!visible md:static md:z-auto"
        data-mobile-open={isMobileOpen}
      >
        <div className="flex flex-col gap-2 my-3 px-3">
          <Button
            variant="secondary"
            className="h-14 md:hidden"
            onClick={() => setMobileOpen(false)}
          >
            Close
          </Button>

          <SearchButton
            runSearch={() => {
              applyChanges()
              setMobileOpen(false)
            }}
          />
        </div>
        <div className="md:space-y-2 pb-4 md:pb-0">
          <ExpandableField label={t('search.includeTags')}>
            <TagsField value={params.include.tags} onChange={setIncludeTags} />
          </ExpandableField>
          <ExpandableField label={t('search.excludeTags')}>
            <TagsField value={params.exclude.tags} onChange={setExcludeTags} />
          </ExpandableField>

          <ExpandableField label={t('search.chapters')}>
            <RangeInput value={params.chapters} range={extremes.chapters} onChange={setChapters} />
          </ExpandableField>
          <ExpandableField label={t('search.words')}>
            <RangeInput value={params.words} range={extremes.words} onChange={setWords} />
          </ExpandableField>
          <ExpandableField label={t('search.wordsPerChapter')}>
            <RangeInput
              value={params.wordsPerChapter}
              range={extremes.wordsPerChapter}
              onChange={setWordsPerChapter}
            />
          </ExpandableField>
          <ExpandableField label={t('search.favorites')}>
            <RangeInput
              value={params.favorites}
              range={extremes.favorites}
              onChange={setFavorites}
            />
          </ExpandableField>
        </div>
      </div>
    </>
  )
}

function SearchButton({ runSearch }: { runSearch: () => void }) {
  const isLoading = useSearchState((s) => s.isLoading)

  return (
    <Button className="w-full h-14 md:h-10" onClick={() => runSearch()}>
      {isLoading ? <ButtonSpinner /> : 'Search'}
    </Button>
  )
}

function ExpandableField({ children, label }: React.PropsWithChildren<{ label: string }>) {
  return (
    <Collapsible defaultOpen={false} className="expandable-field">
      <CollapsibleTrigger asChild>
        <div className="expandable-field__trigger">
          <span className="transition-shadow ">{label}</span>

          <ChevronDown className="absolute top-0 right-3 h-full flex items-center" />
        </div>
      </CollapsibleTrigger>

      <CollapsibleContent className="expandable-field__content">
        <div className="px-6 py-3">{children}</div>
      </CollapsibleContent>
    </Collapsible>
  )
}
