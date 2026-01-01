import { create } from 'zustand/react'
import { ChapterContentEditor } from './editor'
import { EditorElements } from './EditorElements'
import { JSX } from 'preact/jsx-runtime'
import { useSubject } from '@/common/rx'

export type WYSIWYGState = {
  initData: {
    editor: ChapterContentEditor
    elements: EditorElements
  } | null
  initialContent: string
  contentModified: boolean

  init(iframe: HTMLIFrameElement): void
  getContentJSX(): JSX.Element | null
  setContentModified(modified: boolean): void
  setInitialContent(content: string): void
  getContent(): string
  markContentAsFresh(): void
}

export const useWYSIWYG = create<WYSIWYGState>((set, get) => ({
  initData: null,
  contentModified: false,

  init(iframe) {
    const elements = new EditorElements(iframe)

    const initData: WYSIWYGState['initData'] = {
      editor: new ChapterContentEditor(elements),
      elements,
    }
    const { initialContent } = get()
    initData.editor.setContentAndClearHistory(initialContent)
    set({ initData })
  },

  getContentJSX() {
    if (!this.initData) return null

    return this.initData.editor.getContentElement()
  },

  setContentModified(modified) {
    set({ contentModified: modified })
  },

  initialContent: '',

  setInitialContent(content) {
    if (this.initData) {
      this.initData.editor.setContentAndClearHistory(content)
    }
    set({ initialContent: content })
  },

  getContent() {
    if (!this.initData) {
      return ''
    }

    const html = this.initData.editor.getHTML()

    return html
  },

  markContentAsFresh() {
    if (!this.initData) throw new Error('WYSIWYG is not initialized')

    // mark current content as "unchanged"
    this.initData.editor.wasChangedFirstTime.set(false)
  },
}))

export function useWYSIWYGHasChanges() {
  const sub = useWYSIWYG((s) => s.initData?.editor.wasChangedFirstTime)
  return useSubject(sub) ?? false
}
