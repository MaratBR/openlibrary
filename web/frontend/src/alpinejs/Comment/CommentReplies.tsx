import { CommentDto, httpGetCommentReplies } from '@/api/comments'
import UserContent from '@/components/UserContent'
import { render } from 'preact'
import { useEffect, useRef, useState } from 'preact/hooks'

export type CommentRepliesController = {
  close(): void
}

export function initCommentReplies(
  $root: HTMLElement,
  commentId: string,
): CommentRepliesController {
  render(<Replies commentId={commentId} />, $root)

  return {
    close() {
      render(null, $root)
    },
  }
}

function Replies({ commentId }: { commentId: string }) {
  const [replies, setReplies] = useState<CommentDto[]>([])
  const cursorRef = useRef(0)

  function load() {
    httpGetCommentReplies(commentId, cursorRef.current).then((resp) => {
      setReplies((r) => [...r, ...resp.data.comments])
      cursorRef.current = resp.data.nextCursor
    })
  }

  useEffect(load, [commentId])

  return (
    <>
      {replies.map((reply) => {
        return (
          <div key={reply.id} class="chapter-comment chapter-comment--reply">
            <header class="chapter-comment__header">
              <img class="avatar border size-12" src={reply.user.avatar} />
              <strong>{reply.user.name}</strong>
            </header>

            <div class="chapter-comment__content">
              <UserContent value={reply.content} />
            </div>

            <div class="chapter-comment__actions">
              <button>
                <i class="fa-solid fa-thumbs-up" />
              </button>
            </div>
          </div>
        )
      })}
    </>
  )
}
