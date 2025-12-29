import { useState } from 'preact/hooks'
import { createPortal } from 'preact/compat'
import { useWYSIWYG } from './state'

export function EditorIframe() {
  const [loading, setLoading] = useState(true)
  const state = useWYSIWYG()
  const [contentElement, setContentElement] = useState<HTMLElement>()

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
      {contentElement && createPortal(state.editor.getContentElement(), contentElement)}
    </>
  )

  function handleLoad(event: Event) {
    const target = event.target
    if (!(target instanceof HTMLIFrameElement)) return
    if (!target.contentDocument) return

    const contentElement = target.contentDocument.getElementById('BookReaderContent')
    if (!contentElement) return

    setLoading(false)
    setContentElement(contentElement)
  }
}
