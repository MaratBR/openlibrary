import { Editor } from '@tiptap/core'
import { ComponentChildren, MouseEventHandler } from 'preact'
import { EditorState } from './EditorState'

export function Toolbar({ editor }: { editor: Editor }) {
  const { bold, italic, strikethrough, textAlign } = EditorState.useEditorState(editor)

  return (
    <ul class="SimpleEditor__toolbar">
      <ToolbarButton active={bold} onClick={() => editor.chain().toggleBold().focus().run()}>
        <i class="fa-solid fa-bold" />
      </ToolbarButton>
      <ToolbarButton active={italic} onClick={() => editor.chain().toggleItalic().focus().run()}>
        <i class="fa-solid fa-italic" />
      </ToolbarButton>
      <ToolbarButton
        active={strikethrough}
        onClick={() => editor.chain().toggleStrike().focus().run()}
      >
        <i class="fa-solid fa-strikethrough" />
      </ToolbarButton>
      <li class="SimpleEditor__delimiter" aria-hidden="true" />
      <ToolbarButton
        active={textAlign === 'left'}
        onClick={() => editor.chain().focus().setTextAlign('left').run()}
      >
        <i class="fa-solid fa-align-left" />
      </ToolbarButton>
      <ToolbarButton
        active={textAlign === 'center'}
        onClick={() => editor.chain().focus().setTextAlign('center').run()}
      >
        <i class="fa-solid fa-align-center" />
      </ToolbarButton>
      <ToolbarButton
        active={textAlign === 'right'}
        onClick={() => editor.chain().focus().setTextAlign('right').run()}
      >
        <i class="fa-solid fa-align-right" />
      </ToolbarButton>
      <ToolbarButton
        active={textAlign === 'justify'}
        onClick={() => editor.chain().focus().setTextAlign('justify').run()}
      >
        <i class="fa-solid fa-align-justify" />
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
  children: ComponentChildren
}) {
  return (
    <li
      role="button"
      class={`SimpleEditor__btn ${active ? 'SimpleEditor__btn--active' : ''}`}
      onClick={onClick}
    >
      {children}
    </li>
  )
}
