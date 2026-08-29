import { FontsApi } from '@/features/fonts-loader/api'
import { HttpClient } from '@/features/http-client'
import { UserDataApi } from './user-data'
import { Layer } from 'effect'

/** A reusable package of related services and their dependencies. */
export const ApiModuleLive = Layer.mergeAll(FontsApi.layer, UserDataApi.layer).pipe(
  Layer.provideMerge(HttpClient.layer),
)
