import { ManagerBookChapterDto } from '@/api/bm'
import { BMBookAPI } from '@/api/bm/book'
import { FormControl } from '@/components/FormControl'
import { RichTextInput } from '@/components/rte'
import { useOLEditor } from '@/components/rte/RichTextEditor'
import { SlidePanel } from '@/components/SlidePanel'
import Switch from '@/components/Switch'
import { useMutation } from '@tanstack/react-query'
import { FormEvent, useMemo, useState } from 'react'
import { useRevalidator } from 'react-router'

export const CHAPTER_SLIDE_OUT_PARAMETER_NAME = 'chapter'

export function ChapterSlidePanel({
  bookId,
  chapter,
  onClose,
}: {
  bookId: string
  chapter: ManagerBookChapterDto | null
  onClose: () => void
}) {
  return (
    <SlidePanel.Facade open={!!chapter} onClose={onClose}>
      <SlidePanel.Content>
        {chapter && (
          <ChapterEditForm key={chapter.id} bookId={bookId} chapter={chapter} onClose={onClose} />
        )}
      </SlidePanel.Content>
    </SlidePanel.Facade>
  )
}

function ChapterEditForm({
  bookId,
  chapter,
  onClose,
}: {
  bookId: string
  chapter: ManagerBookChapterDto
  onClose: () => void
}) {
  const [name, setName] = useState(chapter.name)
  const [isPubliclyVisible, setIsPubliclyVisible] = useState(chapter.isPubliclyVisible)
  const summaryEditor = useOLEditor({ content: chapter.summary })
  const revalidator = useRevalidator()

  const normalizedName = useMemo(
    () => BMBookAPI.getInstance().normalizeChapterName(name),
    [name],
  )
  const nameError = normalizedName.valid
    ? undefined
    : window._('bookManager.edit.chapterNameInvalid', { max: '70' })
  const willBecomePublic = !chapter.isPubliclyVisible && isPubliclyVisible

  const saveMutation = useMutation({
    mutationFn: () =>
      BMBookAPI.getInstance().updateChapter(bookId, chapter.id, {
        name: normalizedName.value,
        summary: summaryEditor.getHTML(),
        isPubliclyVisible,
      }),
    async onSuccess() {
      await revalidator.revalidate()
      onClose()
    },
    onError(error) {
      window.toast.error(error)
    },
  })

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (normalizedName.valid && !saveMutation.isPending) saveMutation.mutate()
  }

  return (
    <form className="space-y-4" onSubmit={handleSubmit}>
      <div className="flex items-center justify-between gap-4">
        <h2 className="text-2xl font-medium">{chapter.name}</h2>
        <button type="button" className="btn btn--icon btn--ghost" onClick={onClose}>
          <i className="fa-solid fa-xmark" />
        </button>
      </div>

      <FormControl
        htmlFor="chapter-name-input"
        label={window._('bookManager.edit.name')}
        error={nameError}
      >
        <input
          id="chapter-name-input"
          className="input"
          value={name}
          aria-invalid={!!nameError}
          aria-describedby={nameError ? 'chapter-name-input-error' : undefined}
          onChange={(event) => setName(event.currentTarget.value)}
        />
      </FormControl>

      <FormControl label={window._('bookManager.edit.summary')}>
        <RichTextInput editor={summaryEditor} />
      </FormControl>

      <FormControl
        label={window._('bookManager.edit.isPubliclyVisible')}
        description={window._('bookManager.edit.chapterVisibilityDescription')}
      >
        <Switch value={isPubliclyVisible} onChange={setIsPubliclyVisible} />
      </FormControl>

      {willBecomePublic && (
        <div className="alert alert--warning" role="alert">
          {window._('bookManager.edit.chapterPublishWarning')}
        </div>
      )}

      <button
        type="submit"
        className="btn btn--default"
        disabled={!normalizedName.valid || saveMutation.isPending}
      >
        {window._('common.save')}
      </button>

      <div className="mt-6">
        <a href={`/books-manager/book/${bookId}/chapter/${chapter.id}`} className="link">
          {window._('bookManager.edit.goToEditor')}
        </a>
      </div>
    </form>
  )
}
