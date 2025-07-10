import { Editor, EditorOptions } from '@tiptap/core'
import HorizontalRule from '@tiptap/extension-horizontal-rule'
import TextStyle from '@tiptap/extension-text-style'
import Typography from '@tiptap/extension-typography'
import TextAlign from '@tiptap/extension-text-align'
import Image from '@tiptap/extension-image'

import StarterKit from '@tiptap/starter-kit'
import { useEffect, useMemo, useRef } from 'preact/hooks'
import './BookManagerEditor.scss'
import Heading from '@tiptap/extension-heading'
import { createEvent } from '@/lib/event'
import BookContentEditorHeader from './BookContentEditorHeader'
import Underline from '@tiptap/extension-underline'
import { debounce } from '@/common/util/debounce'
import { RefObject } from 'preact'

export type BookContentEditorProps = {
  content: string
  onContentChanged: (editor: Editor) => void
  onBeforeContentChanged?: () => void
  contentChangedDebounce: number
  editorRef?: RefObject<Editor>
}

export default function BookContentEditor({
  content,
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

  const propsRef = useRef({ content })
  propsRef.current.content = content

  useEffect(() => {
    if (root.current === null) throw new Error('root is null')
    editor.current = createEditor(root.current, {
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
    })
  }, [])

  useEffect(() => {
    if (editorRef) editorRef.current = editor.current
  }, [editorRef])

  return (
    <div class="book-editor">
      <BookContentEditorHeader editorRef={editor} editorUpdateEvent={editorUpdateEvent.current} />
      <article class="ol-container">
        <div class="__user-content book-editor__content" ref={root} />
      </article>
    </div>
  )
}

function createEditor(editorElement: HTMLElement, options?: Partial<EditorOptions>) {
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
    ],
    ...options,
  })
}
