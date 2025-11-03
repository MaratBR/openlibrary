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

export type DelayFn<Args extends unknown[]> = {
  (...args: Args): void
  cancel: () => void
}

export function delayFn<Args extends unknown[]>(
  fn: (...args: Args) => void,
  ms: number,
): DelayFn<Args> {
  let timer: number = -1
  let args: Args

  const fn2 = (..._args: Args) => {
    args = _args
    if (timer !== -1) return

    timer = window.setTimeout(() => {
      fn(...args)
      timer = -1
    }, ms)
  }

  ;(fn2 as unknown as { cancel: () => void }).cancel = () => {
    if (timer !== -1) clearTimeout(timer)
  }

  return fn2 as {
    (...args: Args): void
    cancel: () => void
  }
}
