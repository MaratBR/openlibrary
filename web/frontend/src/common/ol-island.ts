import { isAttrTrue } from './html-elements'
import './ol-island.scss'

export interface OLIsland {
  // eslint-disable-next-line no-unused-vars
  mount(el: HTMLElement, data: unknown): () => void
}

function validateIslandName(name: string) {
  if (!name) throw new Error('Island name is not specified')
  if (!/^[a-zA-Z0-9_-]+$/.test(name)) throw new Error('Island name is invalid')
}

class OLIslandsRegistry {
  private readonly _islands: Record<string, OLIsland> = {}

  public static instance: OLIslandsRegistry

  register(name: string, island: OLIsland) {
    validateIslandName(name)
    this._islands[name] = island
    document.dispatchEvent(new CustomEvent(`ol-island:register:${name}`))
  }

  get(name: string): OLIsland | null {
    validateIslandName(name)
    return this._islands[name] || null
  }

  async getAsync(name: string, timeout: number): Promise<OLIsland> {
    const island = this.get(name)
    if (island) {
      return island
    }

    const getPromise = new Promise<OLIsland>((resolve) => {
      let found = false

      const check = () => {
        if (!found) {
          const island = this.get(name)
          if (island) {
            found = true
            document.removeEventListener(`ol-island:register:${name}`, check)
            resolve(island)
          }
        }
      }
      document.addEventListener(`ol-island:register:${name}`, check)
      window.requestAnimationFrame(check)
    })

    const timeoutPromise = new Promise<undefined>((resolve) => setTimeout(resolve, timeout))
    const value = await Promise.race([timeoutPromise, getPromise])
    if (value === undefined) {
      throw new Error('cannot find island within timeout')
    }
    return value
  }
}

OLIslandsRegistry.instance = new OLIslandsRegistry()

abstract class OLIslandElementBase extends HTMLElement {
  private _unmount?: () => void
  private _isCreating: boolean = false

  get active() {
    return isAttrTrue(this.getAttribute('active'))
  }

  set active(value: boolean) {
    if (value) {
      this.setAttribute('active', 'true')
    } else {
      this.removeAttribute('active')
    }
  }

  constructor() {
    super()
    this._handleDestroyRequested = this._handleDestroyRequested.bind(this)
  }

  private _handleDestroyRequested() {
    this.active = false
  }

  private async _create() {
    if (this._isCreating || this._unmount) return
    this._isCreating = true
    this.showLoader()
    const island = await this.getIsland()
    window.requestAnimationFrame(() => {
      if (!this.active) {
        console.warn('[ol-island] by the time island was ready it was already inactive')
        return
      }

      this.childNodes.forEach((node) => {
        if (node instanceof HTMLTemplateElement) {
          return
        }
        node.remove()
      })
      this.dispatchEvent(new CustomEvent('island-before-mount'))
      console.debug('[ol-island] mount', island)
      this._unmount = island.mount(this, undefined)
      this.dispatchEvent(new CustomEvent('island-mount'))
      this._isCreating = false
    })
  }

  showLoader() {
    const template = this.querySelector('template[data-type=loader]')
    if (template instanceof HTMLTemplateElement) {
      const clone = template.content.cloneNode(true)
      this.appendChild(clone)
    }
  }

  connectedCallback() {
    this.addEventListener('island:request-destroy', this._handleDestroyRequested)

    if (this.active) {
      this._create()
    }
  }

  attributeChangedCallback(attribute: string, oldValue: string | null, newValue: string | null) {
    if (attribute === 'active') {
      const old = isAttrTrue(oldValue)
      const new_ = isAttrTrue(newValue)
      if (old !== new_) {
        if (new_) {
          this._create()
        } else {
          this._destroy()
        }

        this.onActiveChanged(new_)
      }
    }
  }

  static get observedAttributes() {
    return ['active']
  }

  private _destroy() {
    if (!this._unmount) return

    this.dispatchEvent(new CustomEvent('island:before-destroy'))
    this._unmount()
    this._unmount = undefined
    this.dispatchEvent(new CustomEvent('island:destroy'))
  }

  disconnectedCallback() {
    this._destroy()
  }

  abstract getIsland(): Promise<OLIsland>

  // eslint-disable-next-line no-unused-vars
  protected onActiveChanged(_active: boolean) {
    // no-op
  }
}

class OLIslandElement extends OLIslandElementBase {
  protected readonly islandName: string
  private static _addedScripts = new Set<string>()

  constructor() {
    super()
    this.islandName = this.getAttribute('name') ?? ''
  }

  async getIsland(): Promise<OLIsland> {
    const name = this.islandName
    if (!name) {
      throw new Error('Island name is not specified')
    }

    this.loadIslandIfNecessary()

    return OLIslandsRegistry.instance.getAsync(name, 30000)
  }
  private loadIslandIfNecessary() {
    const loadMethod = this.getAttribute('load-method')
    if (!loadMethod) {
      return
    }

    if (OLIslandElement._addedScripts.has(this.islandName)) {
      return
    }

    switch (loadMethod) {
      case 'script':
        this.loadFromScript()
        break
      case 'wait':
        break
      default:
        console.warn(`[ol-island] unknown load-method: ${loadMethod}`)
        return
    }
  }

  private loadFromScript() {
    const script = document.createElement('script')
    script.setAttribute('data-island-name', this.islandName)
    script.type = 'module'
    script.src = this.getAttribute('src') || `/_/assets/${this.islandName}.js`
    document.head.appendChild(script)
    OLIslandElement._addedScripts.add(this.islandName)
  }

  connectedCallback() {
    if (!this.islandName) {
      throw new Error('Island name is not specified')
    }

    super.connectedCallback()
  }
}

customElements.define('ol-island', OLIslandElement)

declare global {
  interface Window {
    OLIslandsRegistry: typeof OLIslandsRegistry
  }
}

window.OLIslandsRegistry = OLIslandsRegistry
