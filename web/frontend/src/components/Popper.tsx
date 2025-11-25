import { useEffect, useLayoutEffect, useRef } from 'preact/hooks'
import { createPopper, Options, Instance } from '@popperjs/core'
import { ComponentChild, RefObject } from 'preact'
import { HTMLAttributes } from 'preact/compat'

export type PopperProps = {
  anchorEl?: HTMLElement | RefObject<HTMLElement | null> | null
  children?: ComponentChild
  open?: boolean
} & HTMLAttributes<HTMLDivElement> &
  Partial<Options>

export default function Popper({
  anchorEl,
  children,
  onFirstUpdate,
  placement = 'auto',
  modifiers = [],
  strategy = 'fixed',
  open = true,
  style,
  ...props
}: PopperProps) {
  const ref = useRef<HTMLDivElement | null>(null)
  const instanceRef = useRef<Instance | null>(null)
  const optionsRef = useRef({
    modifiers,
    onFirstUpdate,
    placement,
    strategy,
  })
  optionsRef.current = {
    modifiers,
    onFirstUpdate,
    placement,
    strategy,
  }

  useLayoutEffect(() => {
    if (!ref.current) {
      return
    }

    let el: HTMLElement | null

    if (anchorEl) {
      if (anchorEl instanceof HTMLElement) {
        el = anchorEl
      } else {
        el = anchorEl.current
      }
    } else {
      el = null
    }

    if (!el) {
      return
    }

    const instance = createPopper(el, ref.current, optionsRef.current)
    instanceRef.current = instance

    return () => {
      instanceRef.current = null
      instance.destroy()
    }
  }, [anchorEl])

  useLayoutEffect(() => {
    const { current } = instanceRef
    if (current) {
      current.update()
    }
  }, [open])

  const firstRender = useRef(true)

  useEffect(() => {
    if (firstRender.current) {
      firstRender.current = true
      return
    }

    const { current } = instanceRef
    if (current) {
      current.setOptions({ modifiers, onFirstUpdate, strategy, placement })
    }
  }, [modifiers, onFirstUpdate, strategy, placement])

  return (
    <div style={{ display: open ? 'contents' : 'none' }}>
      <div ref={ref} {...props} data-open={open}>
        {children}
      </div>
    </div>
  )
}
