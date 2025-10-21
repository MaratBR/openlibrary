export function debounce<Args extends unknown[]>(
  fn: (...args: Args) => void,
  ms: number,
): {
  (...args: Args): void
  cancel: () => void
} {
  let timer: number

  const fn2 = (...args: Args) => {
    clearTimeout(timer)

    timer = window.setTimeout(() => {
      fn(...args)
    }, ms)
  }

  ;(fn2 as unknown as { cancel: () => void }).cancel = () => {
    clearTimeout(timer)
  }

  return fn2 as {
    (...args: Args): void
    cancel: () => void
  }
}
