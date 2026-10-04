import { ManagerBookChapterDto } from '@/api/bm'
import { BMBookAPI, ManagerBookDetailsDto } from '@/api/bm/book'
import Popper from '@/components/Popper'
import { useMutation } from '@tanstack/react-query'
import { useRef, useState } from 'react'
import { useRevalidator } from 'react-router'
import { ChapterSlidePanel } from './ChapterSlidePanel'
import { RequestError, useUnsavedChanges } from '../ui'

export function BookChapters({ book }: { book: ManagerBookDetailsDto }) {
  const [chapter, setChapter] = useState<ManagerBookChapterDto | null>(null)
  const [adding, setAdding] = useState(false)
  const trigger = useRef<HTMLButtonElement>(null)
  return (
    <>
      <button ref={trigger} onClick={() => setAdding(true)} className="Btn Btn--primary">
        + {window._('bookManager.edit.addChapter')}
      </button>
      {adding && (
        <AddChapterForm
          bookId={book.id}
          anchor={trigger}
          onClose={() => {
            setAdding(false)
            trigger.current?.focus()
          }}
        />
      )}
      <section className="space-y-4 mt-4">
        {!book.chapters.length && (
          <p className="Card">{window._('bookManager.ui.emptyChapters')}</p>
        )}
        {book.chapters.map((chapter) => (
          <article
            key={chapter.id}
            className="Card flex flex-wrap justify-between gap-4"
            data-testid="BookChapters_Chapter"
          >
            <div className="min-w-0">
              <h2 className="text-xl">{chapter.name}</h2>
              <p className="text-secondary-foreground my-2">
                {window._(
                  chapter.isPubliclyVisible ? 'bookManager.ui.visible' : 'bookManager.ui.hidden',
                )}{' '}
                ·{' '}
                {window._(
                  chapter.words === 1 ? 'bookManager.ui.wordCountOne' : 'bookManager.ui.wordCount',
                  { count: chapter.words.toLocaleString() },
                )}
              </p>
              {chapter.scheduledAt && (
                <p className="text-sm text-secondary-foreground">
                  {window._('bookManager.ui.scheduled', {
                    date: new Date(chapter.scheduledAt).toLocaleString(undefined, {
                      timeZoneName: 'short',
                    }),
                  })}
                </p>
              )}
            </div>
            <div className="flex gap-3 flex-wrap items-center">
              <button
                disabled={adding}
                className="Btn Btn--outline"
                onClick={() => setChapter(chapter)}
              >
                {window._('bookManager.ui.editDetails')}
              </button>
              <a
                className="Link"
                aria-disabled={adding}
                onClick={(e) => {
                  if (adding) e.preventDefault()
                }}
                href={`/books-manager/book/${book.id}/chapter/${chapter.id}`}
              >
                {window._('bookManager.ui.openEditor')}
              </a>
            </div>
          </article>
        ))}
      </section>
      <ChapterSlidePanel bookId={book.id} chapter={chapter} onClose={() => setChapter(null)} />
    </>
  )
}

function AddChapterForm({
  bookId,
  anchor,
  onClose,
}: {
  bookId: string
  anchor: React.RefObject<HTMLButtonElement | null>
  onClose: () => void
}) {
  const [name, setName] = useState('')
  const revalidator = useRevalidator()
  const normalized = BMBookAPI.getInstance().normalizeChapterName(name)
  const createChapter = useMutation({
    mutationFn: async () => {
      const response = await BMBookAPI.getInstance().createChapter(bookId, {
        name: normalized.value,
        summary: '',
        content: '',
      })
      response.throwIfError()
    },
    onSuccess: async () => {
      await revalidator.revalidate()
      onClose()
    },
  })
  const guard = useUnsavedChanges(!!name && !createChapter.isSuccess, createChapter.isPending)
  const close = () => guard.leave(onClose)
  return (
    <>
      <Popper
        className="BM"
        onClose={() => {
          if (!document.querySelector('dialog[open]')) close()
        }}
        open
        placement="bottom-start"
        anchorEl={anchor}
      >
        <div className="Card shadow-2xl" style={{ width: 'min(420px, calc(100vw - 32px))' }}>
          <form
            className="space-y-3"
            onSubmit={(e) => {
              e.preventDefault()
              if (normalized.valid && !createChapter.isPending) createChapter.mutate()
            }}
          >
            <label htmlFor="add-chapter-name" className="label">
              {window._('bookManager.ui.chapterName')}
            </label>
            <input
              autoFocus
              id="add-chapter-name"
              className="input"
              value={name}
              disabled={createChapter.isPending}
              aria-describedby="add-chapter-help"
              aria-invalid={!!name && !normalized.valid}
              onChange={(e) => setName(e.currentTarget.value)}
              onKeyDown={(e) => {
                if (e.key === 'Escape') close()
              }}
            />
            <p id="add-chapter-help" className="text-sm text-secondary-foreground">
              {window._('bookManager.edit.chapterNameInvalid', { max: '70' })}
            </p>
            <p role="status">{createChapter.isPending ? window._('bookManager.ui.saving') : ''}</p>
            {createChapter.isError && <RequestError />}
            <div className="flex flex-wrap gap-2">
              <button
                disabled={!normalized.valid || createChapter.isPending}
                className="Btn Btn--primary"
              >
                {window._(
                  createChapter.isPending ? 'bookManager.ui.saving' : 'bookManager.edit.addChapter',
                )}
              </button>
              <button
                type="button"
                disabled={createChapter.isPending}
                className="Btn Btn--outline"
                onClick={close}
              >
                {window._('common.cancel')}
              </button>
            </div>
          </form>
        </div>
      </Popper>
      {guard.prompt}
    </>
  )
}
