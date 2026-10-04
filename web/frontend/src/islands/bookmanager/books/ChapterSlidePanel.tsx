import { ManagerBookChapterDto } from '@/api/bm'
import { BMBookAPI } from '@/api/bm/book'
import { FormControl } from '@/components/FormControl'
import { RichTextInput } from '@/components/rte'
import { useOLEditor } from '@/components/rte/RichTextEditor'
import Switch from '@/components/Switch'
import { useMutation } from '@tanstack/react-query'
import { useEffect, useState } from 'react'
import { useRevalidator } from 'react-router'
import { ManagerDialog, RequestError, useUnsavedChanges } from '../ui'

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
  return chapter ? (
    <ChapterEditForm key={chapter.id} bookId={bookId} chapter={chapter} onClose={onClose} />
  ) : null
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
  const [, setSummary] = useState(chapter.summary)
  const summaryEditor = useOLEditor({
    content: chapter.summary,
    onUpdate: ({ editor }) => setSummary(editor.getHTML()),
    editorProps: {
      attributes: {
        'aria-label': window._('bookManager.edit.summary'),
        role: 'textbox',
        'aria-multiline': 'true',
      },
    },
  })
  const [initialSummary] = useState(() => summaryEditor.getHTML())
  const [saved, setSaved] = useState(false)
  const dirty =
    !saved &&
    (name !== chapter.name ||
      isPubliclyVisible !== chapter.isPubliclyVisible ||
      summaryEditor.getHTML() !== initialSummary)
  const revalidator = useRevalidator()
  const normalizedName = BMBookAPI.getInstance().normalizeChapterName(name)
  const nameError = normalizedName.valid
    ? undefined
    : window._('bookManager.edit.chapterNameInvalid', { max: '70' })
  const saveMutation = useMutation({
    mutationFn: async () => {
      const response = await BMBookAPI.getInstance().updateChapter(bookId, chapter.id, {
        name: normalizedName.value,
        summary: summaryEditor.getHTML(),
        isPubliclyVisible,
      })
      response.throwIfError()
    },
    async onSuccess() {
      setSaved(true)
      await revalidator.revalidate()
      onClose()
    },
  })
  useEffect(() => {
    summaryEditor.setEditable(!saveMutation.isPending)
  }, [summaryEditor, saveMutation.isPending])
  const guard = useUnsavedChanges(dirty, saveMutation.isPending)
  const close = () => guard.leave(onClose)
  return (
    <>
      <ManagerDialog panel title={window._('bookManager.ui.chapterDetails')} onClose={close}>
        <form
          className="space-y-4"
          onSubmit={(event) => {
            event.preventDefault()
            if (normalizedName.valid && !saveMutation.isPending) saveMutation.mutate()
          }}
        >
          <header className="flex items-start justify-between gap-4">
            <div>
              <p className="text-secondary-foreground">{chapter.name}</p>
              <h2 className="text-2xl">{window._('bookManager.ui.chapterDetails')}</h2>
            </div>
            <button
              type="button"
              className="Btn Btn--icon Btn--ghost"
              aria-label={window._('bookManager.ui.close')}
              onClick={close}
            >
              ×
            </button>
          </header>
          <fieldset disabled={saveMutation.isPending} className="space-y-4">
            <FormControl
              htmlFor="chapter-name-input"
              label={window._('bookManager.ui.chapterName')}
              error={nameError}
            >
              <input
                data-dialog-focus
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
              <Switch
                value={isPubliclyVisible}
                onChange={setIsPubliclyVisible}
                slotProps={{
                  input: { 'aria-label': window._('bookManager.edit.isPubliclyVisible') },
                }}
              />
            </FormControl>
            {!chapter.isPubliclyVisible && isPubliclyVisible && (
              <div className="Alert Alert--warning" role="alert">
                {window._('bookManager.edit.chapterPublishWarning')}
              </div>
            )}
          </fieldset>
          {saveMutation.isError && <RequestError />}
          <p role="status">
            {saveMutation.isPending
              ? window._('bookManager.ui.saving')
              : dirty
                ? window._('bookManager.ui.unsaved')
                : ''}
          </p>
          <div className="flex gap-2">
            <button
              type="submit"
              className="Btn Btn--primary"
              disabled={!normalizedName.valid || saveMutation.isPending}
            >
              {window._('bookManager.ui.saveChanges')}
            </button>
            <button
              type="button"
              className="Btn Btn--outline"
              disabled={saveMutation.isPending}
              onClick={close}
            >
              {window._('common.cancel')}
            </button>
          </div>
          <div className="border-t pt-6">
            <a
              className="Link"
              href={`/books-manager/book/${bookId}/chapter/${chapter.id}`}
              onClick={(e) => {
                e.preventDefault()
                const href = e.currentTarget.href
                guard.leave(() => {
                  window.location.href = href
                })
              }}
            >
              {window._('bookManager.ui.openEditor')}
            </a>
          </div>
        </form>
      </ManagerDialog>
      {guard.prompt}
    </>
  )
}
