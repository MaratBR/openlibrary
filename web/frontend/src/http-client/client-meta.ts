import { getCookie } from './util'

const rndId = () => Math.random().toString(16).substring(2)

function isLocalStorageAvailable() {
  if (!window.localStorage) return false
  const test = rndId()
  try {
    window.localStorage.setItem(test, test)
    if (window.localStorage.getItem(test) !== test) return false
    window.localStorage.removeItem(test)
    return true
  } catch {
    return false
  }
}

const DEVICE_ID_KEY = 'OL_DEVICE_ID'
let deviceId: string = ''

function init() {
  if (isLocalStorageAvailable()) {
    deviceId = window.localStorage.getItem(DEVICE_ID_KEY) || ''
    if (deviceId === '') {
      deviceId = rndId()
      window.localStorage.setItem(DEVICE_ID_KEY, deviceId)
    } else {
      deviceId = getCookie(DEVICE_ID_KEY) || ''
      if (deviceId === '') {
        deviceId = rndId()
      }
    }
  }

  document.cookie = `${DEVICE_ID_KEY}=${deviceId};path=/;max-age=31536000;`
}

window.addEventListener('DOMContentLoaded', init)
