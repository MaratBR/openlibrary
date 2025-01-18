declare global {
  interface Window {
    delay: (ms: number) => Promise<void>;
  }
}

window.delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))