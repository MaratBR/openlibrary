import { createContext } from 'preact'
import { useContext, useRef } from 'preact/hooks'
import { DraftDto } from '../bookmanager/contracts'
import { PropsWithChildren } from 'preact/compat'

export class ChapterState {
  constructor(draft: DraftDto) {}
}

const ChapterStateContext = createContext<ChapterState | null>(null)

export function ChapterStateProvider({
  draft,
  children,
}: PropsWithChildren<{
  draft: DraftDto
}>) {
  const state = useRef<ChapterState | null>(null)
  if (state.current === null) state.current = new ChapterState(draft)

  return (
    <ChapterStateContext.Provider value={state.current}>{children}</ChapterStateContext.Provider>
  )
}

export function useChapterState(): ChapterState {
  const ctx = useContext(ChapterStateContext)
  if (!ctx) {
    throw new Error('ChapterStateContext is not available')
  }
  return ctx
}
