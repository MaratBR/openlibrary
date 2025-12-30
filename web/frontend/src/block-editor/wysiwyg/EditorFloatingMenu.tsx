import { Editor } from '@tiptap/react'
import { FloatingMenu } from '@tiptap/react/menus'

export default function EditorFloatingMenu({ editor }: { editor: Editor }) {
  return (
    <FloatingMenu
      shouldShow={({ state }) => {
        const { $from } = state.selection
        const parent = $from.parent

        return parent.type.name === 'paragraph' && parent.content.size === 0
      }}
      editor={editor}
      class="be-floating-menu"
    >
      {window._('editor.floatingMenu')}
    </FloatingMenu>
  )
}
