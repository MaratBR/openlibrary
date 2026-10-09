import Alpine from 'alpinejs'
import { setCookie } from '@/common/cookies'

const FONT_SIZES = [12, 14, 16, 18, 20, 22, 26, 30, 36, 42, 48]
const PAGE_COLORS = ['background', 'surface'] as const
const READER_THEMES = ['system', 'light', 'dark'] as const

type PageColor = (typeof PAGE_COLORS)[number]
type ReaderTheme = (typeof READER_THEMES)[number]

Alpine.data('BookReader', () => ({
  settingsOpen: false,
  fontSize: 18,
  fontFamily: 'serif',
  contentWidth: 72,
  pageColor: 'background' as PageColor,
  readerTheme: 'system' as ReaderTheme,
  authenticated: false,
  saveTimer: undefined as number | undefined,

  init() {
    this.fontSize = validNumber(this.$el.dataset.fontSize, FONT_SIZES, 18)
    const fontIds = Array.from(this.$el.querySelectorAll<HTMLInputElement>('[data-font-stack]'), (input) => input.value)
    this.fontFamily = validValue(this.$el.dataset.fontFamily, fontIds, fontIds[0] ?? 'serif')
    this.contentWidth = validNumber(this.$el.dataset.contentWidth, Array.from({ length: 27 }, (_, i) => 48 + i * 2), 72)
    this.pageColor = validValue(this.$el.dataset.pageColor, PAGE_COLORS, 'background')
    this.readerTheme = validValue(this.$el.dataset.readerTheme, READER_THEMES, 'system')
    this.authenticated = this.$el.dataset.authenticated === 'true'
    this.applyPreferences(false)
    window.OLTheme.theme.subscribe(() => applyReaderTheme(this.readerTheme))
  },

  changeFontSize(increase: boolean) {
    const currentIndex = FONT_SIZES.indexOf(this.fontSize)
    const nextIndex = Math.max(
      0,
      Math.min(FONT_SIZES.length - 1, currentIndex + (increase ? 1 : -1)),
    )
    this.fontSize = FONT_SIZES[nextIndex]
    this.applyPreferences()
  },

  applyPreferences(persist = true) {
    document.documentElement.style.setProperty('--book-font-size', `${this.fontSize}px`)
    this.$root.setAttribute('data-page-color', this.pageColor)
    this.$root.setAttribute('data-font-family', this.fontFamily)
    this.$root.style.setProperty('--reader-content-width', `${this.contentWidth}ch`)
    const selectedFont = Array.from(this.$root.querySelectorAll<HTMLInputElement>('[data-font-stack]')).find((input) => input.value === this.fontFamily)
    if (selectedFont?.dataset.fontStack) {
      this.$root.style.setProperty('--reader-font-family', selectedFont.dataset.fontStack)
    }
    applyReaderTheme(this.readerTheme)

    setCookie('reader_font_size', String(this.fontSize))
    setCookie('reader_font_family', this.fontFamily)
    setCookie('reader_page_color', this.pageColor)
    setCookie('reader_theme', this.readerTheme)
    setCookie('reader_content_width', String(this.contentWidth))

    if (persist && this.authenticated) {
      window.clearTimeout(this.saveTimer)
      this.saveTimer = window.setTimeout(() => void this.savePreferences(), 250)
    }
  },

  async savePreferences() {
    try {
      await fetch('/_api/reader-preferences', {
        method: 'POST',
        headers: { 'content-type': 'application/json' },
        body: JSON.stringify({
          fontSize: this.fontSize,
          contentWidth: this.contentWidth,
          fontFamily: this.fontFamily,
          pageColor: this.pageColor,
          theme: this.readerTheme,
        }),
      }).then((response) => {
        if (!response.ok) throw new Error(`Reader preferences request failed: ${response.status}`)
      })
    } catch (error) {
      console.error(error)
    }
  },

  toggleButton: {
    '@click'() {
      this.settingsOpen = !this.settingsOpen
    },
    ':aria-expanded'() {
      return String(this.settingsOpen)
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
    ':disabled'() {
      return this.fontSize === FONT_SIZES.at(-1)
    },
  },

  decreaseFont: {
    '@click'() {
      this.changeFontSize(false)
    },
    ':disabled'() {
      return this.fontSize === FONT_SIZES[0]
    },
  },

  closeButton: {
    '@click'() {
      this.settingsOpen = false
    },
  },
}))

function validNumber(
  value: string | undefined,
  values: readonly number[],
  fallback: number,
): number {
  const parsed = Number(value)
  return values.includes(parsed) ? parsed : fallback
}

function validValue<T extends string>(
  value: string | undefined,
  values: readonly T[],
  fallback: T,
): T {
  return values.includes(value as T) ? (value as T) : fallback
}

function applyReaderTheme(theme: ReaderTheme): void {
  let dark = false
  if (theme === 'dark') dark = true
  else if (theme === 'system') dark = window.OLTheme.isDarkThemeActive.get()
  document.documentElement.classList.toggle('dark', dark)
}

export type CurrentPosition = {
  window: { height: number; width: number; scrollY: number }
  nearestElement: { path: string; id: string | null; top: number }
}
