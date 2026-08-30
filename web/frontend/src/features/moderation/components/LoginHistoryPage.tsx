import { getGlobalLoginHistory, searchModerationUsers } from '@/features/moderation/api'
import type { LoginHistoryEntry, ModerationUserListEntry } from '@/features/moderation/api'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { ErrorDisplay } from '@/components/error'
import { FormControl } from '@/components/FormControl'
import { Pagination } from '@/components/Pagination'
import { OLAPIResponse } from '@/features/http-client'
import { useEffect, useState } from 'react'
import { Form, Link, LoaderFunctionArgs, useLoaderData, useRouteError } from 'react-router'
import { UserMultiSelect } from './UserMultiSelect'

const dateTimeFormatter = new Intl.DateTimeFormat(undefined, {
  dateStyle: 'medium',
  timeStyle: 'short',
})

async function responseData<T>(request: Promise<OLAPIResponse<T>>) {
  const response = await request
  response.throwIfError()
  return response.data
}

async function resolveUsers(ids: string[]) {
  const results = await Promise.all(
    ids.map(async (id) => {
      const page = await responseData(searchModerationUsers(id, '', '', 1, 1))
      return page.entries.find((user) => user.id === id)
    }),
  )
  return results.filter((user): user is ModerationUserListEntry => Boolean(user))
}

export async function loginHistoryLoader({ request }: LoaderFunctionArgs) {
  const params = new URL(request.url).searchParams
  const userIDs = (params.get('users') ?? '')
    .split(',')
    .map((id) => id.trim())
    .filter(Boolean)
    .slice(0, 20)
  const filters = {
    users: userIDs.join(','),
    search: params.get('search') ?? '',
    status: params.get('status') ?? '',
    dateFrom: params.get('dateFrom') ?? '',
    dateTo: params.get('dateTo') ?? '',
  }
  const page = Math.max(1, Number(params.get('page')) || 1)
  const [history, selectedUsers] = await Promise.all([
    responseData(getGlobalLoginHistory(page, 20, filters)),
    resolveUsers(userIDs),
  ])
  return { history, selectedUsers, filters }
}

export default function LoginHistoryPage() {
  const {
    history,
    selectedUsers: loadedUsers,
    filters,
  } = useLoaderData<Awaited<ReturnType<typeof loginHistoryLoader>>>()
  const [selectedUsers, setSelectedUsers] = useState(loadedUsers)
  const [status, setStatus] = useState(filters.status || 'all')
  useEffect(() => {
    setSelectedUsers(loadedUsers)
    setStatus(filters.status || 'all')
  }, [loadedUsers, filters.status])
  const paginationFilters = filters
  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader title={window._('moderationPortal.loginHistory')} />
      <div className="w-full max-w-360 mx-auto px-4 md:px-8 py-7">
        <p className="text-secondary-foreground mb-6">
          {window._('moderationPortal.loginHistoryPage.description')}
        </p>
        <section className="Card Card--elevated mb-5">
          <Form method="get" className="grid gap-4">
            <input
              type="hidden"
              name="users"
              value={selectedUsers.map((user) => user.id).join(',')}
            />
            <FormControl label={window._('moderationPortal.loginHistoryPage.users')}>
              <UserMultiSelect value={selectedUsers} onChange={setSelectedUsers} />
            </FormControl>
            <div className="grid md:grid-cols-[minmax(14rem,1fr)_11rem_10rem_10rem_auto] gap-3 items-end">
              <FormControl label={window._('moderationPortal.user.loginSearch')}>
                <input
                  className="input"
                  name="search"
                  defaultValue={filters.search}
                  placeholder={window._('moderationPortal.user.loginSearchPlaceholder')}
                />
              </FormControl>
              <FormControl label={window._('moderationPortal.user.sessionStatus')}>
                <input type="hidden" name="status" value={status === 'all' ? '' : status} />
                <Select value={status} onValueChange={setStatus}>
                  <SelectTrigger className="w-full">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">
                      {window._('moderationPortal.user.allSessions')}
                    </SelectItem>
                    <SelectItem value="active">
                      {window._('moderationPortal.user.activeSessions')}
                    </SelectItem>
                    <SelectItem value="expired">
                      {window._('moderationPortal.user.expiredSessions')}
                    </SelectItem>
                    <SelectItem value="terminated">
                      {window._('moderationPortal.user.terminatedSessions')}
                    </SelectItem>
                  </SelectContent>
                </Select>
              </FormControl>
              <FormControl label={window._('moderationPortal.user.dateFrom')}>
                <input
                  className="input"
                  type="date"
                  name="dateFrom"
                  defaultValue={filters.dateFrom}
                />
              </FormControl>
              <FormControl label={window._('moderationPortal.user.dateTo')}>
                <input className="input" type="date" name="dateTo" defaultValue={filters.dateTo} />
              </FormControl>
              <button className="Btn Btn--primary h-10">
                {window._('moderationPortal.user.applyFilters')}
              </button>
            </div>
          </Form>
        </section>
        <div className="flex items-baseline justify-between gap-4 mb-3 px-1">
          <h2 className="text-lg font-semibold">{window._('moderationPortal.loginHistory')}</h2>
          <span className="text-sm text-secondary-foreground">
            {history.total} {window._('moderationPortal.loginHistoryPage.sessionsFound')}
          </span>
        </div>
        {history.entries.length ? (
          <section className="Card Card--nopad overflow-hidden">
            {history.entries.map((entry) => (
              <LoginRow key={entry.id} entry={entry} />
            ))}
          </section>
        ) : (
          <section className="Card py-16 text-center text-secondary-foreground">
            {window._('moderationPortal.user.noLoginHistory')}
          </section>
        )}
        {history.totalPages > 1 && (
          <div className="flex justify-center mt-7">
            <Pagination.Facade
              page={history.page}
              totalPages={history.totalPages}
              size={7}
              getTo={(page) => ({
                search: new URLSearchParams({
                  ...paginationFilters,
                  page: String(page),
                }).toString(),
              })}
            />
          </div>
        )}
      </div>
    </DashboardContent.Root>
  )
}

function LoginRow({ entry }: { entry: LoginHistoryEntry }) {
  const expired = new Date(entry.expiresAt) <= new Date()
  const state = entry.isTerminated
    ? window._('moderationPortal.user.terminated')
    : expired
      ? window._('moderationPortal.user.expired')
      : window._('moderationPortal.user.sessionActive')
  return (
    <div className="p-4 border-b border-border last:border-0 grid lg:grid-cols-[minmax(10rem,1fr)_11rem_10rem_minmax(10rem,1fr)_minmax(12rem,1fr)_7rem] gap-2 items-center">
      <Link className="Link truncate" to={`/users/${entry.userId}`}>
        {entry.userName}
      </Link>
      <span>{dateTimeFormatter.format(new Date(entry.loggedInAt))}</span>
      <span className="font-mono text-sm">{entry.ipAddress}</span>
      <span className="text-sm">
        {[entry.city, entry.region, entry.country].filter(Boolean).join(', ') ||
          window._('moderationPortal.user.unknownLocation')}
      </span>
      <span className="truncate" title={entry.userAgent}>
        {entry.userAgent}
      </span>
      <span
        className={`chip justify-self-start ${!entry.isTerminated && !expired ? 'Chip--primary' : ''}`}
      >
        {state}
      </span>
    </div>
  )
}

export function LoginHistoryErrorPage() {
  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader title={window._('moderationPortal.loginHistory')} />
      <div className="p-8">
        <div className="Card">
          <ErrorDisplay error={useRouteError()} />
        </div>
      </div>
    </DashboardContent.Root>
  )
}
