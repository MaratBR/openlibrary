import { isAttrTrue } from '../html-elements'
import { OLIsland, OLIslandMounted } from './island'

class OLIslandElement extends HTMLElement {
  private _mountedIsland?: OLIslandMounted
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

  get name() {
    const v = this.getAttribute('name')
    if (!v) throw new Error('Island name is not specified')
    return v
  }

  constructor() {
    super()
    this._handleDestroyRequested = this._handleDestroyRequested.bind(this)

    if (this.getAttribute('preload') === 'true' || this.getAttribute('preload') === '') {
      this._preload()
    }
  }

  //#region html element implementation

  connectedCallback() {
    this.addEventListener('island:request-destroy', this._handleDestroyRequested)

    if (this.active) {
      this._create()
    }
  }

  attributeChangedCallback(attribute: string, oldValue: string | null, newValue: string | null) {
    switch (attribute) {
      case 'active':
        {
          const old = isAttrTrue(oldValue)
          const new_ = isAttrTrue(newValue)
          if (old !== new_) {
            if (new_) {
              this._create()
            } else {
              this._destroy()
            }
          }
        }
        break
      case 'data':
        this._updateComponentData()
        break
    }
  }

  static get observedAttributes() {
    return ['active', 'data']
  }

  //#endregion

  private _handleDestroyRequested() {
    this.active = false
  }

  private async _create() {
    if (this._isCreating || this._mountedIsland) return
    this._isCreating = true

    this._showLoader()
    let island: OLIsland
    try {
      ;[island] = await Promise.all([
        this._fetchIsland(),
        new Promise((resolve) => {
          setTimeout(resolve, 300)
        }),
      ])
    } catch (e: unknown) {
      console.error('[ol-island] getIsland call failed, error:', e)
      this._isCreating = false
      return
    }

    window.requestAnimationFrame(() => {
      if (!this.active) {
        console.warn('[ol-island] by the time island was ready it was already inactive')
        return
      }

      try {
        this._clearContent()
        this.dispatchEvent(new CustomEvent('island:before-mount'))
        this._mountedIsland = island.mount(this, this._getData())
        console.debug('[ol-island] mount', island)
        this.dispatchEvent(new CustomEvent('island:mount'))
      } finally {
        this._isCreating = false
      }
    })
  }

  private async _fetchIsland(): Promise<OLIsland> {
    const src = this.getSrc()

    if (src) {
      const module: unknown = await import(src)
      const name = this.getAttribute('name') || 'default'
      const island = getIslandFromModule(module, name)
      if (island) return island
      console.error(`[ol-island] module ${src}`, module)
      throw new Error(`Island ${name} is not found in ${src}`)
    }
    throw new Error('Island src is not specified')
  }

  private getSrc() {
    const src = this.getAttribute('src')
    if (!src) return null

    if (/(\.js$|\.js[#?].*$)/g.test(src)) {
      return src
    }

    return `/_/assets/islands/${src}.js`
  }

  private _clearContent() {
    this.childNodes.forEach((node) => {
      if (node instanceof HTMLTemplateElement) {
        return
      }
      node.remove()
    })
  }

  private _getData() {
    const dataStr = this.getAttribute('data')
    if (dataStr === null || dataStr === '') {
      return undefined
    }

    return JSON.parse(dataStr)
  }

  private _showLoader() {
    const template = this.querySelector('template[data-type=loader]')
    if (template instanceof HTMLTemplateElement) {
      const clone = template.content.cloneNode(true)
      this.appendChild(clone)
    }
  }

  private _updateComponentData() {
    const mountedIsland = this._mountedIsland
    if (!this.active || !mountedIsland) return

    try {
      const data = this._getData()
      mountedIsland.setData(data)
    } catch (e: unknown) {
      console.error('[ol-island] failed to update component data', e)
    }
  }

  private _destroy() {
    if (!this._mountedIsland) {
      this._clearContent()
      return
    }

    this.dispatchEvent(new CustomEvent('island:before-destroy'))
    this._mountedIsland.dispose()
    this._mountedIsland = undefined
    this.dispatchEvent(new CustomEvent('island:destroy'))
  }

  disconnectedCallback() {
    this._destroy()
  }

  private _preload() {
    const src = this.getAttribute('src')

    if (src) {
      import(src)
    }
  }
}

customElements.define('ol-island', OLIslandElement)

declare global {
  interface HTMLElementTagNameMap {
    'ol-island': OLIslandElement
  }
}

function getIslandFromModule(module: unknown, name: string): OLIsland | undefined {
  if (typeof module !== 'object') throw new Error('module is not an object')
  if (module === null) throw new Error('module is null')
  if (!Object.hasOwnProperty.call(module, name))
    throw new Error(`module has no export with name ${name}`)

  let island: OLIsland | undefined = (module as Record<string, OLIsland>)[name]
  if (island) return island
  island = (module as { default: Record<string, OLIsland> }).default[name]
  return island
}
