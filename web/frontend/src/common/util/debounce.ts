/* eslint-disable no-unused-vars */
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
