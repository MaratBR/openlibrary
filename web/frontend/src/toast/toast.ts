import Toastify from 'toastify-js'
// import 'toastify-js/src/toastify.css'

function toast(message: string) {
  Toastify({
    text: message,
    duration: 3000,
    newWindow: true,
    close: true,
    gravity: 'top',
    position: 'right',
    stopOnFocus: true,
  }).showToast()
}

declare global {
  interface Window {
    Toastify: typeof Toastify
    toast: typeof toast
  }
}

window.Toastify = Toastify
window.toast = toast
