import { useState } from 'preact/hooks'
import { ChapterContentEditor } from './state'
import { useChapterState } from '../state'

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

    const contentRoot = target.contentDocument?.getElementById('BookReaderContent')
    if (!contentRoot) throw new Error('cannot find #BookReader content within iframe')

    new ChapterContentEditor({
      element: contentRoot,
      placeholder: 'Placeholder here!',
      state,
    })

    setLoading(false)
  }
}
