import React from 'react'
import { SlotContent, SlotContext } from './slot'
import { InjectCSS } from './inject-css'

export type EncapsulatedCSSProps = React.PropsWithChildren<{ css?: string | null }>

export class EncapsulatedCSS extends React.Component<
  EncapsulatedCSSProps,
  { root: { shadowRoot: ShadowRoot; slot: HTMLElement } | null }
> {
  private static idCounter = 0

  private hostElement: HTMLElement | null = null
  private ready: boolean = false
  private id: string = `EncapsulatedCSS-${EncapsulatedCSS.idCounter++}`

  constructor(props: Readonly<EncapsulatedCSSProps>) {
    super(props)
    this.state = { root: null }
  }

  componentDidMount(): void {
    if (this.ready) return
    this.ready = true
    if (!this.hostElement) throw new Error('this.hostElement is null')
    const shadowRoot = this.hostElement.attachShadow({ mode: 'open' })
    const slot = document.createElement('div')
    slot.id = 'slot'
    slot.style.display = 'contents'
    shadowRoot.appendChild(slot)
    this.setState({ root: { shadowRoot, slot } })
  }

  render(): React.ReactNode {
    const { children, css } = this.props
    const { root } = this.state

    return (
      <>
        <div id={this.id} ref={(el) => (this.hostElement = el)} />
        <SlotContext.Provider value={root?.slot ?? null}>
          <SlotContent>
            {children}
            {root && css && <InjectCSS css={css} document={root.shadowRoot} />}
          </SlotContent>
        </SlotContext.Provider>
      </>
    )
  }
}
