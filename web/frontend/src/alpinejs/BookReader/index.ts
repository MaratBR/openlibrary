import { IMicrotask, queueMicrotasksWithBursts } from '@/lib/microtasks'
import Alpine from 'alpinejs'

const FONT_SIZES = [10, 12, 14, 16, 18, 22, 30, 36, 48, 60, 70, 80, 90, 99]
const FONT_SIZES_REVERSED = FONT_SIZES.slice().reverse()

Alpine.data('BookReader', () => ({
  settingsOpen: false,
  fontSize: 18,

  init() {
    const value = this.$el.getAttribute('data-font-size')
    if (value && !Number.isNaN(+value)) {
      const fontSize = Math.round(+value)
      if (fontSize >= 12 && fontSize <= 50) {
        this.fontSize = fontSize
      }
    }

    // const chapterContent = document.getElementById('ChapterContent')
    // if (chapterContent) {
    //   initScrollPosition(chapterContent)
    // }
  },

  changeFontSize(increase: boolean) {
    if (increase) {
      if (this.fontSize === FONT_SIZES_REVERSED[0]) return

      const next = FONT_SIZES.find((f) => f > this.fontSize)
      if (next) {
        this.fontSize = next
      } else {
        this.fontSize = 18
      }
      setFontSize(this.fontSize)
    } else {
      if (this.fontSize === FONT_SIZES[0]) return

      const next = FONT_SIZES_REVERSED.find((f) => f < this.fontSize)
      if (next) {
        this.fontSize = next
      } else {
        this.fontSize = 18
      }
      setFontSize(this.fontSize)
    }
  },

  toggleButton: {
    '@click'() {
      this.settingsOpen = !this.settingsOpen
    },
  },

  settings: {
    'x-show'() {
      return this.settingsOpen
    },
  },

  increaseFont: {
    '@click'() {
      this.changeFontSize(true)
    },
  },

  decreaseFont: {
    '@click'() {
      this.changeFontSize(false)
    },
  },

  closeButton: {
    '@click'() {
      this.settingsOpen = false
    },
  },
}))

function setFontSize(fontSize: number) {
  document.cookie = `ifs=${fontSize};path=/;max-age=31536000;`
  document.body.style.setProperty('--book-font-size', `${fontSize}px`)
}

// eslint-disable-next-line no-unused-vars
function initScrollPosition(chapterContent: HTMLElement) {
  const recordCurrentPosition = Alpine.throttle(() => {
    getCurrentPosition(chapterContent).then((pos) => {
      console.log(pos)
    })
  }, 200)

  window.addEventListener('scrollend', recordCurrentPosition)
}

export type CurrentPosition = {
  window: {
    height: number
    width: number
    scrollY: number
  }
  nearestElement: {
    path: string
    id: string | null
    top: number
  }
}

let elementAtOffset2: HTMLElement | null = null

async function getCurrentPosition(root: HTMLElement): Promise<CurrentPosition> {
  const elementAtOffset = await findNodeAtOffset(root, 64)

  if (elementAtOffset2) {
    elementAtOffset2.style.removeProperty('outline')
  }
  elementAtOffset2 = elementAtOffset
  if (elementAtOffset2) {
    elementAtOffset2.style.outline = '2px solid red'
  }

  const { innerHeight, innerWidth, scrollY } = window
  const nearestElement: CurrentPosition['nearestElement'] = {
    path: '',
    id: null,
    top: 0,
  }

  if (elementAtOffset) {
    const pathToElement = getRelativeNodePath(root, elementAtOffset)
    nearestElement.path = serializePathSteps(pathToElement)
    nearestElement.top = elementAtOffset.getBoundingClientRect().top

    if (elementAtOffset.id) {
      nearestElement.id = elementAtOffset.id
    }
  }

  return {
    window: {
      height: innerHeight,
      width: innerWidth,
      scrollY,
    },
    nearestElement,
  }
}

function findNodeAtOffset(root: HTMLElement, offset: number): Promise<HTMLElement | null> {
  if (offset < 1) {
    if (root.firstChild instanceof HTMLElement) {
      return Promise.resolve(root.firstChild)
    }
    return Promise.resolve(root)
  }

  const walker = document.createTreeWalker(root, NodeFilter.SHOW_ELEMENT, null)
  let current = walker.currentNode as HTMLElement | null

  const task: IMicrotask = {
    next() {
      if (!current) {
        return true
      }

      if (current === root) {
        current = walker.nextNode() as HTMLElement | null
        return false
      }

      const rect = current.getBoundingClientRect()
      if (rect.top >= offset) {
        return true
      }
      current = walker.nextNode() as HTMLElement | null
      return false
    },
  }

  return new Promise((resolve, reject) => {
    queueMicrotasksWithBursts(task, 10, 5000, {
      onError: reject,
      onTimeout() {
        reject('timeout')
      },
      onSuccess() {
        resolve(current)
      },
    })
  })
}

type PathStep = {
  from: Node
  to: Node
  idx: number
}

function getRelativeNodePath(parent: Node, child: Node): PathStep[] {
  const steps: PathStep[] = []

  let current: Node | null = child

  while (current !== null && current !== parent) {
    const parentNode: Node | null = current.parentNode
    if (parentNode === null) {
      throw new Error('Child is not a descendant of parent')
    }

    const idx = Array.prototype.indexOf.call(parentNode.childNodes, current)

    if (idx === -1) {
      throw new Error('Invariant violation: node not found in parent.childNodes')
    }

    steps.push({
      from: parentNode,
      to: current,
      idx,
    })

    current = parentNode
  }

  if (current !== parent) {
    throw new Error('Child is not a descendant of parent')
  }

  return steps.reverse()
}

function serializePathSteps(steps: PathStep[]) {
  const s: string[] = []

  for (const step of steps) {
    let ss = `${step.idx}:`

    if (step.to instanceof Element) {
      ss += `e:${JSON.stringify({ tag: step.to.tagName })}`
    } else if (step.to instanceof Text) {
      ss += `t:${JSON.stringify({ l: step.to.textContent.length })}`
    } else {
      ss += '?'
    }

    s.push(ss)
  }

  return s.join(',')
}
