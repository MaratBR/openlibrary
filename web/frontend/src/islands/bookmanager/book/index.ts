import { PreactIsland } from '@/islands/common/preact-island'
import BookEditor from './BookEditor'
import BookCoverComponent from './BookCover'
import Chapters from './Chapters'

export const Book = new PreactIsland(BookEditor)
export const BookCover = new PreactIsland(BookCoverComponent)
export const BookChapters = new PreactIsland(Chapters)
