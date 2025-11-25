import { Editor } from '@tiptap/core'
import HorizontalRule from '@tiptap/extension-horizontal-rule'
import TextStyle from '@tiptap/extension-text-style'
import Typography from '@tiptap/extension-typography'
import TextAlign from '@tiptap/extension-text-align'
import Image from '@tiptap/extension-image'
import StarterKit from '@tiptap/starter-kit'
import Heading from '@tiptap/extension-heading'
import Underline from '@tiptap/extension-underline'
import Placeholder from '@tiptap/extension-placeholder'
import { ChapterState } from '../state'

export type ChapterContentEditorOptions = {
  element: HTMLElement
  placeholder: string
  state: ChapterState
}

export class ChapterContentEditor extends Editor {
  constructor({ element, placeholder }: ChapterContentEditorOptions) {
    super({
      element,
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
          placeholder,
        }),
      ],
      // ...options,
    })
  }
}
