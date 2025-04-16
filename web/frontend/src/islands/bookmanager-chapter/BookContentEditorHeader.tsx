import { InternalEvent } from '@/lib/event'
import { Editor } from '@tiptap/core'
import { useEffect, useState } from 'preact/hooks'
import styles from './BookContentEditorHeader.module.scss'
import { RefObject } from 'preact'
import BookContentEditorHeadingMenu from './BookContentEditorHeadingMenu'
import { Level } from '@tiptap/extension-heading'

type State = {
  bold: boolean
  italic: boolean
  strike: boolean
  underline: boolean
  code: boolean
  typography: null | number
  textAlign: 'left' | 'center' | 'right'
}

const defaultState: State = {
  bold: false,
  italic: false,
  strike: false,
  underline: false,
  code: false,
  typography: null,
  textAlign: 'left',
}

function getTextAlign(editor: Editor): State['textAlign'] {
  if (editor.isActive({ textAlign: 'left' })) {
    return 'left'
  }

  if (editor.isActive({ textAlign: 'center' })) {
    return 'center'
  }

  if (editor.isActive({ textAlign: 'right' })) {
    return 'right'
  }

  return 'left'
}

export default function BookContentEditorHeader({
  editorUpdateEvent,
  editorRef,
}: {
  editorUpdateEvent: InternalEvent<Editor>
  editorRef: RefObject<Editor | null>
}) {
  const [state, setState] = useState<State>(defaultState)

  useEffect(() => {
    const callback = (editor: Editor) => {
      setState({
        bold: editor.isActive('bold'),
        italic: editor.isActive('italic'),
        strike: editor.isActive('strike'),
        underline: editor.isActive('underline'),
        code: editor.isActive('code'),
        typography: editor.isActive('heading') ? editor.getAttributes('heading').level : null,
        textAlign: getTextAlign(editor),
      })
    }

    editorUpdateEvent.subscribe(callback)

    return () => {
      editorUpdateEvent.unsubscribe(callback)
    }
  }, [editorUpdateEvent])

  function toggleBold() {
    if (editorRef.current) editorRef.current.chain().focus().toggleBold().run()
  }

  function toggleItalic() {
    if (editorRef.current) editorRef.current.chain().focus().toggleItalic().run()
  }

  function toggleUnderline() {
    if (editorRef.current) editorRef.current.chain().focus().toggleUnderline().run()
  }

  function toggleTextAlign(textAlign: State['textAlign']) {
    if (editorRef.current)
      editorRef.current
        .chain()
        .focus()
        .setTextAlign(textAlign || 'left')
        .run()
  }

  return (
    <header class="ol-book-editor__header">
      <div class="ol-container flex items-center gap-2 px-0">
        <section class={styles.section}>
          <label role="button" class={`ol-btn ol-btn--ghost ${styles.btn}`}>
            <input type="checkbox" id="bold" checked={state.bold} onInput={toggleBold} />
            <span class="material-symbols-outlined">format_bold</span>
          </label>

          <label role="button" class={`ol-btn ol-btn--ghost ${styles.btn}`}>
            <input type="checkbox" id="italic" checked={state.italic} onInput={toggleItalic} />
            <span class="material-symbols-outlined">format_italic</span>
          </label>

          <label role="button" class={`ol-btn ol-btn--ghost ${styles.btn}`}>
            <input
              type="checkbox"
              id="underline"
              checked={state.underline}
              onInput={toggleUnderline}
            />
            <span class="material-symbols-outlined">format_underlined</span>
          </label>
        </section>

        <div class={styles.divider} />

        <section class={styles.section}>
          <label role="button" class={`ol-btn ol-btn--ghost ${styles.btn}`}>
            <input
              type="checkbox"
              id="left"
              checked={state.textAlign === 'left'}
              onInput={() => toggleTextAlign('left')}
            />
            <span class="material-symbols-outlined">format_align_left</span>
          </label>
          <label role="button" class={`ol-btn ol-btn--ghost ${styles.btn}`}>
            <input
              type="checkbox"
              id="center"
              checked={state.textAlign === 'center'}
              onInput={() => toggleTextAlign('center')}
            />
            <span class="material-symbols-outlined">format_align_center</span>
          </label>
          <label role="button" class={`ol-btn ol-btn--ghost ${styles.btn}`}>
            <input
              type="checkbox"
              id="right"
              checked={state.textAlign === 'right'}
              onInput={() => toggleTextAlign('right')}
            />
            <span class="material-symbols-outlined">format_align_right</span>
          </label>
        </section>

        <div class={styles.divider} />

        <BookContentEditorHeadingMenu
          heading={state.typography}
          onChange={(heading) => {
            if (editorRef.current) {
              if (heading === null) {
                editorRef.current.chain().focus().setParagraph().run()
              } else {
                editorRef.current
                  .chain()
                  .focus()
                  .setHeading({ level: heading as Level })
                  .run()
              }
            }
          }}
        />
      </div>
    </header>
  )
}
