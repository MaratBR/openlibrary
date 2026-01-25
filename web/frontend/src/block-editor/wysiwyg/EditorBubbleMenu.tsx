import { BubbleMenu } from '@tiptap/react/menus'

import './EditorBubbleMenu.scss'
import { ChapterContentEditor, useEditorToolbarState } from './editor'
import EditorToggleButton from './EditorToggleButton'

export default function EditorBubbleMenu({
  editor,
  appendTo,
}: {
  editor: ChapterContentEditor
  appendTo?: HTMLElement
}) {
  const { bold, italic, strikethrough, textAlign } = useEditorToolbarState(editor)

  return (
    <BubbleMenu
      class="be-bubble-menu"
      // getReferencedVirtualElement={() => {
      //   const textElement = getSelectedTextElement(editor)
      //   return textElement
      // }}
      options={{
        placement: 'top-start',
      }}
      editor={editor}
      appendTo={appendTo}
    >
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
      <div class="be-toggle-group">
        <EditorToggleButton
          active={textAlign === 'left'}
          onClick={() => editor.chain().focus().setTextAlign('left').run()}
        >
          <i class="fa-solid fa-align-left" />
        </EditorToggleButton>
        <EditorToggleButton
          active={textAlign === 'center'}
          onClick={() => editor.chain().focus().setTextAlign('center').run()}
        >
          <i class="fa-solid fa-align-center" />
        </EditorToggleButton>
        <EditorToggleButton
          active={textAlign === 'right'}
          onClick={() => editor.chain().focus().setTextAlign('right').run()}
        >
          <i class="fa-solid fa-align-right" />
        </EditorToggleButton>
        <EditorToggleButton
          active={textAlign === 'justify'}
          onClick={() => editor.chain().focus().setTextAlign('justify').run()}
        >
          <i class="fa-solid fa-align-justify" />
        </EditorToggleButton>
      </div>
    </BubbleMenu>
  )
}
