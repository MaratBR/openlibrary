import StarterKit from '@tiptap/starter-kit'
import TextStyle from '@tiptap/extension-text-style'
import Typography from '@tiptap/extension-typography'
import HorizontalRule from '@tiptap/extension-horizontal-rule'
import { Editor, EditorOptions } from '@tiptap/core'
import { httpUpdateReview, ratingSchema, ReviewDto, reviewDtoSchema } from './api'
import { useEffect, useRef, useState } from 'preact/hooks'
import { PreactIslandProps } from '../common'

export default function ReviewEditor({ rootElement }: PreactIslandProps) {
  const rootEl = useRef<HTMLDivElement | null>(null)
  const editor = useRef<Editor | null>()
  const [active, setActive] = useState({ bold: false, italic: false })
  const [rating, setRating] = useState(0)
  const [saving, setSaving] = useState(false)
  const bookId = getBookId()

  useEffect(() => {
    const data = getExistingReviewData()

    if (data) {
      setRating(data.rating)
    }

    editor.current = createEditor(rootEl.current!, {
      content: data?.content ?? '',
      onTransaction: () => {
        setActive({
          bold: editor.current!.isActive('bold'),
          italic: editor.current!.isActive('italic'),
        })
      },
    })
  }, [])

  function handleSave() {
    if (saving) return
    setSaving(true)

    httpUpdateReview(bookId, {
      content: editor.current!.getHTML(),
      rating: ratingSchema.parse(rating),
    })
      .then((review) => {
        rootElement.dispatchEvent(new CustomEvent('review:updated', { detail: review }))
        window.toast({ title: window._('reviews.updated') })
      })
      .finally(() => {
        setSaving(false)
      })
  }

  return (
    <div>
      <RatingInput scale={0.5} value={rating} onInput={setRating} />

      <div class="ol-review-editor__toolbar mt-4">
        <button
          class={active.bold ? 'active' : ''}
          onClick={() => editor.current?.chain().focus().toggleBold().run()}
        >
          <span class="material-symbols-outlined">format_bold</span>
        </button>

        <button
          class={active.italic ? 'active' : ''}
          onClick={() => editor.current?.chain().focus().toggleItalic().run()}
        >
          <span class="material-symbols-outlined">format_italic</span>
        </button>
      </div>

      <div ref={rootEl} class="ol-review-editor__content __user-content __user-content--editor" />

      <button
        class="ol-btn ol-btn--lg ol-btn--primary rounded-full mt-3"
        onClick={() => handleSave()}
      >
        {window._('common.save')}
      </button>
    </div>
  )
}

function RatingInput({
  scale = 1,
  value,
  onInput,
}: {
  scale?: number
  value: number
  // eslint-disable-next-line no-unused-vars
  onInput: (value: number) => void
}) {
  const rootElement = useRef<HTMLDivElement | null>(null)
  const disableHalfPoints = true

  function handleClick(event: MouseEvent) {
    if (!rootElement) return

    const rect = rootElement.current!.getBoundingClientRect()
    const x = event.clientX - rect.left
    const width = rect.width
    let newValue = Math.max(Math.min(Math.ceil((x / width) * 10), 10), 1)

    if (disableHalfPoints && newValue % 2 === 1) {
      newValue += 1
    }

    if (newValue !== value) {
      onInput(newValue)
    }
  }

  return (
    <div
      ref={rootElement}
      class="relative cursor-pointer"
      onClick={handleClick}
      style={`width:${540 * scale}px;height:${100 * scale}px`}
    >
      <div class="ol-star-background h-full w-full opacity-15" />
      <div
        class="absolute left-0 top-0 ol-star-background ol-star-background--filled h-full"
        style={`width:${calcPerc(value)}%;background-size:auto ${scale * 100}px`}
      />
    </div>
  )
}

function calcPerc(value: number): number {
  return ((500 * (value / 10) + Math.floor(value / 2) * 10) / 540) * 100
}

/**
 * Initializes and returns a new instance of the Editor with specified extensions and options.
 *
 * @param editorElement - The HTML element where the editor will be mounted.
 * @param options - Optional configuration settings to customize the editor instance.
 *
 * @returns An instance of Editor configured with a set of extensions and options.
 */
function createEditor(editorElement: HTMLElement, options?: Partial<EditorOptions>) {
  return new Editor({
    element: editorElement,
    content: '',
    extensions: [
      StarterKit.configure({
        horizontalRule: false,
        codeBlock: false,
        heading: false,
        code: { HTMLAttributes: { class: 'inline', spellcheck: 'false' } },
        dropcursor: { width: 2, class: 'ProseMirror-dropcursor border' },
      }),
      TextStyle,
      Typography,
      HorizontalRule,
    ],
    ...options,
  })
}

/**
 * Finds a hidden element containing review data in JSON format and parses it
 * into a ReviewDto object using `reviewDtoSchema`. If the element is not found,
 * null is returned.
 *
 * This function is used to pre-fill the review text area with existing review
 * data when the review editor is rehydrated on the server.
 *
 * @returns {ReviewDto | null}
 */
function getExistingReviewData(): ReviewDto | null {
  const el = document.getElementById('island-review-editor-data')

  if (el instanceof HTMLTemplateElement) {
    const data = JSON.parse(el.content.textContent || '')
    return reviewDtoSchema.parse(data)
  }

  return null
}

function getBookId() {
  const v = window.__server__?.bookId
  if (typeof v === 'string' && v) {
    return v
  }

  throw new Error('could not find bookId, __server__.bookId is not set')
}
