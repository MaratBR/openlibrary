import { Editor } from '@tiptap/core'
import { ReactNode, MouseEventHandler } from 'react'
import { EditorState } from './EditorState'

export function Toolbar({ editor }: { editor: Editor }) {
  const { bold, italic, strikethrough, textAlign } = EditorState.useEditorState(editor)

  return (
    <ul className="OlSimpleEditor-toolbar">
      <ToolbarButton
        label={window._('common.richText.bold')}
        active={bold}
        onClick={() => editor.chain().toggleBold().focus().run()}
      >
        <i className="fa-solid fa-bold" />
      </ToolbarButton>
      <ToolbarButton
        label={window._('common.richText.italic')}
        active={italic}
        onClick={() => editor.chain().toggleItalic().focus().run()}
      >
        <i className="fa-solid fa-italic" />
      </ToolbarButton>
      <ToolbarButton
        label={window._('common.richText.strikethrough')}
        active={strikethrough}
        onClick={() => editor.chain().toggleStrike().focus().run()}
      >
        <i className="fa-solid fa-strikethrough" />
      </ToolbarButton>
      <li className="OlSimpleEditor-delimiter" aria-hidden="true" />
      <ToolbarButton
        label={window._('common.richText.left')}
        active={textAlign === 'left'}
        onClick={() => editor.chain().focus().setTextAlign('left').run()}
      >
        <i className="fa-solid fa-align-left" />
      </ToolbarButton>
      <ToolbarButton
        label={window._('common.richText.center')}
        active={textAlign === 'center'}
        onClick={() => editor.chain().focus().setTextAlign('center').run()}
      >
        <i className="fa-solid fa-align-center" />
      </ToolbarButton>
      <ToolbarButton
        label={window._('common.richText.right')}
        active={textAlign === 'right'}
        onClick={() => editor.chain().focus().setTextAlign('right').run()}
      >
        <i className="fa-solid fa-align-right" />
      </ToolbarButton>
      <ToolbarButton
        label={window._('common.richText.justify')}
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
  label,
  onClick,
  children,
}: {
  active: boolean
  label: string
  onClick: MouseEventHandler<HTMLButtonElement>
  children: ReactNode
}) {
  return (
    <li>
      <button
        type="button"
        aria-label={label}
        title={label}
        aria-pressed={active}
        className={`OlSimpleEditor-btn ${active ? 'OlSimpleEditor-btn--active' : ''}`}
        onClick={onClick}
      >
        {children}
      </button>
    </li>
  )
}
