import { mount, Component, unmount } from "svelte";
import type { OLIsland } from "../common/ol-island";

export class SvelteIsland implements OLIsland {
  private _component: Component<{ data: unknown }>;

  constructor(component: Component<{ data: unknown }>) {
    this._component = component;
  }
  
  mount(el: HTMLElement, data: unknown): () => void {
    const mounted = mount(this._component, {
      target: el,
      props: { data }
    })
    return () => unmount(mounted)
  }
}

