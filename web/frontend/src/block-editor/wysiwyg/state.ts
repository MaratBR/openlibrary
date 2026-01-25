import { create } from 'zustand/react'
import { ChapterContentEditor } from './editor'

export type WYSIWYGState = {
  editor: ChapterContentEditor | null
  contentModified: boolean

  init(editor: ChapterContentEditor): void
  getContent(): string
  markContentAsFresh(): void
}

export const useWYSIWYG = create<WYSIWYGState>((set, get) => ({
  editor: null,
  contentModified: false,

  init(editor) {
    set({ editor })

    const onUpdate = () => {
      if (!get().contentModified) set({ contentModified: true })
    }

    const onDestroy = () => {
      set({ editor: null, contentModified: false })
      editor.off('update', onUpdate)
      editor.off('destroy', onDestroy)
    }

    editor.on('destroy', onDestroy)
    editor.on('update', onUpdate)
  },

  getContent() {
    if (!this.editor) {
      return ''
    }

    const html = this.editor.getHTML()

    return html
  },

  markContentAsFresh() {
    set({ contentModified: false })
  },
}))

export function useWYSIWYGHasChanges() {
  const contentModified = useWYSIWYG((s) => s.contentModified)
  return contentModified
}
