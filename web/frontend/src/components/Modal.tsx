import { AnimationEvent, ModalAnimation, useAnimation } from '@/lib/animate'
import clsx from 'clsx'
import { HTMLAttributes, MouseEvent } from 'react'
import { PropsWithChildren, useCallback, useRef, useState } from 'react'
import { createPortal } from 'react-dom';

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
    (e: MouseEvent<HTMLDivElement>) => {
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
    <div ref={ref} className="modal" onClick={handleClick}>
      <div
        ref={animationRef}
        {...slotProps.content}
        className={clsx('modal__content', slotProps.content?.className)}
      >
        {children}
      </div>
    </div>,
    document.body,
  )
}
