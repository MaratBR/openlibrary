import { useMemo, useState } from 'react'
import { createRoot } from 'react-dom/client'

export type CommentEditorController = {
  close: () => void
}

export function initCommentEditor($root: HTMLElement): CommentEditorController {
  const d = document.createElement('div')
  d.classList = 'chapter-comment-reply'

  const root = createRoot(d)
  root.render(<Editor />)
  $root.prepend(d)

  return {
    close() {
      root.unmount()
      d.remove()
    },
  }
}

function Editor() {
  const [text, setText] = useState('')

  const valid = useMemo(() => {
    const length = text.trim().length
    return length > 0 && length < 1000
  }, [text])

  return (
    <>
      <textarea
        placeholder={window._('common.replyPlaceholder')}
        name="text"
        className="chapter-comment-reply__text"
        value={text}
        onChange={(e) => setText((e.target as HTMLTextAreaElement).value)}
      />
      <button disabled={!valid} className="btn btn--secondary btn--sm chapter-comment-reply__reply">
        {window._('common.reply')}
      </button>
    </>
  )
}
