import { Ref, RefCallback } from 'react'
import { useCallback, useRef } from 'react'

export function useForkRef<T>(...refs: Ref<T>[]): RefCallback<T> {
  const refsRef = useRef(refs)
  refsRef.current = refs

  return useCallback((value: T | null) => {
    refsRef.current.forEach((ref) => {
      if (!ref) return
      if (typeof ref === 'function') {
        ref(value)
      } else {
        ref.current = value
      }
    })
  }, [])
}
