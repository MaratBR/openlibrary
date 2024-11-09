import React from 'react'
import { NavigateOptions } from 'react-router'
import { useSearchParams } from 'react-router-dom'

type SetQueryParam = (value: string | null, navigateOps?: NavigateOptions) => void

export function useQueryParam(param: string): [value: string | null, setValue: SetQueryParam] {
  const [params, setParams] = useSearchParams()
  const paramsRef = React.useRef(params)
  paramsRef.current = params

  const setParam = React.useCallback(
    (value: string | null, navigateOps?: NavigateOptions) => {
      const newParams = new URLSearchParams(paramsRef.current.toString())
      if (value === null) {
        newParams.delete(param)
      } else {
        newParams.set(param, value)
      }
      setParams(newParams, navigateOps)
    },
    [param, setParams],
  )

  return [params.get(param), setParam]
}

export function useQueryParamDefault(
  param: string,
  defaultValue: string,
): [value: string | null, setValue: SetQueryParam] {
  const [value, setValue] = useQueryParam(param)

  React.useEffect(() => {
    if (value === null) {
      setValue(defaultValue, { replace: true })
    }
  }, [defaultValue, setValue, value])

  return [value, setValue]
}
