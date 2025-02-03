import Alpine from 'alpinejs'

const FONT_SIZES = [10, 12, 14, 16, 18, 22, 30, 36, 48, 60, 70, 80, 90, 99]
const FONT_SIZES_REVERSED = FONT_SIZES.slice().reverse()

Alpine.data('bookReader', () => ({
  settingsOpen: false,
  fontSize: 18,

  init() {
    const value = this.$el.getAttribute('data-font-size')
    if (value && !Number.isNaN(+value)) {
      const fontSize = Math.round(+value)
      if (fontSize >= 12 && fontSize <= 50) {
        this.fontSize = fontSize
      }
    }
  },

  changeFontSize(increase: boolean) {
    if (increase) {
      if (this.fontSize === FONT_SIZES_REVERSED[0]) return

      const next = FONT_SIZES.find((f) => f > this.fontSize)
      if (next) {
        this.fontSize = next
      } else {
        this.fontSize = 18
      }
      setFontSizeCookie(this.fontSize)
    } else {
      if (this.fontSize === FONT_SIZES[0]) return

      const next = FONT_SIZES_REVERSED.find((f) => f < this.fontSize)
      if (next) {
        this.fontSize = next
      } else {
        this.fontSize = 18
      }
      setFontSizeCookie(this.fontSize)
    }
  },

  toggleButton: {
    '@click'() {
      this.settingsOpen = !this.settingsOpen
    },
  },

  settings: {
    'x-show'() {
      return this.settingsOpen
    },
  },

  increaseFont: {
    '@click'() {
      this.changeFontSize(true)
    },
  },

  decreaseFont: {
    '@click'() {
      this.changeFontSize(false)
    },
  },
}))

function setFontSizeCookie(fontSize: number) {
  document.cookie = `ifs=${fontSize};path=/;max-age=31536000;`
}
