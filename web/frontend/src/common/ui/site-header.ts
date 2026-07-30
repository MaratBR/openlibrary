import { executeAfterDOMIsReady } from '../dom'
import { delayFn } from '@/common/util/fn'

executeAfterDOMIsReady(() => {
  const header = document.getElementById('site-header')
  if (!header) {
    return
  }

  let isBlurred = false

  window.addEventListener(
    'scroll',
    delayFn(() => {
      const { scrollY } = window
      const newIsBlurred = scrollY > 10
      if (newIsBlurred != isBlurred) {
        header.classList.toggle('site-header--blurred', newIsBlurred)
        isBlurred = newIsBlurred
      }
    }, 100),
  )
})
