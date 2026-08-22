import { DashboardContent } from '@/components/dashboard-layout-components'
import { ErrorDisplay } from '@/components/error'
import { Pagination } from '@/components/Pagination'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components'
import {
  changeBookAgeRating,
  changeBookSummary,
  getBookModerationChapters,
  getBookModerationLog,
  getModerationBook,
  performBookModerationAction,
} from '@/features/moderation/api'
import type {
  ModerationBook,
  ModerationBookChapter,
  BookModerationLog,
} from '@/features/moderation/api'
import { ReactNode, useMemo, useState } from 'react'
import {
  LoaderFunctionArgs,
  NavLink,
  useLoaderData,
  useLocation,
  useRouteError,
} from 'react-router'
import {
  ModerationActionsPanel,
  ModerationConfirmation,
  ModerationReasonActionCard,
  ModerationValueActionCard,
} from './ModerationActions'
import { ModerationLog } from './ModerationLog'

const date = new Intl.DateTimeFormat(undefined, { dateStyle: 'medium' })
async function data<T>(promise: Promise<{ throwIfError(): void; data: T }>) {
  const response = await promise
  response.throwIfError()
  return response.data
}

export async function bookModerationRouteLoader({ params, request }: LoaderFunctionArgs) {
  if (!params.bookId) throw new Error(window._('moderationPortal.book.loadError'))
  const leaf = new URL(request.url).pathname.split('/').at(-1)
  const book = await data(getModerationBook(params.bookId))
  const chapters =
    leaf === 'chapters' ? await data(getBookModerationChapters(params.bookId)) : undefined
  const log =
    leaf === 'activity'
      ? await data(getBookModerationLog(params.bookId, { pageSize: 100 }))
      : undefined
  return { book, chapters, log }
}

export default function BookModerationPage() {
  const loaded = useLoaderData<Awaited<ReturnType<typeof bookModerationRouteLoader>>>()
  const leaf = useLocation().pathname.split('/').at(-1)
  const section =
    leaf === 'actions' || leaf === 'chapters' || leaf === 'activity' ? leaf : 'overview'
  return (
    <DashboardContent.Root>
      <div className="py-8 max-w-360 mx-auto">
        <Header book={loaded.book} />
        <nav>
          <ul className="tabs tabs--primary mt-4">
            {(['overview', 'actions', 'chapters', 'activity'] as const).map((tab) => (
              <li key={tab}>
                <NavLink
                  end={tab === 'overview'}
                  className={`tab ${section === tab ? 'tab--active' : ''}`}
                  to={
                    tab === 'overview'
                      ? `/books/${loaded.book.id}`
                      : `/books/${loaded.book.id}/${tab}`
                  }
                >
                  {window._(`moderationPortal.book.${tab}`)}
                </NavLink>
              </li>
            ))}
          </ul>
        </nav>
        <div className="pt-5">
          {section === 'overview' && <Overview book={loaded.book} />}
          {section === 'actions' && <Actions book={loaded.book} />}
          {section === 'chapters' && <Chapters chapters={loaded.chapters ?? []} />}
          {section === 'activity' && <Activity bookId={loaded.book.id} log={loaded.log} />}
        </div>
      </div>
    </DashboardContent.Root>
  )
}

function Header({ book }: { book: ModerationBook }) {
  return (
    <header className="card card--elevated flex flex-col lg:flex-row lg:items-center gap-5">
      <div className="min-w-0 flex-1">
        <h1 className="text-3xl font-semibold truncate">{book.name}</h1>
        <p className="text-sm text-secondary-foreground">
          {window._('moderationPortal.book.bookId')}: {book.id}
        </p>
        <div className="flex gap-2 mt-2">
          <span className="chip">{book.ageRating}</span>
          <span className={`chip ${book.isBanned ? 'chip--destructive' : 'chip--primary'}`}>
            {book.isBanned
              ? window._('moderationPortal.book.banned')
              : window._('moderationPortal.book.active')}
          </span>
        </div>
      </div>
      <div className="grid grid-cols-3 gap-3">
        <Metric
          label={window._('moderationPortal.book.words')}
          value={book.words.toLocaleString()}
        />
        <Metric label={window._('moderationPortal.book.chapters')} value={book.chapters} />
        <Metric label={window._('moderationPortal.book.reports')} value={book.reportsCount} />
      </div>
    </header>
  )
}
function Metric({ label, value }: { label: string; value: ReactNode }) {
  return (
    <div className="border border-border rounded-lg p-3">
      <div className="text-sm text-secondary-foreground">{label}</div>
      <div className="font-semibold mt-1">{value}</div>
    </div>
  )
}

function Overview({ book }: { book: ModerationBook }) {
  return (
    <div className="grid gap-5">
      <section className="card">
        <h2 className="text-xl font-semibold mb-4">
          {window._('moderationPortal.book.generalInfo')}
        </h2>
        <dl className="grid sm:grid-cols-2 xl:grid-cols-4 gap-5">
          <Info
            label={window._('moderationPortal.book.author')}
            value={
              <NavLink className="link" to={`/users/${book.authorUserId}`}>
                {book.authorUserName}
              </NavLink>
            }
          />
          <Info
            label={window._('moderationPortal.book.released')}
            value={date.format(new Date(book.createdAt))}
          />
          <Info label={window._('moderationPortal.book.ageRating')} value={book.ageRating} />
          <Info
            label={window._('moderationPortal.book.visibility')}
            value={
              book.isPubliclyVisible
                ? window._('moderationPortal.book.public')
                : window._('moderationPortal.book.private')
            }
          />
        </dl>
        <h3 className="font-semibold mt-6">{window._('moderationPortal.book.summary')}</h3>
        <p className="mt-2 whitespace-pre-wrap">{book.summary || '--'}</p>
      </section>
      <section className="card">
        <h2 className="text-xl font-semibold mb-4">
          {window._('moderationPortal.book.currentState')}
        </h2>
        {book.isBanned && (
          <div className="p-3 rounded-lg bg-destructive/10 text-destructive mb-4">
            <strong>{window._('moderationPortal.book.banReason')}:</strong> {book.banReason || '--'}
          </div>
        )}
        <Info
          label={window._('moderationPortal.book.pendingReport')}
          value={
            book.latestPendingReport ? (
              <NavLink className="link" to={`/reports/${book.latestPendingReport.id}`}>
                {book.latestPendingReport.number} · {book.latestPendingReport.reason}
              </NavLink>
            ) : (
              window._('moderationPortal.book.noPendingReports')
            )
          }
        />
      </section>
    </div>
  )
}
function Info({ label, value }: { label: string; value: ReactNode }) {
  return (
    <div>
      <dt className="text-sm text-secondary-foreground">{label}</dt>
      <dd className="mt-1">{value}</dd>
    </div>
  )
}

function Actions({ book }: { book: ModerationBook }) {
  const [action, setAction] = useState('rating')
  const [confirmation, setConfirmation] = useState<ModerationConfirmation>()
  const request = (
    title: string,
    description: string,
    run: () => Promise<unknown>,
    destructive = false,
    onSuccess?: () => void,
  ) => setConfirmation({ title, description, run, destructive, onSuccess })
  let form: ReactNode
  if (action === 'rating')
    form = (
      <ModerationValueActionCard
        title={window._('moderationPortal.book.changeAgeRating')}
        description={window._('moderationPortal.book.changeAgeRatingDescription')}
        initialValue={book.ageRating}
        valueLabel={window._('moderationPortal.book.ageRating')}
        reasonLabel={window._('moderationPortal.book.reason')}
        submitLabel={window._('moderationPortal.book.reviewChange')}
        control={(value, setValue) => (
          <Select value={value} onValueChange={setValue}>
            <SelectTrigger className="w-full">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {['?', 'G', 'PG', 'PG-13', 'R', 'NC-17'].map((rating) => (
                <SelectItem value={rating} key={rating}>
                  {rating}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        )}
        onSubmit={(value, reason, onSuccess) =>
          request(
            window._('moderationPortal.book.confirmAgeRating'),
            window._('moderationPortal.book.confirmChangeDescription'),
            () => changeBookAgeRating(book.id, value, reason),
            false,
            onSuccess,
          )
        }
      />
    )
  else if (action === 'summary')
    form = (
      <ModerationValueActionCard
        title={window._('moderationPortal.book.changeSummary')}
        description={window._('moderationPortal.book.changeSummaryDescription')}
        initialValue={book.summary}
        multiline
        valueLabel={window._('moderationPortal.book.summary')}
        reasonLabel={window._('moderationPortal.book.reason')}
        submitLabel={window._('moderationPortal.book.reviewChange')}
        onSubmit={(value, reason, onSuccess) =>
          request(
            window._('moderationPortal.book.confirmSummary'),
            window._('moderationPortal.book.confirmChangeDescription'),
            () => changeBookSummary(book.id, value, reason),
            false,
            onSuccess,
          )
        }
      />
    )
  else
    form = (
      <ModerationReasonActionCard
        title={
          book.isBanned
            ? window._('moderationPortal.book.unban')
            : window._('moderationPortal.book.ban')
        }
        description={
          book.isBanned
            ? window._('moderationPortal.book.unbanDescription')
            : window._('moderationPortal.book.banDescription')
        }
        submitLabel={
          book.isBanned
            ? window._('moderationPortal.book.unban')
            : window._('moderationPortal.book.ban')
        }
        reasonLabel={window._('moderationPortal.book.reason')}
        destructive={!book.isBanned}
        onSubmit={(reason, onSuccess) =>
          request(
            book.isBanned
              ? window._('moderationPortal.book.confirmUnban')
              : window._('moderationPortal.book.confirmBan'),
            book.isBanned
              ? window._('moderationPortal.book.confirmUnbanDescription')
              : window._('moderationPortal.book.confirmBanDescription'),
            () => performBookModerationAction(book.id, book.isBanned ? 'unban' : 'ban', reason),
            !book.isBanned,
            onSuccess,
          )
        }
      />
    )
  return (
    <ModerationActionsPanel
      action={action}
      onActionChange={setAction}
      options={[
        { value: 'rating', label: window._('moderationPortal.book.changeAgeRating') },
        {
          value: 'restriction',
          label: book.isBanned
            ? window._('moderationPortal.book.unban')
            : window._('moderationPortal.book.ban'),
        },
        { value: 'summary', label: window._('moderationPortal.book.changeSummary') },
      ]}
      selectLabel={window._('moderationPortal.book.selectAction')}
      selectId="moderation-book-action"
      confirmation={confirmation}
      onConfirmationChange={setConfirmation}
      successMessage={window._('moderationPortal.book.actionComplete')}
      confirmLabel={window._('moderationPortal.book.confirm')}
    >
      {form}
    </ModerationActionsPanel>
  )
}

function Chapters({ chapters }: { chapters: ModerationBookChapter[] }) {
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(1)
  const filtered = useMemo(
    () => chapters.filter((c) => c.name.toLocaleLowerCase().includes(search.toLocaleLowerCase())),
    [chapters, search],
  )
  const pages = Math.ceil(filtered.length / 100)
  const shown = filtered.slice((page - 1) * 100, page * 100)
  return (
    <section className="card card--nopad overflow-hidden">
      <div className="p-5 border-b border-border">
        <input
          className="input w-full"
          type="search"
          value={search}
          onChange={(e) => {
            setSearch(e.target.value)
            setPage(1)
          }}
          placeholder={window._('moderationPortal.book.searchChapters')}
        />
      </div>
      {shown.map((chapter) => (
        <NavLink
          key={chapter.id}
          to={`/chapters/${chapter.id}`}
          className="block p-4 border-b border-border hover:bg-foreground/5"
        >
          <div className="flex gap-3 justify-between">
            <span className="font-medium">{chapter.name}</span>
            {chapter.hasPendingReports && (
              <span className="chip chip--destructive">
                {window._('moderationPortal.book.pendingReport')}
              </span>
            )}
          </div>
          <div className="text-sm text-secondary-foreground mt-1">
            {window._('moderationPortal.book.released')} {date.format(new Date(chapter.createdAt))}{' '}
            ·{' '}
            {chapter.updatedAt
              ? `${window._('moderationPortal.book.updated')} ${date.format(new Date(chapter.updatedAt))} · `
              : ''}
            {chapter.words.toLocaleString()} {window._('moderationPortal.book.words').toLowerCase()}
          </div>
        </NavLink>
      ))}
      {!shown.length && (
        <div className="p-8 text-center text-secondary-foreground">
          {window._('moderationPortal.book.noChapters')}
        </div>
      )}
      {pages > 1 && (
        <div className="p-4 flex justify-center">
          <Pagination.Root>
            {Array.from({ length: pages }, (_, index) => index + 1).map((number) => (
              <Pagination.Item
                key={number}
                active={number === page}
                onClick={() => setPage(number)}
              >
                {number}
              </Pagination.Item>
            ))}
          </Pagination.Root>
        </div>
      )}
    </section>
  )
}

function Activity({ bookId, log }: { bookId: string; log?: BookModerationLog }) {
  return (
    <ModerationLog
      entries={(log?.entries ?? []).map((entry, index) => ({
        ...entry,
        id: String(index),
        targetType: 'book',
        targetId: bookId,
        payload: entry.payload,
      }))}
      emptyText={window._('moderationPortal.book.noActivity')}
      page={log?.page}
      totalPages={log?.totalPages}
      showTarget={false}
    />
  )
}

export function BookModerationErrorPage() {
  return (
    <DashboardContent.Root>
      <div className="card">
        <ErrorDisplay error={useRouteError()} />
      </div>
    </DashboardContent.Root>
  )
}
