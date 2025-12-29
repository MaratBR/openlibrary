import { Editor } from '@tiptap/react'
import { FloatingMenu } from '@tiptap/react/menus'

export default function EditorFloatingMenu({ editor }: { editor: Editor }) {
  return <FloatingMenu editor={editor}>Floating menu</FloatingMenu>
}
