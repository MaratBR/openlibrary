import { SyntheticEvent, useEffect, useState } from 'react'
import { EditorElements } from './EditorElements'
import { WYSIWYGEditor } from './wysiwyg'
import { ChapterNameInput } from './ChapterNameInput'
import { createPortal } from 'react-dom'
import { appRuntime } from '@/effect/runtime'
import { FontsLoader } from '@/features/fonts-loader/loader'

const fontsLoader = appRuntime.runSync(FontsLoader)

// loads and iframe inside of which we will have the content of the
// chapter
export function EditorIframe({ initialContent }: { initialContent: string }) {
  const [loading, setLoading] = useState(true)
  const [elements, setElements] = useState<EditorElements | null>(null)

  useEffect(() => {
    if (elements === null) return
    return () => {
      void appRuntime.runPromise(fontsLoader.detachIframe(elements.iframe))
    }
  }, [elements])

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
          <span className="Loader" />
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

  async function handleLoad(event: SyntheticEvent<HTMLIFrameElement>) {
    const iframe = event.target
    if (!(iframe instanceof HTMLIFrameElement)) return
    const elements = new EditorElements(iframe)

    await appRuntime.runPromise(fontsLoader.attachIframe(iframe))
    setElements(elements)
    setLoading(false)
  }
}
