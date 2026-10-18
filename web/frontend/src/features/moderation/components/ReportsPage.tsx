import { LoaderFunctionArgs, useLoaderData } from 'react-router'
import type { ModerationReportsSearchResponse } from '@/backend-types'
import { moderationRead } from '../api/workspace'
import {
  WorkspacePage,
  SearchForm,
  FilterInput,
  FilterSelect,
  PageNavigation,
  filtersFrom,
  t,
  WorkspaceError,
} from './workspace/ui'
import { ReportRecords } from './workspace/ReportRecords'
export async function reportsSearchLoader({ request }: LoaderFunctionArgs) {
  const filters = filtersFrom(request)
  const reports = await moderationRead<ModerationReportsSearchResponse>(
    'reports',
    filters,
    request.signal,
  )
  return { reports, filters }
}
export default function ReportsPage() {
  const { reports, filters: f } = useLoaderData<typeof reportsSearchLoader>()
  return (
    <WorkspacePage title={t('reports')}>
      <SearchForm>
        <FilterInput name="search" value={f.search} label={t('reportSearch')} />
        <FilterSelect
          name="status"
          value={f.status || 'pending'}
          label={t('status')}
          options={['pending', 'unreviewed', 'escalated', 'resolved', 'all'].map((v) => [v, t(v)])}
        />
        <FilterSelect
          name="priority"
          value={f.priority}
          label={t('priority')}
          options={['all', 'critical', 'high', 'medium', 'low'].map((v) => [v, t(v)])}
        />
        <FilterSelect
          name="targetType"
          value={f.targetType}
          label={t('targetType')}
          options={['all', 'user', 'book', 'comment'].map((v) => [v, t(v)])}
        />
        <FilterInput name="targetId" value={f.targetId} label={t('targetId')} />
        {f.chapterId && <input type="hidden" name="chapterId" value={f.chapterId} />}
        <FilterInput name="dateFrom" type="date" value={f.dateFrom} label={t('from')} />
        <FilterInput name="dateTo" type="date" value={f.dateTo} label={t('to')} />
        <FilterSelect
          name="sort"
          value={f.sort || 'priority'}
          label={t('sort')}
          options={['priority', 'oldest', 'newest'].map((v) => [v, t(`sort_${v}`)])}
        />
      </SearchForm>
      <p className="text-sm text-secondary-foreground" role="status">
        {reports.total.toLocaleString()} {t('results')}
      </p>
      <ReportRecords reports={reports.entries} />
      <PageNavigation page={reports.page} totalPages={reports.totalPages} />
    </WorkspacePage>
  )
}
export const ReportsErrorPage = WorkspaceError
