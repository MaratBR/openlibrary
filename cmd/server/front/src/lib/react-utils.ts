import React, { useCallback } from 'react'

export function useForceRender() {
  const [, forceRender] = React.useState(0)

  return useCallback(() => {
    forceRender((i) => i + 1)
  }, [forceRender])
}

export function useDebounce<T>(value: T, milliSeconds: number): T {
  const [debouncedValue, setDebouncedValue] = React.useState(value)

  React.useEffect(() => {
    const handler = window.setTimeout(() => {
      setDebouncedValue(value)
    }, milliSeconds)

    return () => {
      clearTimeout(handler)
    }
  }, [value, milliSeconds])

  return debouncedValue
}
