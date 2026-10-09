import { SyntheticEvent, useEffect, useRef, useState } from 'react'
import { EditorElements } from './EditorElements'
import { WYSIWYGEditor } from './wysiwyg'
import { ChapterNameInput } from './ChapterNameInput'
import { createPortal } from 'react-dom'
import { appRuntime } from '@/effect/runtime'
import { FontsLoader } from '@/features/fonts-loader/loader'
import { useFavoriteFonts } from './fonts/state'
import { EditorFonts } from './fonts/service'

const fontsLoader = appRuntime.runSync(FontsLoader)
const editorFonts = appRuntime.runSync(EditorFonts)

// loads and iframe inside of which we will have the content of the
// chapter
export function EditorIframe({
  initialContent,
  chapterFonts,
}: {
  initialContent: string
  chapterFonts: string[]
}) {
  const favoriteFonts = useFavoriteFonts()
  const loadVersion = useRef(0)
  const attachedIframe = useRef<HTMLIFrameElement | null>(null)
  const [loading, setLoading] = useState(true)
  const [elements, setElements] = useState<EditorElements | null>(null)

  useEffect(() => {
    return () => {
      loadVersion.current++
      if (attachedIframe.current) {
        void appRuntime.runPromise(fontsLoader.detachIframe(attachedIframe.current))
        attachedIframe.current = null
      }
    }
  }, [])

  useEffect(() => {
    if (elements === null) return
    void appRuntime
      .runPromise(fontsLoader.addFonts([...chapterFonts, ...favoriteFonts]))
      .catch(() => undefined)
  }, [elements, chapterFonts, favoriteFonts])

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

    const version = ++loadVersion.current
    setLoading(true)
    attachedIframe.current = iframe

    const loadFonts = async () => {
      await appRuntime.runPromise(fontsLoader.attachIframe(iframe))
      await appRuntime.runPromise(editorFonts.initialize()).catch(() => undefined)
      if (version !== loadVersion.current) return
      const fonts = [...new Set([...chapterFonts, ...editorFonts.getState().favoriteFonts])]
      await appRuntime.runPromise(fontsLoader.addFonts(fonts))
      const iframeDocument = iframe.contentDocument
      if (!iframeDocument) return
      // Stylesheet readiness alone does not mean the font files have loaded.
      await Promise.all(
        fonts.flatMap((font) =>
          [document, iframeDocument].map((target) =>
            target.fonts.load(`16px ${JSON.stringify(font)}`),
          ),
        ),
      )
    }

    let timeout: ReturnType<typeof setTimeout> | undefined
    await Promise.race([
      loadFonts().catch(() => undefined),
      new Promise<void>((resolve) => {
        timeout = setTimeout(resolve, 5000)
      }),
    ])
    clearTimeout(timeout)
    if (version !== loadVersion.current) return
    setElements(elements)
    setLoading(false)
  }
}
