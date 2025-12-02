import { useState } from 'preact/hooks'
import { useChapterState } from '../state'
import { ChapterContentEditor } from './editor'

export function EditorIframe() {
  const [loading, setLoading] = useState(true)
  const state = useChapterState()

  return (
    <>
      <iframe
        onLoad={handleLoad}
        name="editor"
        style={{ width: '100%', height: '100%' }}
        src="/books-manager/__fragment/chapter-content-iframe"
      />
      {loading && (
        <div class="absolute inset-0 flex items-center justify-center">
          <span class="loader" />
        </div>
      )}
    </>
  )

  function handleLoad(event: Event) {
    const target = event.target
    if (!(target instanceof HTMLIFrameElement)) return

    new ChapterContentEditor({
      element: target,
      placeholder: 'Placeholder here!',
      state,
    })

    setLoading(false)
  }
}
