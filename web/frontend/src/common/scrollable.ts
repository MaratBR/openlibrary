import 'overlayscrollbars/overlayscrollbars.css'
;(() => {
  import('overlayscrollbars').then(({ OverlayScrollbars }) => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    ;(window as any).OverlayScrollbars = OverlayScrollbars
    document.dispatchEvent(new CustomEvent('overlayscrollbars-ready'))
  })
})()

function enableDragScroll(element: HTMLElement) {
  let isDown = false
  let startX: number, startY: number, scrollLeft: number, scrollTop: number
  let hasDragged = false

  element.addEventListener('mousedown', (e) => {
    isDown = true
    hasDragged = false
    element.classList.add('dragging')
    startX = e.pageX - element.offsetLeft
    startY = e.pageY - element.offsetTop
    scrollLeft = element.scrollLeft
    scrollTop = element.scrollTop
    e.preventDefault()
  })

  element.addEventListener('mouseleave', () => {
    isDown = false
    element.classList.remove('dragging')
  })

  element.addEventListener('mouseup', () => {
    isDown = false
    element.classList.remove('dragging')
  })

  element.addEventListener('mousemove', (e) => {
    if (!isDown) return
    e.preventDefault()
    const x = e.pageX - element.offsetLeft
    const y = e.pageY - element.offsetTop
    const walkX = x - startX
    const walkY = y - startY

    if (Math.abs(walkX) > 3 || Math.abs(walkY) > 3) {
      hasDragged = true
    }

    element.scrollLeft = scrollLeft - walkX
    element.scrollTop = scrollTop - walkY
  })

  // Prevent clicks on children if a drag occurred
  element.addEventListener(
    'click',
    (e) => {
      if (hasDragged) {
        e.stopPropagation()
        e.preventDefault()
        hasDragged = false
      }
    },
    true,
  ) // use capture so it catches before links
}

declare global {
  interface Window {
    enableDragScroll: typeof enableDragScroll
  }
}

window.enableDragScroll = enableDragScroll
