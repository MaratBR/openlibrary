import { Content, Editor } from '@tiptap/core'
import HorizontalRule from '@tiptap/extension-horizontal-rule'
import { TextStyle } from '@tiptap/extension-text-style'
import Typography from '@tiptap/extension-typography'
import TextAlign from '@tiptap/extension-text-align'
import Image from '@tiptap/extension-image'
import StarterKit from '@tiptap/starter-kit'
import Heading from '@tiptap/extension-heading'
import { BulletList, OrderedList } from '@tiptap/extension-list'
import Underline from '@tiptap/extension-underline'
import Placeholder from '@tiptap/extension-placeholder'
import { EditorContent } from '@tiptap/react'
import EditorFloatingMenu from './EditorFloatingMenu'
import EditorBubbleMenu from './EditorBubbleMenu'
import { Subject, useSubject } from '@/common/rx'
import { createPortal } from 'preact/compat'
import { SlashCommand } from './Suggestions'
import { slashCommands } from './slashCommands'
import { SuggestionsDisplay } from './SuggestionsDisplay'
import { EditorElements } from './EditorElements'
import { createEvent } from '@/lib/event'

export type EditorToolbarState = {
  bold: boolean
  italic: boolean
  strikethrough: boolean
  color: string | null
  textType: 'h1' | 'h2' | 'h3' | 'ol' | 'ul' | 'text'
  font: string | null
  fontSize: string | null
  textAlign: 'left' | 'right' | 'center' | 'justify' | null
}

const DEFAULT_STATE: EditorToolbarState = {
  bold: false,
  italic: false,
  strikethrough: false,
  color: null,
  textType: 'text',
  font: null,
  fontSize: null,
  textAlign: 'left',
}

export class ChapterContentEditor extends Editor {
  private _placeholder = ''
  private elements: EditorElements

  public readonly firstChange = createEvent<void>()
  public readonly wasChangedFirstTime = new Subject<boolean | null>(null)

  constructor(elements: EditorElements) {
    super({
      content: '',
      extensions: [
        StarterKit.configure({
          horizontalRule: false,
          codeBlock: false,
          heading: false,
          code: { HTMLAttributes: { class: 'inline', spellcheck: 'false' } },
          dropcursor: { width: 2, class: 'ProseMirror-dropcursor border' },
        }),
        TextStyle,
        Typography,
        HorizontalRule,
        Heading,
        TextAlign.configure({
          types: ['heading', 'paragraph'],
        }),
        Underline,
        Image.configure({
          inline: true,
        }),
        Placeholder.configure({
          placeholder: ({ editor }) => {
            if (editor instanceof ChapterContentEditor) {
              return editor._placeholder
            }
            return ''
          },
        }),
        BulletList,
        OrderedList,
        SlashCommand.configure({
          suggestionClass: 'be-suggestion',
          commands: slashCommands(),
          displayAdapter: new SuggestionsDisplay(elements, () => this),
        }),
      ],
    })

    this.elements = elements

    const onTransaction = () => {
      this.toolbarState.set(this.getCurrentToolbarState())
    }
    this.toolbarState.set(this.getCurrentToolbarState())

    this.on('transaction', onTransaction)

    const onUpdate = () => {
      if (this.wasChangedFirstTime.get() !== false) {
        return
      }
      this.wasChangedFirstTime.set(true)
    }
    this.on('update', onUpdate)

    const onDestroy = () => {
      this.off('transaction', onTransaction)
      this.off('destroy', onDestroy)
      this.off('update', onUpdate)
    }
    this.on('destroy', onDestroy)
  }

  public setPlaceholder(placeholder: string) {
    this._placeholder = placeholder
  }

  public setContentAndClearHistory(content: Content) {
    this.chain().setMeta('addToHistory', false).insertContent(content).run()
    this.wasChangedFirstTime.set(false)
  }

  public getCurrentToolbarState(): EditorToolbarState {
    const textStyle = this.getAttributes('textStyle')

    let textAlign: EditorToolbarState['textAlign'] = 'left'

    if (this.isActive({ textAlign: 'right' })) {
      textAlign = 'right'
    } else if (this.isActive({ textAlign: 'center' })) {
      textAlign = 'center'
    } else if (this.isActive({ textAlign: 'justify' })) {
      textAlign = 'justify'
    }

    const headerLevel = this.isActive('heading') ? this.getAttributes('heading').level : null

    let textType: EditorToolbarState['textType'] = 'text'

    if (typeof headerLevel === 'number') {
      switch (headerLevel) {
        case 1:
          textType = 'h1'
          break
        case 2:
          textType = 'h2'
          break
        case 3:
          textType = 'h3'
          break
      }
    }

    return {
      bold: this.isActive('bold'),
      italic: this.isActive('italic'),
      strikethrough: this.isActive('strike'),
      color: textStyle.color || null,
      textType,
      font: textStyle.fontFamily || null,
      fontSize: textStyle.fontSize || null,
      textAlign,
    }
  }

  toolbarState = new Subject<EditorToolbarState>(DEFAULT_STATE)

  public getContentElement() {
    return (
      <>
        {createPortal(
          <>
            <EditorFloatingMenu editor={this} />
            <EditorContent editor={this} />
          </>,
          this.elements.content,
        )}
        <EditorBubbleMenu editor={this} appendTo={this.elements.contentWrapper} />
      </>
    )
  }
}

export function useEditorToolbarState(editor: ChapterContentEditor) {
  return useSubject(editor.toolbarState)
}
