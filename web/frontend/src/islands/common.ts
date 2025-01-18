import { mount, Component, unmount } from "svelte";
import type { OLIsland } from "../common/ol-island";

export class SvelteIsland implements OLIsland {
  private _component: Component<{}>;

  constructor(component: Component<{}>) {
    this._component = component;
  }
  
  mount(el: HTMLElement, data: unknown): () => void {
    const mounted = mount(this._component, {
      target: el,
      props: {}
    })
    return () => unmount(mounted)
  }
}

