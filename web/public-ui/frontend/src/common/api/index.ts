import { bookAPI } from "./book-api";
import { reviewsAPI } from './reviews-api'

const api = {
  book: bookAPI,
  reviews: reviewsAPI
}

declare global {
  interface OLGlobal {
    api: typeof api
  }
}

window.OL.api = api;