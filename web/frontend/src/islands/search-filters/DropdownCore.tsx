import { useCallback, useEffect, useRef, useState } from 'preact/hooks'
import { JSX } from 'preact/jsx-runtime'
import clsx from 'clsx'

export type DropdownProps = JSX.HTMLAttributes<HTMLDivElement> & {
  slotProps?: {
    input?: JSX.HTMLAttributes<HTMLInputElement>
    menu?: JSX.HTMLAttributes<HTMLDivElement>
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
          event.target instanceof Element &&
          !rootRef.current.contains(event.target) &&
          // only close if we clicked at element that currently exists in DOM
          // is this a good idea?
          document.body.contains(event.target)
        ) {
          setOpen(false)
        }
      })
    }

    document.addEventListener('click', callback)

    return () => {
      document.removeEventListener('click', callback)
    }
  }, [open])

  const handleInputFocus = useCallback(() => setOpen(true), [])

  return (
    <div ref={rootRef} class="ol-dropdown" data-open={open} {...props}>
      {slots.beforeInput}
      <input class="ol-dropdown__input" onFocus={handleInputFocus} {...slotProps.input} />

      <div
        aria-hidden={!open}
        data-dropdown-content
        onMouseDown={preventDefault}
        {...slotProps.menu}
        class={clsx('ol-dropdown__menu', slotProps.menu?.class, slotProps.menu?.className)}
      />
    </div>
  )
}

function preventDefault(e: Event) {
  e.stopPropagation()
}
