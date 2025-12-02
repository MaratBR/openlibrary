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
import { Dispose } from '@/common/rx'

export type ChapterContentEditorOptions = {
  element: HTMLIFrameElement
  placeholder: string
  state: ChapterState
}

export class ChapterContentEditor extends Editor {
  private readonly chapterState: ChapterState
  private _disposeFn: Dispose[] = []
  private readonly $chapterContent: HTMLElement

  constructor({ element, placeholder, state }: ChapterContentEditorOptions) {
    if (!element.contentDocument)
      throw new Error('iframe element does not have contentDocument yet')
    const chapterContentWrap = element.contentDocument.getElementById('ChapterContent')
    if (!chapterContentWrap) throw new Error('cannot find #ChapterContent content within iframe')
    const userContentContainer = element.contentDocument.getElementById('BookReaderContent')
    if (!userContentContainer)
      throw new Error('cannot find #BookReaderContent content within iframe')

    super({
      element: userContentContainer,
      content: state.draft.content,
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
    })
    this.chapterState = state
    this.$chapterContent = chapterContentWrap
    this._initViewState()
  }

  private _initViewState() {
    this._disposeFn.push(
      this.chapterState.view.subscribe(({ editorWidth }) => {
        this.$chapterContent.style.maxWidth = editorWidth
      }),
    )
  }
}
