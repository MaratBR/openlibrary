import { useMemo, useRef, useState } from 'preact/hooks'
import { httpUpdateChaptersOrder, ManagerBookDetailsDto } from './api'
import { twMerge } from 'tailwind-merge'
import clsx from 'clsx'
import { ErrorDisplay } from '@/components/error'
import Popper from '@/components/Popper'
import { useMutation } from '@tanstack/react-query'
import { httpCreateChapter } from '../chapter/api'

type Chapter = ManagerBookDetailsDto['chapters'][number]

export default function Chapters({ book: data }: { book: ManagerBookDetailsDto }) {
  const [chapters, setChapters] = useState(data.chapters || [])
  const [originalOrder, setOriginalOrder] = useState<string[]>([])
  const [isReordering, setIsReordering] = useState(false)
  const [isSavingOrder, setSavingOrder] = useState(false)
  const [savingOrderError, setSavingOrderError] = useState<unknown>()
  const [openChapterName, setOpenChapterName] = useState(false)
  const [chapterName, setChapterName] = useState('')
  const validChapterName = chapterName.trim().length > 0

  const addChapterButton = useRef<HTMLButtonElement | null>(null)

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

  const createChapterMutation = useMutation({
    mutationFn: async () => {
      if (!validChapterName) {
        throw new Error(`invalid chapter name: ${chapterName}`)
      }

      const response = await httpCreateChapter(data.id, {
        name: chapterName,
        content: '',
        isAdultOverride: false,
        summary: '',
      })

      const chapterId = response.data

      window.location.href = `/books-manager/book/${data.id}/chapter/${chapterId}`
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
      <Popper
        style={openChapterName ? {} : { display: 'none' }}
        anchorEl={addChapterButton}
        placement="bottom-end"
      >
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
                class="btn btn--primary  relative"
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

      {chapters?.length ? (
        <div class="chapters-list">
          {isReordering ? (
            chapters.map((chapter: Chapter, index) => (
              <SortableChapterItem
                key={chapter.id}
                chapter={chapter}
                moveChapterUp={() => moveChapterUp(index)}
                moveChapterDown={() => moveChapterDown(index)}
                isFirst={index === 0}
                isLast={index === chapters.length - 1}
                isModified={chapter.order - 1 !== index}
              />
            ))
          ) : (
            <div class="space-y-2">
              {chapters.map((chapter) => (
                <ChapterCard key={chapter.id} chapter={chapter} bookId={data.id} />
              ))}
            </div>
          )}
        </div>
      ) : (
        <div class="text-center py-8 text-gray-500">{window._('bookManager.edit.noChapters')}</div>
      )}
    </>
  )
}

function SortableChapterItem({
  chapter,
  moveChapterUp,
  moveChapterDown,
  isFirst,
  isLast,
  isModified,
}: {
  chapter: Chapter
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
        <div class="flex flex-col mr-2">
          <button
            disabled={isFirst}
            onClick={moveChapterUp}
            class="flex items-center justify-center size-8 hover:bg-highlight disabled:pointer-events-none disabled:opacity-50"
          >
            <i class="fa-solid fa-arrow-up" />
          </button>
          <button
            disabled={isLast}
            onClick={moveChapterDown}
            class="flex items-center justify-center size-8 hover:bg-highlight disabled:pointer-events-none disabled:opacity-50"
          >
            <i class="fa-solid fa-arrow-down" />
          </button>
        </div>
        <div class="flex justify-between items-center flex-1">
          <div>
            <h3 class="font-medium">{chapter.name}</h3>
            <p class="text-sm text-gray-600">{chapter.summary}</p>
            <div class="text-xs text-gray-500 mt-1">
              {window._('bookManager.edit.words')}: {chapter.words}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

function ChapterCard({ chapter, bookId }: { chapter: Chapter; bookId: string }) {
  return (
    <div class="card p-0 overflow-hidden">
      <header class="bg-background px-4">
        <span class="text-xl font-medium leading-10">{chapter.name}</span>&nbsp;&mdash;&nbsp;
        <span class="text-sm">
          <a class="link" href={`/books-manager/book/${bookId}/chapter/${chapter.id}`}>
            {chapter.draftId
              ? window._('bookManager.edit.editDraft')
              : window._('bookManager.edit.edit')}
          </a>
        </span>
      </header>
      <div class="py-2 px-4">
        {chapter.summary ? (
          <p>{chapter.summary}</p>
        ) : (
          <p>{window._('bookManager.edit.emptyChapterSummary')}</p>
        )}
      </div>
      <div class="px-4 py-2 gap-1 flex">
        {chapter.draftId && (
          <span class="chip border bg-background">
            <i class="fa-solid fa-feather-pointed" />
            {window._('bookManager.edit.pendingChanges')}
          </span>
        )}
        {!chapter.isPubliclyVisible && (
          <span class="chip chip--primary">
            <i class="fa-solid fa-eye-slash" />
            {window._('bookManager.edit.chapterHidden')}
          </span>
        )}
        {chapter.isAdultOverride && (
          <span class="chip bg-red-600 text-white">{window._('bookManager.edit.adult')}</span>
        )}
      </div>
    </div>
  )
}
