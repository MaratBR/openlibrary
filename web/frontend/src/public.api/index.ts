import { bookAPI } from './book-api'
import { reviewsAPI } from './reviews-api'

/** @deprecated Compatibility API for server-rendered templates. Do not add new APIs here. */
const api = {
  book: bookAPI,
  reviews: reviewsAPI,
}

declare global {
  interface OLGlobal {
    /** @deprecated Compatibility API for server-rendered templates. */
    api: typeof api
  }
}

window.OL.api = api
