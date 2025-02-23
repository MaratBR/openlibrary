import { debounce } from './debounce'

declare global {
  interface Window {
    debounce: typeof debounce
  }
}

window.debounce = debounce
