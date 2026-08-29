import { HttpClient, OLAPIResponse } from '@/features/http-client'
import { Context, Effect, Layer, Schema } from 'effect'

export type GetUserDataOptions<T> = {
  readonly schema: Schema.Codec<T, unknown>
}

export class UserDataApi extends Context.Service<
  UserDataApi,
  {
    readonly getUserData: <T>(
      key: string,
      options: GetUserDataOptions<T>,
    ) => Effect.Effect<T | null, unknown>
    readonly setUserData: (key: string, data: unknown) => Effect.Effect<void, unknown>
  }
>()('openlibrary/UserDataApi') {
  static readonly layer = Layer.effect(
    this,
    Effect.gen(function* () {
      const httpClient = yield* HttpClient
      const getUserData = <T>(key: string, { schema }: GetUserDataOptions<T>) =>
        Effect.tryPromise(async () => {
          const response = await httpClient
            .get(`/_api/user-data/${encodeURIComponent(key)}`)
            .then((r) => OLAPIResponse.create(r, Schema.NullOr(schema)))
          response.throwIfError()
          return response.data
        })
      const setUserData = (key: string, data: unknown) =>
        Effect.tryPromise(async () => {
          const response = await httpClient
            .put(`/_api/user-data/${encodeURIComponent(key)}`, { json: data })
            .then((r) => OLAPIResponse.createNoBody(r))
          response.throwIfError()
        })
      return UserDataApi.of({ getUserData, setUserData })
    }),
  )
}
