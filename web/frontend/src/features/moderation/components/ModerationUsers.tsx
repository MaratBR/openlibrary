import { searchModerationUsers } from '@/features/moderation/api'
import type { ModerationUserListEntry } from '@/features/moderation/api'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { ErrorDisplay } from '@/components/error'
import { Pagination } from '@/components/Pagination'
import { OLAPIResponse } from '@/http-client'
import { useEffect, useState } from 'react'
import { Form, Link, LoaderFunctionArgs, useLoaderData, useRouteError } from 'react-router'

const dateFormatter = new Intl.DateTimeFormat(undefined, { dateStyle: 'medium' })
const dateTimeFormatter = new Intl.DateTimeFormat(undefined, {
  dateStyle: 'medium',
  timeStyle: 'short',
})

const EMPTY_FILTER_VALUE = 'all'

async function responseData<T>(request: Promise<OLAPIResponse<T>>) {
  const response = await request
  response.throwIfError()
  return response.data
}

export async function moderationUsersLoader({ request }: LoaderFunctionArgs) {
  const params = new URL(request.url).searchParams
  const search = params.get('search') ?? ''
  const banned = params.get('banned') ?? ''
  const role = params.get('role') ?? ''
  const page = Math.max(1, Number(params.get('page')) || 1)
  const users = await responseData(searchModerationUsers(search, banned, role, page))
  return { users, filters: { search, banned, role } }
}

export default function ModerationUsers({ roles }: { roles: string[] }) {
  const { users, filters } = useLoaderData<Awaited<ReturnType<typeof moderationUsersLoader>>>()

  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader title={window._('moderationPortal.users')} />
      <div className="w-full max-w-360 mx-auto px-4 md:px-8 py-7">
        <header className="mb-6">
          <p className="text-secondary-foreground max-w-180">
            {window._('moderationPortal.usersPage.description')}
          </p>
        </header>

        <DashboardContent.Card className="card--elevated p-4 md:p-5 mb-6">
          <Form
            method="get"
            className="grid grid-cols-1 md:grid-cols-[minmax(18rem,1fr)_12rem_12rem_auto] gap-3 items-end"
          >
            <label className="block">
              <span className="block text-sm font-medium mb-1.5">
                {window._('moderationPortal.usersPage.search')}
              </span>
              <div className="relative">
                <i className="fa-solid fa-magnifying-glass absolute left-3 top-1/2 -translate-y-1/2 text-secondary-foreground" />
                <input
                  className="input w-full pl-10"
                  name="search"
                  defaultValue={filters.search}
                  placeholder={window._('moderationPortal.usersPage.searchPlaceholder')}
                />
              </div>
            </label>
            <FilterSelect
              name="banned"
              label={window._('moderationPortal.usersPage.status')}
              value={filters.banned}
              options={[
                ['', window._('moderationPortal.usersPage.allStatuses')],
                ['active', window._('moderationPortal.usersPage.active')],
                ['banned', window._('moderationPortal.usersPage.banned')],
              ]}
            />
            <FilterSelect
              name="role"
              label={window._('moderationPortal.usersPage.role')}
              value={filters.role}
              options={[
                ['', window._('moderationPortal.usersPage.allRoles')],
                ...roles.map((role): [string, string] => [role, roleLabel(role)]),
              ]}
            />
            <button className="btn btn--primary h-10 px-5" type="submit">
              {window._('moderationPortal.usersPage.apply')}
            </button>
          </Form>
        </DashboardContent.Card>

        <div className="flex items-baseline justify-between gap-4 mb-3 px-1">
          <h2 className="text-lg font-semibold">
            {window._('moderationPortal.usersPage.results')}
          </h2>
          <span className="text-sm text-secondary-foreground">
            {users.total} {window._('moderationPortal.usersPage.usersFound')}
          </span>
        </div>

        {users.entries.length ? (
          <div className="grid gap-3">
            {users.entries.map((user) => (
              <UserRow key={user.id} user={user} />
            ))}
          </div>
        ) : (
          <DashboardContent.Card className="py-16 px-6 text-center">
            <i className="fa-solid fa-user-slash text-3xl text-secondary-foreground mb-4" />
            <h2 className="font-semibold text-lg">
              {window._('moderationPortal.usersPage.empty')}
            </h2>
            <p className="text-secondary-foreground mt-1">
              {window._('moderationPortal.usersPage.emptyDescription')}
            </p>
          </DashboardContent.Card>
        )}

        {users.totalPages > 1 && (
          <div className="flex justify-center mt-7">
            <Pagination.Facade
              page={users.page}
              totalPages={users.totalPages}
              size={7}
              getTo={(page) => ({
                search: new URLSearchParams({ ...filters, page: String(page) }).toString(),
              })}
            />
          </div>
        )}
      </div>
    </DashboardContent.Root>
  )
}

function FilterSelect({
  name,
  label,
  value,
  options,
}: {
  name: string
  label: string
  value: string
  options: [string, string][]
}) {
  const [selectedValue, setSelectedValue] = useState(value || EMPTY_FILTER_VALUE)

  useEffect(() => setSelectedValue(value || EMPTY_FILTER_VALUE), [value])

  return (
    <div className="block">
      <label className="block text-sm font-medium mb-1.5" htmlFor={`users-filter-${name}`}>
        {label}
      </label>
      <input
        type="hidden"
        name={name}
        value={selectedValue === EMPTY_FILTER_VALUE ? '' : selectedValue}
      />
      <Select value={selectedValue} onValueChange={setSelectedValue}>
        <SelectTrigger id={`users-filter-${name}`} className="w-full">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {options.map(([key, text]) => (
            <SelectItem key={key || EMPTY_FILTER_VALUE} value={key || EMPTY_FILTER_VALUE}>
              {text}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  )
}

function UserRow({ user }: { user: ModerationUserListEntry }) {
  return (
    <Link
      to={`/users/${user.id}`}
      className="card card--elevated group block p-4 md:p-5 transition-colors hover:bg-foreground/5"
    >
      <div className="flex flex-col lg:flex-row lg:items-center gap-4 lg:gap-7">
        <div className="flex items-center gap-4 min-w-0 lg:w-72">
          <img
            src={user.avatar}
            alt=""
            className="w-13 h-13 rounded-full object-cover bg-secondary ring-1 ring-border"
          />
          <div className="min-w-0">
            <div className="flex items-center gap-2">
              <h3 className="font-semibold truncate group-hover:text-primary">{user.name}</h3>
              {user.isBanned && (
                <span className="rounded-full bg-red-500/10 text-red-600 dark:text-red-400 px-2 py-0.5 text-xs font-medium">
                  {window._('moderationPortal.usersPage.banned')}
                </span>
              )}
            </div>
            <p className="text-xs text-secondary-foreground truncate mt-1">{user.id}</p>
          </div>
        </div>
        <dl className="grid grid-cols-2 md:grid-cols-3 gap-x-7 gap-y-3 flex-1">
          <Fact label={window._('moderationPortal.usersPage.role')} value={roleLabel(user.role)} />
          <Fact
            label={window._('moderationPortal.usersPage.lastVisit')}
            value={
              user.lastVisitAt
                ? dateTimeFormatter.format(new Date(user.lastVisitAt))
                : window._('moderationPortal.usersPage.never')
            }
          />
          <Fact
            label={window._('moderationPortal.usersPage.joined')}
            value={dateFormatter.format(new Date(user.joinedAt))}
          />
        </dl>
        {user.isBanned && (
          <div className="lg:w-72 lg:border-l border-border lg:pl-6">
            <p className="text-xs font-medium uppercase tracking-wide text-red-600 dark:text-red-400">
              {user.bannedAt
                ? window
                    ._('moderationPortal.usersPage.bannedOn')
                    .replace('{{date}}', dateFormatter.format(new Date(user.bannedAt)))
                : window._('moderationPortal.usersPage.banned')}
            </p>
            <p className="text-sm mt-1 line-clamp-2">
              {user.banReason || window._('moderationPortal.usersPage.noBanReason')}
            </p>
          </div>
        )}
        <i className="fa-solid fa-chevron-right hidden lg:block text-secondary-foreground group-hover:text-primary" />
      </div>
    </Link>
  )
}

function Fact({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-xs text-secondary-foreground mb-1">{label}</dt>
      <dd className="text-sm font-medium">{value}</dd>
    </div>
  )
}

function roleLabel(role: string) {
  return window._(`moderationPortal.usersPage.role${role.charAt(0).toUpperCase()}${role.slice(1)}`)
}

export function ModerationUsersErrorPage() {
  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader title={window._('moderationPortal.users')} />
      <div className="p-8">
        <div className="card">
          <ErrorDisplay error={useRouteError()} />
        </div>
      </div>
    </DashboardContent.Root>
  )
}
