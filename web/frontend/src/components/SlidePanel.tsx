import clsx from 'clsx'
import { ReactNode, HTMLAttributes, MouseEventHandler, MouseEvent } from 'react'
import { createClassComponent } from './util'
import { useRef, useState } from 'react'
import { BinaryAnimation, useAnimation } from '@/lib/animate'
import { ForwardedRef, forwardRef } from 'react'

export type SlidePanelRootProps = {
  size?: '1/2' | '3/4' | '5/6'
} & HTMLAttributes<HTMLDivElement>

const SlidePanel_Root = forwardRef(
  ({ size = '1/2', ...props }: SlidePanelRootProps, ref: ForwardedRef<HTMLDivElement>) => {
    return (
      <div
        ref={ref}
        className={clsx('slide-panel shadow-2xl', {
          'w-1/2': size === '1/2',
          'w-3/4': size === '3/4',
          'w-5/6': size === '5/6',
        })}
        {...props}
      />
    )
  },
)

const SlidePanel_Content = createClassComponent('p-8')

function SlidePanel_Overlay({
  children,
  onClickOutside,
  className,
  ...props
}: {
  onClickOutside?: MouseEventHandler<HTMLElement>
} & HTMLAttributes<HTMLDivElement>) {
  const ref = useRef<HTMLDivElement | null>(null)

  function handleClick(event: MouseEvent<HTMLElement>) {
    if (ref.current === event.target) {
      onClickOutside?.(event)
    }
  }

  return (
    <div
      ref={ref}
      onClick={handleClick}
      className={clsx('fixed inset-0 z-10', className)}
      {...props}
    >
      {children}
    </div>
  )
}

class SlideOutAnimation extends BinaryAnimation {
  private _initialWidth: number = 0

  constructor(element: HTMLElement, duration: number) {
    super(element, duration, {
      onBeforeAnimation: (element, show) => {
        if (show) element.style.removeProperty('display')
        requestAnimationFrame(() => {
          this._initialWidth = element.offsetWidth
        })
      },
      onAfterAnimation(element, show) {
        if (!show) element.style.display = 'none'
      },
      onUpdate: (element, latest) => {
        if (this._initialWidth === 0) return
        element.style.transform = `translateX(${this._initialWidth * (1 - latest)}px)`
      },
    })
  }
}

function SlidePanel_Facade({
  open,
  onClose,
  children,
}: {
  open: boolean
  onClose: () => void
  children: ReactNode
}) {
  const [render, setRender] = useState(open)

  const { ref } = useAnimation({
    show: open,
    animation: (el) => new SlideOutAnimation(el, 300),
    onAnimation({ stage }) {
      switch (stage) {
        case 'exited':
          setRender(false)
          break
        case 'entering':
        case 'entered':
          setRender(true)
          break
      }
    },
  })

  return (
    <SlidePanel_Overlay
      style={render ? undefined : { display: 'none' }}
      onClickOutside={() => {
        onClose()
      }}
    >
      <SlidePanel_Root ref={ref}>{render ? children : null}</SlidePanel_Root>
    </SlidePanel_Overlay>
  )
}

export const SlidePanel = {
  Root: SlidePanel_Root,
  Content: SlidePanel_Content,
  Overlay: SlidePanel_Overlay,
  Facade: SlidePanel_Facade,
}
