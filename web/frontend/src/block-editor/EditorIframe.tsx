import { SyntheticEvent, useState } from 'react'
import { EditorElements } from './EditorElements'
import { WYSIWYGEditor } from './wysiwyg'
import { ChapterNameInput } from './ChapterNameInput'
import { createPortal } from 'react-dom'

// loads and iframe inside of which we will have the content of the
// chapter
export function EditorIframe({ initialContent }: { initialContent: string }) {
  const [loading, setLoading] = useState(true)
  const [elements, setElements] = useState<EditorElements | null>(null)

  return (
    <>
      <iframe
        title={window._('editor.chapterContentEditor')}
        onLoad={handleLoad}
        name="editor"
        style={{ width: '100%', height: '100%' }}
        src="/books-manager/__fragment/chapter-content-iframe"
      />
      {loading && (
        <div className="absolute inset-0 flex items-center justify-center">
          <span className="loader" />
        </div>
      )}
      {!loading && elements && (
        <>
          {createPortal(
            <WYSIWYGEditor
              editorOptions={{
                initialContent,
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

  function handleLoad(event: SyntheticEvent<HTMLIFrameElement>) {
    const iframe = event.target
    if (!(iframe instanceof HTMLIFrameElement)) return
    const elements = new EditorElements(iframe)

    setElements(elements)
    setLoading(false)
  }
}
