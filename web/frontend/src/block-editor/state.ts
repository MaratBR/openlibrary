import { create } from 'zustand/react'
import { useWYSIWYG } from './wysiwyg'
import { httpUpdateAndPublishDraft, httpUpdateDraft, httpUpdateDraftChapterName } from '@/api/bm'
import { DraftDto } from './contracts'
import { useWYSIWYGHasChanges } from './wysiwyg/state'

export type BEState = {
  saving: boolean
  autoSave: boolean
  draft: DraftDto | null
  chapterName: string
  error: unknown | null

  chapterNameWasChanged(): boolean

  init(draft: DraftDto): void
  setChapterName(name: string): void
  saveDraft(): Promise<void>
  saveAndPublishDraft(makePublic: boolean): Promise<void>
}

export const useBEState = create<BEState>((set, get) => ({
  saving: false,
  autoSave: true,
  draft: null,
  chapterName: '',
  error: null,

  chapterNameWasChanged() {
    const { chapterName, draft } = get()
    if (!draft) return false
    return chapterName !== draft.chapterName
  },

  init(draft) {
    set({
      draft,
      chapterName: draft.chapterName,
    })

    const wysiwyg = useWYSIWYG.getState()
    wysiwyg.setInitialContent(draft.content)
  },

  setChapterName(name) {
    set({
      chapterName: name,
    })
  },

  async saveDraft() {
    set({
      saving: true,
    })
    try {
      const { draft, chapterName, chapterNameWasChanged } = get()
      if (!draft) throw new Error('cannot save draft - no draft information is available')

      // first update chapter name if necessary
      if (chapterNameWasChanged()) {
        const response = await httpUpdateDraftChapterName(
          draft.book.id,
          draft.chapterId,
          draft.id,
          chapterName,
        )
        response.throwIfError()
      }

      const wysiwyg = useWYSIWYG.getState()
      const content = wysiwyg.getContent()
      const response = await httpUpdateDraft(draft.book.id, draft.chapterId, draft.id, content)
      response.throwIfError()
      wysiwyg.markContentAsFresh()

      set({ saving: false })
    } catch (error: unknown) {
      set({
        error,
        saving: false,
      })
    }
  },

  async saveAndPublishDraft(makePublic: boolean) {
    set({
      saving: true,
    })
    try {
      const { draft, chapterName, chapterNameWasChanged } = get()
      if (!draft) throw new Error('cannot save draft - no draft information is available')

      // first update chapter name if necessary
      if (chapterNameWasChanged()) {
        const response = await httpUpdateDraftChapterName(
          draft.book.id,
          draft.chapterId,
          draft.id,
          chapterName,
        )
        response.throwIfError()
      }

      const wysiwyg = useWYSIWYG.getState()
      const content = wysiwyg.getContent()
      const response = await httpUpdateAndPublishDraft(
        draft.book.id,
        draft.chapterId,
        draft.id,
        content,
        makePublic,
      )
      response.throwIfError()
      wysiwyg.markContentAsFresh()

      set({ saving: false })
    } catch (error: unknown) {
      set({
        error,
        saving: false,
      })
    }
  },
}))

export function useDraftHasChanges() {
  const contentWasChanged = useWYSIWYGHasChanges()
  const nameChanged = useBEState((s) => s.chapterNameWasChanged())
  return nameChanged || contentWasChanged
}

export function useDraftHasNewerRevision() {
  return true
}
