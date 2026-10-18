import { Link } from 'react-router'
import type { ModerationReportListEntryResponse } from '@/backend-types'
import { t, date, EmptyState, StateBadge } from './ui'
export function ReportRecords({ reports }: { reports: ModerationReportListEntryResponse[] }) {
  return reports.length ? (
    <div>
      {reports.map((report) => (
        <Link key={report.id} to={`/reports/${report.id}`} className="ModWorkspace-record">
          <div className="flex flex-wrap items-center gap-2">
            <span className="text-sm text-secondary-foreground">{report.number}</span>
            <StateBadge value={report.status} />
            <StateBadge value={report.priority} />
          </div>
          <h3 className="font-semibold">{report.reason}</h3>
          <p className="text-sm text-secondary-foreground line-clamp-2">
            {report.description || t('noDescription')}
          </p>
          <div className="flex flex-wrap gap-x-4 gap-y-1 text-sm text-secondary-foreground">
            <span>
              {t(report.targetType)} · {report.targetId}
            </span>
            <span>{report.reporterUserName}</span>
            <time dateTime={report.time}>{date(report.time)}</time>
          </div>
        </Link>
      ))}
    </div>
  ) : (
    <EmptyState>{t('noReports')}</EmptyState>
  )
}
