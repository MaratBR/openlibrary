import { httpClient } from '@/features/http-client'
import { Schema } from 'effect'
import type { recentCollectionDto as RecentCollectionDto } from '@/backend-types'

export async function httpGetRecentCollections(): Promise<RecentCollectionDto[]> {
  return httpClient.get('/_api/collections/recent').then((r) => r.json<RecentCollectionDto[]>())
}

export async function httpGetCollectionsContainingBook(
  bookId: string,
): Promise<RecentCollectionDto[]> {
  return httpClient
    .get('/_api/collections/containingBook', { searchParams: { bookId } })
    .then((r) => r.json<RecentCollectionDto[]>())
}

export type { recentCollectionDto as RecentCollectionDto } from '@/backend-types'

export async function httpCreateCollection(name: string): Promise<string> {
  return httpClient
    .post('/_api/collections', { json: { name } })
    .then((r) => r.json())
    .then(Schema.decodeUnknownSync(Schema.String))
}

export async function httpAddBookToCollections(
  bookId: string,
  collectionIds: string[],
): Promise<void> {
  await httpClient.post('/_api/collections/addBook', {
    searchParams: { bookId },
    json: collectionIds,
  })
}
