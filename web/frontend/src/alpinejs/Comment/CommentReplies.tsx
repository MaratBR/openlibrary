import { CommentDto, httpGetCommentReplies, httpLikeComment } from '@/api/comments'
import UserContent from '@/components/UserContent'

import { useEffect, useRef, useState } from 'react'
import { createRoot } from 'react-dom/client'

export type CommentRepliesController = {
  close(): void
}

export function initCommentReplies(
  $root: HTMLElement,
  commentId: string,
): CommentRepliesController {
  const root = createRoot($root)
  root.render(<Replies commentId={commentId} />)
  return { close() { root.unmount() } }
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
    <div>
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
    <article className="ol-comment">
      <img className="ol-comment__avatar" src={reply.user.avatar} alt="" />
      <div className="ol-comment__body">
        <header className="ol-comment__header"><strong className="ol-comment__author">{reply.user.name}</strong></header>
        <div className="ol-comment__content"><UserContent value={reply.content} /></div>
        <div className="ol-comment__actions">
          <button className="ol-comment-action" onClick={handleLike} aria-pressed={liked}>
            <i className="fa-solid fa-arrow-up" /> {likes > 0 && `${likes}`}
          </button>
          <button className="ol-comment-action" onClick={() => alert('Not implemented yet')}><i className="fa-solid fa-ellipsis" /> {window._('common.more')}</button>
        </div>
      </div>
    </article>
  )
}
