import { clsx, type ClassValue } from 'clsx'
import React from 'react'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function useEffectOnce(callback: () => void) {
  React.useEffect(callback, [])
}

type OnlyComponentTypes<TModule extends object> = {
  [K in keyof TModule]: TModule[K] extends React.ComponentType<object> ? K : never
}
type OnlyComponentTypeKeys<TModule extends object> =
  OnlyComponentTypes<TModule>[keyof OnlyComponentTypes<TModule>]

export function componentsChunk<TModule extends object>(fn: () => Promise<TModule>) {
  function componentType<K extends OnlyComponentTypeKeys<TModule>>(key: K) {
    // @ts-expect-error ts(2345) ignoring this error because I don't know how to fix this and it works fine anyway
    return React.lazy(() => {
      const Component = fn().then((r) => ({ default: r[key] }))
      return Component
    })
  }

  function element<K extends OnlyComponentTypeKeys<TModule>>(key: K): React.JSX.Element {
    const Component = componentType(key)

    return <Component />
  }

  return {
    componentType,
    element,
  }
}

export function setRef<T>(ref: React.MutableRefObject<T> | React.RefCallback<T>, value: T) {
  if (typeof ref === 'function') {
    ref(value)
  } else {
    ref.current = value
  }
}

export function useSetRef<T>(
  ref: React.MutableRefObject<T> | React.RefCallback<T> | undefined | null,
  value: T,
) {
  React.useEffect(() => {
    if (!ref) return
    setRef(ref, value)
  }, [ref, value])
}
