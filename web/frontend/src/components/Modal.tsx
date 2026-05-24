import { AnimationEvent, ModalAnimation, useAnimation } from '@/lib/animate'
import clsx from 'clsx'
import { HTMLAttributes, TargetedMouseEvent } from 'preact'
import { createPortal, PropsWithChildren, useCallback, useRef, useState } from 'preact/compat'

export type ModalProps = PropsWithChildren<{
  open: boolean
  onClose?: () => void
  slotProps?: {
    content?: HTMLAttributes<HTMLDivElement>
  }
}>

export default function Modal({ open, children, onClose, slotProps = {} }: ModalProps) {
  const ref = useRef<HTMLDivElement | null>(null)
  const handleClick = useCallback(
    (e: TargetedMouseEvent<HTMLDivElement>) => {
      if (!ref.current || e.target !== ref.current) return
      if (onClose) onClose()
    },
    [onClose],
  )

  const handleAnimation = useCallback((event: AnimationEvent) => {
    setAnimationInProgress(event.stage !== 'exited')
  }, [])

  const [animationInProgress, setAnimationInProgress] = useState(false)

  const shouldRender = open || animationInProgress

  const { ref: animationRef } = useAnimation({
    show: open,
    animation: ModalAnimation.default,
    onAnimation: handleAnimation,
  })

  if (!shouldRender) return null

  return createPortal(
    <div ref={ref} class="modal" onClick={handleClick}>
      <div
        ref={animationRef}
        {...slotProps.content}
        class={clsx('modal__content', slotProps.content?.class)}
      >
        {children}
      </div>
    </div>,
    document.body,
  )
}
