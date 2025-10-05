import { Editor, EditorOptions } from '@tiptap/core'
import HorizontalRule from '@tiptap/extension-horizontal-rule'
import TextStyle from '@tiptap/extension-text-style'
import Typography from '@tiptap/extension-typography'
import TextAlign from '@tiptap/extension-text-align'
import Image from '@tiptap/extension-image'

import StarterKit from '@tiptap/starter-kit'
import { useEffect, useMemo, useRef, useState } from 'preact/hooks'
import './BookManagerEditor.scss'
import Heading from '@tiptap/extension-heading'
import { createEvent } from '@/lib/event'
import BookContentEditorHeader from './BookContentEditorHeader'
import Underline from '@tiptap/extension-underline'
import Placeholder from '@tiptap/extension-placeholder'
import { debounce } from '@/common/util/debounce'
import { RefObject } from 'preact'
import { MouseEventHandler } from 'preact/compat'
import { useMutation } from '@tanstack/react-query'
import { httpUpdateDraftChapterName } from './api'
import { DraftDto } from '../contracts'

export type BookContentEditorProps = {
  bookId: string
  draft: DraftDto
  onContentChanged: (editor: Editor) => void
  onBeforeContentChanged?: () => void
  contentChangedDebounce: number
  editorRef?: RefObject<Editor>
}

export default function BookContentEditor({
  bookId,
  draft,
  onContentChanged,
  contentChangedDebounce,
  onBeforeContentChanged,
  editorRef,
}: BookContentEditorProps) {
  const editor = useRef<Editor | null>(null)
  const root = useRef<HTMLDivElement | null>(null)

  const editorUpdateEvent = useRef(createEvent<Editor>())

  const handleContentChanged = useMemo(
    () =>
      debounce((editor: Editor) => {
        refs.current.onBeforeContentChangedCalled = false
        refs.current.onContentChanged(editor)
      }, contentChangedDebounce),
    [contentChangedDebounce],
  )
  const refs = useRef({
    handleContentChanged,
    onBeforeContentChangedCalled: false,
    onBeforeContentChanged,
    onContentChanged,
  })
  refs.current.onContentChanged = onContentChanged
  refs.current.onBeforeContentChanged = onBeforeContentChanged
  refs.current.handleContentChanged = handleContentChanged

  const propsRef = useRef({ content: draft.content })
  propsRef.current.content = draft.content

  useEffect(() => {
    if (root.current === null) throw new Error('root is null')
    editor.current = createEditor(
      root.current,
      { placeholder: window._('editor.placeholder') },
      {
        content: propsRef.current.content,
        onUpdate: ({ editor }) => {
          editorUpdateEvent.current.fire(editor)

          refs.current.handleContentChanged(editor)
          if (!refs.current.onBeforeContentChangedCalled) {
            refs.current.onBeforeContentChangedCalled = true
            if (refs.current.onBeforeContentChanged) refs.current.onBeforeContentChanged()
          }
        },
        onTransaction: ({ editor }) => {
          editorUpdateEvent.current.fire(editor)
        },
      },
    )
  }, [])

  useEffect(() => {
    if (editorRef) editorRef.current = editor.current
  }, [editorRef])

  const previousChapterName = useRef(draft.chapterName)
  const [chapterName, setChapterName] = useState(draft.chapterName)

  const updateChapterNameMutation = useMutation({
    mutationFn: async (name: string) => {
      await httpUpdateDraftChapterName(bookId, draft.chapterId, draft.id, name)
    },
  })

  function handleChapterNameBlur() {
    const newName = chapterName.trim()
    if (newName === previousChapterName.current) {
      return
    }

    previousChapterName.current = newName
    updateChapterNameMutation.mutate(newName)
  }

  return (
    <div class="book-editor">
      <BookContentEditorHeader editorRef={editor} editorUpdateEvent={editorUpdateEvent.current} />
      <article class="ol-container">
        <div class="mb-8">
          <ChapterNameInput
            value={chapterName}
            onChange={setChapterName}
            onBlur={handleChapterNameBlur}
          />
        </div>
        <div class="__user-content book-editor__content" ref={root} />
      </article>
    </div>
  )
}

function createEditor(
  editorElement: HTMLElement,
  {
    placeholder,
  }: {
    placeholder: string
  },
  options?: Partial<EditorOptions>,
) {
  return new Editor({
    element: editorElement,
    content: '',
    extensions: [
      StarterKit.configure({
        horizontalRule: false,
        codeBlock: false,
        heading: false,
        code: { HTMLAttributes: { class: 'inline', spellcheck: 'false' } },
        dropcursor: { width: 2, class: 'ProseMirror-dropcursor border' },
      }),
      TextStyle,
      Typography,
      HorizontalRule,
      Heading,
      TextAlign.configure({
        types: ['heading', 'paragraph'],
      }),
      Underline,
      Image.configure({
        inline: true,
      }),
      Placeholder.configure({
        placeholder,
      }),
    ],
    ...options,
  })
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
      // onFocus={(e) => {
      //   ;(e.target as HTMLInputElement).select()
      // }}
      placeholder={window._('editor.chapterNamePlaceholder')}
      onBlur={onBlur}
      class="font-title text-3xl leading-3 bg-transparent focus:outline-none py-2"
      value={value}
      onChange={(e) => onChange((e.target as HTMLInputElement).value)}
    />
  )
}
