import * as ratingAPI from './rating'

/** @deprecated Move APIs to their owning feature and call them from an island. */
const reviewsAPI = {
  ...ratingAPI,
}

export { reviewsAPI }
