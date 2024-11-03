import React from 'react'
import { NavigateOptions } from 'react-router'
import { useSearchParams } from 'react-router-dom'
import { useEffectOnce } from './utils'

type SetQueryParam = (value: string | null, navigateOps?: NavigateOptions) => void

export function useQueryParam(param: string): [value: string | null, setValue: SetQueryParam] {
  const [params, setParams] = useSearchParams()
  const paramsRef = React.useRef(params)
  paramsRef.current = params

  const setParam = React.useCallback(
    (value: string | null) => {
      const newParams = new URLSearchParams(paramsRef.current.toString())
      if (value === null) {
        newParams.delete(param)
      } else {
        newParams.set(param, value)
      }
      setParams(newParams)
    },
    [setParams],
  )

  return [params.get(param), setParam]
}

export function useQueryParamDefault(
  param: string,
  defaultValue: string,
): [value: string | null, setValue: SetQueryParam] {
  const [value, setValue] = useQueryParam(param)

  useEffectOnce(() => {
    if (value === null) {
      setValue(defaultValue, { replace: true })
    }
  })

  return [value, setValue]
}
