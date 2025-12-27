import { useEffect, useRef } from 'preact/hooks'
import { ComponentChild, RefObject } from 'preact'
import { createPortal, HTMLAttributes } from 'preact/compat'
import { useFloating, Placement } from '@floating-ui/react'

export type PopperProps = {
  anchorEl?: HTMLElement | RefObject<HTMLElement | null> | null
  children?: ComponentChild
  open?: boolean
  onClose?: () => void
  placement?: Placement
} & Omit<HTMLAttributes<HTMLDivElement>, 'onClose'>

export default function Popper({
  anchorEl,
  children,
  placement,
  open = true,
  style,
  onClose,
  ...props
}: PopperProps) {
  const { refs, floatingStyles, elements } = useFloating({
    elements: {
      reference: anchorEl instanceof HTMLElement ? anchorEl : anchorEl?.current,
    },
    strategy: 'fixed',
    placement,
  })

  const onCloseRef = useRef(onClose)
  onCloseRef.current = onClose

  useEffect(() => {
    if (!open) return

    const onWindowClick = (event: MouseEvent) => {
      if (!(event.target instanceof Node)) return
      if (
        event.target === elements.floating ||
        event.target === anchorEl ||
        (anchorEl instanceof HTMLElement
          ? anchorEl.contains(event.target)
          : anchorEl?.current?.contains(event.target)) ||
        (elements.floating && elements.floating.contains(event.target))
      ) {
        return
      }

      if (onCloseRef.current) {
        onCloseRef.current()
      }
    }

    requestAnimationFrame(() => {
      window.addEventListener('click', onWindowClick)
    })

    return () => {
      window.removeEventListener('click', onWindowClick)
    }
  }, [elements.floating, anchorEl, open])

  if (!open) return null

  return createPortal(
    <div class="contents" style={{ visibility: open ? 'visible' : 'hidden' }}>
      <div ref={refs.setFloating} {...props} style={floatingStyles} data-open={open}>
        {children}
      </div>
    </div>,
    document.body,
  )
}
