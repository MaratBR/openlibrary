export function executeAfterDOMIsReady(callback: () => void): void {
  // Not in a browser (e.g. SSR or worker) — run async so behavior is predictable.
  if (typeof document === 'undefined') {
    if (typeof queueMicrotask !== 'undefined') queueMicrotask(callback)
    else setTimeout(callback, 0)
    return
  }

  // If the DOM is still loading, wait for DOMContentLoaded.
  if (document.readyState === 'loading') {
    const onReady = () => {
      document.removeEventListener('DOMContentLoaded', onReady)
      callback()
    }
    document.addEventListener('DOMContentLoaded', onReady)
    return
  }

  // readyState is 'interactive' or 'complete' — DOM is available.
  // Run async to avoid surprising synchronous reentrancy.
  if (typeof queueMicrotask !== 'undefined') queueMicrotask(callback)
  else setTimeout(callback, 0)
}
