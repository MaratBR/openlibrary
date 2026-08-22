import { getModerationAuditLog } from '@/features/moderation/api'
import type { ModerationAuditLogPage } from '@/features/moderation/api'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { ErrorDisplay } from '@/components/error'
import { FormControl } from '@/components/FormControl'
import { LoaderFunctionArgs, useLoaderData, useNavigate, useRouteError } from 'react-router'
import { ModerationLog } from './ModerationLog'

export async function auditLogLoader({ request }: LoaderFunctionArgs) {
  const url = new URL(request.url)
  const targetType = url.searchParams.get('targetType') ?? ''
  const page = Math.max(1, Number(url.searchParams.get('page')) || 1)
  const response = await getModerationAuditLog(targetType, page)
  response.throwIfError()
  return { log: response.data, targetType }
}
export default function AuditLogPage() {
  const { log, targetType } = useLoaderData<{ log: ModerationAuditLogPage; targetType: string }>()
  const navigate = useNavigate()
  const change = (value: string) =>
    navigate({ search: value === 'all' ? '' : `?targetType=${value}` })
  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader title={window._('moderationPortal.auditLog')} />
      <div className="grid gap-5 max-w-360 mx-auto">
        <section className="card">
          <FormControl label={window._('moderationPortal.audit.targetType')}>
            <Select value={targetType || 'all'} onValueChange={change}>
              <SelectTrigger className="w-full sm:w-72">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">{window._('moderationPortal.audit.allTargets')}</SelectItem>
                {['user', 'book', 'chapter', 'comment'].map((type) => (
                  <SelectItem key={type} value={type}>
                    {window._(`moderationPortal.audit.${type}`)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </FormControl>
        </section>
        <ModerationLog
          entries={log.entries}
          emptyText={window._('moderationPortal.audit.empty')}
          page={log.page}
          totalPages={log.totalPages}
        />
      </div>
    </DashboardContent.Root>
  )
}
export function AuditLogErrorPage() {
  return (
    <DashboardContent.Root>
      <div className="card">
        <ErrorDisplay error={useRouteError()} />
      </div>
    </DashboardContent.Root>
  )
}
