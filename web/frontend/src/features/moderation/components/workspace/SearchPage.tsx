import { Link, LoaderFunctionArgs, useLoaderData } from 'react-router'
import { moderationRead } from '../../api/workspace'
import type {
  ModerationContentPage,
  ModerationBooksPageResponse,
  ModerationUsersPageResponse,
  ModerationReportsSearchResponse,
} from '@/backend-types'
import { WorkspacePage, SearchForm, FilterInput, EmptyState, t } from './ui'
export async function globalSearchLoader({ request }: LoaderFunctionArgs) {
  const search = new URL(request.url).searchParams.get('search')?.trim() ?? ''
  if (!search) return { search, groups: [] }
  const f = { search, pageSize: '5' }
  const [users, books, chapters, comments, reports] = await Promise.all([
    moderationRead<ModerationUsersPageResponse>('users', f, request.signal),
    moderationRead<ModerationBooksPageResponse>(
      'books',
      { ...f, includeBanned: 'true', includeDeleted: 'true' },
      request.signal,
    ),
    moderationRead<ModerationContentPage>('chapters', f, request.signal),
    moderationRead<ModerationContentPage>('comments', f, request.signal),
    moderationRead<ModerationReportsSearchResponse>(
      'reports',
      { ...f, status: 'all' },
      request.signal,
    ),
  ])
  return {
    search,
    groups: [
      {
        kind: 'users',
        total: users.total,
        entries: users.entries.map((e) => ({ id: e.id, title: e.name })),
      },
      {
        kind: 'books',
        total: books.total,
        entries: books.entries.map((e) => ({ id: e.id, title: e.name })),
      },
      { kind: 'chapters', total: chapters.total, entries: chapters.entries },
      { kind: 'comments', total: comments.total, entries: comments.entries },
      {
        kind: 'reports',
        total: reports.total,
        entries: reports.entries.map((e) => ({ id: e.id, title: `${e.number} · ${e.reason}` })),
      },
    ],
  }
}
export default function SearchPage() {
  const { search, groups } = useLoaderData<typeof globalSearchLoader>()
  return (
    <WorkspacePage title={t('search')}>
      <SearchForm>
        <FilterInput name="search" value={search} label={t('globalSearch')} />
      </SearchForm>
      {!search ? (
        <EmptyState>{t('searchHelp')}</EmptyState>
      ) : groups.every((g) => !g.total) ? (
        <EmptyState>{t('noResults')}</EmptyState>
      ) : (
        <div className="grid lg:grid-cols-2 gap-8">
          {groups.map((group) => (
            <section key={group.kind}>
              <div className="flex flex-wrap justify-between gap-3">
                <h2 className="text-xl font-semibold">
                  {t(group.kind)}{' '}
                  <span className="text-secondary-foreground">
                    ({group.total.toLocaleString()})
                  </span>
                </h2>
                <Link
                  className="Link"
                  to={`/${group.kind}?${new URLSearchParams({ search, ...(group.kind === 'reports' ? { status: 'all' } : {}), ...(group.kind === 'books' ? { includeBanned: 'true', includeDeleted: 'true' } : {}) })}`}
                >
                  {t('viewAll')}
                </Link>
              </div>
              {group.entries.length ? (
                group.entries.map((e) => (
                  <Link className="ModWorkspace-record" key={e.id} to={`/${group.kind}/${e.id}`}>
                    <span className="line-clamp-3">{e.title}</span>
                    <span className="text-sm text-secondary-foreground">{e.id}</span>
                  </Link>
                ))
              ) : (
                <EmptyState>{t('noResults')}</EmptyState>
              )}
            </section>
          ))}
        </div>
      )}
    </WorkspacePage>
  )
}
