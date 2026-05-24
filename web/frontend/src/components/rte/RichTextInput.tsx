import { Editor, EditorContent } from '@tiptap/react'
import { Toolbar } from './Toolbar'

export type RichTextInputProps = {
  editor: Editor
}

export function RichTextInput({ editor }: RichTextInputProps) {
  return (
    <div class="SimpleEditor">
      <Toolbar editor={editor} />
      <EditorContent editor={editor} className="SimpleEditor__content user-content" />
    </div>
  )
}
