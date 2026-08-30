import { atom, useAtomValue } from 'jotai'
import type { Getter, Setter } from 'jotai'
import { Effect } from 'effect'
import { BookManagerApi } from '@/features/book-manager/api'
import { DraftDto } from './contracts'
import { wysiwygContentModifiedAtom, wysiwygEditorAtom } from './wysiwyg/state'

export const draftAtom = atom<DraftDto | null>(null)
export const chapterNameAtom = atom('')
export const savingAtom = atom(false)
export const saveErrorAtom = atom<unknown | null>(null)

export const initializeDraftAtom = atom(null, (_get, set, draft: DraftDto) => {
  set(draftAtom, draft)
  set(chapterNameAtom, draft.chapterName)
  set(savingAtom, false)
  set(saveErrorAtom, null)
})

export const chapterNameWasChangedAtom = atom((get) => {
  const draft = get(draftAtom)
  return draft !== null && get(chapterNameAtom) !== draft.chapterName
})

export const chapterNameIsValidAtom = atom((get) => {
  const name = get(chapterNameAtom).trim()
  return name.length > 0 && Array.from(name).length <= 70
})

export const draftHasPendingChangesAtom = atom((get) => {
  const draft = get(draftAtom)
  if (!draft) return false
  return new Date(draft.updatedAt ?? draft.createdAt) > new Date(draft.chapter.contentUpdatedAt)
})

// TODO: compare the draft's base revision with the current chapter revision.
export const draftHasNewerRevisionAtom = atom(true)

function saveDraft(get: Getter, set: Setter) {
  return Effect.gen(function* () {
    yield* Effect.sync(() => {
      set(savingAtom, true)
      set(saveErrorAtom, null)
    })
    const draft = get(draftAtom)
    if (!draft)
      return yield* Effect.fail(new Error('cannot save draft - no draft information is available'))
    const identity = { bookId: draft.book.id, chapterId: draft.chapter.id, draftId: draft.id }
    const api = yield* BookManagerApi
    const chapterName = get(chapterNameAtom)
    if (get(chapterNameWasChangedAtom)) {
      yield* api.updateDraftChapterName(identity, chapterName)
    }
    const editor = get(wysiwygEditorAtom)
    const updatedDraft = yield* api.updateDraft(identity, editor?.getHTML() ?? '')
    yield* Effect.sync(() => {
      set(wysiwygContentModifiedAtom, false)
      set(draftAtom, updatedDraft)
    })
    return updatedDraft
  }).pipe(
    Effect.tapError((error) => Effect.sync(() => set(saveErrorAtom, error))),
    Effect.ensuring(Effect.sync(() => set(savingAtom, false))),
  )
}

export const saveDraftAtom = atom(null, (get, set) => saveDraft(get, set))

export const publishDraftAtom = atom(null, (get, set, makePublic: boolean) =>
  Effect.gen(function* () {
    yield* Effect.sync(() => {
      set(savingAtom, true)
      set(saveErrorAtom, null)
    })
    const draft = get(draftAtom)
    if (!draft)
      return yield* Effect.fail(new Error('cannot save draft - no draft information is available'))
    const identity = { bookId: draft.book.id, chapterId: draft.chapter.id, draftId: draft.id }
    const api = yield* BookManagerApi
    const chapterName = get(chapterNameAtom)
    if (get(chapterNameWasChangedAtom)) {
      yield* api.updateDraftChapterName(identity, chapterName)
    }
    const editor = get(wysiwygEditorAtom)
    const updatedDraft = yield* api.publishDraft(identity, editor?.getHTML() ?? '', makePublic)
    yield* Effect.sync(() => {
      set(wysiwygContentModifiedAtom, false)
      set(draftAtom, updatedDraft)
    })
    return updatedDraft
  }).pipe(
    Effect.tapError((error) => Effect.sync(() => set(saveErrorAtom, error))),
    Effect.ensuring(Effect.sync(() => set(savingAtom, false))),
  ),
)

export const scheduleDraftAtom = atom(null, (get, set, scheduledAt: Date) =>
  Effect.gen(function* () {
    yield* set(saveDraftAtom)
    yield* Effect.sync(() => set(savingAtom, true))
    const draft = get(draftAtom)
    if (!draft)
      return yield* Effect.fail(
        new Error('cannot schedule draft - no draft information is available'),
      )
    const api = yield* BookManagerApi
    const updatedDraft = yield* api.scheduleDraft(
      { bookId: draft.book.id, chapterId: draft.chapter.id, draftId: draft.id },
      scheduledAt,
    )
    yield* Effect.sync(() => set(draftAtom, updatedDraft))
    return updatedDraft
  }).pipe(
    Effect.tapError((error) => Effect.sync(() => set(saveErrorAtom, error))),
    Effect.ensuring(Effect.sync(() => set(savingAtom, false))),
  ),
)

export function useDraftHasChanges() {
  const contentWasChanged = useAtomValue(wysiwygContentModifiedAtom)
  const chapterNameChanged = useAtomValue(chapterNameWasChangedAtom)
  return chapterNameChanged || contentWasChanged
}
