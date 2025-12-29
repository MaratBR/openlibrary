import { ErrorDisplay } from '@/components/error'
import { PreactIslandProps } from '@/lib/island'
import { useEffect, useMemo, useState } from 'preact/hooks'
import z from 'zod'
import ChaptersList from './ChaptersList'
import { useBookChaptersState } from './state'
import { useShallow } from 'zustand/shallow'
import { ChapterSelectorPopupProvider } from './ChapterSelectorPopup'
import { httpCreateChapter } from '@/api/bm'
import { useMutation } from '@tanstack/react-query'

const schema = z.object({
  bookId: z.string(),
})

export function BookChapters({ data: dataUnknown }: PreactIslandProps) {
  const data = useMemo(() => schema.parse(dataUnknown), [dataUnknown])

  const { showLoader, error, isEmpty } = useBookChaptersState(
    useShallow((s) => ({
      showLoader: s.loading,
      isEmpty: !s.loading && s.chapters.length === 0,
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
      ) : isEmpty ? (
        <NoChaptersYet />
      ) : (
        <ChapterSelectorPopupProvider>
          <ChaptersList />
        </ChapterSelectorPopupProvider>
      )}
      {error && <ErrorDisplay error={error} />}
    </div>
  )
}

function NoChaptersYet() {
  const [showInput, setShowInput] = useState(false)

  const bookId = useBookChaptersState((s) => s.bookId)

  const [name, setName] = useState('')

  const createChapterMutation = useMutation({
    mutationFn: async () => {
      if (name.trim().length === 0) return
      const response = await httpCreateChapter(bookId, {
        name,
        summary: '',
        isAdultOverride: false,
        content: '',
      })

      location.href = `/books-manager/book/${bookId}/chapter/${response.data}?first=1`
    },
  })

  return (
    <div class="py-5 px-4">
      <p class="mb-4">{window._('bookManager.edit.noChapters')}</p>
      {showInput ? (
        <>
          <input
            class="input w-96"
            name="name"
            value={name}
            onChange={(e) => setName((e.target as HTMLInputElement).value)}
          />
          <button
            class="btn btn--ghost"
            onClick={() => createChapterMutation.mutate()}
            disabled={name.trim().length === 0 || createChapterMutation.isPending}
          >
            <i class="fa-solid fa-arrow-right" />
          </button>
        </>
      ) : (
        <button onClick={() => setShowInput(true)} class="btn">
          {window._('bookManager.edit.addChapter')}
        </button>
      )}
    </div>
  )
}
