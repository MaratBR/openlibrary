function getScrollbarWidth() {
  // Creating invisible container
  const outer = document.createElement('div')
  outer.style.visibility = 'hidden'
  outer.style.overflow = 'scroll' // forcing scrollbar to appear
  document.body.appendChild(outer)

  // Creating inner element and placing it in the container
  const inner = document.createElement('div')
  outer.appendChild(inner)

  // Calculating difference between container's full width and the child width
  const scrollbarWidth = outer.offsetWidth - inner.offsetWidth

  // Removing temporary elements from the DOM
  outer.parentNode!.removeChild(outer)

  return scrollbarWidth
}

export function initScrollbarWidth() {
  const update = () => {
    if (document.body.scrollHeight > window.innerHeight) {
      document.body.style.setProperty('--scrollbar-width', `${getScrollbarWidth()}px`)
    } else {
      document.body.style.setProperty('--scrollbar-width', '0px')
    }
  }

  const observer = new ResizeObserver(() => update())
  observer.observe(document.body)

  window.requestAnimationFrame(update)
}
