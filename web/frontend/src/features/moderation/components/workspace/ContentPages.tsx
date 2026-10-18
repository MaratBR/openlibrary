import { Link, LoaderFunctionArgs, useLoaderData } from 'react-router'
import { useState } from 'react'
import type {
  ModerationAuditLogPageResponse,
  ModerationReportsSearchResponse,
} from '@/backend-types'
import { getContent, searchContent, moderationRead } from '../../api/workspace'
import { performChapterModerationAction, performCommentModerationAction } from '../../api/content'
import { ModerationActions, ModerationConfirmation } from '../ModerationActions'
import { ModerationLog } from '../ModerationLog'
import { ReportRecords } from './ReportRecords'
import {
  WorkspacePage,
  SearchForm,
  FilterInput,
  FilterSelect,
  PageNavigation,
  filtersFrom,
  t,
  date,
  StateBadge,
  EmptyState,
} from './ui'

type Kind = 'chapters' | 'comments'
export const contentListLoader =
  (kind: Kind) =>
  async ({ request }: LoaderFunctionArgs) => {
    const filters = filtersFrom(request)
    return { page: await searchContent(kind, filters, request.signal), filters }
  }
export function ContentListPage({ kind }: { kind: Kind }) {
  const { page, filters: f } =
    useLoaderData<Awaited<ReturnType<ReturnType<typeof contentListLoader>>>>()
  return (
    <WorkspacePage title={t(kind)}>
      <SearchForm>
        <FilterInput name="search" value={f.search} label={t('contentSearch')} />
        <FilterInput name="bookId" value={f.bookId} label={t('bookId')} />
        <FilterInput name="authorId" value={f.authorId} label={t('authorId')} />
        {kind === 'comments' && (
          <FilterInput name="chapterId" value={f.chapterId} label={t('chapterId')} />
        )}
        <FilterSelect
          name="state"
          value={f.state}
          label={t('visibility')}
          options={['all', 'visible', 'hidden'].map((v) => [v, t(v)])}
        />
      </SearchForm>
      <p className="text-sm text-secondary-foreground" role="status">
        {page.total.toLocaleString()} {t('results')}
      </p>
      {page.entries.length ? (
        <div>
          {page.entries.map((entry) => (
            <Link to={`/${kind}/${entry.id}`} key={entry.id} className="ModWorkspace-record">
              <div className="flex gap-3 flex-wrap items-center">
                <StateBadge value={entry.visible ? 'visible' : 'hidden'} />
                <span className="text-sm text-secondary-foreground">{entry.id}</span>
              </div>
              <h2 className="font-semibold line-clamp-3">{entry.title}</h2>
              <p className="text-sm text-secondary-foreground">
                {entry.bookName} · {entry.authorName} · {date(entry.createdAt)}
              </p>
            </Link>
          ))}
        </div>
      ) : (
        <EmptyState>{t('noContent')}</EmptyState>
      )}
      <PageNavigation page={page.page} totalPages={page.totalPages} />
    </WorkspacePage>
  )
}
export const contentDetailLoader =
  (kind: Kind) =>
  async ({ request, params }: LoaderFunctionArgs) => {
    const id = params.contentId!
    const [entry, log, reports] = await Promise.all([
      getContent(kind, id, request.signal),
      moderationRead<ModerationAuditLogPageResponse>(
        'audit-log',
        { targetType: kind === 'chapters' ? 'chapter' : 'comment', targetId: id, pageSize: '20' },
        request.signal,
      ),
      moderationRead<ModerationReportsSearchResponse>(
        'reports',
        kind === 'chapters'
          ? { chapterId: id, status: 'all', pageSize: '10' }
          : { targetType: 'comment', targetId: id, status: 'all', pageSize: '10' },
        request.signal,
      ),
    ])
    return { entry, log, reports }
  }
export function ContentDetailPage({ kind }: { kind: Kind }) {
  const { entry, log, reports } =
    useLoaderData<Awaited<ReturnType<ReturnType<typeof contentDetailLoader>>>>()
  const [confirmation, setConfirmation] = useState<ModerationConfirmation>()
  const action = entry.visible ? (kind === 'chapters' ? 'hide' : 'remove') : 'restore'
  const reportQuery =
    kind === 'chapters' ? `chapterId=${entry.id}` : `targetType=comment&targetId=${entry.id}`
  const historyQuery = `targetType=${kind === 'chapters' ? 'chapter' : 'comment'}&targetId=${entry.id}`
  return (
    <WorkspacePage
      title={kind === 'chapters' ? entry.title : `${t('comment')} · ${entry.id}`}
      actions={
        <Link className="Btn Btn--outline" to={`/${kind}`}>
          {t(kind)}
        </Link>
      }
    >
      <div className="flex flex-wrap items-center gap-3">
        <StateBadge value={entry.visible ? 'visible' : 'hidden'} />
        <Link className="Link" to={`/books/${entry.bookId}`}>
          {entry.bookName}
        </Link>
        <Link className="Link" to={`/users/${entry.authorId}`}>
          {entry.authorName}
        </Link>
        {entry.chapterId && (
          <Link className="Link" to={`/chapters/${entry.chapterId}`}>
            {entry.chapterName}
          </Link>
        )}
      </div>
      <p className="text-sm text-secondary-foreground">
        {t('created')}: {date(entry.createdAt)}
        {entry.updatedAt && ` · ${t('updated')}: ${date(entry.updatedAt)}`}
      </p>
      <section>
        <h2 className="text-xl font-semibold mb-3">{t('currentContent')}</h2>
        <ContentPreview content={entry.content} />
      </section>
      {entry.parentId && (
        <Link className="Link" to={`/comments/${entry.parentId}`}>
          {t('parentComment')}
        </Link>
      )}
      {entry.replies.length > 0 && (
        <section>
          <h2 className="font-semibold">{t('replies')}</h2>
          <ul className="flex flex-wrap gap-4 mt-3">
            {entry.replies.map((id) => (
              <li key={id}>
                <Link className="Link" to={`/comments/${id}`}>
                  {t('comment')} · {id}
                </Link>
              </li>
            ))}
          </ul>
        </section>
      )}
      <ModerationActions.Root
        action={action}
        onActionChange={() => {}}
        options={[{ value: action, label: t(action) }]}
        selectLabel={t('action')}
        selectId="content-action"
        confirmation={confirmation}
        onConfirmationChange={setConfirmation}
        successMessage={t('saved')}
        confirmLabel={t('confirm')}
      >
        <ModerationActions.ReasonCard
          title={t(action)}
          description={t(`${action}Description`)}
          reasonLabel={t('reason')}
          submitLabel={t(action)}
          destructive={action !== 'restore'}
          onSubmit={(reason, onSuccess) =>
            setConfirmation({
              title: t(action),
              description: `${entry.title.slice(0, 160)} · ${entry.id}. ${t(`${action}Description`)}`,
              destructive: action !== 'restore',
              onSuccess,
              run: () =>
                kind === 'chapters'
                  ? performChapterModerationAction(entry.id, action as 'hide' | 'restore', reason)
                  : performCommentModerationAction(
                      entry.id,
                      action as 'remove' | 'restore',
                      reason,
                    ),
            })
          }
        />
      </ModerationActions.Root>
      <section>
        <div className="flex justify-between gap-4">
          <h2 className="text-xl font-semibold">{t('relatedReports')}</h2>
          <Link className="Link" to={`/reports?status=all&${reportQuery}`}>
            {t('viewAll')}
          </Link>
        </div>
        <ReportRecords reports={reports.entries} />
      </section>
      <section>
        <div className="flex justify-between gap-4 mb-4">
          <h2 className="text-xl font-semibold">{t('history')}</h2>
          <Link className="Link" to={`/audit-log?${historyQuery}`}>
            {t('viewAll')}
          </Link>
        </div>
        <ModerationLog entries={log.entries} emptyText={t('noActivity')} />
      </section>
    </WorkspacePage>
  )
}
// Text-only investigation avoids executing author markup inside the staff portal.
// DOMParser does not attach the supplied document to the live DOM.
function ContentPreview({ content }: { content: string }) {
  const doc = new DOMParser().parseFromString(content, 'text/html')
  doc.querySelectorAll('script, style').forEach((node) => node.remove())
  doc.querySelectorAll('p, div, br, h1, h2, h3, li').forEach((node) => node.append('\n'))
  return <div className="ModWorkspace-content">{doc.body.textContent || t('noContent')}</div>
}
