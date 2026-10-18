import { Link, LoaderFunctionArgs, useLoaderData } from 'react-router'
import { getWorkspaceOverview, moderationRead } from '../../api/workspace'
import type {
  ModerationReportsSearchResponse,
  ModerationAuditLogPageResponse,
} from '@/backend-types'
import { WorkspacePage, t, date, EmptyState } from './ui'
import { ReportRecords } from './ReportRecords'
import { ModerationLog } from '../ModerationLog'

export async function overviewLoader({ request }: LoaderFunctionArgs) {
  const [overview, reports, log] = await Promise.all([
    getWorkspaceOverview(request.signal),
    moderationRead<ModerationReportsSearchResponse>(
      'reports',
      { status: 'pending', sort: 'priority', pageSize: '5' },
      request.signal,
    ),
    moderationRead<ModerationAuditLogPageResponse>(
      'audit-log',
      { action: 'report_decision', pageSize: '5' },
      request.signal,
    ),
  ])
  return { overview, reports, log }
}
export default function OverviewPage() {
  const { overview, reports, log } = useLoaderData<typeof overviewLoader>()
  return (
    <WorkspacePage
      title={t('overview')}
      actions={
        <Link className="Btn Btn--primary" to="/reports">
          {t('openQueue')}
        </Link>
      }
    >
      <div className="grid lg:grid-cols-[minmax(0,2fr)_minmax(15rem,1fr)] gap-8 items-start">
        <section>
          <h2 className="text-xl font-semibold">{t('pendingWork')}</h2>
          <p className="text-secondary-foreground mt-2 mb-4">
            {overview.oldestPending
              ? `${t('oldestPending')}: ${date(overview.oldestPending)}`
              : t('queueClear')}
          </p>
          <ReportRecords reports={reports.entries} />
        </section>
        <section>
          <h2 className="text-xl font-semibold">{t('queues')}</h2>
          {(['unreviewed', 'escalated', 'resolved'] as const).map((status) => (
            <Link className="ModWorkspace-queue" key={status} to={`/reports?status=${status}`}>
              <span>{t(status)}</span>
              <strong>{overview[status].toLocaleString()}</strong>
            </Link>
          ))}
          <p className="text-sm text-secondary-foreground mt-5">{t('manualReview')}</p>
          <div className="flex flex-wrap gap-4 mt-3">
            <Link className="Link" to="/triage">
              {t('triage')}
            </Link>
            <Link className="Link" to="/notifications">
              {t('notifications')}
            </Link>
          </div>
        </section>
      </div>
      <section>
        <h2 className="text-xl font-semibold mb-4">{t('recentDecisions')}</h2>
        {log.entries.length ? (
          <ModerationLog entries={log.entries} emptyText={t('noActivity')} />
        ) : (
          <EmptyState>{t('noActivity')}</EmptyState>
        )}
      </section>
    </WorkspacePage>
  )
}
export function IntegrationPage({ feature }: { feature: 'triage' | 'notifications' }) {
  return (
    <WorkspacePage title={t(feature)}>
      <section className="Card Card--dashboard">
        <span className="ModWorkspace-status">{t('notAvailable')}</span>
        <h2 className="text-xl font-semibold mt-4">{t(`${feature}Title`)}</h2>
        <p className="text-secondary-foreground mt-3 max-w-160">{t(`${feature}Description`)}</p>
        <Link className="Btn Btn--primary mt-6" to="/reports">
          {t('openQueue')}
        </Link>
      </section>
    </WorkspacePage>
  )
}
