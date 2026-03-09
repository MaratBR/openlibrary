import { httpBmCreateChapter, ManagerBookChapterDto } from '@/api/bm'
import { ManagerBookDetailsDto } from '@/api/bm/book'
import Popper from '@/components/Popper'
import { formatNumberK } from '@/util/fmt'
import { useMutation } from '@tanstack/react-query'
import { useRef, useState } from 'preact/hooks'
import { NavLink } from 'react-router'

export function BookChapters({ book }: { book: ManagerBookDetailsDto }) {
  return (
    <>
      <AddChapterButton bookId={book.id} />

      <div class="card mt-4 px-0">
        <table class="table">
          <thead>
            <tr>
              <th class="w-8 text-muted-foreground">#</th>
              <th>{''}</th>
              <th>{''}</th>
            </tr>
          </thead>
          <tbody>
            {book.chapters.map((chapter, index) => (
              <ChapterRow key={chapter.id} index={index} book={book} chapter={chapter} />
            ))}
          </tbody>
        </table>
      </div>
    </>
  )
}

function ChapterRow({
  book,
  chapter,
  index,
}: {
  book: ManagerBookDetailsDto
  chapter: ManagerBookChapterDto
  index: number
}) {
  return (
    <tr>
      <td class="text-muted-foreground text-sm">{index + 1}</td>
      <td>
        <span class="text-xl font-medium">{chapter.name}</span>

        <div class="flex gap-2 mt-2">
          <NavLink to={`/books/${book.id}/chapters/${chapter.id}`} class="btn btn--lg">
            <i class="fa-solid fa-pen mr-2" />
            {window._('common.edit')}
          </NavLink>
        </div>
      </td>
      <td>
        <div class="flex gap-1">
          {chapter.isAdultOverride && <AdultChip />}
          {chapter.isPubliclyVisible && <HiddenChip />}
          <div class="chip">{window._('book.words', { count: formatNumberK(chapter.words) })}</div>
        </div>
      </td>
    </tr>
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

  const createChapter = useMutation({
    mutationFn: async () => {
      const response = await httpBmCreateChapter(bookId, {
        name,
        summary: '',
        isAdultOverride: false,
        content: '',
      })

      const chapterId = response.data
    },
  })

  return (
    <>
      <button onClick={() => setOpen(true)} ref={ref} class="btn primary btn--outline btn--lg mt-4">
        <i class="fa-solid fa-plus mr-2" />
        {window._('bookManager.edit.addChapter')}
      </button>
      <Popper onClose={() => setOpen(false)} open={open} placement="bottom-start" anchorEl={ref}>
        <div class="card max-w-128 shadow-2xl">
          <form action="#" onSubmit={handleSubmit}>
            <div class="flex gap-1">
              <input
                class="input"
                placeholder={window._('bookManager.edit.chapterNamePlaceholder')}
              />
              <button class="btn primary">{window._('bookManager.edit.addChapter')}</button>
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
