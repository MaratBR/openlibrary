import { atom } from 'jotai'
import { Effect } from 'effect'
import { ChapterContentEditor } from './editor'

export const wysiwygEditorAtom = atom<ChapterContentEditor | null>(null)
export const wysiwygContentModifiedAtom = atom(false)

export const mountWysiwygEditorAtom = atom(null, (_get, set, editor: ChapterContentEditor) =>
  Effect.sync(() => {
    const handleUpdate = () => set(wysiwygContentModifiedAtom, true)

    editor.on('update', handleUpdate)
    set(wysiwygEditorAtom, editor)
    set(wysiwygContentModifiedAtom, false)

    return () => {
      editor.off('update', handleUpdate)
      set(wysiwygEditorAtom, null)
      set(wysiwygContentModifiedAtom, false)
      editor.destroy()
    }
  }),
)
