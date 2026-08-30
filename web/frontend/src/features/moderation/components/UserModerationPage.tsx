import { ErrorDisplay } from '@/components/error'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { FormControl } from '@/components/FormControl'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components'
import {
  banUser,
  changeUserAbout,
  getModerationUser,
  getModerationUserBooks,
  getModerationUserComments,
  getModerationUserHistory,
  getModerationUserReports,
  getUserLoginHistory,
  getUserLoginLocations,
  permanentlyBanUser,
  renameUser,
  unbanUser,
} from '@/features/moderation/api'
import type {
  LoginHistoryEntry,
  LoginLocation,
  ModerationUser,
  ModerationUserBooksPage,
  ModerationUserCommentsPage,
  ModerationUserHistoryPage,
  ModerationUserReportsPage,
} from '@/features/moderation/api'
import { Pagination } from '@/components/Pagination'
import { FormEvent, ReactNode, useEffect, useState } from 'react'
import {
  LoaderFunctionArgs,
  NavLink,
  useLoaderData,
  useLocation,
  useRouteError,
} from 'react-router'
import {
  ModerationActionCard,
  ModerationActionsPanel,
  ModerationConfirmation,
  ModerationReasonActionCard,
  ModerationValueActionCard,
} from './ModerationActions'
import { ModerationLog, ModerationLogEntry } from './ModerationLog'
import SanitizeHTML from '@/common/SanitizeHTML'
import { queryClient } from '@/react/queryCache'
import { OLAPIResponse } from '@/features/http-client'

const dateFormatter = new Intl.DateTimeFormat(undefined, { dateStyle: 'medium' })
const dateTimeFormatter = new Intl.DateTimeFormat(undefined, {
  dateStyle: 'medium',
  timeStyle: 'short',
})

type PageResult<T> = {
  entries: T[]
  page: number
  pageSize: number
  total: number
  totalPages: number
}

async function responseData<T>(request: Promise<OLAPIResponse<T>>) {
  const response = await request
  response.throwIfError()
  return response.data
}

export const userModerationRouteLoader = async ({ params, request }: LoaderFunctionArgs) => {
  const userId = params.userId
  if (!userId) throw new Error(window._('moderationPortal.user.loadError'))

  const url = new URL(request.url)
  const leaf = url.pathname.split('/').at(-1) ?? ''
  const page = Math.max(1, Number(url.searchParams.get('page')) || 1)
  const user = await responseData(
    queryClient.fetchQuery({
      queryFn: () => getModerationUser(userId),
      queryKey: ['mod', 'user', userId],
      staleTime: 5000,
    }),
  )

  let books: ModerationUserBooksPage | undefined
  let comments: ModerationUserCommentsPage | undefined
  let history: ModerationUserHistoryPage | undefined
  let reports: ModerationUserReportsPage | undefined
  let logins: PageResult<LoginHistoryEntry> | undefined
  let locations: LoginLocation[] | undefined

  if (leaf === userId) {
    ;[books, comments, locations] = await Promise.all([
      responseData(getModerationUserBooks(userId, 1, 4)),
      responseData(getModerationUserComments(userId, 1, 4)),
      responseData(getUserLoginLocations(userId)),
    ])
  } else if (leaf === 'activity') {
    ;[history, reports, logins] = await Promise.all([
      responseData(getModerationUserHistory(userId, 1, 4)),
      responseData(getModerationUserReports(userId, 1, 4)),
      responseData(getUserLoginHistory(userId, 1, 4)),
    ])
  } else if (leaf === 'books') books = await responseData(getModerationUserBooks(userId, page))
  else if (leaf === 'comments')
    comments = await responseData(getModerationUserComments(userId, page))
  else if (leaf === 'history') history = await responseData(getModerationUserHistory(userId, page))
  else if (leaf === 'reports') reports = await responseData(getModerationUserReports(userId, page))

  return { user, books, comments, history, reports, logins, locations }
}

export default function UserModerationPage() {
  const data = useLoaderData<Awaited<ReturnType<typeof userModerationRouteLoader>>>()
  const { user } = data
  const location = useLocation()

  const leaf = location.pathname.split('/').at(-1) ?? ''
  const section = ['activity', 'history', 'reports'].includes(leaf)
    ? 'activity'
    : leaf === 'actions'
      ? 'actions'
      : 'overview'

  return (
    <DashboardContent.Root>
      <div className="py-8 max-w-360 mx-auto">
        <UserHeader user={user} />
        <nav aria-label={window._('moderationPortal.user.navigation')}>
          <ul className="tabs tabs--primary mt-4">
            <UserTab active={section === 'overview'} end to={`/users/${user.id}`}>
              {window._('moderationPortal.user.overview')}
            </UserTab>
            <UserTab active={section === 'activity'} to={`/users/${user.id}/activity`}>
              {window._('moderationPortal.user.activity')}
            </UserTab>
            <UserTab active={section === 'actions'} to={`/users/${user.id}/actions`}>
              {window._('moderationPortal.user.actions')}
            </UserTab>
          </ul>
        </nav>
        <div className="pt-5">
          {leaf === user.id && <Overview user={user} data={data} />}
          {leaf === 'activity' && <Activity userId={user.id} data={data} />}
          {['history', 'reports', 'books', 'comments'].includes(leaf) && (
            <ResourcePage resource={leaf} data={data} />
          )}
          {leaf === 'actions' && <Actions user={user} />}
        </div>
      </div>
    </DashboardContent.Root>
  )
}

export function UserModerationErrorPage() {
  const error = useRouteError()
  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader title={window._('moderationPortal.user.title')} />
      <div className="card">
        <ErrorDisplay error={error} />
      </div>
    </DashboardContent.Root>
  )
}

function UserTab({
  children,
  to,
  end = false,
  active,
}: {
  children: ReactNode
  to: string
  end?: boolean
  active: boolean
}) {
  return (
    <li>
      <NavLink end={end} className={`tab ${active ? 'tab--active' : ''}`} to={to}>
        {children}
      </NavLink>
    </li>
  )
}

function UserHeader({ user }: { user: ModerationUser }) {
  return (
    <header className="card card--elevated">
      <div className="flex flex-col xl:flex-row xl:items-center gap-6">
        <div className="flex items-center gap-4 min-w-0 xl:min-w-96">
          <img className="avatar flex-none" src={user.avatar} alt="" />
          <div className="min-w-0">
            <h1 className="text-3xl font-semibold truncate">{user.name}</h1>
            <p className="text-sm text-secondary-foreground break-all">
              {window._('moderationPortal.user.userId')}: {user.id}
            </p>
            <div className="flex flex-wrap gap-2 mt-2">
              <span className="chip">{user.role}</span>
              <span className={`chip ${user.isBanned ? 'chip--destructive' : 'chip--primary'}`}>
                {user.isBanned
                  ? window._('moderationPortal.user.banned')
                  : window._('moderationPortal.user.active')}
              </span>
            </div>
          </div>
        </div>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 flex-1">
          <Metric
            icon="fa-calendar"
            label={window._('moderationPortal.user.joined')}
            value={dateFormatter.format(new Date(user.joinedAt))}
          />
          <Metric
            icon="fa-book"
            label={window._('moderationPortal.user.books')}
            value={user.booksTotal}
          />
          <Metric
            icon="fa-comments"
            label={window._('moderationPortal.user.comments')}
            value={user.commentsTotal}
          />
          <Metric
            icon="fa-users"
            label={window._('moderationPortal.user.followers')}
            value={user.followersTotal}
          />
        </div>
      </div>
    </header>
  )
}

function Metric({ icon, label, value }: { icon: string; label: string; value: ReactNode }) {
  return (
    <div className="border border-border rounded-lg p-3 min-w-0">
      <div className="flex items-center gap-2 text-secondary-foreground text-sm">
        <i className={`fa-solid ${icon}`} />
        <span>{label}</span>
      </div>
      <div className="font-semibold mt-1 truncate">{value}</div>
    </div>
  )
}

function Overview({
  user,
  data,
}: {
  user: ModerationUser
  data: Awaited<ReturnType<typeof userModerationRouteLoader>>
}) {
  const email = user.email.trim() || '--'
  const emailDomain = user.email.includes('@') ? user.email.split('@').at(-1) || '--' : '--'
  const accountAge = formatAccountAge(user.joinedAt)
  return (
    <div className="grid gap-5">
      <section className="card">
        <h2 className="text-xl font-semibold mb-4">{window._('moderationPortal.user.account')}</h2>
        <dl className="grid sm:grid-cols-2 xl:grid-cols-4 gap-5">
          <Signal
            icon="fa-envelope"
            label={window._('moderationPortal.user.email')}
            value={email}
          />
          <Signal
            icon="fa-at"
            label={window._('moderationPortal.user.emailDomain')}
            value={emailDomain}
          />
          <Signal
            icon="fa-shield-halved"
            label={window._('moderationPortal.user.emailVerification')}
            value={
              user.isEmailVerified
                ? window._('moderationPortal.user.verified')
                : window._('moderationPortal.user.notVerified')
            }
          />
          <Signal
            icon="fa-hourglass-half"
            label={window._('moderationPortal.user.accountAge')}
            value={accountAge}
          />
        </dl>
      </section>
      <section className="card">
        <h2 className="text-xl font-semibold mb-4">
          {window._('moderationPortal.user.accountSignals')}
        </h2>
        <div className="grid sm:grid-cols-2 xl:grid-cols-4 gap-3">
          <LoginLocations locations={data.locations ?? []} />
          {['lastActivity', 'passwordChanges', 'emailChanges', 'twoFactor'].map((key) => (
            <UnavailableSignal key={key} label={window._(`moderationPortal.user.${key}`)} />
          ))}
        </div>
      </section>
      <RecentContent user={user} books={data.books} comments={data.comments} />
    </div>
  )
}

function RecentContent({
  user,
  books,
  comments,
}: {
  user: ModerationUser
  books?: ModerationUserBooksPage
  comments?: ModerationUserCommentsPage
}) {
  return (
    <div className="grid lg:grid-cols-2 gap-5">
      <PreviewCard
        title={window._('moderationPortal.user.recentBooks')}
        to={`/users/${user.id}/books`}
        link={`${window._('moderationPortal.user.viewAll')} (${user.booksTotal})`}
      >
        {books?.entries.map((book) => (
          <div key={book.id} className="py-3 border-b border-border last:border-0">
            <div className="flex items-center justify-between gap-3">
              <a
                className="link font-medium truncate"
                href={`/book/${book.id}`}
                target="_blank"
                rel="noreferrer"
              >
                {book.name}
              </a>
              <NavLink className="link text-sm shrink-0" to={`/books/${book.id}`}>
                {window._('moderationPortal.user.moderate')}
              </NavLink>
            </div>
            <div className="text-sm text-secondary-foreground">
              {dateFormatter.format(new Date(book.createdAt))}
            </div>
          </div>
        ))}
      </PreviewCard>
      <PreviewCard
        title={window._('moderationPortal.user.recentComments')}
        to={`/users/${user.id}/comments`}
        link={window._('moderationPortal.user.viewAll')}
      >
        {comments?.entries.map((comment) => (
          <div
            key={comment.id}
            className="relative py-3 border-b border-border last:border-0 hover:text-primary"
          >
            <NavLink
              className="absolute inset-0"
              to={`/comments/${comment.id}`}
              aria-label={window._('moderationPortal.user.moderateComment')}
            />
            <SanitizeHTML className="truncate pointer-events-none" value={comment.content} />
            <div className="text-sm text-secondary-foreground">
              {comment.bookName} · {comment.chapterName}
            </div>
          </div>
        ))}
      </PreviewCard>
    </div>
  )
}

function PreviewCard({
  title,
  to,
  link,
  children,
}: {
  title: string
  to: string
  link: string
  children: ReactNode
}) {
  return (
    <section className="card">
      <div className="flex justify-between items-center gap-4">
        <h2 className="text-xl font-semibold">{title}</h2>
        <NavLink className="link text-sm" to={to}>
          {link}
        </NavLink>
      </div>
      <div className="mt-2">{children}</div>
    </section>
  )
}

function Signal({ icon, label, value }: { icon: string; label: string; value: ReactNode }) {
  return (
    <div className="min-w-0">
      <dt className="text-sm text-secondary-foreground flex items-center gap-2">
        <i className={`fa-solid ${icon}`} />
        {label}
      </dt>
      <dd className="mt-1 break-words">{value}</dd>
    </div>
  )
}

function UnavailableSignal({ label }: { label: string }) {
  return (
    <div className="border border-border rounded-lg p-3 bg-secondary/50">
      <div className="font-medium">{label}</div>
      <div className="text-sm text-secondary-foreground mt-1">
        {window._('moderationPortal.user.notAvailable')}
      </div>
    </div>
  )
}

function LoginLocations({ locations }: { locations: LoginLocation[] }) {
  return (
    <div className="border border-border rounded-lg p-3 bg-secondary/50 sm:col-span-2">
      <div className="font-medium mb-2">{window._('moderationPortal.user.loginLocations')}</div>
      {locations.length ? (
        <ul className="grid gap-2">
          {locations.map((location) => (
            <li
              key={`${location.country}-${location.region}-${location.city}`}
              className="flex items-center gap-2 text-sm"
            >
              <i className="fa-solid fa-location-dot text-secondary-foreground" />
              <span className="font-medium">
                {[location.city, location.region, location.country].filter(Boolean).join(', ')}
              </span>
              <span className="text-secondary-foreground ml-auto">
                {dateFormatter.format(new Date(location.lastSeenAt))}
              </span>
            </li>
          ))}
        </ul>
      ) : (
        <div className="text-sm text-secondary-foreground">
          {window._('moderationPortal.user.noLoginLocations')}
        </div>
      )}
    </div>
  )
}

function Activity({
  userId,
  data,
}: {
  userId: string
  data: Awaited<ReturnType<typeof userModerationRouteLoader>>
}) {
  const { history, reports, logins } = data
  return (
    <div className="grid gap-5">
      <>
        <PreviewCard
          title={window._('moderationPortal.user.moderationHistory')}
          to={`/users/${userId}/history`}
          link={window._('moderationPortal.user.showAll')}
        >
          {history?.entries.length ? (
            history.entries.map((entry) => (
              <ModerationLogEntry
                key={entry.id}
                entry={{ ...entry, action: entry.type }}
                showTarget={false}
              />
            ))
          ) : (
            <EmptyState text={window._('moderationPortal.user.noModerationHistory')} />
          )}
        </PreviewCard>
        <PreviewCard
          title={window._('moderationPortal.user.reports')}
          to={`/users/${userId}/reports`}
          link={window._('moderationPortal.user.showAll')}
        >
          {reports?.entries.map((report) => (
            <div key={report.id} className="py-3 border-b border-border last:border-0">
              <NavLink className="link font-medium" to={`/reports/${report.id}`}>
                {report.number} · {report.reason}
              </NavLink>
              <div className="text-sm text-secondary-foreground">
                {report.targetType} · {dateTimeFormatter.format(new Date(report.time))}
              </div>
            </div>
          ))}
        </PreviewCard>
        <PreviewCard
          title={window._('moderationPortal.user.loginHistory')}
          to={`/users/${userId}/login-history`}
          link={window._('moderationPortal.user.showAll')}
        >
          {logins?.entries.length ? (
            logins.entries.map((entry, index) => (
              <div
                key={`${entry.loggedInAt}-${index}`}
                className="py-3 border-b border-border last:border-0 grid sm:grid-cols-[12rem_10rem_1fr] gap-2"
              >
                <span>{dateTimeFormatter.format(new Date(entry.loggedInAt))}</span>
                <span className="font-mono text-sm">{entry.ipAddress}</span>
                <span className="truncate text-secondary-foreground">{entry.userAgent}</span>
              </div>
            ))
          ) : (
            <EmptyState text={window._('moderationPortal.user.noLoginHistory')} />
          )}
        </PreviewCard>
      </>
    </div>
  )
}

function EmptyState({ text }: { text: string }) {
  return <div className="p-8 text-center text-secondary-foreground">{text}</div>
}

function ResourcePage({
  resource,
  data,
}: {
  resource: string
  data: Awaited<ReturnType<typeof userModerationRouteLoader>>
}) {
  switch (resource) {
    case 'books':
      return (
        <PagedList
          title={window._('moderationPortal.user.books')}
          data={data.books}
          render={(book: ModerationUserBooksPage['entries'][number]) => (
            <div>
              <div className="flex items-center justify-between gap-3">
                <a
                  className="link font-medium truncate"
                  href={`/book/${book.id}`}
                  target="_blank"
                  rel="noreferrer"
                >
                  {book.name}
                </a>
                <NavLink className="link text-sm shrink-0" to={`/books/${book.id}`}>
                  {window._('moderationPortal.user.moderate')}
                </NavLink>
              </div>
              <div className="text-sm text-secondary-foreground">
                {dateFormatter.format(new Date(book.createdAt))}
              </div>
            </div>
          )}
        />
      )
    case 'comments':
      return (
        <PagedList
          title={window._('moderationPortal.user.comments')}
          data={data.comments}
          render={(comment: ModerationUserCommentsPage['entries'][number]) => (
            <div className="relative hover:text-primary">
              <NavLink
                className="absolute inset-0"
                to={`/comments/${comment.id}`}
                aria-label={window._('moderationPortal.user.moderateComment')}
              />
              <SanitizeHTML className="truncate pointer-events-none" value={comment.content} />
              <div className="text-sm text-secondary-foreground">
                {comment.bookName} · {comment.chapterName}
              </div>
            </div>
          )}
        />
      )
    case 'history':
      return (
        <ModerationLog
          entries={(data.history?.entries ?? []).map((entry) => ({ ...entry, action: entry.type }))}
          emptyText={window._('moderationPortal.user.noModerationHistory')}
          page={data.history?.page}
          totalPages={data.history?.totalPages}
          showTarget={false}
        />
      )
    case 'reports':
      return (
        <PagedList
          title={window._('moderationPortal.user.reports')}
          data={data.reports}
          render={(report: ModerationUserReportsPage['entries'][number]) => (
            <div>
              <NavLink className="link font-medium" to={`/reports/${report.id}`}>
                {report.number} · {report.reason}
              </NavLink>
              <div>{report.description}</div>
              <div className="text-sm text-secondary-foreground">
                {report.targetType} · {dateTimeFormatter.format(new Date(report.time))}
              </div>
            </div>
          )}
        />
      )
    default:
      return null
  }
}

function PagedList<T>({
  title,
  data,
  render,
}: {
  title: string
  data?: PageResult<T>
  render: (entry: T) => ReactNode
}) {
  return (
    <section className="card card--nopad overflow-hidden">
      <div className="p-6 border-b border-border">
        <h2 className="text-xl font-semibold">{title}</h2>
      </div>
      {!data ? (
        <div className="min-h-32 grid place-items-center">
          <span className="loader" />
        </div>
      ) : data.entries.length === 0 ? (
        <EmptyState text={window._('moderationPortal.user.noEntries')} />
      ) : (
        <div>
          {data.entries.map((entry, index) => (
            <div className="p-4 border-b border-border last:border-0" key={index}>
              {render(entry)}
            </div>
          ))}
          <div className="p-4 flex justify-center">
            <Pagination.Facade page={data.page} totalPages={data.totalPages} size={7} />
          </div>
        </div>
      )}
    </section>
  )
}

function Actions({ user }: { user: ModerationUser }) {
  const [confirmation, setConfirmation] = useState<ModerationConfirmation>()
  const [action, setAction] = useState(user.isBanned ? 'unban' : 'temporary-ban')

  useEffect(() => {
    if (user.isBanned && action !== 'unban' && action !== 'rename' && action !== 'change-about')
      setAction('unban')
    if (!user.isBanned && action === 'unban') setAction('temporary-ban')
  }, [action, user.isBanned])

  let actionForm: ReactNode
  if (action === 'temporary-ban')
    actionForm = <TemporaryBanCard userId={user.id} request={setConfirmation} />
  else if (action === 'permanent-ban')
    actionForm = (
      <ModerationReasonActionCard
        title={window._('moderationPortal.user.permanentBan')}
        description={window._('moderationPortal.user.permanentBanDescription')}
        destructive
        submitLabel={window._('moderationPortal.user.permanentBan')}
        reasonLabel={window._('moderationPortal.user.reason')}
        onSubmit={(reason, onSuccess) =>
          setConfirmation({
            title: window._('moderationPortal.user.confirmPermanentBan'),
            description: window._('moderationPortal.user.confirmPermanentBanDescription'),
            destructive: true,
            run: () => permanentlyBanUser(user.id, reason),
            onSuccess,
          })
        }
      />
    )
  else if (action === 'unban')
    actionForm = (
      <ModerationReasonActionCard
        title={window._('moderationPortal.user.unban')}
        description={window._('moderationPortal.user.unbanDescription')}
        submitLabel={window._('moderationPortal.user.unban')}
        reasonLabel={window._('moderationPortal.user.reason')}
        onSubmit={(reason, onSuccess) =>
          setConfirmation({
            title: window._('moderationPortal.user.confirmUnban'),
            description: window._('moderationPortal.user.confirmUnbanDescription'),
            run: () => unbanUser(user.id, reason),
            onSuccess,
          })
        }
      />
    )
  else if (action === 'rename')
    actionForm = (
      <ModerationValueActionCard
        title={window._('moderationPortal.user.rename')}
        description={window._('moderationPortal.user.renameDescription')}
        initialValue={user.name}
        multiline={false}
        valueLabel={window._('moderationPortal.user.newValue')}
        reasonLabel={window._('moderationPortal.user.reason')}
        submitLabel={window._('moderationPortal.user.reviewChange')}
        onSubmit={(value, reason, onSuccess) =>
          setConfirmation({
            title: window._('moderationPortal.user.confirmRename'),
            description: window._('moderationPortal.user.confirmProfileChangeDescription'),
            run: () => renameUser(user.id, value, reason),
            onSuccess,
          })
        }
      />
    )
  else
    actionForm = (
      <ModerationValueActionCard
        title={window._('moderationPortal.user.changeAbout')}
        description={window._('moderationPortal.user.changeAboutDescription')}
        initialValue={user.about}
        multiline
        valueLabel={window._('moderationPortal.user.newValue')}
        reasonLabel={window._('moderationPortal.user.reason')}
        submitLabel={window._('moderationPortal.user.reviewChange')}
        onSubmit={(value, reason, onSuccess) =>
          setConfirmation({
            title: window._('moderationPortal.user.confirmChangeAbout'),
            description: window._('moderationPortal.user.confirmProfileChangeDescription'),
            run: () => changeUserAbout(user.id, value, reason),
            onSuccess,
          })
        }
      />
    )

  return (
    <ModerationActionsPanel
      action={action}
      onActionChange={setAction}
      options={[
        ...(user.isBanned
          ? [{ value: 'unban', label: window._('moderationPortal.user.unban') }]
          : [
              { value: 'temporary-ban', label: window._('moderationPortal.user.temporaryBan') },
              { value: 'permanent-ban', label: window._('moderationPortal.user.permanentBan') },
            ]),
        { value: 'rename', label: window._('moderationPortal.user.rename') },
        { value: 'change-about', label: window._('moderationPortal.user.changeAbout') },
      ]}
      selectLabel={window._('moderationPortal.user.selectAction')}
      selectId="moderation-user-action"
      confirmation={confirmation}
      onConfirmationChange={setConfirmation}
      successMessage={window._('moderationPortal.user.actionComplete')}
      confirmLabel={window._('moderationPortal.user.confirm')}
    >
      {actionForm}
    </ModerationActionsPanel>
  )
}

function TemporaryBanCard({
  userId,
  request,
}: {
  userId: string
  request: (confirmation: ModerationConfirmation) => void
}) {
  const [days, setDays] = useState('7')
  const [reason, setReason] = useState('')
  const submit = (event: FormEvent) => {
    event.preventDefault()
    if (!reason.trim()) return
    const until = new Date(Date.now() + Number(days) * 24 * 60 * 60 * 1000)
    request({
      title: window._('moderationPortal.user.confirmTemporaryBan'),
      description: window._('moderationPortal.user.confirmTemporaryBanDescription'),
      destructive: true,
      run: () => banUser(userId, reason.trim(), until),
      onSuccess: () => setReason(''),
    })
  }
  return (
    <ModerationActionCard
      title={window._('moderationPortal.user.temporaryBan')}
      description={window._('moderationPortal.user.temporaryBanDescription')}
    >
      <form onSubmit={submit} className="grid gap-4">
        <FormControl
          label={window._('moderationPortal.user.duration')}
          htmlFor="temporary-ban-duration"
        >
          <Select value={days} onValueChange={setDays}>
            <SelectTrigger id="temporary-ban-duration" className="w-full">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1">{window._('moderationPortal.user.oneDay')}</SelectItem>
              <SelectItem value="7">{window._('moderationPortal.user.sevenDays')}</SelectItem>
              <SelectItem value="30">{window._('moderationPortal.user.thirtyDays')}</SelectItem>
            </SelectContent>
          </Select>
        </FormControl>
        <FormControl label={window._('moderationPortal.user.reason')}>
          <textarea
            className="input min-h-24"
            required
            value={reason}
            onChange={(event) => setReason(event.target.value)}
          />
        </FormControl>
        <button className="btn btn--destructive justify-self-start">
          {window._('moderationPortal.user.temporaryBan')}
        </button>
      </form>
    </ModerationActionCard>
  )
}

function formatAccountAge(joinedAt: string) {
  const days = Math.max(0, Math.floor((Date.now() - new Date(joinedAt).getTime()) / 86_400_000))
  if (days < 30) return `${days} ${window._('moderationPortal.user.days')}`
  const months = Math.floor(days / 30)
  if (months < 24) return `${months} ${window._('moderationPortal.user.months')}`
  return `${Math.floor(months / 12)} ${window._('moderationPortal.user.years')}`
}
