import BookContentEditor from './BookContentEditor'
import { useMemo, useRef, useState } from 'preact/hooks'
import { PreactIslandProps } from '@/islands/common/preact-island'
import { DraftDtoSchema } from '../contracts'
import { Editor } from '@tiptap/core'
import { z } from 'zod'
import { httpUpdateAndPublishDraft, httpUpdateDraft, httpUpdateDraftChapterName } from './api'
import { createPortal, MouseEventHandler } from 'preact/compat'
import clsx from 'clsx'
import { useMutation } from '@tanstack/react-query'

const dataSchema = z.object({
  bookId: z.string(),
  draft: DraftDtoSchema,
})

export default function BookManagerEditor({ data }: PreactIslandProps) {
  const { bookId, draft } = useMemo(() => dataSchema.parse(data), [data])

  const editorRef = useRef<Editor | null>(null)

  const previousChapterName = useRef(draft.chapterName)
  const [chapterName, setChapterName] = useState(draft.chapterName)
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
      await httpUpdateAndPublishDraft(bookId, draft.chapterId, draft.id, content)
      await d
    },
  })

  const updateChapterNameMutation = useMutation({
    mutationFn: async (name: string) => {
      await httpUpdateDraftChapterName(bookId, draft.chapterId, draft.id, name)
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

  function handleChapterNameBlur() {
    const newName = chapterName.trim()
    if (newName === previousChapterName.current) {
      return
    }

    previousChapterName.current = newName
    updateChapterNameMutation.mutate(newName)
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
        <ChapterNameInput
          value={chapterName}
          onChange={setChapterName}
          onBlur={handleChapterNameBlur}
        />,
        document.getElementById('slot:header-text')!,
      )}
      {createPortal(
        <>
          <button
            disabled={saveAndPublishMutation.isPending}
            onClick={() => saveAndPublish()}
            id="actions:saveAndPublish"
            class={clsx('btn btn--secondary rounded-full', {
              'with-loader': saveAndPublishMutation.isPending,
            })}
          >
            {window._('editor.publishDraft')}
          </button>
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

function ChapterNameInput({
  value,
  onChange,
  onBlur,
}: {
  value: string
  onChange: (value: string) => void
  onBlur: MouseEventHandler<HTMLInputElement>
}) {
  return (
    <input
      onFocus={(e) => {
        ;(e.target as HTMLInputElement).select()
      }}
      onBlur={onBlur}
      class="h-full outline-none w-full focus:bg-muted mr-12 pl-2 -ml-2"
      value={value}
      onChange={(e) => onChange((e.target as HTMLInputElement).value)}
    />
  )
}
