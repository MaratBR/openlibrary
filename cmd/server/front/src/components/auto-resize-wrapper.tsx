import { useEffect, useRef, useState, ReactNode } from 'react'

type AutoResizeWrapperProps = {
  children: ReactNode
  timingFunction: (v: number) => number
  duration: number
  className?: string
}

type Dimensions = {
  width: number
  height: number
}

type AnimationState = {
  startWidth: number
  startHeight: number
  targetWidth: number
  targetHeight: number
  startTime: number
}

const AutoResizeWrapper = ({
  children,
  timingFunction = (v: number) => v,
  duration = 300,
  className = '',
}: AutoResizeWrapperProps) => {
  const wrapperRef = useRef<HTMLDivElement>(null)
  const contentRef = useRef<HTMLDivElement>(null)
  const animationRef = useRef<number>()
  const animationState = useRef<AnimationState>()
  const [dimensions, setDimensions] = useState<Dimensions>({ width: 0, height: 0 })
  const [isFirstRender, setIsFirstRender] = useState(true)

  const animate = (currentTime: number) => {
    if (!animationState.current || !wrapperRef.current) return

    const { startTime, startWidth, startHeight, targetWidth, targetHeight } = animationState.current
    const elapsed = currentTime - startTime
    const progress = Math.min(elapsed / duration, 1)
    const eased = timingFunction(progress)

    const currentWidth = Math.ceil(startWidth + (targetWidth - startWidth) * eased)
    const currentHeight = Math.ceil(startHeight + (targetHeight - startHeight) * eased)

    wrapperRef.current.style.setProperty('--inner-width', `${currentWidth}px`)
    wrapperRef.current.style.setProperty('--inner-height', `${currentHeight}px`)

    if (progress < 1) {
      animationRef.current = requestAnimationFrame(animate)
    }
  }

  const startAnimation = (newDimensions: Dimensions) => {
    if (!wrapperRef.current) return

    const currentWidth = wrapperRef.current.offsetWidth
    const currentHeight = wrapperRef.current.offsetHeight

    if (animationRef.current) {
      cancelAnimationFrame(animationRef.current)
    }

    animationState.current = {
      startWidth: currentWidth,
      startHeight: currentHeight,
      targetWidth: newDimensions.width,
      targetHeight: newDimensions.height,
      startTime: performance.now(),
    }

    animationRef.current = requestAnimationFrame(animate)
  }

  useEffect(() => {
    const observer = new ResizeObserver((entries: ResizeObserverEntry[]) => {
      const entry = entries[0]
      if (entry) {
        const newDimensions = {
          width: entry.contentRect.width,
          height: entry.contentRect.height,
        }

        if (isFirstRender) {
          setDimensions(newDimensions)
          setIsFirstRender(false)
        } else {
          startAnimation(newDimensions)
          setDimensions(newDimensions)
        }
      }
    })

    if (contentRef.current) {
      observer.observe(contentRef.current)
    }

    return () => {
      observer.disconnect()
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current)
      }
    }
  }, [isFirstRender, duration])

  return (
    <div
      ref={wrapperRef}
      className={`overflow-hidden ${className}`}
      style={
        {
          '--inner-width': isFirstRender ? 'auto' : `${dimensions.width}px`,
          '--inner-height': isFirstRender ? 'auto' : `${dimensions.height}px`,
        } as React.CSSProperties
      }
    >
      <div ref={contentRef}>{children}</div>
    </div>
  )
}

export default AutoResizeWrapper
