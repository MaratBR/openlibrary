function tween(
  from: number,
  to: number,
  duration: number,
  // eslint-disable-next-line no-unused-vars
  onUpdate: (value: number) => void,
  // eslint-disable-next-line no-unused-vars
  timingFunction: (v: number) => number,
) {
  const start = Date.now()
  const delta = to - from
  const step = () => {
    const elapsed = Date.now() - start
    const value = timingFunction(from + (delta * elapsed) / duration)
    onUpdate(value)
    if (elapsed < duration) {
      requestAnimationFrame(step)
    }
  }
  requestAnimationFrame(step)
}

function easeInOutQuint(x: number): number {
  return x < 0.5 ? 16 * x * x * x * x * x : 1 - Math.pow(-2 * x + 2, 5) / 2
}

class Toast {
  private readonly _root: HTMLElement

  constructor(element: HTMLElement) {
    this._root = document.createElement('div')
    this._root.appendChild(element)
  }

  mount(root: HTMLElement) {
    window.requestAnimationFrame(() => {
      this._root.style.opacity = '0'
      this._root.style.overflow = 'hidden'
      this._root.style.maxHeight = '0px'
      root.appendChild(this._root)

      const height = this._root.clientHeight
      const duration = 2000

      tween(
        0,
        1,
        duration,
        (v) => {
          this._root.style.maxHeight = `${height * v}px`
          this._root.style.opacity = `${v}`
          this._root.style.marginBottom = `${8 * v}px`
        },
        easeInOutQuint,
      )
    })
  }
}

class Toasts {
  private readonly _root: HTMLElement

  private readonly _toasts: Toast[] = []

  constructor(root: HTMLElement) {
    this._root = root
  }

  add(element: HTMLElement) {
    const toast = new Toast(element)
    toast.mount(this._root)
    this._toasts.push(toast)
  }
}

const toastsElement = document.createElement('section')
toastsElement.style.padding = '16px'
toastsElement.style.position = 'fixed'
toastsElement.style.left = '0px'
toastsElement.style.bottom = '0px'
toastsElement.style.zIndex = '10000'

document.body.appendChild(toastsElement)
const toasts = new Toasts(toastsElement)

setInterval(() => {
  toast('Hello world!')
}, 2000)

setTimeout(() => {
  window.location.reload()
}, 40000)

function toast(message: string) {
  const element = document.createElement('div')
  element.classList.add('ol-toast')
  element.classList.add('ol-toast--simple')
  element.textContent = message

  toasts.add(element)
}

declare global {
  interface Window {
    toast: typeof toast
  }
}

window.toast = toast
