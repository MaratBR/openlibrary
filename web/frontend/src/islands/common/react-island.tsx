import { OLIsland, OLIslandMounted } from '@/lib/island'
import { ReactIslandSetup } from './setup'
import { ComponentType, ReactElement, useState } from 'react'
import { createRoot } from 'react-dom/client';

export type ReactIslandProps = { data?: unknown; rootElement: HTMLElement }

abstract class ReactIslandBase implements OLIsland {
  private _component: ComponentType<ReactIslandProps>

  constructor(component: ComponentType<ReactIslandProps>) {
    this._component = component
  }

  abstract wrap(element: ReactElement): ReactElement

  mount(el: HTMLElement, data: unknown): OLIslandMounted {
    let setData: (data: unknown) => void = () => {}

    const StateProxy = () => {
      const [innerData, setInnerData] = useState(data)

      setData = setInnerData

      return <this._component rootElement={el} data={innerData} />
    }

    const root = createRoot(el)
    root.render(
      this.wrap(<StateProxy />)
    )

    return {
      dispose() {
        root.unmount()
      },
      setData,
    }
  }
}

export class ReactIsland extends ReactIslandBase {
  wrap(element: ReactElement): ReactElement {
    return <ReactIslandSetup>{element}</ReactIslandSetup>
  }
}
