// Keep IDs as decimal strings: book IDs can exceed JavaScript's safe integer range.
export function parseBookReference(value: string, origin: string): string | null {
  let id = value.trim()
  if (!/^\d+$/.test(id)) {
    try {
      const url = new URL(id, origin)
      if (url.origin !== origin) return null
      const match = /^\/book\/(\d+)(?:-[^/]*)?\/?$/.exec(url.pathname)
      if (!match) return null
      id = match[1]
    } catch {
      return null
    }
  }
  id = id.replace(/^0+/, '')
  if (!id || id.length > 19 || (id.length === 19 && id > '9223372036854775807')) return null
  return id
}

export function loadBookPreview(
  container: HTMLElement,
  bookID: string,
  initialHTML?: string,
): () => void {
  const document = container.ownerDocument
  const controller = new AbortController()
  const status = document.createElement('p')
  status.setAttribute('role', 'status')
  status.textContent = window._('editor.bookPreview.loading')
  const content = document.createElement('div')
  const link = document.createElement('a')
  link.className = 'BookPreviewWidget-open'
  link.href = `/book/${bookID}`
  link.textContent = window._('editor.bookPreview.open')
  const arrow = document.createElement('i')
  arrow.className = 'fa-solid fa-arrow-right text-xs'
  arrow.setAttribute('aria-hidden', 'true')
  link.append(arrow)
  const retry = document.createElement('button')
  retry.type = 'button'
  retry.className = 'Btn Btn--sm Btn--outline'
  retry.textContent = window._('editor.bookPreview.retry')
  retry.hidden = true
  const actions = document.createElement('div')
  actions.className = 'BookPreviewWidget-actions'
  actions.append(link, retry)
  container.replaceChildren(status, content, actions)

  async function load() {
    retry.hidden = true
    status.hidden = false
    status.textContent = window._('editor.bookPreview.loading')
    content.replaceChildren()
    container.setAttribute('aria-busy', 'true')
    try {
      const response = await fetch(`/book/${bookID}/__fragment/preview-card`, {
        signal: controller.signal,
        credentials: 'same-origin',
      })
      if (!response.ok) throw new Error('Preview unavailable')
      const html = await response.text()
      if (controller.signal.aborted) return
      // Only use the same-origin server fragment; never persisted widget HTML.
      const template = document.createElement('template')
      template.innerHTML = html
      if (!template.content.querySelector('.BookCardPreview')) throw new Error('Invalid preview')
      content.replaceChildren(template.content)
      status.hidden = true
    } catch {
      if (controller.signal.aborted) return
      status.textContent = window._('editor.bookPreview.unavailable')
      retry.hidden = false
    } finally {
      if (!controller.signal.aborted) container.removeAttribute('aria-busy')
    }
  }
  retry.addEventListener('click', () => void load())
  if (initialHTML !== undefined) {
    const template = document.createElement('template')
    template.innerHTML = initialHTML
    content.replaceChildren(template.content)
    status.hidden = true
    retry.hidden = !!content.querySelector('.BookCardPreview')
  } else {
    void load()
  }
  return () => controller.abort()
}
