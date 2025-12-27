import { useContext, useEffect, useLayoutEffect, useRef, useState } from 'preact/hooks'
import ChapterSelector, { ChapterSelectorProps } from '../ChapterSelector'
import { useBookChaptersState } from './state'
import { createContext, createPortal, PropsWithChildren } from 'preact/compat'
import { computePosition } from '@floating-ui/react'
import { getErrorMessage } from '@/common/error'

export function ChapterSelectorPopup({
  element,
  onSelected,
  ActionDescriptionComponent,
  onClose,
}: {
  element: HTMLElement
  onClose: () => void
} & Pick<ChapterSelectorProps, 'onSelected' | 'ActionDescriptionComponent'>) {
  const chapters = useBookChaptersState((s) => s.chapters)

  const containerRef = useRef<HTMLDivElement | null>(null)

  useLayoutEffect(() => {
    const { current: container } = containerRef
    if (!container) return
    computePosition(element, container, { strategy: 'absolute', placement: 'left' })
      .then((position) => {
        container.style.left = `${position.x}px`
        container.style.top = `${Math.max(20, position.y)}px`
      })
      .catch((err) => {
        alert(`failed to position popup: ${getErrorMessage(err)}`)
      })
  }, [element])

  const onCloseRef = useRef(onClose)
  onCloseRef.current = onClose
  useEffect(() => {
    const onWindowClick = (event: Event) => {
      const { current: container } = containerRef
      if (!container || !(event.target instanceof Node)) return

      if (container.contains(event.target) || container === event.target) {
        return
      }

      onCloseRef.current?.()
    }

    window.addEventListener('click', onWindowClick)
  }, [element])

  return createPortal(
    <div class="card shadow-2xl rounded-2xl absolute" ref={containerRef}>
      <ChapterSelector
        chapters={chapters}
        onSelected={onSelected}
        ActionDescriptionComponent={ActionDescriptionComponent}
      />
    </div>,
    document.body,
  )
}

type OpenChapterSelectorPopupProps = Pick<
  ChapterSelectorProps,
  'ActionDescriptionComponent' | 'onSelected'
> & {
  element: HTMLElement
}

const ChapterSelectorPopupContext = createContext({
  open(_props: OpenChapterSelectorPopupProps) {},
  close() {},
})

export function useChapterSelectorPopup() {
  return useContext(ChapterSelectorPopupContext)
}

export function ChapterSelectorPopupProvider({ children }: PropsWithChildren) {
  const [props, setProps] = useState<OpenChapterSelectorPopupProps | null>(null)

  const ctx = useRef({
    open: setProps,
    close: () => setProps(null),
  })

  return (
    <ChapterSelectorPopupContext.Provider value={ctx.current}>
      {children}
      {props && <ChapterSelectorPopup onClose={() => setProps(null)} {...props} />}
    </ChapterSelectorPopupContext.Provider>
  )
}
