export type ToastOptions = {
  render: (element: HTMLElement, controller: ToastController) => () => void
  duration: number
  meta?: Record<string, string>
  animationDuration?: number
}

export type ToastController = {
  close: () => void
}

export type ToastImplementation = (toast: ToastOptions) => void

export type ToastFunction = {
  (options: {
    title?: string
    text?: string
    type?: string
    duration?: number
    close?: boolean
  }): void
  impl: ToastImplementation
}

declare global {
  interface Window {
    toast: ToastFunction
  }

  var toast: ToastFunction
}

const toasts: ToastOptions[] = []

let toast: ToastFunction = function (
  this: ToastFunction,
  { text, title, type = 'info', duration = 500000, close = true },
) {
  this.impl({
    render(element, { close }) {
      const d = (cls: string) => {
        const div = document.createElement('div')
        div.className = cls
        return div
      }
      const $layout = d('toast-layout')
      const $icon = d('toast-layout__icon')
      const $content = d('toast-layout__content')
      if (title) {
        const $title = d('toast-layout__title')
        $title.innerText = title || ''
        $content.appendChild($title)
      }
      if (text) {
        const $text = d('toast-layout__text')
        $text.innerText = text || ''
        $content.appendChild($text)
      }

      switch (type) {
        case 'success':
          $icon.innerHTML = '<i class="fa-solid fa-circle-check"></i>'
          break
        case 'error':
          $icon.innerHTML = '<i class="fa-solid fa-triangle-exclamation"></i>'
          break
        case 'warning':
        case 'warn':
          type = 'warn'
          $icon.innerHTML = '<i class="fa-solid fa-circle-exclamation"></i>'
          break
        case 'info':
        default:
          type = 'info'
          $icon.innerHTML = '<i class="fa-solid fa-circle-info"></i>'
          break
      }
      $icon.setAttribute('data-color', type)

      $layout.append($icon, $content)

      if (close) {
        const $close = document.createElement('div')
        $close.innerHTML = '<i class="fa-solid fa-xmark"></i>'
        $close.classList = 'toast-layout__close'
        $layout.appendChild($close)
        let closed = false
        $close.addEventListener('click', () => {
          if (closed) {
            return
          }
          closed = true
          close()
        })
      }

      element.appendChild($layout)
      return () => {}
    },
    duration,
  })
} as ToastFunction

toast.impl = (options) => {
  toasts.push(options)
}

toast = toast.bind(toast)

window.toast = toast

import('./toastImplementation').then((m) => {
  window.toast.impl = m.toastImplementation
  toasts.forEach(m.toastImplementation)
})
