import { ToastImplementation, ToastOptions } from './toast'
import { animate } from 'popmotion'

async function toaster(options: ToastOptions, parent: HTMLElement) {
  const element = document.createElement('div')
  element.classList.add('toast')
  const cleanup: (() => void)[] = []
  const renderCleanup = options.render(element, {
    close() {
      close()
    },
    addCleanup(cb) {
      cleanup.push(cb)
    },
  })
  parent.appendChild(element)

  const toasterRect = element.getBoundingClientRect()

  let animationProgress = 0

  function onUpdate(latest: number) {
    animationProgress = latest

    let scaleY = 1,
      translateX = 0

    const maxTranslateX = window.innerWidth - toasterRect.x
    const maxTranslateY = window.innerHeight - toasterRect.y

    if (latest < 50) {
      translateX = maxTranslateX
    } else {
      translateX = maxTranslateX * (1 - (latest - 50) / 50)
    }

    if (latest < 50) {
      scaleY = latest / 50
    }

    element.style.transform = `translateX(${translateX}px)`
    element.style.marginBottom = `calc(-${maxTranslateY * (1 - scaleY)}px + var(--spacing) * 2)`
  }

  const animationDuration = options.animationDuration ?? 300
  const appearAnimation = animate({
    from: 0,
    to: 100,
    onUpdate,
    duration: animationDuration,
  })

  let closed = false

  const close = () => {
    if (closed) return
    closed = true
    appearAnimation.stop()

    animate({
      from: animationProgress,
      to: 0,
      duration: animationDuration,
      onUpdate,
      onComplete() {
        cleanup.reverse().forEach((cb) => cb())
        renderCleanup()
        element.remove()
      },
    })
  }

  setTimeout(close, options.duration)
}

let toastContainer: HTMLElement | undefined

export const toastImplementation: ToastImplementation = (options: ToastOptions) => {
  if (!toastContainer) {
    toastContainer = document.createElement('div')
    toastContainer.className = 'toasts'
    document.body.appendChild(toastContainer)
  }

  toaster(options, toastContainer)
}
