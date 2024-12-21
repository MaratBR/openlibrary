import React, { useRef } from 'react'
import { useForkRef } from 'rooks'

export type VisibilityTriggerProps = React.HTMLAttributes<HTMLDivElement> & {
  onVisibilityChange?: (visible: boolean) => void
  onAppear?: () => void
}

const VisibilityTrigger = React.forwardRef(
  ({ ...props }: VisibilityTriggerProps, fRef: React.ForwardedRef<HTMLDivElement>) => {
    const lRef = React.useRef<HTMLDivElement | null>(null)
    const ref = useForkRef(fRef, lRef)
    const visibleRef = useRef<boolean>()
    const appearedRef = useRef(false)

    React.useEffect(() => {
      const observer = new IntersectionObserver((entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting !== visibleRef.current) {
            visibleRef.current = entry.isIntersecting
            props.onVisibilityChange?.(entry.isIntersecting)

            if (!appearedRef.current) {
              appearedRef.current = true
              props.onAppear?.()
            }
          }
        })
      })

      if (lRef.current) {
        observer.observe(lRef.current)
      }

      return () => observer.disconnect()
    }, [lRef, props])

    return <div ref={ref} {...props} />
  },
)

export default VisibilityTrigger
