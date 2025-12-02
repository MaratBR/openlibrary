import { createContext } from 'preact'
import { useContext, useRef } from 'preact/hooks'
import { PropsWithChildren } from 'preact/compat'
import { DraftDto } from '../contracts'
import { ViewState } from './view'

export class ChapterState {
  private readonly _draft: DraftDto
  public view = new ViewState('100%')

  constructor(draft: DraftDto) {
    this._draft = draft
  }

  get draft() {
    return this._draft
  }
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
