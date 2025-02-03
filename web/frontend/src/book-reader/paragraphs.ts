import styles from './styles.module.css'

class Paragraphs {
  private readonly _root: HTMLElement
  private _active: HTMLParagraphElement | null = null

  constructor(root: HTMLElement) {
    this._root = root
    this._onClick = this._onClick.bind(this)
  }

  private _onClick(e: MouseEvent) {
    if (!(e.target instanceof HTMLParagraphElement)) {
      return
    }

    if (this._active) {
      this._active.classList.toggle(styles.active, false)
      this._active = null
    }

    this._active = e.target
    this._active.classList.toggle(styles.active, true)
  }

  init() {
    const paragraphs = this._root.querySelectorAll('p')
    paragraphs.forEach((p) => {
      if (!p.textContent?.trim()) return
      p.classList.add(styles.p)
      p.addEventListener('click', this._onClick)
    })
  }
}

export function initParagraphs(root: HTMLElement) {
  new Paragraphs(root).init()
}
