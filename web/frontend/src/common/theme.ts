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
  const HTML = document.getElementsByTagName('html')[0]

  const isDarkCb = (isDark: boolean) => {
    HTML.classList.toggle('dark', isDark)
  }
  OLTheme.isDarkThemeActive.subscribe(isDarkCb)
  isDarkCb(OLTheme.isDarkThemeActive.get())

  const themeCb = (theme: Theme) => {
    HTML.setAttribute('data-theme', theme)
  }
  OLTheme.theme.subscribe(themeCb)
  themeCb(OLTheme.theme.get())
})
