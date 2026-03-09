import { useCallback, useMemo } from 'preact/hooks'
import { useSearchParams } from 'react-router'

export type UseQueryParameter = [
  value: string,
  setValue: (value: string | ((prev: string) => string), replace?: boolean) => void,
]

export type QueryParameterOptions = {
  defaultValue: string
  name: string
  map?: (value: string) => string
}

function getValue(parameters: URLSearchParams, options: QueryParameterOptions): string {
  let value = (parameters.get(options.name) || '').trim() || options.defaultValue

  if (options.map) {
    value = options.map(value)
  }

  return value
}

export function createQueryParameter(options: QueryParameterOptions): () => UseQueryParameter {
  return () => {
    const [parameters, setParameters] = useSearchParams()

    const value = useMemo(() => getValue(parameters, options), [parameters])

    const setValue = useCallback(
      (val: string | ((prev: string) => string), replace: boolean = false) => {
        if (typeof val === 'string') {
          let valStr = val
          if (options.map) {
            valStr = options.map(valStr)
          }
          setParameters(
            (prev) => {
              prev.set(options.name, valStr)
              return prev
            },
            { replace },
          )
        } else {
          setParameters((prev) => {
            const oldVal = getValue(prev, options)
            let newVal = val(oldVal)
            if (options.map) newVal = options.map(newVal)
            prev.set(options.name, newVal)
            return prev
          })
        }
      },
      [setParameters],
    )

    return [value, setValue]
  }
}

export function createEnumParameter(name: string, values: string[]) {
  if (values.length === 0)
    return createQueryParameter({
      name,
      defaultValue: '',
    })

  return createQueryParameter({
    name,
    defaultValue: values[0],
    map(value) {
      return values.includes(value) ? value : values[0]
    },
  })
}
