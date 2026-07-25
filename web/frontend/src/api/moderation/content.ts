import { httpClient, OLAPIResponse } from '@/http-client'

export type ChapterModerationAction = 'hide' | 'restore'
export type CommentModerationAction = 'remove' | 'restore'

export function performChapterModerationAction(
  chapterId: string,
  action: ChapterModerationAction,
  reason: string,
) {
  return httpClient
    .post(`/_api/moderation/chapters/${encodeURIComponent(chapterId)}/actions/${action}`, {
      json: { reason },
    })
    .then((response) => OLAPIResponse.createNoBody(response))
}

export function performCommentModerationAction(
  commentId: string,
  action: CommentModerationAction,
  reason: string,
) {
  return httpClient
    .post(`/_api/moderation/comments/${encodeURIComponent(commentId)}/actions/${action}`, {
      json: { reason },
    })
    .then((response) => OLAPIResponse.createNoBody(response))
}
