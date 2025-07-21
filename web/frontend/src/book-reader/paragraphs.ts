import styles from './styles.module.css'

class Paragraphs {
  private readonly _root: HTMLElement
  private _active: HTMLParagraphElement | null = null

  constructor(root: HTMLElement) {
    this._root = root
    this._onClick = this._onClick.bind(this)
  }

  private _onClick(_e: MouseEvent, p: HTMLParagraphElement) {
    if (this._active) {
      this._active.classList.toggle(styles.active, false)
      this._active = null
    }

    this._active = p
    this._active.classList.toggle(styles.active, true)
  }

  init() {
    const paragraphs = this._root.querySelectorAll(':scope > p')
    paragraphs.forEach((el) => {
      const p = el as HTMLParagraphElement
      if (!containsOnlyTextElements(p)) return
      p.classList.add(styles.p)
      p.addEventListener('click', (e) => {
        this._onClick(e, p)
      })
    })
  }
}

export function initParagraphs(root: HTMLElement) {
  new Paragraphs(root).init()
}

function containsOnlyTextElements(node: Node): boolean {
  // List of allowed text-related elements
  const allowedTextElements = new Set([
    'b',
    'strong',
    'em',
    'i',
    'u',
    'span',
    'sub',
    'sup',
    'small',
    'mark',
    'del',
    'ins',
    'a',
    'p',
    'br',
    'code',
    'abbr',
    'cite',
    'q',
    's',
    'time',
  ])

  // Loop through all child nodes
  for (let i = 0; i < node.childNodes.length; i++) {
    const child = node.childNodes[i]

    // If the child is a text node, continue
    if (child.nodeType === Node.TEXT_NODE) {
      continue
    }

    // If the child is an element node, check if it's allowed
    if (child.nodeType === Node.ELEMENT_NODE) {
      const tagName = (child as Element).tagName.toLowerCase()

      // If the element is not in the allowed list, return false
      if (!allowedTextElements.has(tagName)) {
        return false
      }

      // Recursively check the child element
      if (!containsOnlyTextElements(child)) {
        return false
      }
    }
  }

  // If all child nodes are text or allowed elements, return true
  return true
}
