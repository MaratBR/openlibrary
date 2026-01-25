import { SuggestionProps } from '@tiptap/suggestion'
import { SlashCommandDisplayAdapter, SlashCommandItem } from './Suggestions'
import { Subject, useSubject } from '@/common/rx'
import { RefObject, render } from 'preact'
import { computePosition, flip, offset } from '@floating-ui/react'
import { useEffect, useLayoutEffect, useRef } from 'preact/hooks'
import { wrapVirtualElement } from '@/lib/iframe'
import { getNextElement, getPrevElement } from '@/lib/html-elements'
import { Editor } from '@tiptap/core'

export interface SuggestionsElements {
  iframe: HTMLIFrameElement
  contentWrapper: HTMLElement
}

export class SuggestionsDisplay implements SlashCommandDisplayAdapter {
  private readonly root: HTMLElement
  private readonly elements: SuggestionsElements
  private readonly props = new Subject<
    SuggestionProps<SlashCommandItem, SlashCommandItem> | undefined
  >(undefined)
  private readonly getEditor: () => Editor
  private readonly focusCallbackRef: RefObject<(arrowUp: boolean) => void> = { current: null }

  constructor(elements: SuggestionsElements, getEditor: () => Editor) {
    this.elements = elements
    this.getEditor = getEditor
    this.root = document.createElement('div')
    document.body.appendChild(this.root)
  }

  show(props: SuggestionProps<SlashCommandItem, SlashCommandItem>): void {
    this.props.set(props)
    render(
      <Suggestions
        focusCallbackRef={this.focusCallbackRef}
        props={this.props}
        elements={this.elements}
        getEditor={this.getEditor}
      />,
      this.root,
    )
  }
  update(props: SuggestionProps<SlashCommandItem, SlashCommandItem>): void {
    this.props.set(props)
  }
  hide(): void {
    render(null, this.root)
  }
  focus(key: 'ArrowUp' | 'ArrowDown'): void {
    this.focusCallbackRef.current?.(key === 'ArrowUp')
  }
}

function Suggestions({
  props,
  elements,
  focusCallbackRef,
  getEditor,
}: {
  props: Subject<SuggestionProps<SlashCommandItem, SlashCommandItem> | undefined>
  elements: SuggestionsElements
  focusCallbackRef: RefObject<(arrowUp: boolean) => void>
  getEditor: () => Editor
}) {
  const propsValue = useSubject(props)
  const { items = [], command = () => {}, clientRect } = propsValue ?? {}
  const ref = useRef<HTMLDivElement | null>(null)

  useLayoutEffect(() => {
    const update = () => {
      const { current: modalElement } = ref
      if (!modalElement || !clientRect) {
        return
      }
      const domRect = clientRect()
      if (!domRect) return

      computeSuggestionsModalPosition(domRect, modalElement, elements.iframe).then((position) => {
        const wrappedElement = wrapVirtualElement(elements.iframe, elements.contentWrapper)
        const contentWrapperRect = wrappedElement.getBoundingClientRect()
        modalElement.style.transform = `translate(${contentWrapperRect.x + 8}px, ${position.y}px)`
      })
    }

    update()

    // TODO implement a propert recomputation only when necessary
    const interval = setInterval(update, 250)
    return () => clearInterval(interval)
  }, [clientRect, elements])

  const ulRef = useRef<HTMLUListElement | null>(null)

  useEffect(() => {
    focusCallbackRef.current = (arrowUp) => {
      const ul = ulRef.current
      if (!ul) return

      let li: HTMLLIElement | null

      if (arrowUp) {
        li = ul.querySelector('li:last-child') as HTMLLIElement | null
      } else {
        li = ul.querySelector('li:first-child') as HTMLLIElement | null
      }

      li?.focus()
    }
  }, [focusCallbackRef])

  function handleKeyDown(event: KeyboardEvent) {
    if (!(event.target instanceof HTMLLIElement)) {
      return
    }

    if (event.key === 'ArrowDown') {
      const nextLI = getNextElement(
        event.target,
        (el) => el instanceof HTMLLIElement,
      ) as HTMLElement
      nextLI.focus()
    } else if (event.key === 'ArrowUp') {
      const prevLI = getPrevElement(
        event.target,
        (el) => el instanceof HTMLLIElement,
      ) as HTMLElement
      prevLI.focus()
    } else if (event.key === 'Enter') {
      const commandName = event.target.dataset.command
      const item = items.find((x) => x.name === commandName)
      if (item) {
        command(item)
      }
    } else {
      const editor = getEditor()
      if (event.key.length === 1) {
        editor.chain().focus().insertContent(event.key).run()
      } else if (event.key === 'Backspace') {
        const { from, to } = editor.view.state.selection
        if (from === to) {
          editor
            .chain()
            .focus()
            .deleteRange({ from: from - 1, to: from })
            .run()
        } else {
          editor.chain().focus().deleteSelection().run()
        }
      }

      const li = event.target
      window.requestAnimationFrame(() => {
        li.focus()
      })
    }
  }

  return (
    <div class="be-suggestions-modal" ref={ref}>
      <ul ref={ulRef}>
        {items.map((item, index) => (
          <li
            key={item.name}
            tabIndex={0}
            role="button"
            onClick={(e) => {
              e.preventDefault()
              command(item)
            }}
            onKeyDown={handleKeyDown}
            data-command={item.name}
            className="be-suggestion-item"
          >
            <div class="be-suggestion-item__name">
              {item.icon && <div class="be-suggestion-item__icon">{item.icon}</div>}
              {item.name}
            </div>

            {item.description && (
              <div class="be-suggestion-item__description">{item.description}</div>
            )}
          </li>
        ))}
      </ul>
    </div>
  )
}

function computeSuggestionsModalPosition(
  clientRect: DOMRect,
  modalElement: HTMLElement,
  iframe: HTMLIFrameElement,
) {
  const virtualElement = wrapVirtualElement(iframe, {
    getBoundingClientRect() {
      return clientRect
    },
  })

  return computePosition(virtualElement, modalElement, {
    strategy: 'fixed',
    placement: 'bottom',
    middleware: [
      flip({
        fallbackPlacements: ['top', 'bottom'],
      }),
      offset({
        mainAxis: clientRect.height * 0,
      }),
    ],
  })
}
