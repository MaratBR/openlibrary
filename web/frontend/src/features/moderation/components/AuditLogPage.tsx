import { LoaderFunctionArgs, useLoaderData } from 'react-router'
import type { ModerationAuditLogPageResponse } from '@/backend-types'
import { moderationRead } from '../api/workspace'
import { ModerationLog } from './ModerationLog'
import {
  WorkspacePage,
  SearchForm,
  FilterInput,
  FilterSelect,
  PageNavigation,
  filtersFrom,
  t,
  targetOptions,
  WorkspaceError,
} from './workspace/ui'
export async function auditLogLoader({ request }: LoaderFunctionArgs) {
  const filters = filtersFrom(request)
  return {
    filters,
    log: await moderationRead<ModerationAuditLogPageResponse>('audit-log', filters, request.signal),
  }
}
export default function AuditLogPage() {
  const { log, filters: f } = useLoaderData<typeof auditLogLoader>()
  return (
    <WorkspacePage title={t('auditLog')}>
      <SearchForm>
        <FilterSelect
          name="targetType"
          value={f.targetType}
          label={t('targetType')}
          options={targetOptions()}
        />
        <FilterInput name="targetId" value={f.targetId} label={t('targetId')} />
        <FilterInput name="actorId" value={f.actorId} label={t('actorId')} />
        <FilterInput name="action" value={f.action} label={t('actionCode')} />
        <FilterInput name="dateFrom" value={f.dateFrom} label={t('from')} type="date" />
        <FilterInput name="dateTo" value={f.dateTo} label={t('to')} type="date" />
      </SearchForm>
      <p className="text-sm text-secondary-foreground">
        {log.total.toLocaleString()} {t('results')}
      </p>
      <ModerationLog entries={log.entries} emptyText={t('noActivity')} />
      <PageNavigation page={log.page} totalPages={log.totalPages} />
    </WorkspacePage>
  )
}
export const AuditLogErrorPage = WorkspaceError
