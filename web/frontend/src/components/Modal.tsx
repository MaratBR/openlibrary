import { createPortal, MouseEvent, PropsWithChildren, useCallback, useRef } from 'preact/compat'

export type ModalProps = PropsWithChildren<{
  open: boolean
  onClose?: () => void
}>

export default function Modal({ open, children, onClose }: ModalProps) {
  const ref = useRef<HTMLDivElement | null>(null)
  const handleClick = useCallback(
    (e: MouseEvent<HTMLDivElement>) => {
      if (!ref.current || e.target !== ref.current) return
      if (onClose) onClose()
    },
    [onClose],
  )

  if (!open) return null

  return createPortal(
    <div ref={ref} class="modal" onClick={handleClick}>
      <div class="modal__content">{children}</div>
    </div>,
    document.body,
  )
}
