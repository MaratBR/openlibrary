import { CommentDto, httpGetCommentReplies, httpLikeComment } from '@/api/comments'
import UserContent from '@/components/UserContent'

import { useEffect, useRef, useState } from 'react'

export type CommentRepliesController = {
  close(): void
}

export function initCommentReplies(
  _$root: HTMLElement,
  _commentId: string,
): CommentRepliesController {
  alert('Comments are not finished')
  throw new Error('not implemented')
  // render(<Replies commentId={commentId} />, $root)

  // return {
  //   close() {
  //     render(null, $root)
  //   },
  // }
}

export function Replies({ commentId }: { commentId: string }) {
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
    <div className="mt-4 border-b pb-2">
      {replies.map((reply) => {
        return <Reply key={reply.id} reply={reply} />
      })}
    </div>
  )
}

function Reply({ reply }: { reply: CommentDto }) {
  const [liked, setLiked] = useState(!!reply.likedAt)
  const [likes, setLikes] = useState(
    reply.likes +
      (reply.likedAt && new Date(reply.likedAt) > new Date(reply.likesUpdatedAt) ? 1 : 0),
  )
  const opIdRef = useRef(0)

  function handleLike() {
    const opId = opIdRef.current++
    setLiked(!liked)
    setLikes(liked ? likes - 1 : likes + 1)
    httpLikeComment(reply.id, !liked).catch((err) => {
      if (opId === opIdRef.current) {
        setLiked(liked)
      }
      window.toast.error(err)
    })
  }

  return (
    <div className="chapter-comment chapter-comment--reply">
      <header className="chapter-comment__header">
        <img className="avatar border size-12" src={reply.user.avatar} />
        <strong>{reply.user.name}</strong>
      </header>

      <div className="chapter-comment__content">
        <UserContent value={reply.content} />
      </div>

      <div className="chapter-comment__actions">
        <button onClick={handleLike} data-set={`${liked}`}>
          <i className="fa-solid fa-thumbs-up" />
          {likes > 0 && `${likes}`}
        </button>
      </div>
    </div>
  )
}
