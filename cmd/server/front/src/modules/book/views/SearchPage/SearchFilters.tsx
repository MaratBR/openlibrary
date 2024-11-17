import { useSearchState } from './state'
import './SearchFilters.css'
import { useBookSearchParams } from './search-params'
import { Button } from '@/components/ui/button'
import { ButtonSpinner } from '@/components/spinner'
import { httpGetBookExtremes } from '../../api'
import { useQuery } from '@tanstack/react-query'
import React from 'react'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import { ChevronDown, Filter } from 'lucide-react'
import TagsField from '../../components/tags-field'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { RangeInput } from './RangeInput'

export default function SearchFilters() {
  const [isMobileOpen, setMobileOpen] = React.useState(false)
  const [autoApply, setAutoApply] = React.useState(false)

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
    hasChanges,
  } = useBookSearchParams()

  useFiltersAutoApply(autoApply, applyChanges, hasChanges)

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
          <div className="mb-3 hidden md:flex">
            <Label htmlFor="auto-apply-filters-switch">Auto-apply filters once they change</Label>
            <Switch
              checked={autoApply}
              onCheckedChange={setAutoApply}
              id="auto-apply-filters-switch"
            />
          </div>
          <SearchButton
            runSearch={() => {
              applyChanges()
              setMobileOpen(false)
            }}
          />
        </div>
        <div className="md:space-y-2">
          <ExpandableField label="Include tags">
            <TagsField value={params.include.tags} onChange={setIncludeTags} />
          </ExpandableField>
          <ExpandableField label="Exclude tags">
            <TagsField value={params.exclude.tags} onChange={setExcludeTags} />
          </ExpandableField>

          <ExpandableField label="Chapters">
            <RangeInput value={params.chapters} range={extremes.chapters} onChange={setChapters} />
          </ExpandableField>
          <ExpandableField label="Words">
            <RangeInput value={params.words} range={extremes.words} onChange={setWords} />
          </ExpandableField>
          <ExpandableField label="Words per chapter">
            <RangeInput
              value={params.wordsPerChapter}
              range={extremes.wordsPerChapter}
              onChange={setWordsPerChapter}
            />
          </ExpandableField>
          <ExpandableField label="Favorites">
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

function useFiltersAutoApply(
  enabled: boolean,
  applyChanges: () => void,
  hasChanges: () => boolean,
) {
  React.useEffect(() => {
    if (enabled) {
      const interval = setInterval(() => {
        if (hasChanges()) {
          applyChanges()
        }
      }, 1000)
      return () => clearInterval(interval)
    }
  }, [enabled, applyChanges, hasChanges])
}

function SearchButton({ runSearch }: { runSearch: () => void }) {
  const isLoading = useSearchState((s) => s.loading)

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
        <div className="expandable-field__trigger mx-2 pl-6 py-3 relative">
          <span className="transition-shadow md:group-hover:shadow-[0_2px_0px_currentColor]">
            {label}
          </span>

          <ChevronDown className="absolute top-0 right-3 h-full flex items-center" />
        </div>
      </CollapsibleTrigger>

      <CollapsibleContent className="px-6 py-3">{children}</CollapsibleContent>
    </Collapsible>
  )
}
