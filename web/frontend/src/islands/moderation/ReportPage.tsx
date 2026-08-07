import { getModerationReport } from '@/api/moderation'
import type { ModerationReportDetail } from '@/api/moderation'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { ErrorDisplay } from '@/components/error'
import { OLAPIResponse } from '@/http-client'
import { LoaderFunctionArgs, NavLink, useLoaderData, useRouteError } from 'react-router'

const dateTimeFormatter = new Intl.DateTimeFormat(undefined, {
  dateStyle: 'medium',
  timeStyle: 'short',
})

async function responseData<T>(request: Promise<OLAPIResponse<T>>) {
  const response = await request
  response.throwIfError()
  return response.data
}

export async function reportRouteLoader({ params }: LoaderFunctionArgs) {
  if (!params.reportId) throw new Error(window._('moderationPortal.report.loadError'))
  return responseData(getModerationReport(params.reportId))
}

export default function ReportPage() {
  const report = useLoaderData<ModerationReportDetail>()
  const target = targetPresentation(report)

  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader
        title={
          <div className="flex items-center gap-3">
            <span>{window._('moderationPortal.report.title')}</span>
            <span className="font-mono text-sm font-normal text-secondary-foreground">
              {report.number}
            </span>
          </div>
        }
      >
        <span className="rounded-full bg-amber-500/12 text-amber-700 dark:text-amber-300 px-3 py-1 text-sm font-semibold capitalize">
          {report.status}
        </span>
      </DashboardContent.StickyHeader>

      <div className="w-full max-w-360 mx-auto px-4 md:px-8 py-7">
        <div className="grid xl:grid-cols-[minmax(0,1fr)_22rem] gap-5 items-start">
          <main className="grid gap-5">
            <section className="card card--elevated overflow-hidden">
              <div className="flex items-start gap-4 border-b border-border pb-5">
                <div className="w-11 h-11 rounded-xl bg-red-500/10 text-red-600 dark:text-red-400 grid place-items-center shrink-0">
                  <i className={`fa-solid ${target.icon}`} />
                </div>
                <div className="min-w-0">
                  <div className="flex flex-wrap items-center gap-2 mb-1">
                    <span className="text-xs font-semibold uppercase tracking-wider text-secondary-foreground">
                      {target.label}
                    </span>
                    <span className="rounded-full bg-red-500/10 text-red-600 dark:text-red-400 px-2 py-0.5 text-xs font-medium capitalize">
                      {report.priority} {window._('moderationPortal.report.priority')}
                    </span>
                  </div>
                  <h1 className="text-2xl font-semibold leading-tight">{report.reason}</h1>
                  <p className="text-sm text-secondary-foreground mt-2">
                    {window._('moderationPortal.report.reportedBy')}{' '}
                    <NavLink className="link" to={`/users/${report.reporterUserId}`}>
                      {report.reporterUserName ||
                        window._('moderationPortal.report.unknownReporter')}
                    </NavLink>{' '}
                    · {dateTimeFormatter.format(new Date(report.time))}
                  </p>
                </div>
              </div>
              <div className="pt-5">
                <h2 className="text-sm font-semibold uppercase tracking-wider text-secondary-foreground mb-2">
                  {window._('moderationPortal.report.description')}
                </h2>
                <p className="leading-7 whitespace-pre-wrap">
                  {report.description || window._('moderationPortal.report.noDescription')}
                </p>
              </div>
            </section>

            <section className="card">
              <div className="flex items-center justify-between gap-4 mb-4">
                <h2 className="text-xl font-semibold">
                  {window._('moderationPortal.report.reportedContent')}
                </h2>
                <NavLink className="btn btn--outline" to={target.to}>
                  {window._('moderationPortal.report.openTarget')}
                  <i className="fa-solid fa-arrow-up-right-from-square ml-2" />
                </NavLink>
              </div>
              <div className="rounded-xl border border-border bg-secondary/45 p-4 flex items-center gap-4">
                <div className="w-10 h-10 rounded-lg bg-surface grid place-items-center text-secondary-foreground">
                  <i className={`fa-solid ${target.icon}`} />
                </div>
                <div>
                  <div className="font-semibold">{target.label}</div>
                  <div className="font-mono text-sm text-secondary-foreground">
                    {report.targetId}
                  </div>
                </div>
              </div>
              <p className="text-sm text-secondary-foreground mt-3">{target.help}</p>
            </section>

            <section className="card">
              <h2 className="text-xl font-semibold mb-5">
                {window._('moderationPortal.report.activity')}
              </h2>
              <ol className="grid gap-0">
                {report.activities.map((activity, index) => (
                  <li
                    key={`${activity.time}-${activity.kind}`}
                    className="grid grid-cols-[2rem_1fr] gap-3"
                  >
                    <div className="flex flex-col items-center">
                      <span className="w-8 h-8 rounded-full bg-secondary grid place-items-center text-secondary-foreground">
                        <i className={`fa-solid ${activityIcon(activity.kind)} text-xs`} />
                      </span>
                      {index < report.activities.length - 1 && (
                        <span className="w-px flex-1 bg-border min-h-8" />
                      )}
                    </div>
                    <div className="pb-5">
                      <p>
                        <span className="font-semibold">
                          {activity.actor || window._('moderationPortal.report.system')}
                        </span>{' '}
                        {activity.description}
                      </p>
                      <time className="text-sm text-secondary-foreground">
                        {dateTimeFormatter.format(new Date(activity.time))}
                      </time>
                    </div>
                  </li>
                ))}
              </ol>
            </section>
          </main>

          <aside className="grid gap-5 xl:sticky xl:top-22">
            <section className="card">
              <h2 className="text-lg font-semibold mb-4">
                {window._('moderationPortal.report.details')}
              </h2>
              <dl className="grid gap-4">
                <Detail
                  label={window._('moderationPortal.report.status')}
                  value={capitalize(report.status)}
                />
                <Detail
                  label={window._('moderationPortal.report.priorityLabel')}
                  value={capitalize(report.priority)}
                />
                <Detail
                  label={window._('moderationPortal.report.assignedTo')}
                  value={report.assignedTo}
                  subvalue={report.assignedTeam}
                />
                <Detail
                  label={window._('moderationPortal.report.channel')}
                  value={report.channel}
                />
                <Detail
                  label={window._('moderationPortal.report.sla')}
                  value={dateTimeFormatter.format(new Date(report.slaDeadline))}
                />
              </dl>
            </section>
            <section className="card">
              <h2 className="text-lg font-semibold mb-3">
                {window._('moderationPortal.report.tags')}
              </h2>
              <div className="flex flex-wrap gap-2">
                {report.tags.map((tag) => (
                  <span key={tag} className="rounded-full bg-secondary px-3 py-1 text-sm">
                    {tag}
                  </span>
                ))}
              </div>
            </section>
            <section className="rounded-xl border border-dashed border-border p-4 text-sm text-secondary-foreground">
              <i className="fa-solid fa-flask mr-2" />
              {window._('moderationPortal.report.mockNotice')}
            </section>
          </aside>
        </div>
      </div>
    </DashboardContent.Root>
  )
}

function Detail({ label, value, subvalue }: { label: string; value: string; subvalue?: string }) {
  return (
    <div>
      <dt className="text-xs text-secondary-foreground mb-1">{label}</dt>
      <dd className="font-medium">{value}</dd>
      {subvalue && <dd className="text-sm text-secondary-foreground">{subvalue}</dd>}
    </div>
  )
}

function targetPresentation(report: ModerationReportDetail) {
  if (report.targetType === 'user')
    return {
      icon: 'fa-user',
      label: window._('moderationPortal.report.targetUser'),
      to: `/users/${report.targetId}`,
      help: window._('moderationPortal.report.userHelp'),
    }
  if (report.targetType === 'book')
    return {
      icon: 'fa-book',
      label: window._('moderationPortal.report.targetBook'),
      to: `/books/${report.targetId}`,
      help: window._('moderationPortal.report.bookHelp'),
    }
  return {
    icon: 'fa-comment',
    label: window._('moderationPortal.report.targetComment'),
    to: `/comments/${report.targetId}`,
    help: window._('moderationPortal.report.commentHelp'),
  }
}

function activityIcon(kind: string) {
  if (kind === 'assignment') return 'fa-user-check'
  if (kind === 'priority') return 'fa-arrow-up'
  return 'fa-flag'
}

function capitalize(value: string) {
  return value.charAt(0).toUpperCase() + value.slice(1)
}

export function ReportErrorPage() {
  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader title={window._('moderationPortal.report.title')} />
      <div className="p-8">
        <div className="card">
          <ErrorDisplay error={useRouteError()} />
        </div>
      </div>
    </DashboardContent.Root>
  )
}
