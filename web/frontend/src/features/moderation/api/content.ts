import { httpClient, OLAPIResponse } from '@/features/http-client'
import type { ModerationReasonRequest } from '@/backend-types'

export type ChapterModerationAction = 'hide' | 'restore'
export type CommentModerationAction = 'remove' | 'restore'

export function performChapterModerationAction(
  chapterId: string,
  action: ChapterModerationAction,
  reason: string,
) {
  const body: ModerationReasonRequest = { reason }
  return httpClient
    .post(`/_api/moderation/chapters/${encodeURIComponent(chapterId)}/actions/${action}`, {
      json: body,
    })
    .then((response) => OLAPIResponse.createNoBody(response))
}

export function performCommentModerationAction(
  commentId: string,
  action: CommentModerationAction,
  reason: string,
) {
  const body: ModerationReasonRequest = { reason }
  return httpClient
    .post(`/_api/moderation/comments/${encodeURIComponent(commentId)}/actions/${action}`, {
      json: body,
    })
    .then((response) => OLAPIResponse.createNoBody(response))
}
