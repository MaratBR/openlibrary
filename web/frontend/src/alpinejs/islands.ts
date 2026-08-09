import BookCardPreviewIsland from '@/islands/book-card-preview'
import { Islands } from './island'
import SearchFiltersIsland from '@/islands/search-filters'
import UserMenuIsland from '@/islands/header-user-menu'

Islands.register('bookmanager/BM', () => import('@/islands/bookmanager').then((r) => r.BMIsland))

Islands.register('moderation/Portal', () =>
  import('@/islands/moderation').then((r) => r.ModerationPortalIsland),
)

Islands.register('search/filters', () => Promise.resolve(SearchFiltersIsland))

Islands.register('book-card-preview', () => Promise.resolve(BookCardPreviewIsland))

Islands.register('header-user-menu', () => Promise.resolve(UserMenuIsland))

Islands.register('comments', () => import('@/islands/comments').then((module) => module.default))
Islands.register('report-button', () => import('@/islands/report-button').then((module) => module.default))
