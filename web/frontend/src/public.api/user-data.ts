import { httpClient, OLAPIResponse } from '@/features/http-client'
import { Effect } from 'effect'
import { type ZodType } from 'zod'

export namespace UserDataApi {
  export type GetUserDataOptions<T> = {
    schema: ZodType<T>
  }

  export function getUserData<T>(key: string, { schema }: GetUserDataOptions<T>) {
    return Effect.tryPromise(async () => {
      const response = await httpClient
        .get(`/_api/user-data/${encodeURIComponent(key)}`)
        .then((r) => OLAPIResponse.create(r, schema.nullable()))
      response.throwIfError()
      return response.data
    })
  }

  export function setUserData(key: string, data: unknown) {
    return Effect.tryPromise(async () => {
      const response = await httpClient
        .put(`/_api/user-data/${encodeURIComponent(key)}`, { json: data })
        .then((r) => OLAPIResponse.createNoBody(r))
      response.throwIfError()
    })
  }
}
