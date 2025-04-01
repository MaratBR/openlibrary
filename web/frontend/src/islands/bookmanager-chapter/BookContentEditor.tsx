import { Editor, EditorOptions } from '@tiptap/core'
import HorizontalRule from '@tiptap/extension-horizontal-rule'
import TextStyle from '@tiptap/extension-text-style'
import Typography from '@tiptap/extension-typography'
import TextAlign from '@tiptap/extension-text-align'

import StarterKit from '@tiptap/starter-kit'
import { useEffect, useRef } from 'preact/hooks'
import './BookManagerEditor.scss'
import Heading from '@tiptap/extension-heading'
import { createEvent } from '@/lib/event'
import BookContentEditorHeader from './BookContentEditorHeader'
import Underline from '@tiptap/extension-underline'

export type BookContentEditorProps = {
  content: string
}

export default function BookContentEditor({ content }: BookContentEditorProps) {
  const editor = useRef<Editor | null>(null)
  const root = useRef<HTMLDivElement | null>(null)

  const editorUpdateEvent = useRef(createEvent<Editor>())

  const propsRef = useRef({ content })
  propsRef.current.content = content

  useEffect(() => {
    if (root.current === null) throw new Error('root is null')
    editor.current = createEditor(root.current, {
      content: propsRef.current.content,
      onUpdate: ({ editor }) => {
        editorUpdateEvent.current.fire(editor)
      },
      onTransaction: ({ editor }) => {
        editorUpdateEvent.current.fire(editor)
      },
    })
  }, [])

  return (
    <div class="ol-book-editor">
      <BookContentEditorHeader editorRef={editor} editorUpdateEvent={editorUpdateEvent.current} />
      <div class="__user-content ol-book-editor__content ol-container" ref={root} />
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
    ],
    ...options,
  })
}
