import { httpClient } from '@/modules/common/api'

export async function httpUploadBookCover(bookId: string, file: File): Promise<{ url: string }> {
  const formData = new FormData()
  formData.append('file', file)
  return httpClient
    .post(`/api/manager/books/${bookId}/cover`, { body: formData })
    .then((r) => r.json())
}
