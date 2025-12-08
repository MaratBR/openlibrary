import { Editor, EditorEvents } from '@tiptap/core'
import HorizontalRule from '@tiptap/extension-horizontal-rule'
import { TextStyle, FontSize, FontFamily } from '@tiptap/extension-text-style'
import Typography from '@tiptap/extension-typography'
import Color from '@tiptap/extension-color'
import StarterKit from '@tiptap/starter-kit'

import TextAlign from '@tiptap/extension-text-align'
import { ComponentChildren, render } from 'preact'
import { Subject, useSubject } from '@/common/rx'
import { MouseEventHandler } from 'preact/compat'

export type State = {
  bold: boolean
  italic: boolean
  strikethrough: boolean
  color: string | null
  header: number | null
  font: string | null
  fontSize: string | null
  textAlign: 'left' | 'right' | 'center' | 'justify' | null
}

const DEFAULT_STATE: State = {
  bold: false,
  italic: false,
  strikethrough: false,
  color: null,
  header: null,
  font: null,
  fontSize: null,
  textAlign: null,
}

export class SimpleEditor extends Editor {
  tiptapState = new Subject<State>(DEFAULT_STATE)

  constructor(element: HTMLElement) {
    const html = element.innerHTML
    element.classList.add('SimpleEditor')
    element.classList.add('user-content')

    const contentElement = document.createElement('div')
    contentElement.classList.add('SimpleEditor__content')

    super({
      element: contentElement,
      content: html,
      extensions: [
        StarterKit.configure({
          horizontalRule: false,
          codeBlock: false,
          heading: false,
          code: { HTMLAttributes: { class: 'inline', spellcheck: 'false' } },
          dropcursor: { width: 2, class: 'ProseMirror-dropcursor border' },
        }),
        TextStyle,
        FontSize,
        FontFamily,
        Typography,
        HorizontalRule,
        Color,
        TextAlign,
      ],
    })

    const toolbarWrapper = document.createElement('div')

    window.requestAnimationFrame(() => {
      element.innerHTML = ''
      element.appendChild(toolbarWrapper)
      element.appendChild(contentElement)
      render(<Toolbar editor={this} />, toolbarWrapper)
    })

    this.on('update', this._onUpdate.bind(this))
  }

  private _getCurrentState(): State {
    const textStyle = this.getAttributes('textStyle')

    let textAlign: State['textAlign'] = 'left'

    if (this.isActive({ textAlign: 'right' })) {
      textAlign = 'right'
    } else if (this.isActive({ textAlign: 'center' })) {
      textAlign = 'center'
    } else if (this.isActive({ textAlign: 'justify' })) {
      textAlign = 'justify'
    }

    return {
      bold: this.isActive('bold'),
      italic: this.isActive('italic'),
      strikethrough: this.isActive('strikethrough'),
      color: textStyle.color || null,
      header: this.isActive('heading') ? this.getAttributes('heading').level : null,
      font: textStyle.fontFamily || null,
      fontSize: textStyle.fontSize || null,
      textAlign,
    }
  }

  private _onUpdate(_props: EditorEvents['update']) {
    this.tiptapState.set(this._getCurrentState())
  }
}

function Toolbar({ editor }: { editor: SimpleEditor }) {
  const { bold, italic, strikethrough } = useSubject(editor.tiptapState)

  return (
    <ul class="SimpleEditor__toolbar">
      <ToolbarButton active={bold} onClick={() => editor.chain().setBold().focus().run()}>
        <i class="fa-solid fa-bold" />
      </ToolbarButton>
      <ToolbarButton active={bold} onClick={() => editor.chain().setItalic().focus().run()}>
        <i class="fa-solid fa-italic" />
      </ToolbarButton>
      <ToolbarButton active={bold} onClick={() => editor.chain().setStrike().focus().run()}>
        <i class="fa-solid fa-strikethrough" />
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
      class={`SimpleEditor__btn ${active ? 'SimpleEditor__btn--actuive' : ''}`}
      onClick={onClick}
    >
      {children}
    </li>
  )
}
