import { OLIsland, OLIslandMounted } from '@/lib/island'
import { ComponentChild, ComponentType, render } from 'preact'
import { PreactIslandSetup } from './setup'
import { useState } from 'preact/hooks'

export type PreactIslandProps = { data?: unknown; rootElement: HTMLElement }

abstract class PreactIslandBase implements OLIsland {
  private _component: ComponentType<PreactIslandProps>

  constructor(component: ComponentType<PreactIslandProps>) {
    this._component = component
  }

  abstract wrap(element: ComponentChild): ComponentChild

  mount(el: HTMLElement, data: unknown): OLIslandMounted {
    let setData: (data: unknown) => void = () => {}

    const StateProxy = () => {
      const [innerData, setInnerData] = useState(data)

      setData = setInnerData

      return <this._component rootElement={el} data={innerData} />
    }

    render(this.wrap(<StateProxy />), el)

    return {
      dispose() {
        render(null, el)
      },
      setData,
    }
  }
}

export class PreactIsland extends PreactIslandBase {
  wrap(element: ComponentChild): ComponentChild {
    return <PreactIslandSetup>{element}</PreactIslandSetup>
  }
}
