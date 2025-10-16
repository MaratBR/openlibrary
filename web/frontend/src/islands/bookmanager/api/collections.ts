import { httpClient } from '@/http-client'
import { z } from 'zod'

const recentCollectionDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
})

export type RecentCollectionDto = z.infer<typeof recentCollectionDtoSchema>

export async function httpGetRecentCollections(): Promise<RecentCollectionDto[]> {
  return httpClient
    .get('/_api/collections/recent')
    .then((r) => r.json())
    .then(z.array(recentCollectionDtoSchema).parse)
}

export async function httpGetCollectionsContainingBook(
  bookId: string,
): Promise<RecentCollectionDto[]> {
  return httpClient
    .get('/_api/collections/containingBook', { searchParams: { bookId } })
    .then((r) => r.json())
    .then(z.array(recentCollectionDtoSchema).parse)
}

export async function httpCreateCollection(name: string): Promise<string> {
  return httpClient
    .post('/_api/collections', { json: { name } })
    .then((r) => r.json())
    .then(z.string().parse)
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
