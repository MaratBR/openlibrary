import { create } from 'zustand/react'
import {
  httpScheduleDraft,
  httpUpdateAndPublishDraft,
  httpUpdateDraft,
  httpUpdateDraftChapterName,
} from '@/api/bm'
import { DraftDto } from './contracts'
import { useWYSIWYG, useWYSIWYGHasChanges } from './wysiwyg/state'

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
  saveAndScheduleDraft(scheduledAt: Date): Promise<void>
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
      saving: false,
      error: null,
    })
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
          draft.chapter.id,
          draft.id,
          chapterName,
        )
        response.throwIfError()
      }

      const wysiwyg = useWYSIWYG.getState()
      const content = wysiwyg.getContent()
      const response = await httpUpdateDraft(draft.book.id, draft.chapter.id, draft.id, content)
      response.throwIfError()
      wysiwyg.markContentAsFresh()

      set({ saving: false, draft: response.data })
    } catch (error: unknown) {
      set({
        error,
        saving: false,
      })
      throw error
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
          draft.chapter.id,
          draft.id,
          chapterName,
        )
        response.throwIfError()
      }

      const wysiwyg = useWYSIWYG.getState()
      const content = wysiwyg.getContent()
      const response = await httpUpdateAndPublishDraft(
        draft.book.id,
        draft.chapter.id,
        draft.id,
        content,
        makePublic,
      )
      response.throwIfError()
      wysiwyg.markContentAsFresh()

      set({ saving: false, draft: response.data })
    } catch (error: unknown) {
      set({
        error,
        saving: false,
      })
      throw error
    }
  },

  async saveAndScheduleDraft(scheduledAt: Date) {
    await get().saveDraft()
    set({ saving: true })
    try {
      const { draft } = get()
      if (!draft) throw new Error('cannot schedule draft - no draft information is available')
      const response = await httpScheduleDraft(
        draft.book.id,
        draft.chapter.id,
        draft.id,
        scheduledAt.toISOString(),
      )
      response.throwIfError()
      set({ saving: false, error: null, draft: response.data })
    } catch (error: unknown) {
      set({ error, saving: false })
      throw error
    }
  },
}))

export function useDraftHasChanges() {
  const contentWasChanged = useWYSIWYGHasChanges()
  const chapterNameChanged = useBEState((s) => s.chapterNameWasChanged())
  return chapterNameChanged || contentWasChanged
}

export function useChapterNameIsValid() {
  return useBEState((s) => {
    const name = s.chapterName.trim()
    return name.length > 0 && Array.from(name).length <= 70
  })
}

export function useDraftHasPendingChanges() {
  return useBEState((s) => {
    const draft = s.draft
    if (!draft) return false
    return new Date(draft.updatedAt ?? draft.createdAt) > new Date(draft.chapter.contentUpdatedAt)
  })
}

export function useDraftHasNewerRevision() {
  return true
}
