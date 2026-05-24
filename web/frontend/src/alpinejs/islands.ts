import BookCardPreviewIsland from '@/islands/book-card-preview'
import { Islands } from './island'
import SearchFiltersIsland from '@/islands/search-filters'
import UserMenuIsland from '@/islands/header-user-menu'

Islands.register('bookmanager/BM', () => import('@/islands/bookmanager').then((r) => r.BMIsland))

Islands.register('search/filters', () => Promise.resolve(SearchFiltersIsland))

Islands.register('book-card-preview', () => Promise.resolve(BookCardPreviewIsland))

Islands.register('header-user-menu', () => Promise.resolve(UserMenuIsland))
