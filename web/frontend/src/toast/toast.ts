import anime from 'animejs'

class Toast {
  private readonly _root: HTMLElement
  private _activeAnimation?: anime.AnimeInstance
  private _unmountTimeout?: number
  private readonly _options: ToastOptions

  constructor(element: HTMLElement, options: ToastOptions) {
    this._root = document.createElement('div')
    this._root.appendChild(element)
    this._options = options
  }

  mount(root: HTMLElement) {
    const closeButton = this._root.querySelector('[data-toast-close]')

    if (closeButton) {
      closeButton.addEventListener('click', () => {
        this.unmount()
      })
    }

    window.requestAnimationFrame(() => {
      this._root.style.opacity = '0'
      this._root.style.position = 'fixed'

      window.requestAnimationFrame(() => {
        root.appendChild(this._root)
        const height = this._root.clientHeight
        this._root.style.height = '0px'
        this._root.style.removeProperty('position')

        this._activeAnimation = anime({
          targets: this._root,
          duration: this._options.animationDuration,
          height,
          opacity: 1,
          marginBottom: 8,
        })
      })
    })

    if (this._options.duration >= 0) {
      this._unmountTimeout = window.setTimeout(() => {
        this.unmount()
      }, this._options.duration)
    }
  }

  unmount() {
    if (this._unmountTimeout !== undefined) {
      clearTimeout(this._unmountTimeout)
      this._unmountTimeout = undefined
    }

    if (this._activeAnimation) {
      this._activeAnimation.pause()
      this._activeAnimation = undefined
    }

    this._activeAnimation = anime({
      targets: this._root,
      opacity: 0,
      duration: this._options.animationDuration,
      height: 0,
      marginBottom: 0,
      elasticity: 0,
      complete: () => {
        this._root.remove()
        this._activeAnimation = undefined
      },
    })
  }
}

export type ToastOptions = {
  duration: number
  animationDuration: number
}

class Toasts {
  private readonly _root: HTMLElement

  private readonly _toasts: Toast[] = []

  constructor(root: HTMLElement) {
    this._root = root
  }

  add(element: HTMLElement, options: ToastOptions) {
    const toast = new Toast(element, options)
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

export type ToastType = 'success' | 'error' | 'info'

function createBasicToast(title: Node, content: Node, type?: ToastType): HTMLElement {
  const div = document.createElement('div')
  div.classList.add('ol-toast')
  div.classList.add('ol-toast--simple')

  {
    const close = document.createElement('div')
    close.role = 'button'
    close.classList.add('ol-toast__close')
    close.setAttribute('data-toast-close', 'true')
    close.innerHTML = '<i style="font-size:22px" class="fa-solid fa-xmark" />'
    div.appendChild(close)
  }

  if (type) {
    const iconWrapper = document.createElement('div')
    iconWrapper.classList.add('ol-toast__icon')

    const iconElement = document.createElement('i')
    iconElement.classList.add('fa-solid')

    iconWrapper.dataset.type = type
    switch (type) {
      case 'info':
        iconElement.classList.add('fa-info')
        break
      case 'success':
        iconElement.classList.add('fa-circle-check')
        break
      case 'error':
        iconElement.classList.add('fa-triangle-exclamation')
        break
    }

    iconWrapper.appendChild(iconElement)
    div.appendChild(iconWrapper)
    div.classList.add('ol-toast--withIcon')
  }

  {
    const mainSection = document.createElement('div')
    mainSection.classList.add('ol-toast__main')

    const contentWrapper = document.createElement('div')
    contentWrapper.classList.add('ol-toast__content')
    contentWrapper.appendChild(content)
    const titleElement = document.createElement('div')
    titleElement.classList.add('ol-toast__title')
    titleElement.appendChild(title)

    mainSection.appendChild(titleElement)
    mainSection.appendChild(contentWrapper)
    div.appendChild(mainSection)
  }

  return div
}

function toast({
  options,
  title,
  text,
  type,
}: {
  options?: Partial<ToastOptions>
  title?: string
  text?: string
  type?: ToastType
}) {
  const toastElement = createBasicToast(
    document.createTextNode(title ?? ''),
    document.createTextNode(text ?? ''),
    type,
  )

  toasts.add(toastElement, {
    duration: 5000,
    animationDuration: 500,
    ...options,
  })
}

declare global {
  interface Window {
    toast: typeof toast
  }
}

window.toast = (..._args) => {}

// const btn = document.createElement('button')
// btn.textContent = 'Toast'
// btn.setAttribute(
//   'style',
//   'position: fixed; right: 32px; bottom: 10px; z-index:10000; background-color: #222; color: #fff; padding: 10px; border-radius: 4px; cursor: pointer;',
// )
// btn.addEventListener('click', () => {
//   toast({ title: 'Title', text: 'Lorem ipsum dolor sit amet', options: { duration: 5000 } })
// })
// document.addEventListener('DOMContentLoaded', () => {
//   document.body.appendChild(btn)
// })
