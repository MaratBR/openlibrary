import { httpClient, OLAPIResponse } from '@/http-client'
import type { CommentDto } from '@/backend-types'

export function httpLikeComment(commentId: string, like: boolean) {
  return httpClient.post('/_api/comments/like', {
    searchParams: {
      like,
      commentId,
      ts: Date.now(),
    },
  })
}

export function httpAddComment(chapterId: string, content: string) {
  return httpClient.post('/_api/comments/add', { json: { chapterId, content } })
}

export function httpGetCommentReplies(commentId: string, cursor: number) {
  return httpClient
    .get('/_api/comments/replies', {
      searchParams: {
        commentId,
        cursor,
      },
    })
    .then((r) =>
      OLAPIResponse.create<{ cursor: number; nextCursor: number; comments: CommentDto[] }>(r),
    )
}

export type { CommentDto } from '@/backend-types'
