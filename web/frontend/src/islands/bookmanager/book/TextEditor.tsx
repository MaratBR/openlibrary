import { Editor } from '@tiptap/core'
import { RefObject } from 'preact'
import Paragraph from '@tiptap/extension-paragraph'
import Document from '@tiptap/extension-document'
import Bold from '@tiptap/extension-bold'
import Italic from '@tiptap/extension-italic'
import Strike from '@tiptap/extension-strike'
import Underline from '@tiptap/extension-underline'
import Text from '@tiptap/extension-text'

import { useRef, useEffect } from 'preact/hooks'

interface TextEditorProps {
  value: RefObject<string>
}

const TextEditor = ({ value }: TextEditorProps) => {
  const rootRef = useRef<HTMLDivElement>(null)
  const editor = useRef<Editor>()

  useEffect(() => {
    editor.current = new Editor({
      content: value.current,
      element: rootRef.current!,
      extensions: [Document, Paragraph, Bold, Italic, Strike, Underline, Text],
      onTransaction(props) {
        value.current = props.editor.getHTML()
      },
    })

    return () => editor.current?.destroy()
  }, [value])

  return (
    <div className="border-gray-400 dark:border-gray-600 border rounded-lg shadow-sm">
      <div className="border-b p-2 bg-background flex gap-1 rounded-t-lg">
        <button
          type="button"
          onClick={() => editor.current?.chain().toggleBold().run()}
          className="size-8 rounded hover:bg-secondary"
          title="Bold"
        >
          <strong>B</strong>
        </button>
        <button
          type="button"
          onClick={() => editor.current?.chain().focus().toggleItalic().run()}
          className="size-8 rounded hover:bg-secondary"
          title="Italic"
        >
          <em>I</em>
        </button>
        <button
          type="button"
          onClick={() => editor.current?.chain().focus().toggleStrike().run()}
          className="size-8 rounded hover:bg-secondary"
          title="Strike-through"
        >
          <span className="line-through">S</span>
        </button>
        <button
          type="button"
          onClick={() => editor.current?.chain().focus().toggleUnderline().run()}
          className="size-8 rounded hover:bg-secondary"
          title="Underline"
        >
          <span className="underline">U</span>
        </button>
      </div>
      <div
        ref={rootRef}
        className={`user-content [&>.tiptap]:min-h-32 [&>.tiptap]:p-4 [&>.tiptap]:!outline-none`}
      />
    </div>
  )
}

export default TextEditor
