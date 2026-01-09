import { render } from 'preact'
import { useMemo, useState } from 'preact/hooks'

export type CommentEditorController = {
  close: () => void
}

export function initCommentEditor($root: HTMLElement): CommentEditorController {
  const d = document.createElement('div')
  d.classList = 'chapter-comment-reply'
  render(<Editor />, d)
  $root.prepend(d)

  return {
    close() {
      render(null, d)
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
        class="chapter-comment-reply__text"
        value={text}
        onChange={(e) => setText((e.target as HTMLTextAreaElement).value)}
      />
      <button disabled={!valid} class="chapter-comment-reply__reply btn btn--secondary btn--sm">
        {window._('common.reply')}
      </button>
    </>
  )
}
