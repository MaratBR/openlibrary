import Alpine from 'alpinejs'

Alpine.data('TagsAutocomplete', () => ({
  _id: Math.random().toString().substring(2),
  _dispose: null as null | (() => void),

  init() {
    import('./TagsAutocompleteSearcher').then((m) => {
      this._dispose = m.initTagsAutocomplete(this.$root)
    })
  },

  destroy() {
    if (this._dispose) this._dispose()
  },
}))
