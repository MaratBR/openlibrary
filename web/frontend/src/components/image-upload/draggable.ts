export type DragEvent = {
  x: number
  y: number
}

export type DraggableInit = {
  element: HTMLElement
  bounds: ResizableImageBounds
  // eslint-disable-next-line no-unused-vars
  onDrag: (event: DragEvent) => void
}

export type ResizableImageBounds = {
  x0: number
  y0: number
  x1: number
  y1: number
}

export function initDraggable({ element, bounds, onDrag }: DraggableInit) {
  let isDragging = false
  let startX: number
  let startY: number

  const onMouseDown = (event: MouseEvent) => {
    // Prevent text selection during drag
    event.preventDefault()

    if (event.button !== 0) return

    // Store the initial mouse position relative to the element
    startX = event.clientX - element.offsetLeft
    startY = event.clientY - element.offsetTop

    // Add styles to indicate dragging
    element.style.cursor = 'grabbing'
    element.style.userSelect = 'none'

    // Start dragging
    isDragging = true

    // Add move and up event listeners to document to handle drag outside element
    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
  }

  const onMouseMove = (event: MouseEvent) => {
    if (!isDragging) return

    // Calculate new position
    let newX = event.clientX - startX,
      newY = event.clientY - startY

    if (newX < bounds.x0) newX = bounds.x0
    else if (newX > bounds.x1) newX = bounds.x1

    if (newY < bounds.y0) newY = bounds.y0
    else if (newY > bounds.y1) newY = bounds.y1

    onDrag({
      x: newX,
      y: newY,
    })

    // Update element position
    element.style.left = `${newX}px`
    element.style.top = `${newY}px`
  }

  const onMouseUp = () => {
    // Reset dragging state
    isDragging = false

    // Restore styles
    element.style.cursor = 'grab'
    element.style.userSelect = 'auto'

    // Remove event listeners
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
  }

  // Initial cursor style
  element.style.cursor = 'grab'

  // Add mousedown listener
  element.addEventListener('mousedown', onMouseDown)

  window.requestAnimationFrame(() => {
    element.style.position = 'absolute'
  })

  // Return cleanup function
  return () => {
    element.removeEventListener('mousedown', onMouseDown)
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
  }
}
