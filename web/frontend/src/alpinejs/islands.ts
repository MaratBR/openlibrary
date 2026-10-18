import BookCardPreviewIsland from '@/islands/book-card-preview'
import { Islands } from './island'
import SearchFiltersIsland from '@/islands/search-filters'

Islands.register('bookmanager/BM', () => import('@/islands/bookmanager').then((r) => r.BMIsland))

Islands.register('moderation/Portal', () =>
  import('@/features/moderation').then((r) => r.ModerationPortalIsland),
)

Islands.register('search/filters', () => Promise.resolve(SearchFiltersIsland))

Islands.register('book-card-preview', () => Promise.resolve(BookCardPreviewIsland))

// Keep this lazy: Vite loads the shared panel CSS with the island chunk.
Islands.register('header-actions', () =>
  import('@/islands/header-user-menu').then((module) => module.default),
)

Islands.register('comments', () => import('@/islands/comments').then((module) => module.default))
Islands.register('report-button', () =>
  import('@/features/moderation/components/report-button').then((module) => module.default),
)
