import { getCookie, setCookie } from './cookies'
import { executeAfterDOMIsReady } from './dom'
import { PREFERS_DARK_MODE } from './media'
import { Derived, Subject } from './rx'

const SYSTEM_THEME = 'system'
const LIGHT_THEME = 'light'
const DARK_THEME = 'dark'
type Theme = typeof SYSTEM_THEME | typeof DARK_THEME | typeof LIGHT_THEME

namespace OLTheme {
  export const theme = new Subject<Theme>(SYSTEM_THEME)

  export function toggle() {
    const value = theme.get()

    switch (value) {
      case 'system':
        theme.set('light')
        break
      case 'light':
        theme.set('dark')
        break
      case 'dark':
        theme.set('system')
        break
    }
  }

  export const isDarkThemeActive = new Derived(
    [theme, PREFERS_DARK_MODE],
    (theme, prefersDarkMode): boolean => {
      switch (theme) {
        case SYSTEM_THEME:
          return prefersDarkMode
        case DARK_THEME:
          return true
        case LIGHT_THEME:
          return false
        default:
          return false
      }
    },
  )

  // Reader choices affect the current page without changing the website preference.
  export function applyPageTheme(readerTheme?: string): void {
    const dark = readerTheme === 'dark' || readerTheme === 'oled'
      || (readerTheme !== 'light' && isDarkThemeActive.get())
    const html = document.documentElement
    html.classList.toggle('dark', dark)
    html.dataset.theme = readerTheme === 'oled' || readerTheme === 'dark'
      ? 'dark' : readerTheme === 'light' ? 'light' : theme.get()
    if (readerTheme) html.dataset.readerTheme = readerTheme
    else delete html.dataset.readerTheme
  }

  document.dispatchEvent(new CustomEvent('OLTheme:ready'))

  const THEME_COOKIE = '_theme'

  const themeFromCookie = getCookie(THEME_COOKIE)

  if (themeFromCookie) {
    if (
      themeFromCookie === SYSTEM_THEME ||
      themeFromCookie == LIGHT_THEME ||
      themeFromCookie === DARK_THEME
    ) {
      theme.set(themeFromCookie)
    } else {
      setCookie(THEME_COOKIE, theme.get())
    }
  }

  theme.subscribe((t) => setCookie(THEME_COOKIE, t))
}

declare global {
  interface Window {
    OLTheme: typeof OLTheme
    OL_DEFAULT_THEME?: string
  }
}

window.OLTheme = OLTheme

executeAfterDOMIsReady(() => {
  const applyTheme = () => {
    OLTheme.applyPageTheme(document.querySelector<HTMLElement>('.BookReader')?.dataset.readerTheme)
  }
  OLTheme.isDarkThemeActive.subscribe(applyTheme)
  OLTheme.theme.subscribe(applyTheme)
  applyTheme()
})
