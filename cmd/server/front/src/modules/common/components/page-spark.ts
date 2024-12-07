import './page-spark.css'

const THICKNESS = 25

export function spark(
  stepWidth: number,
  appearDuration: number,
  disappearDuration: number,
  stillDuration: number,
  segmentDelay: number,
  initialPosition: number = 0,
) {
  const svg = document.getElementById('page-spark')
  if (!(svg instanceof SVGSVGElement)) throw new Error('spark not found')
  const g = document.createElementNS('http://www.w3.org/2000/svg', 'g')
  g.classList.add('page-spark__instance')
  g.style.setProperty('--spark-appear-duration', `${appearDuration}ms`)
  g.style.setProperty('--spark-disappear-duration', `${disappearDuration}ms`)
  g.style.setProperty('--spark-segment-delay', `${segmentDelay}ms`)
  g.dataset.spark = 'false'
  generateSVGSteps(
    g,
    window.innerWidth + THICKNESS,
    window.innerHeight + THICKNESS,
    stepWidth,
    initialPosition,
    COLORS,
    THICKNESS,
  )

  svg.appendChild(g)

  window.requestAnimationFrame(() => {
    g.dataset.spark = 'true'

    setTimeout(() => {
      window.requestAnimationFrame(() => {
        g.dataset.spark = 'false'
      })

      setTimeout(() => {
        g.remove()
      }, disappearDuration)
    }, appearDuration + stillDuration)
  })
}

function createSparkRoot(thickness: number) {
  const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg')
  svg.id = 'page-spark'
  svg.setAttribute('width', '100%')
  svg.setAttribute('height', '100%')
  svg.setAttribute('filter', 'url(#blur)')
  svg.setAttribute(
    'viewBox',
    `0 0 ${window.innerWidth + thickness} ${window.innerHeight + thickness}`,
  )
  svg.innerHTML = `
  <defs>
    <filter id="blur">
      <feGaussianBlur in="SourceGraphic" stdDeviation="15" />
    </filter>
  </defs>`

  const root = document.createElement('div')
  root.style.setProperty('--thickness', `${thickness}px`)

  root.classList.add('page-spark')
  root.appendChild(svg)

  document.body.appendChild(root)
}

const MAX_COLORS = 21
const COLORS = Array.from({ length: MAX_COLORS }).map((_, index) => {
  return `hsl(${Math.round((index / MAX_COLORS) * 360)}, 100%, 50%)`
})

function generateSVGSteps(
  element: SVGElement,
  viewportWidth: number,
  viewportHeight: number,
  suggestedStepWidth: number,
  startingPosition: number,
  colors: string[],
  thickness: number = 20,
) {
  const totalLength = viewportHeight * 2 + viewportWidth * 2
  const stepWidth = Math.round(totalLength / colors.length / 3)

  let remainingLength = totalLength
  let colorIndex = 0
  let backwards = false
  let elementIndex = 0
  const totalElements = Math.ceil(totalLength / stepWidth)

  function Rect(): SVGRectElement {
    return document.createElementNS('http://www.w3.org/2000/svg', 'rect')
  }

  function getPosition() {
    let position = totalLength - remainingLength + startingPosition
    if (position > viewportWidth * 2 + viewportHeight * 2) {
      position = position - viewportWidth * 2 - viewportHeight * 2
    }
    return position
  }

  while (remainingLength > 0) {
    let rect2: SVGRectElement | null = null
    const rect = Rect()

    const color = colors[colorIndex]
    rect.setAttribute('fill', color)
    if (backwards) {
      colorIndex--
      if (colorIndex < 0) {
        colorIndex = 1
        backwards = false
      }
    } else {
      colorIndex++
      if (colorIndex > colors.length - 1) {
        colorIndex = colors.length - 2
        backwards = true
      }
    }

    const position = getPosition()

    if (position < viewportWidth) {
      let seg1Size = stepWidth

      if (position + seg1Size > viewportWidth) {
        // second segment is required
        seg1Size = viewportWidth - position
        rect2 = Rect()
        rect2.setAttribute('fill', color)
        rect2.classList.add('position-right')
        rect2.setAttribute('y', '0')
        rect2.setAttribute('x', (viewportWidth - thickness).toString())
        rect2.setAttribute('width', thickness.toString())
        rect2.setAttribute('height', (stepWidth - seg1Size).toString())
      }

      // top
      rect.classList.add('position-top')
      rect.setAttribute('y', '0')
      rect.setAttribute('x', position.toString())
      rect.setAttribute('width', seg1Size.toString())
      rect.setAttribute('height', thickness.toString())
    } else if (position < viewportWidth + viewportHeight) {
      let seg1Size = stepWidth

      if (position + seg1Size > viewportWidth + viewportHeight) {
        // second segment is required
        seg1Size = viewportWidth + viewportHeight - position
        const seg2Size = stepWidth - seg1Size
        rect2 = Rect()
        rect2.setAttribute('fill', color)
        rect2.classList.add('position-bottom')
        rect2.setAttribute('y', (viewportHeight - thickness).toString())
        rect2.setAttribute('x', (viewportWidth - seg2Size).toString())
        rect2.setAttribute('height', thickness.toString())
        rect2.setAttribute('width', seg2Size.toString())
      }

      // right
      rect.classList.add('position-right')
      rect.setAttribute('y', (position - viewportWidth).toString())
      rect.setAttribute('x', (viewportWidth - thickness).toString())
      rect.setAttribute('height', seg1Size.toString())
      rect.setAttribute('width', thickness.toString())
    } else if (position < viewportWidth * 2 + viewportHeight) {
      let seg1Size = stepWidth

      if (position + seg1Size > viewportWidth * 2 + viewportHeight) {
        seg1Size = viewportWidth * 2 + viewportHeight - position
        const seg2Size = stepWidth - seg1Size
        rect2 = Rect()
        rect2.setAttribute('fill', color)
        rect2.classList.add('position-left')
        rect2.setAttribute('y', (viewportHeight - seg2Size).toString())
        rect2.setAttribute('x', '0')
        rect2.setAttribute('height', seg2Size.toString())
        rect2.setAttribute('width', thickness.toString())
      }

      // bottom
      rect.classList.add('position-bottom')
      rect.setAttribute('y', (viewportHeight - thickness).toString())
      rect.setAttribute(
        'x',
        (viewportWidth - (position - viewportWidth - viewportHeight) - seg1Size).toString(),
      )
      rect.setAttribute('width', seg1Size.toString())
      rect.setAttribute('height', thickness.toString())
    } else if (position < viewportWidth * 2 + viewportHeight * 2) {
      let seg1Size = stepWidth

      if (position + seg1Size > viewportWidth * 2 + viewportHeight * 2) {
        seg1Size = viewportWidth * 2 + viewportHeight * 2 - position
        const seg2Size = stepWidth - seg1Size
        rect2 = Rect()
        rect2.setAttribute('fill', color)
        rect2.classList.add('position-top')
        rect2.setAttribute('y', '0')
        rect2.setAttribute('x', '0')
        rect2.setAttribute('height', thickness.toString())
        rect2.setAttribute('width', seg2Size.toString())
      }

      rect.classList.add('position-left')
      rect.setAttribute(
        'y',
        (viewportHeight - (position - viewportHeight - viewportWidth * 2) - stepWidth).toString(),
      )
      rect.setAttribute('x', '0')
      rect.setAttribute('width', thickness.toString())
      rect.setAttribute('height', stepWidth.toString())
    }

    remainingLength -= stepWidth
    const actualIndex = Math.min(elementIndex, totalElements - elementIndex)
    rect.style.setProperty('--index', actualIndex + '')
    element.appendChild(rect)
    if (rect2) {
      rect2.style.setProperty('--index', actualIndex + '')
      element.appendChild(rect2)
    }
    elementIndex++
  }
}

export function initPageSpark() {
  if (self !== top) return

  const init = () => {
    createSparkRoot(25)
  }

  if (document.readyState === 'complete') {
    init()
  } else {
    document.addEventListener('DOMContentLoaded', init)
  }
}
