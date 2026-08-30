import { Editor } from '@tiptap/core'
import { ReactNode, MouseEventHandler } from 'react'
import { EditorState } from './EditorState'

export function Toolbar({ editor }: { editor: Editor }) {
  const { bold, italic, strikethrough, textAlign } = EditorState.useEditorState(editor)

  return (
    <ul className="OlSimpleEditor-toolbar">
      <ToolbarButton active={bold} onClick={() => editor.chain().toggleBold().focus().run()}>
        <i className="fa-solid fa-bold" />
      </ToolbarButton>
      <ToolbarButton active={italic} onClick={() => editor.chain().toggleItalic().focus().run()}>
        <i className="fa-solid fa-italic" />
      </ToolbarButton>
      <ToolbarButton
        active={strikethrough}
        onClick={() => editor.chain().toggleStrike().focus().run()}
      >
        <i className="fa-solid fa-strikethrough" />
      </ToolbarButton>
      <li className="OlSimpleEditor-delimiter" aria-hidden="true" />
      <ToolbarButton
        active={textAlign === 'left'}
        onClick={() => editor.chain().focus().setTextAlign('left').run()}
      >
        <i className="fa-solid fa-align-left" />
      </ToolbarButton>
      <ToolbarButton
        active={textAlign === 'center'}
        onClick={() => editor.chain().focus().setTextAlign('center').run()}
      >
        <i className="fa-solid fa-align-center" />
      </ToolbarButton>
      <ToolbarButton
        active={textAlign === 'right'}
        onClick={() => editor.chain().focus().setTextAlign('right').run()}
      >
        <i className="fa-solid fa-align-right" />
      </ToolbarButton>
      <ToolbarButton
        active={textAlign === 'justify'}
        onClick={() => editor.chain().focus().setTextAlign('justify').run()}
      >
        <i className="fa-solid fa-align-justify" />
      </ToolbarButton>
    </ul>
  )
}

function ToolbarButton({
  active,
  onClick,
  children,
}: {
  active: boolean
  onClick: MouseEventHandler<HTMLLIElement>
  children: ReactNode
}) {
  return (
    <li
      role="button"
      className={`OlSimpleEditor-btn ${active ? 'OlSimpleEditorBtn--active' : ''}`}
      onClick={onClick}
    >
      {children}
    </li>
  )
}
