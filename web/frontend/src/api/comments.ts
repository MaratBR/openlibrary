import { httpClient, OLAPIResponse } from '@/http-client'
import z from 'zod'

export function httpLikeComment(commentId: string, like: boolean) {
  return httpClient.post('/_api/comments/like', {
    searchParams: {
      like,
      commentId,
      ts: Date.now(),
    },
  })
}

export const commentUserDto = z.object({
  id: z.string(),
  name: z.string(),
  avatar: z.string(),
})

export const commentDtoSchema = z.object({
  id: z.string(),
  content: z.string(),
  user: commentUserDto,
  createdAt: z.string(),
  updatedAt: z.string().nullable(),
  likedAt: z.string().nullable(),
  likes: z.number(),
  likesUpdatedAt: z.string(),
  subcomments: z.number(),
})

export type CommentDto = z.infer<typeof commentDtoSchema>

export function httpGetCommentReplies(commentId: string, cursor: number) {
  return httpClient
    .get('/_api/comments/replies', {
      searchParams: {
        commentId,
        cursor,
      },
    })
    .then((r) =>
      OLAPIResponse.create(
        r,
        z.object({
          cursor: z.number(),
          nextCursor: z.number(),
          comments: commentDtoSchema.array(),
        }),
      ),
    )
}
