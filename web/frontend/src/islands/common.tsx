import { OLIsland } from '@/lib/island'
import { ComponentType, render } from 'preact'

export type PreactIslandProps = { data?: unknown; rootElement: HTMLElement }

export class PreactIsland implements OLIsland {
  private _component: ComponentType<PreactIslandProps>

  constructor(component: ComponentType<PreactIslandProps>) {
    this._component = component
  }

  mount(el: HTMLElement, data: unknown): () => void {
    render(<this._component rootElement={el} data={data} />, el)

    return () => {
      render(null, el)
    }
  }
}
