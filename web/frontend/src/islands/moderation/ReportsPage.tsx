import { searchModerationReports } from '@/api/moderation'
import type { ModerationReportsSearch } from '@/api/moderation'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { ErrorDisplay } from '@/components/error'
import { Pagination } from '@/components/Pagination'
import { OLAPIResponse } from '@/http-client'
import { useState } from 'react'
import { Form, LoaderFunctionArgs, NavLink, useLoaderData, useRouteError } from 'react-router'

const dateTimeFormatter = new Intl.DateTimeFormat(undefined, {
  dateStyle: 'medium',
  timeStyle: 'short',
})

const ALL_TARGET_TYPES = 'all'

async function responseData<T>(request: Promise<OLAPIResponse<T>>) {
  const response = await request
  response.throwIfError()
  return response.data
}

export async function reportsSearchLoader({ request }: LoaderFunctionArgs) {
  const params = new URL(request.url).searchParams
  const search = params.get('search') ?? ''
  const targetType = params.get('targetType') ?? ''
  const page = Math.max(1, Number(params.get('page')) || 1)
  const reports = await responseData(searchModerationReports(search, targetType, page))
  return { reports, filters: { search, targetType } }
}

export default function ReportsPage({ view }: { view: 'overview' | 'search' }) {
  const data = useLoaderData() as Awaited<ReturnType<typeof reportsSearchLoader>> | undefined

  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader title={window._('moderationPortal.reports')} />
      <div className="w-full max-w-360 mx-auto px-4 md:px-8 py-7">
        <nav aria-label={window._('moderationPortal.reportsPage.navigation')}>
          <ul className="tabs tabs--primary mb-6">
            <li>
              <NavLink
                className={({ isActive }) => `tab ${isActive ? 'tab--active' : ''}`}
                end
                to="/reports"
              >
                {window._('moderationPortal.reportsPage.overview')}
              </NavLink>
            </li>
            <li>
              <NavLink
                className={({ isActive }) => `tab ${isActive ? 'tab--active' : ''}`}
                to="/reports/search"
              >
                {window._('moderationPortal.reportsPage.search')}
              </NavLink>
            </li>
          </ul>
        </nav>
        {view === 'overview' ? <ReportsOverview /> : data && <ReportsSearch data={data} />}
      </div>
    </DashboardContent.Root>
  )
}

function ReportsOverview() {
  // TODO: Add report queues and summary metrics to the overview.
  return null
}

function ReportsSearch({
  data,
}: {
  data: { reports: ModerationReportsSearch; filters: { search: string; targetType: string } }
}) {
  const { reports, filters } = data
  return (
    <div>
      <DashboardContent.Card className="card--elevated p-4 md:p-5 mb-6">
        <Form
          method="get"
          className="grid md:grid-cols-[minmax(18rem,1fr)_14rem_auto] gap-3 items-end"
        >
          <label>
            <span className="block text-sm font-medium mb-1.5">
              {window._('moderationPortal.reportsPage.searchLabel')}
            </span>
            <div className="relative">
              <i className="fa-solid fa-magnifying-glass absolute left-3 top-1/2 -translate-y-1/2 text-secondary-foreground" />
              <input
                className="input w-full pl-10"
                name="search"
                defaultValue={filters.search}
                placeholder={window._('moderationPortal.reportsPage.searchPlaceholder')}
              />
            </div>
          </label>
          <div>
            <label className="block text-sm font-medium mb-1.5" htmlFor="reports-target-type">
              {window._('moderationPortal.reportsPage.targetType')}
            </label>
            <TargetTypeSelect key={filters.targetType} targetType={filters.targetType} />
          </div>
          <button className="btn btn--default h-10 px-5" type="submit">
            {window._('moderationPortal.reportsPage.apply')}
          </button>
        </Form>
      </DashboardContent.Card>

      <div className="flex items-baseline justify-between mb-3 px-1">
        <h2 className="text-lg font-semibold">
          {window._('moderationPortal.reportsPage.results')}
        </h2>
        <span className="text-sm text-secondary-foreground">
          {reports.total} {window._('moderationPortal.reportsPage.found')}
        </span>
      </div>
      {reports.entries.length ? (
        <div className="grid gap-3">
          {reports.entries.map((report) => (
            <NavLink
              key={report.id}
              to={`/reports/${report.id}`}
              className="card card--elevated group block p-4 md:p-5 hover:bg-foreground/5"
            >
              <div className="flex items-start gap-4">
                <div className="w-10 h-10 rounded-xl bg-red-500/10 text-red-600 dark:text-red-400 grid place-items-center shrink-0">
                  <i className={`fa-solid ${targetIcon(report.targetType)}`} />
                </div>
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-x-3 gap-y-1">
                    <span className="font-mono text-sm text-secondary-foreground">
                      {report.number}
                    </span>
                    <span className="text-xs uppercase tracking-wider text-secondary-foreground">
                      {report.targetType}
                    </span>
                  </div>
                  <h3 className="font-semibold text-lg mt-1 group-hover:text-primary">
                    {report.reason}
                  </h3>
                  <p className="text-sm text-secondary-foreground mt-1 line-clamp-2">
                    {report.description || window._('moderationPortal.report.noDescription')}
                  </p>
                  <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs text-secondary-foreground mt-3">
                    <span>{dateTimeFormatter.format(new Date(report.time))}</span>
                    <span>
                      {window._('moderationPortal.reportsPage.target')} {report.targetId}
                    </span>
                    <span>
                      {report.reporterUserName ||
                        window._('moderationPortal.report.unknownReporter')}
                    </span>
                  </div>
                </div>
                <i className="fa-solid fa-chevron-right text-secondary-foreground group-hover:text-primary mt-3" />
              </div>
            </NavLink>
          ))}
        </div>
      ) : (
        <DashboardContent.Card className="py-14 text-center text-secondary-foreground">
          {window._('moderationPortal.reportsPage.empty')}
        </DashboardContent.Card>
      )}
      {reports.totalPages > 1 && (
        <div className="flex justify-center mt-7">
          <Pagination.Facade
            page={reports.page}
            totalPages={reports.totalPages}
            size={7}
            getTo={(page) => ({
              search: new URLSearchParams({ ...filters, page: String(page) }).toString(),
            })}
          />
        </div>
      )}
    </div>
  )
}

function TargetTypeSelect({ targetType }: { targetType: string }) {
  const [value, setValue] = useState(targetType || ALL_TARGET_TYPES)

  return (
    <>
      <input type="hidden" name="targetType" value={value === ALL_TARGET_TYPES ? '' : value} />
      <Select value={value} onValueChange={setValue}>
        <SelectTrigger id="reports-target-type" className="w-full">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value={ALL_TARGET_TYPES}>
            {window._('moderationPortal.reportsPage.allTargets')}
          </SelectItem>
          <SelectItem value="user">{window._('moderationPortal.report.targetUser')}</SelectItem>
          <SelectItem value="book">{window._('moderationPortal.report.targetBook')}</SelectItem>
          <SelectItem value="comment">
            {window._('moderationPortal.report.targetComment')}
          </SelectItem>
        </SelectContent>
      </Select>
    </>
  )
}

function targetIcon(targetType: string) {
  if (targetType === 'user') return 'fa-user'
  if (targetType === 'book') return 'fa-book'
  return 'fa-comment'
}

export function ReportsErrorPage() {
  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader title={window._('moderationPortal.reports')} />
      <div className="p-8">
        <div className="card">
          <ErrorDisplay error={useRouteError()} />
        </div>
      </div>
    </DashboardContent.Root>
  )
}
