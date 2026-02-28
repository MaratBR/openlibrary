import { AnimationEvent, AnimationWrapper, ModalAnimation } from '@/lib/animate'
import { TargetedMouseEvent } from 'preact'
import { createPortal, PropsWithChildren, useCallback, useRef, useState } from 'preact/compat'

export type ModalProps = PropsWithChildren<{
  open: boolean
  onClose?: () => void
}>

export default function Modal({ open, children, onClose }: ModalProps) {
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

  if (!shouldRender) return null

  return createPortal(
    <div ref={ref} class="modal" onClick={handleClick}>
      <AnimationWrapper
        onAnimation={handleAnimation}
        show={open}
        animation={ModalAnimation.default}
      >
        <div class="modal__content">{children}</div>
      </AnimationWrapper>
    </div>,
    document.body,
  )
}
