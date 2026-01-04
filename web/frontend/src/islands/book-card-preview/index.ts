import { Dispose } from '@/common/rx'
import { debounce } from '@/common/util/fn'
import { ModalAnimation } from '@/lib/animate'
import { OLIsland, OLIslandMounted } from '@/lib/island'
import { computePosition } from '@floating-ui/react'
import z from 'zod'

const dataSchema = z.object({
  selector: z.string().nullable().optional(),
})

const DUMMY_ISLAND: OLIslandMounted = {
  setData(_data) {},
  dispose() {},
}

class BookCardPreviewIsland implements OLIsland {
  mount(_el: HTMLElement, data: unknown): OLIslandMounted {
    const { selector } = dataSchema.parse(data)
    if (!selector) return DUMMY_ISLAND
    const $root = document.querySelector(selector)
    if (!$root) return DUMMY_ISLAND

    const popover = createPopover()

    const disposables: Dispose[] = []
    const $elements = $root.querySelectorAll('[data-book-card-preview]')
    $elements.forEach(($el) => {
      const bookId = $el.getAttribute('data-book-card-preview')
      if (!bookId) {
        return
      }

      const onMouseOver = debounce(() => {
        popover.animation.setShow(false, 0)
        popover.loadBookContent(bookId, $el).then(() => {
          popover.animation.setShow(true)
        })
      }, 200)

      const onMouseLeave = () => {
        onMouseOver.cancel()
        popover.animation.setShow(false)
      }

      $el.addEventListener('mouseover', onMouseOver)
      $el.addEventListener('mouseleave', onMouseLeave)

      const dispose = () => {
        $el.removeEventListener('mouseover', onMouseOver)
        $el.removeEventListener('mouseleave', onMouseLeave)
      }
      disposables.push(dispose)
    })

    return {
      setData(_data) {},
      dispose() {
        popover.dispose()

        disposables.reverse().forEach((cb) => cb())
        disposables.splice(0, disposables.length)
      },
    }
  }
}

function createPopover() {
  const div = document.createElement('div')
  div.className = 'bg-card absolute z-10 shadow-2xl rounded-2xl overflow-hidden'
  div.style.display = 'none'
  div.style.maxWidth = '400px'
  document.body.appendChild(div)

  const animation = new ModalAnimation(div, 150)
  const cache: Record<string, string> = {}

  return {
    animation,
    dispose() {
      animation.dispose()
      div.remove()
    },
    async loadBookContent(bookId: string, anchorEl: Element) {
      if (cache[bookId]) {
        div.innerHTML = cache[bookId]
      } else {
        const res = await fetch(`/book/${bookId}/__fragment/preview-card`)
        const html = await res.text()
        div.innerHTML = html
        cache[bookId] = html
      }
      const pos = await computePosition(anchorEl, div, { placement: 'top-end' })
      div.style.left = `${pos.x}px`
      div.style.top = `${pos.y}px`
    },
  }
}

export default new BookCardPreviewIsland()
