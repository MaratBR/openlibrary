import { z } from 'zod'
import { BroadcastChannel } from 'broadcast-channel'
import kefir, { Emitter, Observable } from 'kefir'

function subscribeToSystemTheme(callback: (isDark: boolean) => void): () => void {
  if (!window.matchMedia) {
    console.warn('browser does not support matchMedia')
    return () => {}
  }

  const query = window.matchMedia('(prefers-color-scheme: dark)')
  callback(query.matches)
  const onChange = () => {
    callback(query.matches)
  }
  query.addEventListener('change', onChange)

  return () => {
    query.removeEventListener('change', onChange)
  }
}

function createPrefersColorSchemeObservable(): Observable<boolean, never> {
  const subject = kefir.stream<boolean, never>((emitter) => {
    return subscribeToSystemTheme((isDark) => {
      emitter.emit(isDark)
    })
  })
  return subject
}

const themeSchema = z.enum(['light', 'dark', 'system'])
export type Theme = z.infer<typeof themeSchema>

function createThemeSettingsObservable(): [Observable<Theme, never>, (theme: Theme) => void] {
  function getSavedThemeOrDefault(): Theme {
    const theme = localStorage['theme']
    if (themeSchema.safeParse(theme).success) {
      return theme
    }
    return 'system'
  }

  const themeBroadcastChannel = new BroadcastChannel<Theme>('theme')
  let emitterOuter: Emitter<Theme, never>
  const subject = kefir.stream<Theme, never>((emitter) => {
    emitterOuter = emitter
    emitter.emit(getSavedThemeOrDefault())
    const onMessage = (theme: Theme) => {
      emitter.emit(theme)
    }
    themeBroadcastChannel.addEventListener('message', onMessage)
    return () => {
      themeBroadcastChannel.removeEventListener('message', onMessage)
    }
  })

  const set = (theme: Theme) => {
    if (!themeSchema.safeParse(theme).success) {
      throw new Error('invalid theme value: ' + theme)
    }
    localStorage['theme'] = theme
    themeBroadcastChannel.postMessage(theme)
    emitterOuter.emit(theme)
  }

  return [subject, set]
}

const prefersColorScheme = createPrefersColorSchemeObservable()
const [themeSettings, setTheme] = createThemeSettingsObservable()

const theme: Observable<Theme, never> = kefir
  .combine({ theme: themeSettings, isDark: prefersColorScheme })
  .map(({ theme, isDark }) => {
    if (theme === 'system') {
      return isDark ? 'dark' : 'light'
    } else {
      return theme
    }
  })
  .ignoreErrors()

theme.onValue((theme) => {
  document.body.classList.toggle('dark', theme === 'dark')
})

if (import.meta.env.DEV) {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ;(window as any).setTheme = setTheme
}
