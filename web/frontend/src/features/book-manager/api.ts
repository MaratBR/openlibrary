import type { DraftDto } from '@/backend-types'
import { HttpClient, OLAPIResponse } from '@/features/http-client'
import { Context, Effect, Layer } from 'effect'

type DraftIdentity = {
  readonly bookId: string
  readonly chapterId: string
  readonly draftId: string
}

export class BookManagerApi extends Context.Service<
  BookManagerApi,
  {
    readonly updateDraftChapterName: (
      identity: DraftIdentity,
      chapterName: string,
    ) => Effect.Effect<void, unknown>
    readonly updateDraft: (
      identity: DraftIdentity,
      content: string,
    ) => Effect.Effect<DraftDto, unknown>
    readonly publishDraft: (
      identity: DraftIdentity,
      content: string,
      makePublic: boolean,
    ) => Effect.Effect<DraftDto, unknown>
    readonly scheduleDraft: (
      identity: DraftIdentity,
      scheduledAt: Date,
    ) => Effect.Effect<DraftDto, unknown>
  }
>()('openlibrary/BookManagerApi') {
  static readonly layer = Layer.effect(
    this,
    Effect.gen(function* () {
      const httpClient = yield* HttpClient
      const draftPath = ({ bookId, chapterId, draftId }: DraftIdentity) =>
        `/_api/books-manager/book/${bookId}/${chapterId}/${draftId}`

      const updateDraftChapterName = Effect.fn('BookManagerApi.updateDraftChapterName')(function* (
        identity: DraftIdentity,
        chapterName: string,
      ) {
        const response = yield* Effect.tryPromise(() =>
          httpClient.post(`${draftPath(identity)}/chapterName`, {
            body: chapterName,
            headers: { 'Content-Type': 'text/plain' },
          }),
        )
        yield* Effect.tryPromise(() =>
          OLAPIResponse.createNoBody(response).then((r) => r.throwIfError()),
        )
      })

      const updateDraft = Effect.fn('BookManagerApi.updateDraft')(function* (
        identity: DraftIdentity,
        content: string,
      ) {
        const response = yield* Effect.tryPromise(() =>
          httpClient.post(draftPath(identity), { json: { content } }),
        )
        const result = yield* Effect.tryPromise(() => OLAPIResponse.create<DraftDto>(response))
        return yield* Effect.try({
          try: () => {
            result.throwIfError()
            return result.data
          },
          catch: (error) => error,
        })
      })

      const publishDraft = Effect.fn('BookManagerApi.publishDraft')(function* (
        identity: DraftIdentity,
        content: string,
        makePublic: boolean,
      ) {
        const response = yield* Effect.tryPromise(() =>
          httpClient.post(`${draftPath(identity)}/publish`, {
            json: { content },
            searchParams: { makePublic },
          }),
        )
        const result = yield* Effect.tryPromise(() => OLAPIResponse.create<DraftDto>(response))
        return yield* Effect.try({
          try: () => {
            result.throwIfError()
            return result.data
          },
          catch: (error) => error,
        })
      })

      const scheduleDraft = Effect.fn('BookManagerApi.scheduleDraft')(function* (
        identity: DraftIdentity,
        scheduledAt: Date,
      ) {
        const response = yield* Effect.tryPromise(() =>
          httpClient.post(`${draftPath(identity)}/schedule`, {
            json: { scheduledAt: scheduledAt.toISOString() },
          }),
        )
        const result = yield* Effect.tryPromise(() => OLAPIResponse.create<DraftDto>(response))
        return yield* Effect.try({
          try: () => {
            result.throwIfError()
            return result.data
          },
          catch: (error) => error,
        })
      })

      return BookManagerApi.of({
        updateDraftChapterName,
        updateDraft,
        publishDraft,
        scheduleDraft,
      })
    }),
  )
}
