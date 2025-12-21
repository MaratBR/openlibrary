import { Islands } from './island'

Islands.register('bookmanager/book/BookChaptersIsland', () =>
  import('@/islands/bookmanager/book').then((m) => m.BookChaptersIsland),
)
