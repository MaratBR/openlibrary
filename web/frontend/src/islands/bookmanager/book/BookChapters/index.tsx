import { ErrorDisplay } from '@/components/error'
import { PreactIslandProps } from '@/lib/island'
import { useEffect, useMemo } from 'preact/hooks'
import z from 'zod'
import ChaptersList from './ChaptersList'
import { useBookChaptersState } from './state'
import { useShallow } from 'zustand/shallow'
import { ChapterSelectorPopupProvider } from './ChapterSelectorPopup'

const schema = z.object({
  bookId: z.string(),
})

export function BookChapters({ data: dataUnknown }: PreactIslandProps) {
  const data = useMemo(() => schema.parse(dataUnknown), [dataUnknown])

  const { showLoader, error } = useBookChaptersState(
    useShallow((s) => ({
      showLoader: s.loading && s.chapters.length === 0,
      error: s.error,
    })),
  )

  useEffect(() => {
    useBookChaptersState.getState().loadChapters(data.bookId)
  }, [data.bookId])

  return (
    <div>
      {showLoader ? (
        <div style="height:120px" class="flex items-center justify-center">
          <span class="loader" />
        </div>
      ) : (
        <ChapterSelectorPopupProvider>
          <ChaptersList />
        </ChapterSelectorPopupProvider>
      )}
      {error && <ErrorDisplay error={error} />}
    </div>
  )
}
