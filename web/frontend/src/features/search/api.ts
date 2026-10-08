import { HttpClient } from '@/features/http-client'
import { Context, Effect, Layer, Schema } from 'effect'

export const TagsCategory = Schema.Literals([
  'other',
  'warning',
  'fandom',
  'rel',
  'reltype',
  'genre',
  'unknown',
])

export type TagsCategory = Schema.Schema.Type<typeof TagsCategory>

export const DefinedTagDto = Schema.Struct({
  id: Schema.String,
  name: Schema.String,
  desc: Schema.String,
  adult: Schema.Boolean,
  spoiler: Schema.Boolean,
  cat: TagsCategory,
})

export type DefinedTagDto = Schema.Schema.Type<typeof DefinedTagDto>

export class SearchApi extends Context.Service<
  SearchApi,
  {
    readonly searchTags: (query: string) => Effect.Effect<ReadonlyArray<DefinedTagDto>, unknown>
  }
>()('openlibrary/SearchApi') {
  static readonly layer = Layer.effect(
    this,
    Effect.gen(function* () {
      const httpClient = yield* HttpClient
      const fetchTags = Effect.fn('SearchApi.fetchTags')(function* (query: string) {
        const response = yield* Effect.tryPromise(() =>
          httpClient
            .get(
              window.location.pathname.startsWith('/books-manager')
                ? '/_api/books-manager/tags'
                : '/_api/tags',
              {
                searchParams: {
                  q: query,
                  ...(new URLSearchParams(window.location.search).get('admin.link') === '1' ||
                  window.location.pathname.startsWith('/admin')
                    ? { 'admin.link': '1' }
                    : {}),
                  ...(new URLSearchParams(window.location.search).get('admin.override') === '1'
                    ? { 'admin.override': '1' }
                    : {}),
                },
              },
            )
            .json(),
        )
        return yield* Schema.decodeUnknownEffect(Schema.Array(DefinedTagDto))(response)
      })
      // Results depend on current account preferences and management mode.
      const searchTags = fetchTags

      return SearchApi.of({ searchTags })
    }),
  )
}
