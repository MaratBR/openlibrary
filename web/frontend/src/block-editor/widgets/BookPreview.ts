import { Node, mergeAttributes } from '@tiptap/core'
import { loadBookPreview, parseBookReference } from '@/features/book-preview'

function storedBookID(value: unknown): string | null {
  if (typeof value !== 'string') return null
  try {
    const data = JSON.parse(value)
    return typeof data.bookId === 'string' && /^\d+$/.test(data.bookId)
      ? parseBookReference(data.bookId, window.location.origin)
      : null
  } catch {
    return null
  }
}

export const BookPreview = Node.create<{ previews: Record<string, string> }>({
  name: 'bookPreview',
  group: 'block',
  atom: true,
  draggable: true,

  addOptions() {
    return { previews: {} }
  },

  addAttributes() {
    return {
      bookId: {
        default: null,
        rendered: false,
        parseHTML: (element) => storedBookID(element.getAttribute('data')),
      },
    }
  },

  parseHTML() {
    return [{ tag: 'ol-widget[name="book-preview"]' }]
  },

  renderHTML({ node, HTMLAttributes }) {
    return [
      'ol-widget',
      mergeAttributes(HTMLAttributes, {
        name: 'book-preview',
        data: JSON.stringify({ bookId: node.attrs.bookId }),
      }),
    ]
  },

  addNodeView() {
    const previews = this.options.previews
    return ({ node, editor, view, getPos }) => {
      const document = view.dom.ownerDocument
      const dom = document.createElement('div')
      dom.className = 'BookPreviewWidget'
      dom.contentEditable = 'false'
      let currentNode = node
      let disposePreview: (() => void) | undefined

      function configure() {
        disposePreview?.()
        const form = document.createElement('form')
        form.className = 'BookPreviewWidget-form'
        const label = document.createElement('label')
        label.textContent = window._('editor.bookPreview.bookLink')
        const input = document.createElement('input')
        input.className = 'Input'
        input.type = 'text'
        input.required = true
        input.value = currentNode.attrs.bookId ?? ''
        label.append(input)
        const error = document.createElement('p')
        error.setAttribute('role', 'alert')
        error.hidden = true
        const button = document.createElement('button')
        button.type = 'submit'
        button.className = 'Btn Btn--primary'
        button.textContent = window._('editor.bookPreview.insert')
        form.append(label, error, button)
        form.addEventListener('submit', (event) => {
          event.preventDefault()
          const id = parseBookReference(input.value, window.location.origin)
          if (!id) {
            input.setAttribute('aria-invalid', 'true')
            error.textContent = window._('editor.bookPreview.invalid')
            error.hidden = false
            return
          }
          const pos = getPos()
          if (typeof pos !== 'number') return
          const unchanged = id === currentNode.attrs.bookId
          editor.view.dispatch(editor.state.tr.setNodeMarkup(pos, undefined, { bookId: id }))
          if (unchanged) render()
          editor.commands.focus()
        })
        dom.replaceChildren(form)
      }

      function render() {
        disposePreview?.()
        if (!currentNode.attrs.bookId) {
          configure()
          return
        }
        const preview = document.createElement('div')
        const change = document.createElement('button')
        change.type = 'button'
        change.className = 'Btn Btn--sm Btn--ghost'
        change.textContent = window._('editor.bookPreview.change')
        change.addEventListener('click', configure)
        dom.replaceChildren(preview)
        disposePreview = loadBookPreview(
          preview,
          currentNode.attrs.bookId,
          previews[currentNode.attrs.bookId],
        )
        preview.querySelector('.BookPreviewWidget-actions')?.append(change)
      }
      render()

      return {
        dom,
        update(updatedNode) {
          if (updatedNode.type !== currentNode.type) return false
          const changed = updatedNode.attrs.bookId !== currentNode.attrs.bookId
          currentNode = updatedNode
          if (changed) render()
          return true
        },
        stopEvent: (event) => !!(event.target as Element).closest('form, button, a'),
        ignoreMutation: () => true,
        destroy: () => disposePreview?.(),
      }
    }
  },
})
