import { animate, easeInOut } from 'popmotion'
import React, { useLayoutEffect, useRef } from 'react'

export function EditButtonMenuExpandableSection({
  children,
  speedPixelsPerSecond = 400,
}: React.PropsWithChildren<{ speedPixelsPerSecond?: number }>) {
  const rootRef = useRef<HTMLDivElement | null>(null)
  const innerRef = useRef<HTMLDivElement | null>(null)

  const propsRef = useRef({ speedPixelsPerSecond })
  propsRef.current.speedPixelsPerSecond = speedPixelsPerSecond

  useLayoutEffect(() => {
    const $root = rootRef.current
    if (!$root) return

    $root.style.overflow = 'hidden'
    $root.style.height = '0px'

    const $inner = innerRef.current
    if (!$inner) return

    let animation: { stop: () => void } | null = null

    const observer = new ResizeObserver(() => {
      if (animation) animation.stop()

      const { clientHeight: rootHeight } = $root
      const { clientHeight: innerHeight } = $inner

      if (Math.round(innerHeight) === Math.round(rootHeight)) {
        $root.style.height = `${innerHeight}px`
        return
      }

      animation = animate({
        from: rootHeight,
        to: innerHeight,
        duration:
          (Math.abs(rootHeight - innerHeight) / propsRef.current.speedPixelsPerSecond) * 1000,
        ease: easeInOut,
        onUpdate(latest) {
          $root.style.height = `${latest}px`
        },
      })
    })

    observer.observe($inner)

    return () => observer.disconnect()
  }, [])

  return (
    <div ref={rootRef} className="be-bubble-menu__section">
      <div className="w-full min-w-0 max-w-full" ref={innerRef}>
        {children}
      </div>
    </div>
  )
}
