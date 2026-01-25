import { useState } from 'preact/hooks'
import { EditorElements } from './EditorElements'
import { createPortal } from 'preact/compat'
import { WYSIWYGEditor } from './wysiwyg'
import { ChapterNameInput } from './ChapterNameInput'

// loads and iframe inside of which we will have the content of the
// chapter
export function EditorIframe() {
  const [loading, setLoading] = useState(true)
  const [elements, setElements] = useState<EditorElements | null>(null)

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
      {!loading && elements && (
        <>
          {createPortal(
            <WYSIWYGEditor
              editorOptions={{
                contentElement: elements.content,
                contentWrapperElement: elements.contentWrapper,
                iframe: elements.iframe,
              }}
            />,
            elements.content,
          )}
          {createPortal(<ChapterNameInput />, elements.contentWrapperHeader)}
        </>
      )}
    </>
  )

  function handleLoad(event: Event) {
    const iframe = event.target
    if (!(iframe instanceof HTMLIFrameElement)) return
    const elements = new EditorElements(iframe)

    setElements(elements)
    setLoading(false)
  }
}
