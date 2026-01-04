import BookCardPreviewIsland from '@/islands/book-card-preview'
import { Islands } from './island'
import SearchFiltersIsland from '@/islands/search-filters'

Islands.register('bookmanager/book/BookChaptersIsland', () =>
  import('@/islands/bookmanager/book').then((m) => m.BookChaptersIsland),
)

Islands.register('search/filters', () => Promise.resolve(SearchFiltersIsland))

Islands.register('book-card-preview', () => Promise.resolve(BookCardPreviewIsland))
