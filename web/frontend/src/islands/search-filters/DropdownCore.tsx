import { HTMLAttributes, JSX, useCallback, useEffect, useRef, useState } from 'react'
import clsx from 'clsx'

export type DropdownProps = HTMLAttributes<HTMLDivElement> & {
  slotProps?: {
    input?: HTMLAttributes<HTMLInputElement>
    menu?: HTMLAttributes<HTMLDivElement>
  }
  slots?: {
    beforeInput?: JSX.Element
  }
}

export function DropdownCore({ slotProps = {}, slots = {}, ...props }: DropdownProps) {
  const [open, setOpen] = useState(false)
  const rootRef = useRef<HTMLDivElement | null>(null)

  useEffect(() => {
    if (!open) return

    const callback = (event: MouseEvent) => {
      window.requestAnimationFrame(() => {
        if (!rootRef.current) return

        if (
          (event.target instanceof Element &&
            !rootRef.current.contains(event.target) &&
            // only close if we clicked at element that currently exists in DOM
            // is this a good idea?
            document.body.contains(event.target)) ||
          event.target === document.documentElement
        ) {
          setOpen(false)
        }
      })
    }

    window.addEventListener('click', callback)

    return () => {
      window.removeEventListener('click', callback)
    }
  }, [open])

  const handleInputFocus = useCallback(() => setOpen(true), [])

  return (
    <div ref={rootRef} className="input-dropdown" data-open={open} {...props}>
      {slots.beforeInput}
      <input className="Dropdown-input" onFocus={handleInputFocus} {...slotProps.input} />

      <div
        aria-hidden={!open}
        onMouseDown={preventDefault}
        {...slotProps.menu}
        className={clsx('Dropdown-menu', slotProps.menu?.className, slotProps.menu?.className)}
      />
    </div>
  )
}

function preventDefault(e: { stopPropagation: () => void }) {
  e.stopPropagation()
}
