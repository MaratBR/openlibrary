import { initParagraphs } from './paragraphs'

function init(...args: unknown[]) {
  if (args.length === 0) {
    throw new Error('init function should have arguments')
  }

  if (!(args[0] instanceof HTMLElement)) {
    throw new Error('First argument should be HTMLElement')
  }

  initParagraphs(args[0])
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
;(window as any).__initBookReader = init
document.dispatchEvent(new CustomEvent('ol:book-reader:ready'))
