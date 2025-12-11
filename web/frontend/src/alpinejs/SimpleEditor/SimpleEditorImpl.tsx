import { Editor } from '@tiptap/core'
import HorizontalRule from '@tiptap/extension-horizontal-rule'
import { TextStyle, FontSize, FontFamily } from '@tiptap/extension-text-style'
import Typography from '@tiptap/extension-typography'
import Color from '@tiptap/extension-color'
import Paragraph from '@tiptap/extension-paragraph'
import Document from '@tiptap/extension-document'
import Bold from '@tiptap/extension-bold'
import Italic from '@tiptap/extension-italic'
import Strike from '@tiptap/extension-strike'
import Underline from '@tiptap/extension-underline'
import Text from '@tiptap/extension-text'

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
  textAlign: 'left',
}

export class SimpleEditor extends Editor {
  tiptapState = new Subject<State>(DEFAULT_STATE)

  private readonly _toolbarWrapper: HTMLElement

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
        Document,
        Paragraph,
        Bold,
        Italic,
        Strike,
        Underline,
        Text,
        TextStyle,
        FontSize,
        FontFamily,
        Typography,
        HorizontalRule,
        Color,
        TextAlign.configure({
          types: ['heading', 'paragraph'],
        }),
      ],
    })

    this._toolbarWrapper = document.createElement('div')

    window.requestAnimationFrame(() => {
      element.innerHTML = ''
      element.appendChild(this._toolbarWrapper)
      element.appendChild(contentElement)
      render(<Toolbar editor={this} />, this._toolbarWrapper)
    })

    const onUpdate = this._onUpdate.bind(this)

    this.on('update', onUpdate)
    this.on('transaction', onUpdate)

    this.on('destroy', () => {
      this.off('update', onUpdate)
      this.off('transaction', onUpdate)
    })
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

  private _onUpdate() {
    this.tiptapState.set(this._getCurrentState())
  }

  destroy(): void {
    super.destroy()
    render(null, this._toolbarWrapper)
  }
}

function Toolbar({ editor }: { editor: SimpleEditor }) {
  const { bold, italic, strikethrough, textAlign } = useSubject(editor.tiptapState)

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
