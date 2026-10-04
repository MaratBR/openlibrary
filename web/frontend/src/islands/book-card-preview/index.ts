import { Dispose } from '@/common/rx'
import { debounce } from '@/common/util/fn'
import { ModalAnimation } from '@/lib/animate'
import { OLIsland, OLIslandMounted } from '@/lib/island'
import './style.scss'
import { autoUpdate, flip, offset, shift, computePosition } from '@floating-ui/react'
import { Schema } from 'effect'

const dataSchema = Schema.Struct({
  selector: Schema.optional(Schema.NullOr(Schema.String)),
})

const DUMMY_ISLAND: OLIslandMounted = {
  setData(_data) {},
  dispose() {},
}

function findBySelector($el: HTMLElement, selector: string) {
  if (selector === ':parent') {
    return $el.parentElement
  }
  return document.querySelector(selector)
}

class BookCardPreviewIsland implements OLIsland {
  mount(el: HTMLElement, data: unknown): OLIslandMounted {
    const { selector } = Schema.decodeUnknownSync(dataSchema)(data)
    if (!selector) return DUMMY_ISLAND
    const $root = findBySelector(el, selector)
    if (!$root) return DUMMY_ISLAND

    let request = 0
    const dismiss = () => {
      ++request
      popover.hide()
    }
    const popover = createPopover(dismiss)
    const disposables: Dispose[] = []
    const $elements = $root.querySelectorAll('[data-book-card-preview]')
    $elements.forEach(($el) => {
      const bookId = $el.getAttribute('data-book-card-preview')
      if (!bookId) {
        return
      }

      if (($el as { _bookCardPreviewReady?: boolean })._bookCardPreviewReady === true) {
        return
      }

      ;($el as { _bookCardPreviewReady?: boolean })._bookCardPreviewReady = true

      const show = async (mobile: boolean) => {
        const current = ++request
        popover.hide()
        try {
          await popover.loadBookContent(bookId, $el, () => current === request, mobile)
        } catch {
          if (current === request) popover.hide()
        }
      }
      const onMouseOver = debounce(() => { void show(false) }, 200)
      const onPointerEnter = (event: Event) => {
        if ((event as PointerEvent).pointerType === 'mouse') onMouseOver()
      }
      const onPointerLeave = (event: Event) => {
        if ((event as PointerEvent).pointerType !== 'mouse') return
        onMouseOver.cancel()
        dismiss()
      }

      let pressTimer: ReturnType<typeof setTimeout> | undefined
      let pressOrigin: { x: number; y: number } | undefined
      let suppressClickUntil = 0
      const cancelPress = () => {
        clearTimeout(pressTimer)
        pressTimer = undefined
        pressOrigin = undefined
      }
      const onPointerDown = (event: Event) => {
        const pointer = event as PointerEvent
        if (pointer.pointerType === 'mouse' || !pointer.isPrimary) return
        cancelPress()
        suppressClickUntil = 0
        pressOrigin = { x: pointer.clientX, y: pointer.clientY }
        pressTimer = setTimeout(() => {
          pressTimer = undefined
          suppressClickUntil = Date.now() + 1500
          void show(true)
        }, 500)
      }
      const onPointerMove = (event: Event) => {
        const pointer = event as PointerEvent
        if (pressOrigin && Math.hypot(pointer.clientX - pressOrigin.x, pointer.clientY - pressOrigin.y) > 10) cancelPress()
      }
      const onPointerUp = () => {
        if (suppressClickUntil) suppressClickUntil = Date.now() + 1500
        cancelPress()
      }
      const onClick = (event: Event) => {
        if (Date.now() < suppressClickUntil) {
          event.preventDefault()
          event.stopImmediatePropagation()
          suppressClickUntil = 0
        }
      }
      const onContextMenu = (event: Event) => {
        if (pressOrigin || Date.now() < suppressClickUntil) event.preventDefault()
      }
      const listeners: [string, EventListener][] = [
        ['pointerenter', onPointerEnter], ['pointerleave', onPointerLeave],
        ['pointerdown', onPointerDown], ['pointermove', onPointerMove],
        ['pointerup', onPointerUp], ['pointercancel', cancelPress],
        ['click', onClick], ['contextmenu', onContextMenu],
      ]
      listeners.forEach(([type, listener]) => $el.addEventListener(type, listener, true))
      const dispose = () => {
        onMouseOver.cancel()
        cancelPress()
        delete ($el as { _bookCardPreviewReady?: boolean })._bookCardPreviewReady
        listeners.forEach(([type, listener]) => $el.removeEventListener(type, listener, true))
      }
      disposables.push(dispose)
    })

    return {
      setData(_data) {},
      dispose() {
        ++request
        popover.dispose()

        disposables.reverse().forEach((cb) => cb())
        disposables.splice(0, disposables.length)
      },
    }
  }
}

function createPopover(dismiss: () => void) {
  const layer = document.createElement('div')
  layer.className = 'BookCardPreviewLayer'
  layer.hidden = true
  let backdropPressed = false
  layer.addEventListener('pointerdown', () => { backdropPressed = true })
  layer.addEventListener('click', () => {
    // The release click from the opening long press may target the new backdrop.
    if (backdropPressed) dismiss()
    backdropPressed = false
  })
  const onKeyDown = (event: KeyboardEvent) => {
    if (!layer.hidden && event.key === 'Escape') dismiss()
  }
  document.addEventListener('keydown', onKeyDown)
  document.body.appendChild(layer)
  const div = document.createElement('div')
  div.className = 'BookCardPreviewPopover'
  div.setAttribute('aria-hidden', 'true')
  div.style.display = 'none'
  layer.appendChild(div)

  const animation = new ModalAnimation(div, 150)
  const cache: Record<string, string> = {}
  let stopPositioning: Dispose | undefined
  let previousOverflow: string | undefined

  const hide = () => {
    layer.hidden = true
    if (previousOverflow !== undefined) {
      document.body.style.overflow = previousOverflow
      previousOverflow = undefined
    }
    stopPositioning?.()
    stopPositioning = undefined
    animation.stop()
    animation.setState(false, { duration: 0, force: true })
  }

  return {
    hide,
    dispose() {
      hide()
      animation.dispose()
      document.removeEventListener('keydown', onKeyDown)
      layer.remove()
    },
    async loadBookContent(bookId: string, anchorEl: Element, isCurrent: () => boolean, mobile: boolean) {
      let html = cache[bookId]
      if (!html) {
        const res = await fetch(`/book/${bookId}/__fragment/preview-card`)
        if (!res.ok) throw new Error('Book preview unavailable')
        html = await res.text()
        cache[bookId] = html
      }
      if (!isCurrent()) return
      div.innerHTML = html
      backdropPressed = false
      layer.classList.toggle('BookCardPreviewLayer--mobile', mobile)
      layer.hidden = false
      if (mobile) {
        previousOverflow = document.body.style.overflow
        document.body.style.overflow = 'hidden'
        div.style.removeProperty('left')
        div.style.removeProperty('top')
        div.style.visibility = 'visible'
        animation.setState(true)
        return
      }
      // Floating UI needs a measurable element before the entrance animation.
      div.style.display = 'block'
      div.style.visibility = 'hidden'
      div.style.transform = 'none'
      const updatePosition = async () => {
        const pos = await computePosition(anchorEl, div, {
          strategy: 'fixed',
          placement: 'right',
          middleware: [offset(12), flip({ fallbackPlacements: ['left', 'top', 'bottom'], padding: 12 }), shift({ padding: 12 })],
        })
        if (!isCurrent()) return
        div.style.left = `${pos.x}px`
        div.style.top = `${pos.y}px`
      }
      await updatePosition()
      if (!isCurrent()) return
      div.style.visibility = 'visible'
      animation.setState(true)
      stopPositioning = autoUpdate(anchorEl, div, updatePosition)
    },
  }
}

export default new BookCardPreviewIsland()
