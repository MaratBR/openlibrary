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
import History from '@tiptap/extension-history'
import TextAlign from '@tiptap/extension-text-align'
import { useEditor } from '@tiptap/react'

export type OLTiptapEditorOptions = {
  element?: HTMLElement
  content?: string
}

export type UseOLEditorOptions = {
  content?: string
}

export function useOLEditor(options: UseOLEditorOptions = {}) {
  return useEditor({
    content: options.content,
    extensions: [
      Document,
      History,
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
}
