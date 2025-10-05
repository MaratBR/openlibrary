import BookContentEditor from './BookContentEditor'
import { useMemo, useRef, useState } from 'preact/hooks'
import { PreactIslandProps } from '@/islands/common/preact-island'
import { DraftDtoSchema } from '../contracts'
import { Editor } from '@tiptap/core'
import { z } from 'zod'
import { httpUpdateAndPublishDraft, httpUpdateDraft } from './api'
import { createPortal } from 'preact/compat'
import clsx from 'clsx'
import { useMutation } from '@tanstack/react-query'
import Switch from '@/components/Switch'

const dataSchema = z.object({
  bookId: z.string(),
  draft: DraftDtoSchema,
})

export default function BookManagerEditor({ data }: PreactIslandProps) {
  const { bookId, draft } = useMemo(() => dataSchema.parse(data), [data])

  const editorRef = useRef<Editor | null>(null)

  const [makeChapterVisible, setMakeChapterVisible] = useState(true)
  const [beforeSaving, setBeforeSaving] = useState(false)
  const [publishPopupOpen, setPublishPopupOpen] = useState(false)

  const savingMutation = useMutation({
    mutationFn: async (content: string) => {
      await httpUpdateDraft(bookId, draft.chapterId, draft.id, content)
      setPublishPopupOpen(false)
    },
    onSettled() {
      setBeforeSaving(false)
    },
  })

  const saveAndPublishMutation = useMutation({
    mutationFn: async (content: string) => {
      await httpUpdateAndPublishDraft(
        bookId,
        draft.chapterId,
        draft.id,
        content,
        makeChapterVisible,
      )
      setPublishPopupOpen(false)
    },
  })

  const refs = useRef({ content: '' })

  async function handleContentChange(editorParam?: Editor) {
    const editor = editorParam ?? editorRef.current
    if (!editor) throw new Error('cannot find editor')
    if (savingMutation.isPending) return

    refs.current.content = editor.getHTML()
    savingMutation.mutate(refs.current.content)
  }

  function save() {
    if (!editorRef.current) return
    const content = editorRef.current.getHTML()
    savingMutation.mutate(content)
  }

  function saveAndPublish() {
    if (!editorRef.current) return
    const content = editorRef.current.getHTML()
    saveAndPublishMutation.mutate(content)
  }

  const publishButtonRef = useRef<HTMLButtonElement | null>(null)

  return (
    <>
      <BookContentEditor
        editorRef={editorRef}
        contentChangedDebounce={1000}
        onContentChanged={handleContentChange}
        onBeforeContentChanged={handleBeforeContentChanged}
        draft={draft}
        bookId={bookId}
      />
      {createPortal(
        <>
          <button
            ref={publishButtonRef}
            disabled={saveAndPublishMutation.isPending}
            onClick={() => setPublishPopupOpen(true)}
            id="actions:saveAndPublish"
            class="btn btn--secondary rounded-full"
          >
            {window._('editor.publishDraft')}
          </button>
          <div class="relative">
            <div
              data-open={publishPopupOpen}
              class="card p-4 shadow-2xl absolute right-0 top-0 rounded-2xl min-w-48 transition-opacity data-[open=true]:opacity-1 data-[open=false]:opacity-0 data-[open=false]:pointer-events-none"
            >
              <strong>{window._('editor.publishAreYouSure')}</strong>
              <p>{window._('editor.publishWarning')}</p>
              {!draft.isChapterPubliclyAvailable && (
                <div class="form-control form-control--horizontal bg-muted rounded-xl p-2 -m-2 mt-1">
                  <div class="form-control__label p-0">{window._('editor.makeChapterVisible')}</div>
                  <div class="form-control__value">
                    <Switch value={makeChapterVisible} onChange={setMakeChapterVisible} />
                  </div>
                </div>
              )}
              <div class="flex -ml-2 gap-1 mt-4">
                <button
                  class="btn btn--destructive rounded-full"
                  onClick={() => setPublishPopupOpen(false)}
                >
                  {window._('common.cancel')}
                </button>
                <button
                  ref={publishButtonRef}
                  disabled={saveAndPublishMutation.isPending}
                  onClick={() => saveAndPublish()}
                  class={clsx('btn btn--secondary rounded-full', {
                    'with-loader': saveAndPublishMutation.isPending,
                  })}
                >
                  {window._('editor.publishDraft')}
                </button>
              </div>
            </div>
          </div>
          <button
            disabled={savingMutation.isPending}
            onClick={() => save()}
            id="actions:save"
            class={clsx('w-28 btn btn--primary rounded-full', {
              'with-loader': savingMutation.isPending || beforeSaving,
            })}
          >
            {window._('editor.saveDraft')}
          </button>
        </>,
        document.getElementById('slot:actions')!,
      )}
    </>
  )

  function handleBeforeContentChanged() {
    setBeforeSaving(true)
  }
}
