import { useMemo, useRef, useState } from 'preact/hooks'
import { httpUpdateChaptersOrder, managerBookDetailsSchema } from './api'
import { ErrorDisplay } from '@/components/error'
import Popper from '@/components/Popper'
import { useMutation } from '@tanstack/react-query'
import { PreactIslandProps } from '@/lib/island'
import { httpCreateChapter } from '@/api/bm'

// type Chapter = ManagerBookDetailsDto['chapters'][number]

export default function Chapters({ data }: PreactIslandProps) {
  const book = useMemo(() => managerBookDetailsSchema.parse(data), [data])
  const [chapters, setChapters] = useState(book.chapters || [])
  const [originalOrder, setOriginalOrder] = useState<string[]>([])
  const [isReordering, setIsReordering] = useState(false)
  const [isSavingOrder, setSavingOrder] = useState(false)
  const [savingOrderError, setSavingOrderError] = useState<unknown>()
  const [openChapterName, setOpenChapterName] = useState(false)
  const [chapterName, setChapterName] = useState('')
  const validChapterName = chapterName.trim().length > 0

  const addChapterButton = useRef<HTMLButtonElement | null>(null)

  const handleStartReordering = () => {
    setOriginalOrder(chapters.map((x) => x.id))
    setIsReordering(true)
  }

  const handleCancelReordering = () => {
    setIsReordering(false)
    setChapters(originalOrder.map((id) => chapters.find((x) => x.id === id)!))
  }

  const handleSaveOrder = () => {
    setSavingOrder(true)

    const newOrder = chapters.map((x) => x.id)

    httpUpdateChaptersOrder(book.id, newOrder)
      .then(() => {
        setSavingOrderError(undefined)
        setIsReordering(false)
        setOriginalOrder(newOrder)
        setChapters((chapters) =>
          chapters.map((c, index) => ({
            ...c,
            order: index + 1,
          })),
        )
      })
      .catch((error) => setSavingOrderError(error))
      .finally(() => {
        setSavingOrder(false)
      })
  }

  const createChapterMutation = useMutation({
    mutationFn: async () => {
      if (!validChapterName) {
        throw new Error(`invalid chapter name: ${chapterName}`)
      }

      const response = await httpCreateChapter(book.id, {
        name: chapterName,
        content: '',
        isAdultOverride: false,
        summary: '',
      })

      const chapterId = response.data

      window.location.href = `/books-manager/book/${book.id}/chapter/${chapterId}`
    },
  })

  const numberOfUpdates = useMemo(() => {
    if (!isReordering) return 0

    let count = 0
    for (let i = 0; i < chapters.length; i++) {
      const chapter = chapters[i]
      if (i !== originalOrder.indexOf(chapter.id)) {
        count++
      }
    }
    return count
  }, [isReordering, chapters, originalOrder])

  return (
    <>
      <Popper open={openChapterName} anchorEl={addChapterButton} placement="bottom-end">
        <div class="card shadow-md mt-1 flex items-center gap-1">
          <input
            value={chapterName}
            onInput={(e) => setChapterName((e.target as HTMLInputElement).value)}
            class="input w-64"
            placeholder={window._('bookManager.edit.chapterNamePlaceholder')}
          />
          <button
            disabled={createChapterMutation.isPending || !validChapterName}
            class="btn btn--ghost"
            onClick={() => createChapterMutation.mutate()}
          >
            <i class="fa-solid fa-check" />
          </button>
        </div>
      </Popper>

      {savingOrderError && (
        <div class="mb-4">
          <ErrorDisplay error={savingOrderError} />
        </div>
      )}

      <div class="flex justify-between items-center mb-4">
        <span />
        <div class="flex gap-2">
          {!isReordering ? (
            <>
              <button class="btn btn--secondary " onClick={handleStartReordering}>
                {window._('bookManager.edit.reorder')}
              </button>
              <button
                ref={addChapterButton}
                class="btn btn--primary relative"
                onClick={() => setOpenChapterName((x) => !x)}
              >
                {window._('bookManager.edit.addChapter')}
              </button>
            </>
          ) : (
            <>
              <button class="btn btn--secondary " onClick={handleCancelReordering}>
                {window._('bookManager.edit.cancel')}
              </button>
              <button class="btn btn--primary" onClick={handleSaveOrder}>
                {isSavingOrder ? (
                  <span class="loader loader--dark" />
                ) : (
                  window._('bookManager.edit.save')
                )}
              </button>
            </>
          )}
        </div>
      </div>

      {isReordering && (
        <div class="mb-4 ol-alert ol-alert--warning">
          {window._('bookManager.edit.changesPending', { count: `${numberOfUpdates}` })}
        </div>
      )}

      <table class="table">
        <thead>
          <tr>
            <th>{window._('common.name')}</th>
            <th>{window._('common.info')}</th>
          </tr>
        </thead>
        <tbody>
          {chapters.map((chapter) => {
            return (
              <tr key={chapter.id}>
                <td>
                  {chapter.name} <br />
                  <div class="flex gap-3">
                    <a class="link" href={`/book/${book.id}/chapter/${chapter.id}`}>
                      {window._('common.view')}
                    </a>
                    |
                    <a class="link" href={`/books-manager/book/${book.id}/chapter/${chapter.id}`}>
                      {window._('common.edit')}
                    </a>
                  </div>
                </td>
                <td>{chapter.isAdultOverride}</td>
              </tr>
            )
          })}
        </tbody>
        <tfoot>
          <tr>
            <th>{window._('common.name')}</th>
            <th>{window._('common.info')}</th>
          </tr>
        </tfoot>
      </table>

      {!chapters?.length && (
        <div class="text-center py-8 text-gray-500">{window._('bookManager.edit.noChapters')}</div>
      )}
    </>
  )
}
