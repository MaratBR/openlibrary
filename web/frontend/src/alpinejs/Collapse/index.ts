import Alpine from 'alpinejs'

Alpine.data('Collapse', () => ({
  can: false,
  expand: false,
  resizeObserver: null as ResizeObserver | null,

  init() {
    this.can = this.$root.hasAttribute('data-collapsible-init')

    window.requestAnimationFrame(() => {
      this.resizeObserver = new ResizeObserver(() => {
        this.recalculate()
      })
      this.resizeObserver.observe(this.$refs.content)
    })
  },

  destroy() {
    this.resizeObserver?.disconnect()
    this.resizeObserver = null
  },

  recalculate() {
    const el = this.$refs.content
    const maxHeight = +(this.$root.dataset.collapsibleHeight || '')
    this.can = el.clientHeight + 20 > maxHeight
  },

  content: {
    'x-ref': 'content',
  },

  button: {
    '@click'() {
      this.expand = !this.expand
    },

    'x-show'() {
      return this.can
    },
  },

  buttonLabel: {
    'x-text'() {
      return this.expand ? window._('common.less') : window._('common.more')
    },
  },

  buttonIcon: {
    ':style'() {
      return this.expand ? 'transform:rotate(180deg)' : ''
    },
  },
}))
