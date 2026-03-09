import NewBook from './new-book'
import BookCollectionsPopup from './BookCollectionsPopup'
import { PreactIsland } from '../common/preact-island'
import BM from './BM'

export { NewBook, BookCollectionsPopup }

export const BMIsland = new PreactIsland(BM)
