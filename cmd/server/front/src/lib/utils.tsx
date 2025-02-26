import { clsx, type ClassValue } from 'clsx'
import React from 'react'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function useEffectOnce(callback: () => void) {
  // eslint-disable-next-line react-hooks/exhaustive-deps
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

export function debounce<Args extends unknown[]>(
  fn: (...args: Args) => void,
  ms: number,
): (...args: Args) => void {
  let timer: number

  return (...args: Args) => {
    clearTimeout(timer)

    timer = window.setTimeout(() => {
      fn(...args)
    }, ms)
  }
}

export function isNotFalsy<T>(value: T | null | undefined): value is T {
  return !!value
}

export function toDictionaryByProperty<T extends object>(
  array: T[],
  key: keyof T,
): Record<string, T> {
  const d: Record<string, T> = {}

  for (let i = 0; i < array.length; i++) {
    const item = array[i]
    d[item[key] as string] = item
  }

  return d
}

export function delayMs(ms: number): Promise<void> {
  // if (ms <= 0) return Promise.resolve()
  return new Promise((resolve) => window.setTimeout(resolve, ms))
}
