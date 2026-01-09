import { httpClient } from '@/http-client'

export function httpLikeComment(commentId: string, like: boolean) {
  return httpClient.post('/_api/comments/like', {
    searchParams: {
      like,
      commentId,
      ts: Date.now(),
    },
  })
}
