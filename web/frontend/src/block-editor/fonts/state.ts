import { useCallback, useEffect, useSyncExternalStore } from 'react'
import { appRuntime } from '@/effect/runtime'
import { EditorFonts } from './service'

const editorFonts = appRuntime.runSync(EditorFonts)

function useEditorFontsState() {
  return useSyncExternalStore(editorFonts.subscribe, editorFonts.getState, editorFonts.getState)
}

export function useInitializeEditorFonts() {
  useEffect(() => {
    void appRuntime.runPromise(editorFonts.initialize()).catch(() => undefined)
  }, [])
}

export function useFonts() {
  const state = useEditorFontsState()

  return {
    ...state,
    init: editorFonts.initialize,
    favorite: editorFonts.setFavorite,
  }
}

export function useFavoriteFonts() {
  const { favoriteFonts } = useFonts()
  return favoriteFonts
}

export function useFavoriteFontState(font: string) {
  const selected = useEditorFontsState().favoriteFonts.includes(font)
  const select = useCallback(
    (favorite: boolean) => appRuntime.runPromise(editorFonts.setFavorite(font, favorite)),
    [font],
  )

  return { selected, select }
}
