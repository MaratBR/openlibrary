import { useMemo, useState } from 'preact/hooks'
import { httpUpdateChaptersOrder, managerBookDetailsSchema } from './api'
import { twMerge } from 'tailwind-merge'
import clsx from 'clsx'
import { ErrorDisplay } from '@/components/error'
import { PreactIslandProps } from '@/islands/common/preact-island'

type Chapter = {
  id: string
  name: string
  words: number
  order: number
  createdAt: string
  summary: string
}

function SortableChapterItem({
  chapter,
  bookId,
  isReordering,
  moveChapterUp,
  moveChapterDown,
  isFirst,
  isLast,
  isModified,
}: {
  chapter: Chapter
  bookId: string
  isReordering: boolean
  moveChapterUp: () => void
  moveChapterDown: () => void
  isFirst: boolean
  isLast: boolean
  isModified: boolean
}) {
  return (
    <div
      className={twMerge(
        clsx(
          'card relative p-4 mb-2 before:absolute before:block before:bg-primary before:w-2 before:h-full before:left-0 before:top-0 before:invisible overflow-hidden',
          {
            'before:visible': isModified,
          },
        ),
      )}
    >
      <div class="flex items-center">
        {isReordering && (
          <div class="flex flex-col mr-2">
            <button
              disabled={isFirst}
              onClick={moveChapterUp}
              class="flex items-center justify-center size-8 hover:bg-highlight disabled:pointer-events-none disabled:opacity-50"
            >
              <span class="material-symbols-outlined">arrow_upward</span>
            </button>
            <button
              disabled={isLast}
              onClick={moveChapterDown}
              class="flex items-center justify-center size-8 hover:bg-highlight disabled:pointer-events-none disabled:opacity-50"
            >
              <span class="material-symbols-outlined">arrow_downward</span>
            </button>
          </div>
        )}
        <div class="flex justify-between items-center flex-1">
          <div>
            <h3 class="font-medium">{chapter.name}</h3>
            <p class="text-sm text-gray-600">{chapter.summary}</p>
            <div class="text-xs text-gray-500 mt-1">
              {window._('bookManager.edit.words')}: {chapter.words}
            </div>
          </div>
          <div class="flex gap-2" style={isReordering ? { display: 'none' } : undefined}>
            <a
              target="_blank"
              href={`/books-manager/book/${bookId}/chapter/${chapter.id}`}
              class="btn btn--secondary"
              rel="noreferrer"
            >
              {window._('bookManager.edit.edit')}
              <span class="material-symbols-outlined">open_in_new</span>
            </a>
            <button disabled={isReordering} class="btn btn--destructive">
              {window._('bookManager.edit.delete')}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

export default function Chapters({ data: dataUnknown }: PreactIslandProps) {
  const data = useMemo(() => managerBookDetailsSchema.parse(dataUnknown), [dataUnknown])
  const [chapters, setChapters] = useState(data.chapters || [])
  const [originalOrder, setOriginalOrder] = useState<string[]>([])
  const [isReordering, setIsReordering] = useState(false)
  const [isSavingOrder, setSavingOrder] = useState(false)
  const [savingOrderError, setSavingOrderError] = useState<unknown>()

  const moveChapterUp = (index: number) => {
    if (index > 0) {
      setChapters((prevChapters) => {
        const newChapters = [...prevChapters]
        const temp = newChapters[index - 1]
        newChapters[index - 1] = newChapters[index]
        newChapters[index] = temp
        return newChapters
      })
    }
  }

  const moveChapterDown = (index: number) => {
    if (index < chapters.length - 1) {
      setChapters((prevChapters) => {
        const newChapters = [...prevChapters]
        const temp = newChapters[index + 1]
        newChapters[index + 1] = newChapters[index]
        newChapters[index] = temp
        return newChapters
      })
    }
  }

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

    httpUpdateChaptersOrder(data.id, newOrder)
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
              <button class="btn btn--secondary rounded-full" onClick={handleStartReordering}>
                {window._('bookManager.edit.reorder')}
              </button>
              <button class="btn btn--primary rounded-full">
                {window._('bookManager.edit.addChapter')}
              </button>
            </>
          ) : (
            <>
              <button
                class="btn btn--secondary rounded-full"
                onClick={handleCancelReordering}
              >
                {window._('bookManager.edit.cancel')}
              </button>
              <button class="btn btn--primary  rounded-full" onClick={handleSaveOrder}>
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

      {chapters?.length ? (
        <div class="chapters-list">
          {chapters.map((chapter: Chapter, index) => (
            <SortableChapterItem
              key={chapter.id}
              bookId={data.id}
              chapter={chapter}
              isReordering={isReordering}
              moveChapterUp={() => moveChapterUp(index)}
              moveChapterDown={() => moveChapterDown(index)}
              isFirst={index === 0}
              isLast={index === chapters.length - 1}
              isModified={chapter.order - 1 !== index}
            />
          ))}
        </div>
      ) : (
        <div class="text-center py-8 text-gray-500">{window._('bookManager.edit.noChapters')}</div>
      )}
    </>
  )
}
