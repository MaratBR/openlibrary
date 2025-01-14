import { ISLANDS } from "./islands";

const loadedIslands: Record<string, OLIsland> = {};

customElements.define('ol-island', class extends HTMLElement {
  private _unmount?: () => void;

  constructor() {
    super();
    this._onDisposeRequested = this._onDisposeRequested.bind(this);
  }

  private _getData() {
    const data = this.getAttribute('data');
    if (data) {
      return JSON.parse(data);
    }

    const dataSource = this.getAttribute('data-source');
    if (dataSource) {
      const el = document.querySelector('#' + dataSource);
      if (el) {
        const dataFromSource = el.textContent;
        if (dataFromSource) {
          return JSON.parse(dataFromSource);
        }
      }
    }

    return null;
  }

  private async _loadAsync(name: string) {
    if (!ISLANDS[name]) {
      throw new Error('Island not found: ' + name);
    }
    await new Promise((resolve) => {
      setTimeout(resolve, 600)
    })
    const island = (await ISLANDS[name]()).default;
    loadedIslands[name] = island;
    this._start(island);
  }

  private _onDisposeRequested() {
    this.dispose()
  }

  private _start(island: OLIsland) {
    this.dispatchEvent(new CustomEvent('island-before-mount'))
    this._unmount = island.mount(this, this._getData())
    this.dispatchEvent(new CustomEvent('island-mount'))
  }

  connectedCallback() {
    this.addEventListener('island-request-dispose', this._onDisposeRequested);

    const name = this.getAttribute('name')
    if (!name) {
      throw new Error('Island element must have a name attribute')
    }
    const island = loadedIslands[name];
    if (!island) {
      const loader = document.createElement('div');
      loader.classList.add('island-loader');
      this.appendChild(loader);
      this._loadAsync(name);
      return;
    }
    
    window.requestAnimationFrame(() => {
      this._start(island);
    })
  }

  dispose() {
    if (!this._unmount) return;

    this.dispatchEvent(new CustomEvent('island-before-dispose'));
    this._unmount();
    this._unmount = undefined;
  }

  disconnectedCallback() {
    this.dispose();
  }
})

export interface OLIsland {
  mount(el: HTMLElement, data: unknown): () => void;
}