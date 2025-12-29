import { BubbleMenu } from '@tiptap/react/menus'

import './EditorBubbleMenu.scss'
import { ChapterContentEditor, useEditorToolbarState } from './editor'
import EditorToggleButton from './EditorToggleButton'
import TextFeatureSelector from './TextFeatureSelector'

export default function EditorBubbleMenu({
  editor,
  appendTo,
}: {
  editor: ChapterContentEditor
  appendTo: HTMLElement
}) {
  const { bold, italic, strikethrough } = useEditorToolbarState(editor)

  return (
    <BubbleMenu class="be-bubble-menu" editor={editor} appendTo={appendTo}>
      <div class="be-toggle-group">
        <EditorToggleButton active={bold} onClick={() => editor.chain().focus().toggleBold().run()}>
          <i class="fa-solid fa-bold" />
        </EditorToggleButton>
        <EditorToggleButton
          active={italic}
          onClick={() => editor.chain().focus().toggleItalic().run()}
        >
          <i class="fa-solid fa-italic" />
        </EditorToggleButton>
        <EditorToggleButton
          active={strikethrough}
          onClick={() => editor.chain().focus().toggleStrike().run()}
        >
          <i class="fa-solid fa-strikethrough" />
        </EditorToggleButton>
      </div>

      <div class="be-bubble-menu__delimiter" />

      <TextFeatureSelector editor={editor} />
    </BubbleMenu>
  )
}
