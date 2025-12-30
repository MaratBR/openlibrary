import { useState } from 'preact/hooks'
import { useWYSIWYG } from './state'

export function EditorIframe() {
  const [loading, setLoading] = useState(true)
  const state = useWYSIWYG()

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
      {state.renderContent()}
    </>
  )

  function handleLoad(event: Event) {
    const target = event.target
    if (!(target instanceof HTMLIFrameElement)) return

    state.init(target)
    setLoading(false)
  }
}
