import './logger'

import { httpUpdateReadingListStartReading, httpUpdateReadingListStatus } from './book-api'
import './common.css'

document.addEventListener('DOMContentLoaded', () => {
  document.querySelectorAll('[data-hidden-initially]').forEach(el => {
    el.removeAttribute('data-hidden-initially')
  })
})

const bookAPI = {
  httpUpdateReadingListStatus,
  httpUpdateReadingListStartReading
}

interface OLGlobals {
  bookAPI: typeof bookAPI
}

const globals = (window as unknown) as OLGlobals
globals.bookAPI = bookAPI