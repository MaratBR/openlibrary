import { create } from 'zustand/react'
import { ChapterContentEditor } from './editor'
import { EditorElements } from './EditorElements'
import { JSX } from 'preact/jsx-runtime'

export type WYSIWYGState = {
  initData: {
    editor: ChapterContentEditor
    elements: EditorElements
  } | null
  initialContent: string

  init(iframe: HTMLIFrameElement): void
  renderContent(): JSX.Element | null

  setInitialContent(content: string): void
}

export const useWYSIWYG = create<WYSIWYGState>(() => ({
  initData: null,

  init(iframe) {
    const elements = new EditorElements(iframe)

    this.initData = {
      editor: new ChapterContentEditor(elements),
      elements,
    }
    if (this.initialContent) {
      this.initData.editor.setContentAndClearHistory(this.initialContent)
    }
  },

  renderContent() {
    if (!this.initData) return null

    return this.initData.editor.getContentElement()
  },

  initialContent: '',

  setInitialContent(content) {
    if (this.initData) {
      this.initData.editor.setContentAndClearHistory(content)
    }
    this.initialContent = content
  },
}))
