import { bookAPI } from './book-api'
import { reviewsAPI } from './reviews-api'
import * as userDataAPI from './user-data'

const api = {
  book: bookAPI,
  reviews: reviewsAPI,
  userData: userDataAPI,
}

declare global {
  interface OLGlobal {
    api: typeof api
  }
}

window.OL.api = api
