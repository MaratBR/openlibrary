import { useDebounce, useForceRender } from '@/lib/react-utils'
import React, { useRef } from 'react'
import LoadingBar from './loading-bar'
import ReactDOM from 'react-dom'

interface ILoadingBarContext {
  set(): void
  release(): void
}

const LoadingBarContext = React.createContext<ILoadingBarContext>({
  set: () => {},
  release: () => {},
})

class LoadingBarContextImpl implements ILoadingBarContext {
  private _counterRef: React.MutableRefObject<number>
  private _update: () => void

  constructor(update: () => void, counterRef: React.MutableRefObject<number>) {
    this._update = update
    this._counterRef = counterRef
  }

  set(): void {
    this._counterRef.current++
    if (this._counterRef.current === 1) {
      this._update()
    }
  }
  release(): void {
    this._counterRef.current--
    if (this._counterRef.current === 0) {
      this._update()
    }
  }
}

export function useHeaderLoading(loading: boolean) {
  const ctx = React.useContext(LoadingBarContext)

  const loadingDebounced = useDebounce(loading, 300)

  React.useEffect(() => {
    if (!loadingDebounced) return

    ctx.set()
    return () => ctx.release()
  }, [loadingDebounced, ctx])
}

export function LoadingBarProvider({ children }: React.PropsWithChildren) {
  const counterRef = useRef(0)
  const fr = useForceRender()

  const ctx = useRef<ILoadingBarContext | null>(null)
  if (ctx.current === null) {
    ctx.current = new LoadingBarContextImpl(fr, counterRef)
  }

  return (
    <LoadingBarContext.Provider value={ctx.current}>
      {children}
      {counterRef.current > 0 &&
        ReactDOM.createPortal(
          <div className="fixed top-0 w-screen z-50">
            <LoadingBar />
          </div>,
          document.body,
        )}
    </LoadingBarContext.Provider>
  )
}
