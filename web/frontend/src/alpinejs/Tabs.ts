import Alpine from 'alpinejs'

Alpine.data('Tabs', () => ({
  tab: '',
  tabParam: 'tab',
  fallbackTab: '',
  _dispose: undefined as undefined | (() => void),

  init() {
    this.tabParam = this.$root.dataset.queryParam || this.tabParam
    this.fallbackTab = this.$root.dataset.fallbackTab || this.fallbackTab
    this._update()

    const onPopState = () => {
      this._update()
    }
    window.addEventListener('popstate', onPopState)

    this._dispose = () => {
      window.removeEventListener('popstate', onPopState)
    }
  },

  destroy() {
    this._dispose?.()
  },

  setTab(tab: string) {
    this.tab = tab

    const url = new URL(window.location.toString())
    url.searchParams.set('tab', tab)

    history.pushState({ tab }, '', url)
  },

  _update() {
    const newTab =
      new URLSearchParams(window.location.search).get(this.tabParam) || this.fallbackTab
    console.log(window.location.search)
    console.log(newTab)
    this.tab = newTab
  },
}))
