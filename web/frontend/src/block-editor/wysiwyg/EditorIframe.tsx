import { useState } from 'preact/hooks'
import { useWYSIWYG } from './state'
import { EditorElements } from './editor'

export function EditorIframe() {
  const [loading, setLoading] = useState(true)
  const state = useWYSIWYG()
  const [elements, setElements] = useState<EditorElements>()

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
      {elements && state.editor.getContentElement(elements)}
    </>
  )

  function handleLoad(event: Event) {
    const target = event.target
    if (!(target instanceof HTMLIFrameElement)) return
    if (!target.contentDocument) return

    const editorWrapElement = target.contentDocument.getElementById('BlockEditorWrap')
    const contentElement = target.contentDocument.getElementById('ChapterContent')

    if (!contentElement) return
    if (!editorWrapElement) return

    setLoading(false)
    setElements({
      content: contentElement,
      wrapper: editorWrapElement,
    })
  }
}
