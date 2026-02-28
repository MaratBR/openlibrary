import { animate } from 'popmotion'
import {
  cloneElement,
  forwardRef,
  JSX,
  useCallback,
  useEffect,
  useLayoutEffect,
  useRef,
} from 'preact/compat'
import { createEvent, Unsubscribe } from './event'

export type SetShowOptions = {
  duration?: number
  force?: boolean
  onComplete?: (cancelled: boolean) => void
}

export type AnimationEvent = {
  stage: 'entering' | 'entered' | 'exiting' | 'exited'
}

type AnimationCallback = (stage: AnimationEvent) => void

export type AnimationDefinition = {
  onUpdate: (element: HTMLElement, latest: number) => void
  onBeforeAnimation: (element: HTMLElement, show: boolean) => void
  onAfterAnimation: (Elementlement: HTMLElement, show: boolean) => void
}

export class BinaryAnimation {
  private readonly element: HTMLElement
  private progress = 0
  private _stop: (() => void) | undefined = undefined
  private readonly callbacks: AnimationDefinition
  private readonly duration: number
  private readonly _event = createEvent<AnimationEvent>()

  constructor(element: HTMLElement, duration: number, callbacks: AnimationDefinition) {
    this.element = element
    this.callbacks = callbacks
    this.duration = duration
  }

  setState(show: boolean, options?: SetShowOptions) {
    const { duration = this.duration, onComplete, force = false } = options ?? {}

    if (this.progress === (show ? 1 : 0) && !force) {
      return
    }

    if (duration <= 0) {
      // instantly transition to desired state
      this.callbacks.onBeforeAnimation(this.element, show)
      this.callbacks.onUpdate(this.element, show ? 1 : 0)
      this.callbacks.onAfterAnimation(this.element, show)

      this._event.fire({ stage: show ? 'entered' : 'exited' })

      return
    }

    if (this._stop) {
      this._stop()
      this._stop = undefined
    }

    let callbackFired = false

    this._event.fire({ stage: show ? 'entering' : 'exiting' })
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
        this.callbacks.onBeforeAnimation(this.element, show)
      },
      onComplete: () => {
        this._stop = undefined
        this.callbacks.onAfterAnimation(this.element, show)

        if (onComplete && !callbackFired) {
          callbackFired = true
          onComplete(false)
        }

        this._event.fire({ stage: show ? 'entered' : 'exited' })
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

  subscribe(callback: AnimationCallback): Unsubscribe {
    return this._event.subscribe(callback)
  }
}

export class ModalAnimation extends BinaryAnimation {
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
      onBeforeAnimation(element, show) {
        if (show) {
          element.style.removeProperty('display')
        }
      },
    })
  }

  static factory = (duration: number) => (element: HTMLElement) =>
    new ModalAnimation(element, duration)

  static default = this.factory(150)
}

export type AnimationProps = {
  animation: (element: HTMLElement) => BinaryAnimation
  children: JSX.Element
  show: boolean
  onAnimation?: AnimationCallback
}

// this feeels clunky somehow but will do for now
export const AnimationWrapper = forwardRef(
  ({ animation, children, show, onAnimation }: AnimationProps, ref) => {
    const animationInstanceRef = useRef<BinaryAnimation | null>(null)

    useLayoutEffect(() => {
      const { current: instance } = animationInstanceRef
      if (!instance) return
      instance.setState(show)
    }, [show])

    useEffect(() => {
      const { current: instance } = animationInstanceRef
      if (!instance || !onAnimation) return
      return instance.subscribe(onAnimation)
    }, [onAnimation])

    const initAnimation = useCallback((element: unknown) => {
      if (typeof ref === 'function') {
        ref(element)
      } else if (ref) {
        ref.current = element
      }

      if (!(element instanceof HTMLElement)) return

      const animationInstance = animation(element)
      animationInstance.setState(show, {
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
  },
)
