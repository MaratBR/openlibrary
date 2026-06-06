import { ManagerBookChapterDto } from '@/api/bm'
import { BMBookAPI, ManagerBookDetailsDto } from '@/api/bm/book'
import Popper from '@/components/Popper'
import { formatNumberK } from '@/util/fmt'
import { useMutation } from '@tanstack/react-query'
import { useMemo, useRef, useState } from 'preact/hooks'
import { NavLink, useRevalidator } from 'react-router'
import { CHAPTER_SLIDE_OUT_PARAMETER_NAME, ChapterSlidePanel } from './chapter-page'

export function BookChapters({ book }: { book: ManagerBookDetailsDto }) {
  return (
    <>
      <AddChapterButton bookId={book.id} />

      <section class="space-y-4 mt-4">
        {book.chapters.map((chapter, index) => (
          <Chapter key={chapter.id} index={index} book={book} chapter={chapter} />
        ))}
      </section>

      <ChapterSlidePanel />
    </>
  )
}

function Chapter({
  book,
  chapter,
  index,
}: {
  book: ManagerBookDetailsDto
  chapter: ManagerBookChapterDto
  index: number
}) {
  return (
    <div data-testid="BookChapters_Chapter" class="card rounded-lg grid cols-2">
      <div>
        <span class="text-xl font-medium">{chapter.name}</span>
        <div class="flex gap-2 mt-2">
          <NavLink
            to={`/books/${book.id}?t=chapters&${CHAPTER_SLIDE_OUT_PARAMETER_NAME}=${chapter.id}`}
            className="btn btn--lg"
          >
            <i class="fa-solid fa-pen mr-2" />
            {window._('common.edit')}
          </NavLink>
        </div>
      </div>

      <div>
        <div class="flex gap-1">
          {chapter.isAdultOverride && <AdultChip />}
          {chapter.isPubliclyVisible && <HiddenChip />}
          <div class="chip">{window._('book.words', { count: formatNumberK(chapter.words) })}</div>
        </div>
      </div>
    </div>
  )
}

function AdultChip() {
  return <div class="chip chip--destructive">{window._('common.adult')}</div>
}

function HiddenChip() {
  return (
    <div class="chip chip--secondary">
      <i class="fa-solid fa-eye-slash mr-1" />
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
      <button onClick={() => setOpen(true)} ref={ref} class="btn btn--outline btn--lg mt-4">
        <i class="fa-solid fa-plus mr-2" />
        {window._('bookManager.edit.addChapter')}
      </button>
      <Popper onClose={() => setOpen(false)} open={open} placement="bottom-start" anchorEl={ref}>
        <div class="card max-w-128 shadow-2xl">
          <form action="#" onSubmit={handleSubmit}>
            <div class="flex gap-1">
              <input
                class="input"
                value={name}
                onChange={(e) => setName((e.target as HTMLInputElement).value)}
                placeholder={window._('bookManager.edit.chapterNamePlaceholder')}
              />
              <button disabled={!valid} class="btn">
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
