import { animate } from 'popmotion'
import { cloneElement, forwardRef, JSX, useCallback, useLayoutEffect, useRef } from 'preact/compat'

export type SetShowOptions = {
  duration?: number
  force?: boolean
  onComplete?: (cancelled: boolean) => void
}

export interface AnimationController {
  setShow(show: boolean, options?: SetShowOptions): void
  dispose(): void
}

export type AnimationCallbacks = {
  onUpdate: (element: HTMLElement, latest: number) => void
  onBeforeAnimationm: (element: HTMLElement, show: boolean) => void
  onAfterAnimation: (Elementlement: HTMLElement, show: boolean) => void
}

export class ShowAnimation implements AnimationController {
  private readonly element: HTMLElement
  private progress = 0
  private _stop: (() => void) | undefined = undefined
  private readonly callbacks: AnimationCallbacks
  private readonly duration: number

  constructor(element: HTMLElement, duration: number, callbacks: AnimationCallbacks) {
    this.element = element
    this.callbacks = callbacks
    this.duration = duration
  }

  setShow(show: boolean, options?: SetShowOptions) {
    const { duration = this.duration, onComplete, force = false } = options ?? {}

    if (this.progress === (show ? 1 : 0) && !force) {
      return
    }

    if (duration <= 0) {
      // instantly transition to desired state
      this.callbacks.onBeforeAnimationm(this.element, show)
      this.callbacks.onUpdate(this.element, show ? 1 : 0)
      this.callbacks.onAfterAnimation(this.element, show)

      return
    }

    if (this._stop) {
      this._stop()
      this._stop = undefined
    }

    let callbackFired = false

    this._stop = animate({
      duration: duration * (show ? 1 - this.progress : this.progress),
      from: this.progress,
      to: show ? 1 : 0,
      onStop() {
        if (onComplete && !callbackFired) {
          callbackFired = true
          onComplete(true)
        }
      },
      onPlay: () => {
        this.callbacks.onBeforeAnimationm(this.element, show)
      },
      onComplete: () => {
        this._stop = undefined
        this.callbacks.onAfterAnimation(this.element, show)

        if (onComplete && !callbackFired) {
          callbackFired = true
          onComplete(false)
        }
      },
      onUpdate: (latest) => {
        this.progress = latest
        this.callbacks.onUpdate(this.element, latest)
      },
    }).stop
  }

  stop() {
    if (this._stop) {
      this._stop()
      this._stop = undefined
    }
  }

  dispose() {
    this.stop()
  }
}

export class ModalAnimation extends ShowAnimation {
  constructor(element: HTMLElement, duration: number) {
    super(element, duration, {
      onUpdate(element, latest) {
        element.style.opacity = `${latest}`
        element.style.transform = `scale(${0.95 + 0.05 * latest})`
      },
      onAfterAnimation(element, show) {
        if (!show) {
          element.style.display = 'none'
        }
      },
      onBeforeAnimationm(element, show) {
        if (show) {
          element.style.removeProperty('display')
        }
      },
    })
  }

  static factory = (duration: number) => (element: HTMLElement) =>
    new ModalAnimation(element, duration)
}

export type AnimationProps = {
  animation: (element: HTMLElement) => AnimationController
  children: JSX.Element
  show: boolean
}

// this feeels clunky somehow but will do for now
export const AnimationWrapper = forwardRef(({ animation, children, show }: AnimationProps, ref) => {
  const animationInstanceRef = useRef<AnimationController | null>(null)

  useLayoutEffect(() => {
    const { current: instance } = animationInstanceRef
    if (!instance) return
    instance.setShow(show)
  }, [show])

  const initAnimation = useCallback((element: unknown) => {
    if (typeof ref === 'function') {
      ref(element)
    } else if (ref) {
      ref.current = element
    }

    if (!(element instanceof HTMLElement)) return

    const animationInstance = animation(element)
    animationInstance.setShow(show, {
      duration: 0,
      force: true,
    })
    animationInstanceRef.current = animationInstance
    return () => {
      animationInstance.dispose()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return cloneElement(children, {
    ref: initAnimation,
  })
})
