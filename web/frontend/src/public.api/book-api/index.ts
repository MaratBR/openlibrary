import * as readingListAPI from './reading-list'

/** @deprecated Move APIs to their owning feature and call them from an island. */
const bookAPI = {
  ...readingListAPI,
}

export { bookAPI }
