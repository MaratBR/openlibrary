import { HttpClient } from '@/features/http-client'
import { Cache, Context, Effect, Layer, Schema } from 'effect'

export const TagsCategory = Schema.Literals([
  'other',
  'warning',
  'fandom',
  'rel',
  'reltype',
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
          httpClient.get('/_api/tags', { searchParams: { q: query } }).json(),
        )
        return yield* Schema.decodeUnknownEffect(Schema.Array(DefinedTagDto))(response)
      })
      const tagsCache = yield* Cache.make({
        capacity: 100,
        lookup: fetchTags,
        timeToLive: '5 minutes',
      })
      const searchTags = Effect.fn('SearchApi.searchTags')((query: string) =>
        Cache.get(tagsCache, query),
      )

      return SearchApi.of({ searchTags })
    }),
  )
}
