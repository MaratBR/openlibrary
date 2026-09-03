import { BubbleMenu } from '@tiptap/react/menus'

import './EditorBubbleMenu.scss'
import { ChapterContentEditor, useEditorToolbarState } from '../editor'
import EditorToggleButton from './EditorToggleButton'
import { ColorPallette, ColorPalletteSection } from './ColorPallette'
import { FontPicker, Fonts } from './Fonts'

export function EditorBubbleMenu({
  editor,
  appendTo,
}: {
  editor: ChapterContentEditor
  appendTo?: HTMLElement
}) {
  const { bold, italic, strikethrough, textAlign } = useEditorToolbarState(editor)

  return (
    <BubbleMenu
      className="BeBubbleMenu"
      // getReferencedVirtualElement={() => {
      //   const textElement = getSelectedTextElement(editor)
      //   return textElement
      // }}
      options={{
        placement: 'top-start',
        flip: { padding: 8 },
        shift: { padding: 8 },
        size: {
          padding: 8,
          apply({ availableWidth, elements }) {
            elements.floating.style.maxWidth = `${availableWidth}px`
          },
        },
      }}
      editor={editor}
      appendTo={appendTo}
    >
      <div className="BeBubbleMenu-main">
        <div className="BeToggleGroup">
          <EditorToggleButton
            active={bold}
            onClick={() => editor.chain().focus().toggleBold().run()}
          >
            <i className="fa-solid fa-bold" />
          </EditorToggleButton>
          <EditorToggleButton
            active={italic}
            onClick={() => editor.chain().focus().toggleItalic().run()}
          >
            <i className="fa-solid fa-italic" />
          </EditorToggleButton>
          <EditorToggleButton
            active={strikethrough}
            onClick={() => editor.chain().focus().toggleStrike().run()}
          >
            <i className="fa-solid fa-strikethrough" />
          </EditorToggleButton>
        </div>
        <div className="BeBubbleMenu-delimiter" />
        <div className="BeToggleGroup">
          <EditorToggleButton
            active={textAlign === 'left'}
            onClick={() => editor.chain().focus().setTextAlign('left').run()}
          >
            <i className="fa-solid fa-align-left" />
          </EditorToggleButton>
          <EditorToggleButton
            active={textAlign === 'center'}
            onClick={() => editor.chain().focus().setTextAlign('center').run()}
          >
            <i className="fa-solid fa-align-center" />
          </EditorToggleButton>
          <EditorToggleButton
            active={textAlign === 'right'}
            onClick={() => editor.chain().focus().setTextAlign('right').run()}
          >
            <i className="fa-solid fa-align-right" />
          </EditorToggleButton>
          <EditorToggleButton
            active={textAlign === 'justify'}
            onClick={() => editor.chain().focus().setTextAlign('justify').run()}
          >
            <i className="fa-solid fa-align-justify" />
          </EditorToggleButton>
        </div>
        <div className="BeToggleGroup">
          <ColorPallette />
          <Fonts />
        </div>
      </div>

      <ColorPalletteSection editor={editor} />
      <FontPicker editor={editor} />
    </BubbleMenu>
  )
}
