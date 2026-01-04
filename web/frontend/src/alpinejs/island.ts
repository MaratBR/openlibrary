import { OLIsland, OLIslandMounted } from '@/lib/island'
import Alpine from 'alpinejs'

export namespace Islands {
  const registry: Record<string, OLIsland | (() => Promise<OLIsland>)> = {}

  export function register(name: string, island: OLIsland | (() => Promise<OLIsland>)) {
    registry[name] = island
  }

  export function getByName(name: string): OLIsland | (() => Promise<OLIsland>) | undefined {
    return registry[name]
  }
}

Alpine.data('Island', ({ name, data }: { name: string; data: unknown }) => ({
  _mounted: null as null | OLIslandMounted,
  _loading: false,

  init() {
    const islandOrProvider = Islands.getByName(name)

    if (!islandOrProvider) {
      console.error(`[Islands] island ${name} not found`)
    } else if (typeof islandOrProvider === 'function') {
      this._loading = true
      islandOrProvider()
        .then((island) => {
          this._loading = false
          this._mount(island)
        })
        .catch((err) => {
          this._loading = false
          this._error = err
        })
    } else {
      this._mount(islandOrProvider)
    }
  },

  _mount(island: OLIsland) {
    this._unmount()

    const container = document.createElement('div')
    container.style.display = 'contents'
    this._mounted = island.mount(container, data)
    this.$root.appendChild(container)

    this.$el.dispatchEvent(new CustomEvent('island:mount'))
  },

  _unmount() {
    if (this._mounted) {
      this._mounted.dispose()
      this._mounted = null
    }

    if (this.$refs._container instanceof HTMLElement) {
      this.$refs._container.remove()
    }
  },

  _error(_err: unknown) {
    this._mounted?.dispose()
    this._mounted = null
  },

  destroy() {
    this._unmount()
  },
}))
