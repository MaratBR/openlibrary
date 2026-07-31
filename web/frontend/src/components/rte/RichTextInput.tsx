import { Editor, EditorContent } from '@tiptap/react'
import { Toolbar } from './Toolbar'
import type { ReactNode } from 'react'

export type RichTextInputProps = {
  editor: Editor
  className?: string
  toolbar?: ReactNode | ((editor: Editor) => ReactNode)
}

export function RichTextInput({ editor, className = '', toolbar }: RichTextInputProps) {
  return (
    <div className={`ol-simple-editor ${className}`.trim()}>
      {typeof toolbar === 'function' ? toolbar(editor) : (toolbar ?? <Toolbar editor={editor} />)}
      <EditorContent editor={editor} className="ol-simple-editor__content user-content" />
    </div>
  )
}
