import { ManagerBookChapterDto } from '@/api/bm'
import { BMBookAPI, ManagerBookDetailsDto } from '@/api/bm/book'
import Popper from '@/components/Popper'
import { formatNumberK } from '@/util/fmt'
import { useMutation } from '@tanstack/react-query'
import { SubmitEvent, useMemo, useRef, useState } from 'react'
import { useRevalidator } from 'react-router'
import { ChapterSlidePanel } from './ChapterSlidePanel'

export function BookChapters({ book }: { book: ManagerBookDetailsDto }) {
  const [chapter, setChapter] = useState<ManagerBookChapterDto | null>(null)

  return (
    <>
      <AddChapterButton bookId={book.id} />

      <section className="space-y-4 mt-4">
        {book.chapters.map((chapter) => (
          <Chapter key={chapter.id} chapter={chapter} onOpenChapter={setChapter} />
        ))}
      </section>

      <ChapterSlidePanel chapter={chapter} onClose={() => setChapter(null)} />
    </>
  )
}

function Chapter({
  chapter,
  onOpenChapter,
}: {
  chapter: ManagerBookChapterDto
  onOpenChapter: (chapter: ManagerBookChapterDto) => void
}) {
  return (
    <div data-testid="BookChapters_Chapter" className="card grid cols-2">
      <div>
        <span className="text-xl font-medium">{chapter.name}</span>
        <div className="flex gap-2 mt-2">
          <button onClick={() => onOpenChapter(chapter)} className="btn btn--lg btn--default">
            <i className="fa-solid fa-pen mr-2" />
            {window._('common.edit')}
          </button>
        </div>
      </div>

      <div>
        <div className="flex gap-1">
          {chapter.isAdultOverride && <AdultChip />}
          {chapter.isPubliclyVisible && <HiddenChip />}
          <div className="chip">
            {window._('book.words', { count: formatNumberK(chapter.words) })}
          </div>
        </div>
      </div>
    </div>
  )
}

function AdultChip() {
  return <div className="chip chip--destructive">{window._('common.adult')}</div>
}

function HiddenChip() {
  return (
    <div className="chip chip--secondary">
      <i className="fa-solid fa-eye-slash mr-1" />
      {window._('bookManager.edit.chapterHidden')}
    </div>
  )
}

function AddChapterButton({ bookId }: { bookId: string }) {
  const ref = useRef<HTMLButtonElement | null>(null)
  const [open, setOpen] = useState(false)
  const [name, setName] = useState('')

  const revalidator = useRevalidator()

  const valid = useMemo(() => {
    const { valid } = BMBookAPI.getInstance().normalizeChapterName(name)
    return valid
  }, [name])

  const createChapter = useMutation({
    mutationFn: async () => {
      await BMBookAPI.getInstance().createChapter(bookId, {
        name,
        summary: '',
        isAdultOverride: false,
        content: '',
      })

      setOpen(false)
      await revalidator.revalidate()
      setName('')
    },
  })

  return (
    <>
      <button onClick={() => setOpen(true)} ref={ref} className="btn btn--outline btn--lg mt-4">
        <i className="fa-solid fa-plus mr-2" />
        {window._('bookManager.edit.addChapter')}
      </button>
      <Popper onClose={() => setOpen(false)} open={open} placement="bottom-start" anchorEl={ref}>
        <div className="card max-w-128 shadow-2xl">
          <form action="#" onSubmit={handleSubmit}>
            <div className="flex gap-1">
              <input
                className="input"
                value={name}
                onChange={(e) => setName((e.target as HTMLInputElement).value)}
                placeholder={window._('bookManager.edit.chapterNamePlaceholder')}
              />
              <button disabled={!valid} className="btn btn--default">
                {window._('bookManager.edit.addChapter')}
              </button>
            </div>
          </form>
        </div>
      </Popper>
    </>
  )

  function handleSubmit(e: SubmitEvent) {
    e.preventDefault()
    createChapter.mutate()
  }
}
