import BookContentEditor from './BookContentEditor'
import { useMemo, useRef, useState } from 'preact/hooks'
import { PreactIslandProps } from '@/islands/common/preact-island'
import { DraftDtoSchema } from '../contracts'
import { Editor } from '@tiptap/core'
import { z } from 'zod'
import { httpUpdateDraft } from './api'
import { createPortal } from 'preact/compat'
import clsx from 'clsx'
import { useMutation } from '@tanstack/react-query'

const dataSchema = z.object({
  bookId: z.string(),
  draft: DraftDtoSchema,
})

export default function BookManagerEditor({ data }: PreactIslandProps) {
  const { bookId, draft } = useMemo(() => dataSchema.parse(data), [data])

  const editorRef = useRef<Editor | null>(null)

  const [beforeSaving, setBeforeSaving] = useState(false)

  const savingMutation = useMutation({
    mutationFn: async (content: string) => {
      const t = performance.now()
      const d = window.delay(500)
      await httpUpdateDraft(bookId, draft.chapterId, draft.id, content)
      await d
      console.log(`saving took ${performance.now() - t} ms`)
    },
    onSettled() {
      setBeforeSaving(false)
    },
  })

  const saveAndPublishMutation = useMutation({
    mutationFn: async (content: string) => {
      const d = window.delay(500)
      await httpUpdateDraft(bookId, draft.chapterId, draft.id, content)
      await d
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

  return (
    <>
      <BookContentEditor
        editorRef={editorRef}
        contentChangedDebounce={1000}
        onContentChanged={handleContentChange}
        onBeforeContentChanged={handleBeforeContentChanged}
        content={draft.content}
      />
      {createPortal(
        <>
          <button
            disabled={saveAndPublishMutation.isPending}
            onClick={() => saveAndPublish()}
            id="actions:saveAndPublish"
            class={clsx('ol-btn ol-btn--secondary rounded-full', {
              'with-loader': saveAndPublishMutation.isPending,
            })}
          >
            {window._('editor.publishDraft')}
          </button>
          <button
            disabled={savingMutation.isPending}
            onClick={() => save()}
            id="actions:save"
            class={clsx('w-28 ol-btn ol-btn--primary rounded-full', {
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
