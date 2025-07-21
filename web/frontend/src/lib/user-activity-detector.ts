declare global {
  interface Window {
    FIRST_USER_ACTIVITY?: boolean
  }
}

let FIRST_USER_ACTIVITY = false

export function initFirstActivityEvent() {
  window.FIRST_USER_ACTIVITY = FIRST_USER_ACTIVITY
  if (FIRST_USER_ACTIVITY) return

  let timeout = 0,
    count = 0

  const cb = () => {
    count++
    if (timeout === 0) {
      timeout = window.setTimeout(() => {
        if (count > 5) {
          window.FIRST_USER_ACTIVITY = FIRST_USER_ACTIVITY = true
          document.dispatchEvent(new CustomEvent('first_user_activity'))
          window.removeEventListener('scroll', cb)
          document.removeEventListener('mousemove', cb)
          document.removeEventListener('keydown', cb)
        }
        timeout = count = 0
      }, 1000)
    }
  }

  window.addEventListener('scroll', cb)
  document.addEventListener('mousemove', cb)
  document.addEventListener('keydown', cb)
}
