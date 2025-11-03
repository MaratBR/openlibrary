import { debounce } from './fn'

declare global {
  interface Window {
    debounce: typeof debounce
  }
}

window.debounce = debounce
