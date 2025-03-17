import { createPortal, PropsWithChildren } from 'preact/compat'

export type ModalProps = PropsWithChildren<{
  open: boolean
}>

export default function Modal({ open, children }: ModalProps) {
  if (!open) return null

  return createPortal(
    <div class="ol-modal">
      <div class="ol-modal__content">{children}</div>
    </div>,
    document.body,
  )
}
