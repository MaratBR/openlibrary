export type ToastOptions = {
  render: (element: HTMLElement, controller: ToastController) => () => void
  duration: number
  meta?: Record<string, string>
  animationDuration?: number
}

export type ToastController = {
  close: () => void
  addCleanup: (cb: () => void) => void
}

export type ToastImplementation = (toast: ToastOptions) => void

export type ToastFunction = {
  (options: {
    title?: string
    text?: string
    type?: string
    duration?: number
    close?: boolean
    customContent?: (element: HTMLElement) => undefined | (() => void)
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

window.toast = function (
  this: ToastFunction,
  { text, title, type = 'info', duration = 5000, close = true, customContent },
) {
  window.toast.impl({
    render(element, { close: closeCb, addCleanup }) {
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

      if (customContent) {
        const $customContent = d('toast-layout__customContent')
        $content.appendChild($customContent)
        const customContentCleanup = customContent($customContent)
        if (customContentCleanup) addCleanup(customContentCleanup)
      }

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
          closeCb()
        })
      }

      element.appendChild($layout)
      return () => {}
    },
    duration,
  })
} as ToastFunction

window.toast.impl = (options) => {
  toasts.push(options)
}
window.toast = window.toast.bind(window.toast)

import('./toastImplementation').then((m) => {
  window.toast.impl = m.toastImplementation
  toasts.forEach(m.toastImplementation)
})
