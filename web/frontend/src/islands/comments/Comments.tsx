import type { CommentDto } from '@/api/comments'
import {
  httpAddComment,
  httpGetCommentReplies,
  httpGetComments,
  httpLikeComment,
  httpReplyToComment,
  httpUpdateComment,
} from '@/api/comments'
import UserContent from '@/components/UserContent'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components'
import { RichTextInput, useOLEditor } from '@/components/rte'
import type { ReactIslandProps } from '@/islands/common/react-island'
import { SelfUserDtoSchema } from '@/api/auth/user'
import { FormEvent, useEffect, useRef, useState } from 'react'

type CommentSort = 'newest' | 'oldest' | 'popular'

export function Comments(_: ReactIslandProps) {
  const chapterId = getChapterId()

  const userResult = SelfUserDtoSchema.safeParse(window.__server__.user)
  const user = userResult.success ? userResult.data : null
  const [comments, setComments] = useState<CommentDto[]>([])
  const [nextCursor, setNextCursor] = useState(0)
  const [total, setTotal] = useState(0)
  const [sort, setSort] = useState<CommentSort>(() => parseSort(location.search))
  const [loading, setLoading] = useState(true)
  const initialized = useRef(false)

  useEffect(() => {
    if (initialized.current) return
    initialized.current = true
    httpGetComments(chapterId, sort, 0)
      .then((response) => {
        setComments(response.data.comments)
        setNextCursor(response.data.nextCursor)
        setTotal(response.data.total)
      })
      .catch((error) => window.toast.error(error))
      .finally(() => setLoading(false))
  }, [chapterId, sort])

  async function loadMore() {
    setLoading(true)
    try {
      const response = await httpGetComments(chapterId, sort, nextCursor)
      setComments((current) => [...current, ...response.data.comments])
      setNextCursor(response.data.nextCursor)
    } catch (error) {
      window.toast.error(error)
    } finally {
      setLoading(false)
    }
  }

  async function changeSort(value: CommentSort) {
    setSort(value)
    setLoading(true)
    try {
      const response = await httpGetComments(chapterId, value, 0)
      setComments(response.data.comments)
      setNextCursor(response.data.nextCursor)
      setTotal(response.data.total)
      const url = new URL(window.location.href)
      url.searchParams.set('sort', value)
      url.searchParams.delete('cursor')
      window.history.replaceState(null, '', url)
    } catch (error) {
      window.toast.error(error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      {user ? (
        <Composer chapterId={chapterId} avatar={user.avatar.md} />
      ) : (
        <div className="ol-comments__signin card">
          <a className="btn btn--default" href="/login">
            {window._('comments.signInToComment')}
          </a>
        </div>
      )}
      <div className="ol-comments__toolbar">
        <div className="ol-comments__toolbar-label">
          <span className="ol-comments__toolbar-title">{window._('comments.discussion')}</span>
          {(!loading || total > 0) && (
            <span className="ol-comments__count">
              {window._('common.commentsCount', { Count: `${total}` })}
            </span>
          )}
        </div>
        <label className="sr-only" htmlFor="ChapterCommentsSort">
          {window._('comments.sort')}
        </label>
        <Select
          value={sort}
          disabled={loading}
          onValueChange={(value) => changeSort(value as CommentSort)}
        >
          <SelectTrigger id="ChapterCommentsSort" className="select--sm ol-comments__sort w-64">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="newest">{window._('comments.newest')}</SelectItem>
            <SelectItem value="oldest">{window._('comments.oldest')}</SelectItem>
            <SelectItem value="popular">{window._('comments.popular')}</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="ol-comment-list">
        {!loading && comments.length === 0 && (
          <p className="ol-comments__empty">{window._('comments.empty')}</p>
        )}
        {comments.map((comment) => (
          <Comment
            key={comment.id}
            comment={comment}
            chapterId={chapterId}
            authenticated={!!user}
            currentUserId={user?.id}
          />
        ))}
        {nextCursor > 0 && (
          <button className="btn btn--ghost btn--lg mt-4" disabled={loading} onClick={loadMore}>
            {window._('common.loadMore')}
          </button>
        )}
      </div>
    </>
  )
}

function Composer({ chapterId, avatar }: { chapterId: string; avatar: string }) {
  const [content, setContent] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const editor = useOLEditor({ onUpdate: ({ editor }) => setContent(editor.getHTML()) })

  async function submit(event: FormEvent) {
    event.preventDefault()
    setSubmitting(true)
    try {
      await httpAddComment(chapterId, content)
      window.location.reload()
    } catch (error) {
      window.toast.error(error)
      setSubmitting(false)
    }
  }

  return (
    <form className="ol-comment-composer" onSubmit={submit}>
      <img className="ol-comment__avatar" src={avatar} alt="" />
      <div className="ol-comment-composer__body">
        {editor && <RichTextInput editor={editor} />}
        <div className="ol-comment-composer__footer">
          <p className="ol-comment-composer__hint">{window._('comments.kindHint')}</p>
          <div className="ol-comment-composer__actions">
            <button
              className="btn btn--default"
              disabled={submitting || !editor || editor.isEmpty || content.length > 2000}
            >
              {window._('comments.post')}
            </button>
          </div>
        </div>
      </div>
    </form>
  )
}

function Comment({
  comment,
  chapterId,
  authenticated,
  currentUserId,
  allowReply = true,
}: {
  comment: CommentDto
  chapterId: string
  authenticated: boolean
  currentUserId?: string
  allowReply?: boolean
}) {
  const [liked, setLiked] = useState(!!comment.likedAt)
  const [likes, setLikes] = useState(realLikes(comment))
  const [replying, setReplying] = useState(false)
  const [editing, setEditing] = useState(false)
  const [content, setContent] = useState(comment.content)
  const [updatedAt, setUpdatedAt] = useState(comment.updatedAt)
  const operation = useRef(0)

  async function toggleLike() {
    const previous = liked
    const id = ++operation.current
    setLiked(!previous)
    setLikes((value) => value + (previous ? -1 : 1))
    try {
      await httpLikeComment(comment.id, !previous)
    } catch (error) {
      if (id === operation.current) {
        setLiked(previous)
        setLikes((value) => value + (previous ? 1 : -1))
      }
      window.toast.error(error)
    }
  }

  return (
    <article className={`ol-comment${comment.deleted ? ' ol-comment--deleted' : ''}`}>
      <img className="ol-comment__avatar" src={comment.user.avatar} alt="" />
      <div className="ol-comment__body">
        <header className="ol-comment__header">
          <a href={`/users/${comment.user.id}`} className="ol-comment__author">
            {comment.user.name}
          </a>
          <Time value={comment.createdAt} />
          {isEdited(comment.createdAt, updatedAt) && (
            <span className="ol-comment__meta">{window._('comments.edited')}</span>
          )}
        </header>
        {editing ? (
          <CommentEditor
            commentId={comment.id}
            content={content}
            onCancel={() => setEditing(false)}
            onUpdated={(updated) => {
              setContent(updated.content)
              setUpdatedAt(updated.updatedAt)
              setEditing(false)
            }}
          />
        ) : (
          <div className="ol-comment__content">
            <UserContent value={content} />
          </div>
        )}
        <div className="ol-comment__actions">
          <button
            className="ol-comment-action"
            disabled={!authenticated}
            onClick={toggleLike}
            aria-pressed={liked}
          >
            <i className="fa-solid fa-arrow-up" aria-hidden="true" /> {likes > 0 && likes}
          </button>
          {authenticated && allowReply && (
            <button className="ol-comment-action" onClick={() => setReplying((value) => !value)}>
              <i className="fa-solid fa-reply" aria-hidden="true" /> {window._('common.reply')}
            </button>
          )}
          {currentUserId === comment.user.id && !comment.deleted && !editing && (
            <button className="ol-comment-action" onClick={() => setEditing(true)}>
              <i className="fa-solid fa-pencil" aria-hidden="true" /> {window._('common.edit')}
            </button>
          )}
        </div>
        {replying && <ReplyComposer chapterId={chapterId} parentId={comment.id} />}
        {comment.subcomments > 0 && (
          <Replies
            commentId={comment.id}
            authenticated={authenticated}
            currentUserId={currentUserId}
          />
        )}
      </div>
    </article>
  )
}

function Replies({
  commentId,
  authenticated,
  currentUserId,
}: {
  commentId: string
  authenticated: boolean
  currentUserId?: string
}) {
  const [replies, setReplies] = useState<CommentDto[]>([])
  const [cursor, setCursor] = useState(0)
  const initialized = useRef(false)

  useEffect(() => {
    if (initialized.current) return
    initialized.current = true
    httpGetCommentReplies(commentId, 0)
      .then((response) => {
        setReplies(response.data.comments)
        setCursor(response.data.nextCursor)
      })
      .catch((error) => window.toast.error(error))
  }, [commentId])

  async function load() {
    try {
      const response = await httpGetCommentReplies(commentId, cursor)
      setReplies((value) => [...value, ...response.data.comments])
      setCursor(response.data.nextCursor)
    } catch (error) {
      window.toast.error(error)
    }
  }

  return (
    <div className="ol-comment-replies">
      {replies.map((reply) => (
        <Comment
          key={reply.id}
          comment={reply}
          chapterId=""
          authenticated={authenticated}
          currentUserId={currentUserId}
          allowReply={false}
        />
      ))}
      {cursor > 0 && (
        <button className="btn btn--ghost mt-2" onClick={load}>
          {window._('common.loadMore')}
        </button>
      )}
    </div>
  )
}

function ReplyComposer({ chapterId, parentId }: { chapterId: string; parentId: string }) {
  const [content, setContent] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const editor = useOLEditor({ onUpdate: ({ editor }) => setContent(editor.getHTML()) })
  async function submit(event: FormEvent) {
    event.preventDefault()
    setSubmitting(true)
    try {
      await httpReplyToComment(chapterId, parentId, content)
      window.location.reload()
    } catch (error) {
      window.toast.error(error)
      setSubmitting(false)
    }
  }
  return (
    <form className="ol-comment-reply-editor" onSubmit={submit}>
      {editor && <RichTextInput editor={editor} />}
      <div className="ol-comment-reply-editor__actions">
        <button
          className="btn btn--default"
          disabled={submitting || !editor || editor.isEmpty || content.length > 2000}
        >
          {window._('comments.post')}
        </button>
      </div>
    </form>
  )
}

function CommentEditor({
  commentId,
  content: initialContent,
  onCancel,
  onUpdated,
}: {
  commentId: string
  content: string
  onCancel: () => void
  onUpdated: (comment: CommentDto) => void
}) {
  const [content, setContent] = useState(initialContent)
  const [saving, setSaving] = useState(false)
  const editor = useOLEditor({
    content: initialContent,
    onUpdate: ({ editor }) => setContent(editor.getHTML()),
  })

  async function submit(event: FormEvent) {
    event.preventDefault()
    if (!editor || saving) return
    setSaving(true)
    try {
      const response = await httpUpdateComment(commentId, content)
      onUpdated(response.data)
    } catch (error) {
      window.toast.error(error)
      setSaving(false)
    }
  }

  return (
    <form className="ol-comment-edit" onSubmit={submit}>
      {editor && <RichTextInput editor={editor} />}
      <div className="ol-comment-edit__actions">
        <button type="button" className="btn btn--ghost" disabled={saving} onClick={onCancel}>
          {window._('common.cancel')}
        </button>
        <button
          className="btn btn--default"
          disabled={saving || !editor || editor.isEmpty || content.length > 2000}
        >
          {window._('common.save')}
        </button>
      </div>
    </form>
  )
}

function Time({ value }: { value: string }) {
  return (
    <time className="ol-comment__meta" dateTime={value}>
      {relativeTime(value)}
    </time>
  )
}

function realLikes(comment: CommentDto) {
  return (
    comment.likes +
    (comment.likedAt && new Date(comment.likedAt) > new Date(comment.likesUpdatedAt) ? 1 : 0)
  )
}

function isEdited(createdAt: string, updatedAt: string | null) {
  return !!updatedAt && new Date(updatedAt).getTime() > new Date(createdAt).getTime()
}

function relativeTime(value: string) {
  const seconds = Math.round((new Date(value).getTime() - Date.now()) / 1000)
  const ranges: [Intl.RelativeTimeFormatUnit, number][] = [
    ['year', 31_536_000],
    ['month', 2_592_000],
    ['week', 604_800],
    ['day', 86_400],
    ['hour', 3_600],
    ['minute', 60],
  ]
  const formatter = new Intl.RelativeTimeFormat(document.documentElement.lang || undefined, {
    numeric: 'auto',
  })
  for (const [unit, size] of ranges)
    if (Math.abs(seconds) >= size) return formatter.format(Math.round(seconds / size), unit)
  return formatter.format(seconds, 'second')
}

function parseSort(search: string): CommentSort {
  const value = new URLSearchParams(search).get('sort')
  return value === 'oldest' || value === 'popular' ? value : 'newest'
}

function getChapterId(): string {
  const chapterId = window.__server__?.chapterId
  if (!chapterId) throw new Error('could not find chapterId, __server__.chapterId is not set')
  return chapterId
}
