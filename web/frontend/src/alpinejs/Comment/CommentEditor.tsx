import { useMemo, useState } from 'react'
import { createRoot } from 'react-dom/client'

export type CommentEditorController = {
  close: () => void
}

export function initCommentEditor($root: HTMLElement, commentId: string): CommentEditorController {
  const d = document.createElement('div')
  d.classList = 'ol-comment-reply-editor'

  const root = createRoot(d)
  root.render(<Editor commentId={commentId} />)
  $root.prepend(d)

  return {
    close() {
      root.unmount()
      d.remove()
    },
  }
}

function Editor({ commentId }: { commentId: string }) {
  const [text, setText] = useState('')
	const [submitting, setSubmitting] = useState(false)

  const valid = useMemo(() => {
    const length = text.trim().length
    return length > 0 && length < 1000
  }, [text])

  return (
    <form onSubmit={async (event) => {
		event.preventDefault()
		setSubmitting(true)
		try {
			const response = await fetch('/_api/comments/add', {
				method: 'POST', headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ chapterId: `${window.__server__.chapterId}`, parentCommentId: commentId, content: text }),
			})
			if (!response.ok) throw new Error('Could not post reply')
			window.location.reload()
		} catch (error) {
			window.toast.error(error)
			setSubmitting(false)
		}
	}}>
      <textarea
        placeholder={window._('common.replyPlaceholder')}
        name="text"
        className="ol-comment-reply-editor__input"
        value={text}
        onChange={(e) => setText((e.target as HTMLTextAreaElement).value)}
      />
      <div className="ol-comment-reply-editor__actions">
		<button disabled={!valid || submitting} className="btn btn--default btn--sm">{window._('common.reply')}</button>
      </div>
    </form>
  )
}
