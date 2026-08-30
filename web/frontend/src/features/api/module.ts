import { FontsApi } from '@/features/fonts-loader/api'
import { HttpClient } from '@/features/http-client'
import { UserDataApi } from './user-data'
import { SearchApi } from '@/features/search/api'
import { Layer } from 'effect'
import { BookManagerApi } from '@/features/book-manager/api'

/** A reusable package of related services and their dependencies. */
export const ApiModuleLive = Layer.mergeAll(
  FontsApi.layer,
  SearchApi.layer,
  UserDataApi.layer,
  BookManagerApi.layer,
).pipe(Layer.provideMerge(HttpClient.layer))
