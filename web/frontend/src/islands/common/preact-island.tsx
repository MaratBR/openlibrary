import { OLIsland } from '@/lib/island'
import { ComponentChild, ComponentType, render } from 'preact'
import { PreactIslandSetup } from './setup'

export type PreactIslandProps = { data?: unknown; rootElement: HTMLElement }

abstract class PreactIslandBase implements OLIsland {
  private _component: ComponentType<PreactIslandProps>

  constructor(component: ComponentType<PreactIslandProps>) {
    this._component = component
  }

  abstract wrap(element: ComponentChild): ComponentChild

  mount(el: HTMLElement, data: unknown): () => void {
    render(this.wrap(<this._component rootElement={el} data={data} />), el)

    return () => {
      render(null, el)
    }
  }
}

export class PreactIsland extends PreactIslandBase {
  wrap(element: ComponentChild): ComponentChild {
    return <PreactIslandSetup>{element}</PreactIslandSetup>
  }
}
