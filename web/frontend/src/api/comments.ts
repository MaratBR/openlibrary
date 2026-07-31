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
  return httpClient
    .post('/_api/comments/add', { json: { chapterId, content } })
    .then((r) => OLAPIResponse.create<CommentDto>(r))
}

export function httpReplyToComment(chapterId: string, parentCommentId: string, content: string) {
  return httpClient
    .post('/_api/comments/add', { json: { chapterId, parentCommentId, content } })
    .then((r) => OLAPIResponse.create<CommentDto>(r))
}

export function httpUpdateComment(commentId: string, content: string) {
  return httpClient
    .put(`/_api/comments/${encodeURIComponent(commentId)}`, { json: { content } })
    .then((r) => OLAPIResponse.create<CommentDto>(r))
}

export function httpGetComments(chapterId: string, sort: string, cursor: number) {
  return httpClient.get('/_api/comments', { searchParams: { chapterId, sort, cursor } }).then((r) =>
    OLAPIResponse.create<{
      cursor: number
      nextCursor: number
      comments: CommentDto[]
      total: number
    }>(r),
  )
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
