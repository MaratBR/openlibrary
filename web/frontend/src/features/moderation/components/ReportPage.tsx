import { decideModerationReport, getModerationReport } from '@/features/moderation/api'
import type { ModerationReportDecision, ModerationReportDetail } from '@/features/moderation/api'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
  Timeline,
} from '@/components'
import { Checkbox } from '@/components'
import { FormControl } from '@/components/FormControl'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { ErrorDisplay } from '@/components/error'
import { OLAPIResponse } from '@/http-client'
import { FormEvent, useState } from 'react'
import {
  LoaderFunctionArgs,
  NavLink,
  useLoaderData,
  useRevalidator,
  useRouteError,
} from 'react-router'
import { ModerationConfirmation, ModerationConfirmationDialog } from './ModerationActions'

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
          <div>
            <div className="text-sm font-normal text-secondary-foreground mb-1">
              <NavLink className="link" to="/reports/search">
                {window._('moderationPortal.reports')}
              </NavLink>
              <span className="mx-2">/</span>
              {report.number}
            </div>
            <span>{report.reason}</span>
          </div>
        }
      >
        <span className="btn btn--outline pointer-events-none">
          <i className="fa-solid fa-circle text-secondary-foreground text-xs mr-2" />
          {capitalize(report.status)}
        </span>
      </DashboardContent.StickyHeader>

      <div className="w-full max-w-360 mx-auto px-4 md:px-8 py-5">
        <div className="grid xl:grid-cols-[minmax(0,1fr)_22rem] gap-4 items-start">
          <main className="grid gap-4">
            <ReportSummary report={report} />
            <ReportedContent report={report} target={target} />
            {report.bookContext && <BookContext report={report} />}
            <Activity report={report} />
          </main>
          <DecisionPanel report={report} />
        </div>
      </div>
    </DashboardContent.Root>
  )
}

function ReportSummary({ report }: { report: ModerationReportDetail }) {
  return (
    <section className="card card--elevated">
      <dl className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-5">
        <Detail label={window._('moderationPortal.report.reason')} value={report.reason} />
        <Detail
          label={window._('moderationPortal.report.severity')}
          value={capitalize(report.priority)}
          dot
        />
        <Detail
          label={window._('moderationPortal.report.submitted')}
          value={dateTimeFormatter.format(new Date(report.time))}
        />
        <Detail
          label={window._('moderationPortal.report.reporter')}
          value={report.reporterUserName || window._('moderationPortal.report.unknownReporter')}
        />
      </dl>
      <div>
        <div className="text-sm text-secondary-foreground">
          {window._('moderationPortal.report.reporterStatement')}
        </div>
        <p className="mt-1 whitespace-pre-wrap">
          {report.description || window._('moderationPortal.report.noDescription')}
        </p>
      </div>
    </section>
  )
}

function ReportedContent({
  report,
  target,
}: {
  report: ModerationReportDetail
  target: ReturnType<typeof targetPresentation>
}) {
  const book = report.bookContext
  return (
    <section className="card card--elevated">
      <h2 className="text-lg font-semibold mb-4">
        {window._('moderationPortal.report.reportedContent')}
        <span className="ml-3 text-sm font-normal text-secondary-foreground">
          <i className="fa-solid fa-lock mr-2" />
          {window._('moderationPortal.report.snapshot')}
        </span>
      </h2>
      {book ? (
        <div className="flex gap-4">
          <img
            className="w-24 sm:w-30 aspect-[2/3] object-cover rounded-lg border border-border shrink-0"
            src={book.coverUrl}
            alt=""
          />
          <div className="min-w-0">
            <span className="inline-flex items-center rounded-full bg-secondary px-2.5 py-1 text-xs font-medium text-secondary-foreground mb-2">
              <i
                className={`fa-solid ${book.scope === 'text' ? 'fa-quote-left' : book.scope === 'chapter' ? 'fa-file-lines' : 'fa-book'} mr-2`}
              />
              {window._(`moderationPortal.report.scope${capitalize(book.scope)}`)}
            </span>
            <h3 className="text-lg font-semibold">{book.title}</h3>
            <p className="text-sm">
              {window._('moderationPortal.report.by')}{' '}
              <span className="text-primary">{book.author}</span>
            </p>
            {book.scope === 'chapter' && (
              <div className="mt-3 rounded-lg border border-border bg-secondary/45 px-4 py-3 flex items-center gap-3">
                <i className="fa-regular fa-file-lines text-secondary-foreground" />
                <div>
                  <div className="text-xs text-secondary-foreground">
                    {window._('moderationPortal.report.reportedChapter')}
                  </div>
                  <div className="font-medium">{book.chapter}</div>
                </div>
              </div>
            )}
            {book.scope === 'text' && (
              <div className="mt-3 border-l-3 border-amber-400 pl-4">
                <div className="text-sm font-medium mb-1">{book.chapter}</div>
                <p className="leading-6 line-clamp-4">
                  <mark className="bg-amber-100 dark:bg-amber-900/40 px-1">{book.excerpt}</mark>
                </p>
              </div>
            )}
            {book.scope === 'book' && (
              <p className="text-sm text-secondary-foreground mt-3">
                {window._('moderationPortal.report.wholeBookReported')}
              </p>
            )}
            <NavLink className="link inline-block mt-3" to={target.to}>
              <i className="fa-solid fa-arrow-up-right-from-square mr-2" />
              {window._('moderationPortal.report.openTarget')}
            </NavLink>
          </div>
        </div>
      ) : report.userContext ? (
        <div className="grid md:grid-cols-[minmax(0,1fr)_auto] gap-5 items-start">
          <div>
            <div className="flex flex-wrap items-center gap-2 mb-3">
              <span className="chip">{report.userContext.role}</span>
              <span
                className={`chip ${report.userContext.isBanned ? 'chip--destructive' : 'chip--primary'}`}
              >
                {report.userContext.isBanned
                  ? window._('moderationPortal.report.userBanned')
                  : window._('moderationPortal.report.userActive')}
              </span>
            </div>
            <h3 className="text-xl font-semibold">{report.userContext.name}</h3>
            <p className="text-sm text-secondary-foreground break-all">
              {report.userContext.email}
            </p>
            <p className="mt-4 whitespace-pre-wrap">
              {report.userContext.about || window._('moderationPortal.report.noProfileDescription')}
            </p>
            <dl className="grid sm:grid-cols-2 gap-3 mt-4 text-sm">
              <Detail
                label={window._('moderationPortal.report.joined')}
                value={dateTimeFormatter.format(new Date(report.userContext.joinedAt))}
              />
              <Detail
                label={window._('moderationPortal.report.relatedReports')}
                value={String(report.userContext.relatedReports)}
              />
            </dl>
          </div>
          <NavLink className="btn btn--outline" to={target.to}>
            {window._('moderationPortal.report.openTarget')}
          </NavLink>
        </div>
      ) : report.commentContext ? (
        <div>
          <div className="flex flex-wrap gap-2 mb-3">
            <span
              className={`chip ${report.commentContext.deletedAt ? 'chip--destructive' : 'chip--primary'}`}
            >
              {report.commentContext.deletedAt
                ? window._('moderationPortal.report.commentRemoved')
                : window._('moderationPortal.report.commentVisible')}
            </span>
            <span className="chip">
              {window._('moderationPortal.report.relatedReportCount', {
                count: String(report.commentContext.relatedReports),
              })}
            </span>
          </div>
          <blockquote className="rounded-lg border-l-4 border-primary bg-secondary/45 px-5 py-4 whitespace-pre-wrap">
            {report.commentContext.content}
          </blockquote>
          <div className="grid sm:grid-cols-2 gap-4 mt-4 text-sm">
            <div>
              <span className="text-secondary-foreground">
                {window._('moderationPortal.report.commentAuthor')}
              </span>
              <NavLink className="link block" to={`/users/${report.commentContext.authorId}`}>
                {report.commentContext.authorName}
              </NavLink>
            </div>
            <div>
              <span className="text-secondary-foreground">
                {window._('moderationPortal.report.commentLocation')}
              </span>
              <NavLink className="link block" to={`/books/${report.commentContext.bookId}`}>
                {report.commentContext.bookName} · {report.commentContext.chapterName}
              </NavLink>
            </div>
          </div>
          <div className="flex flex-wrap items-center justify-between gap-3 mt-4">
            <time className="text-sm text-secondary-foreground">
              {dateTimeFormatter.format(new Date(report.commentContext.createdAt))}
            </time>
            <NavLink className="btn btn--outline" to={target.to}>
              {window._('moderationPortal.report.openTarget')}
            </NavLink>
          </div>
        </div>
      ) : (
        <div className="rounded-xl border border-border bg-secondary/45 p-4 flex items-center gap-4">
          <i className={`fa-solid ${target.icon} text-secondary-foreground`} />
          <div className="flex-1">
            <div className="font-semibold">{target.label}</div>
            <div className="font-mono text-sm text-secondary-foreground">{report.targetId}</div>
          </div>
          <NavLink className="btn btn--outline" to={target.to}>
            {window._('moderationPortal.report.openTarget')}
          </NavLink>
        </div>
      )}
    </section>
  )
}

function BookContext({ report }: { report: ModerationReportDetail }) {
  const book = report.bookContext!
  return (
    <>
      <div className="rounded-xl border border-amber-400 bg-amber-50 dark:bg-amber-950/25 px-5 py-3 flex items-center gap-3">
        <i className="fa-regular fa-clock text-amber-700" />
        <span className="font-medium flex-1">
          TODO: Finish this (NOT IMPLEMENTED){' '}
          {window._('moderationPortal.report.editedAfter', {
            minutes: String(book.editedAfterMinutes),
          })}
        </span>
        <button className="link">
          {window._('moderationPortal.report.compareChanges')}{' '}
          <i className="fa-solid fa-chevron-right ml-2" />
        </button>
      </div>
      <div className="grid md:grid-cols-2 gap-4">
        <section className="card">
          <h2 className="text-lg font-semibold mb-4">
            {window._('moderationPortal.report.contentContext')}
          </h2>
          <dl className="grid grid-cols-[auto_1fr] gap-x-6 gap-y-2 text-sm">
            <dt>{window._('moderationPortal.report.rating')}</dt>
            <dd>
              <Badge>{book.rating}</Badge>
            </dd>
            <dt>{window._('moderationPortal.report.warnings')}</dt>
            <dd>
              {book.warnings?.length
                ? book.warnings.join(', ')
                : window._('moderationPortal.report.none')}
            </dd>
            <dt>{window._('moderationPortal.report.publicationStatus')}</dt>
            <dd>
              <Badge>{book.publicationState}</Badge>
            </dd>
            <dt>{window._('moderationPortal.report.lastUpdated')}</dt>
            <dd>{dateTimeFormatter.format(new Date(book.lastUpdated))}</dd>
          </dl>
        </section>
        <section className="card">
          <h2 className="text-lg font-semibold mb-4">
            {window._('moderationPortal.report.history')}
          </h2>
          <p>
            <i className="fa-solid fa-users mr-3" />
            {window._('moderationPortal.report.relatedReportCount', {
              count: String(book.relatedReports),
            })}
          </p>
          <p className="mt-3">
            <i className="fa-solid fa-shield-halved mr-3" />
            {window._('moderationPortal.report.noPreviousEnforcement')}
          </p>
        </section>
      </div>
    </>
  )
}

const dispositions = ['no_violation', 'request_changes', 'action_taken', 'escalated'] as const

function DecisionPanel({ report }: { report: ModerationReportDetail }) {
  const revalidator = useRevalidator()
  const [disposition, setDisposition] = useState<(typeof dispositions)[number]>('no_violation')
  const [action, setAction] = useState('')
  const [policyReason, setPolicyReason] = useState('')
  const [internalNote, setInternalNote] = useState('')
  const [notifyTarget, setNotifyTarget] = useState(false)
  const [value, setValue] = useState('')
  const [until, setUntil] = useState('')
  const [confirmation, setConfirmation] = useState<ModerationConfirmation>()
  const [pending, setPending] = useState(false)
  const actions = report.availableActions
  const resolved = report.status === 'resolved'
  const needsValue = [
    'user.rename',
    'user.change_about',
    'book.change_age_rating',
    'book.change_summary',
  ].includes(action)
  const submit = (event: FormEvent) => {
    event.preventDefault()
    const payload: Record<string, unknown> = {}
    if (needsValue) payload.value = value
    if (action === 'user.temporary_ban') payload.until = new Date(until).toISOString()
    const body: ModerationReportDecision = {
      disposition,
      action: disposition === 'action_taken' ? action : '',
      policyReason: policyReason.trim(),
      internalNote: internalNote.trim(),
      notifyTarget,
      payload,
    }
    setConfirmation({
      title: window._('moderationPortal.report.confirmDecision'),
      description: window._('moderationPortal.report.confirmDecisionHelp'),
      destructive: disposition === 'action_taken',
      run: () => decideModerationReport(report.id, body),
    })
  }
  const confirm = async () => {
    if (!confirmation) return
    setPending(true)
    try {
      const response = (await confirmation.run()) as { throwIfError?: () => void }
      response.throwIfError?.()
      setConfirmation(undefined)
      window.toast({ title: window._('moderationPortal.report.decisionRecorded') })
      await revalidator.revalidate()
    } catch (error) {
      window.toast.error(error)
    } finally {
      setPending(false)
    }
  }
  return (
    <aside className="card xl:sticky xl:top-22">
      <h2 className="text-xl font-semibold mb-4">{window._('moderationPortal.report.decision')}</h2>
      {resolved ? (
        <p className="rounded-lg bg-secondary p-4 text-secondary-foreground">
          {window._('moderationPortal.report.alreadyResolved')}
        </p>
      ) : (
        <form onSubmit={submit} className="grid gap-4">
          <FormControl label={window._('moderationPortal.report.outcome')}>
            <Select
              value={disposition}
              onValueChange={(next) => {
                setDisposition(next as typeof disposition)
                setAction('')
              }}
            >
              <SelectTrigger className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {dispositions.map((item) => (
                  <SelectItem key={item} value={item}>
                    {window._(`moderationPortal.report.disposition_${item}`)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </FormControl>
          {disposition === 'action_taken' && (
            <FormControl label={window._('moderationPortal.report.enforcementAction')}>
              <Select value={action} onValueChange={setAction}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder={window._('moderationPortal.report.chooseAction')} />
                </SelectTrigger>
                <SelectContent>
                  {actions.map((item) => (
                    <SelectItem key={item} value={item}>
                      {window._(`moderationPortal.report.action_${item.replace(/\./g, '_')}`)}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </FormControl>
          )}
          {action === 'user.temporary_ban' && (
            <FormControl label={window._('moderationPortal.report.banUntil')}>
              <input
                className="input"
                type="datetime-local"
                required
                value={until}
                onChange={(event) => setUntil(event.target.value)}
              />
            </FormControl>
          )}
          {action === 'book.change_age_rating' && (
            <FormControl label={window._('moderationPortal.report.newValue')}>
              <Select value={value} onValueChange={setValue}>
                <SelectTrigger className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {['?', 'G', 'PG', 'PG-13', 'R', 'NC-17'].map((rating) => (
                    <SelectItem key={rating} value={rating}>
                      {rating}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </FormControl>
          )}
          {needsValue && action !== 'book.change_age_rating' && (
            <FormControl label={window._('moderationPortal.report.newValue')}>
              <textarea
                className="input min-h-20"
                required
                value={value}
                onChange={(event) => setValue(event.target.value)}
              />
            </FormControl>
          )}
          <FormControl label={window._('moderationPortal.report.policyReason')}>
            <textarea
              className="input min-h-20"
              required
              value={policyReason}
              onChange={(event) => setPolicyReason(event.target.value)}
            />
          </FormControl>
          <FormControl label={window._('moderationPortal.report.internalNote')}>
            <textarea
              className="input min-h-24"
              maxLength={1000}
              value={internalNote}
              onChange={(event) => setInternalNote(event.target.value)}
              placeholder={window._('moderationPortal.report.notePlaceholder')}
            />
          </FormControl>
          <label className="flex items-start gap-3 cursor-pointer" htmlFor="report-notify-target">
            <Checkbox
              id="report-notify-target"
              className="mt-0.5"
              checked={notifyTarget}
              onCheckedChange={(checked) => setNotifyTarget(checked === true)}
            />
            <span>
              {window._('moderationPortal.report.notifyAuthor')}
              <small className="block text-secondary-foreground mt-1">
                {window._('moderationPortal.report.notifyUnavailable')}
              </small>
            </span>
          </label>
          <button
            disabled={disposition === 'action_taken' && !action}
            className="btn btn--primary w-full"
          >
            <i className="fa-solid fa-shield-halved mr-2" />
            {window._('moderationPortal.report.applyAction')}
          </button>
        </form>
      )}
      <p className="text-xs italic text-center text-secondary-foreground mt-5">
        {window._('moderationPortal.report.recordedNotice')}
      </p>
      <ModerationConfirmationDialog
        confirmation={confirmation}
        pending={pending}
        onClose={() => setConfirmation(undefined)}
        onConfirm={() => void confirm()}
        confirmLabel={window._('moderationPortal.report.applyAction')}
      />
    </aside>
  )
}

function Activity({ report }: { report: ModerationReportDetail }) {
  return (
    <section className="card">
      <h2 className="text-lg font-semibold mb-4">{window._('moderationPortal.report.activity')}</h2>
      <Timeline.Root>
        {report.activities.map((activity) => (
          <Timeline.Item
            key={`${activity.time}-${activity.kind}`}
            marker={
              <i
                className={`fa-solid ${activity.kind === 'decision' ? 'fa-gavel' : 'fa-flag'} text-xs`}
              />
            }
          >
            <p>
              <b>{activity.actor}</b>{' '}
              {activity.kind === 'decision'
                ? window._(`moderationPortal.report.disposition_${activity.disposition}`)
                : activity.description}
            </p>
            {activity.action && (
              <p className="text-sm mt-1">
                {window._(`moderationPortal.report.action_${activity.action.replace(/\./g, '_')}`)}
              </p>
            )}
            {activity.policyReason && (
              <p className="text-sm text-secondary-foreground mt-1">{activity.policyReason}</p>
            )}
            {activity.internalNote && (
              <p className="mt-2 whitespace-pre-wrap">{activity.internalNote}</p>
            )}
            {activity.notifyTarget && (
              <p className="text-xs text-secondary-foreground mt-1">
                {window._('moderationPortal.report.notificationRequested')}
              </p>
            )}
            <time className="text-sm text-secondary-foreground">
              {dateTimeFormatter.format(new Date(activity.time))}
            </time>
          </Timeline.Item>
        ))}
      </Timeline.Root>
    </section>
  )
}

function Detail({ label, value, dot }: { label: string; value: string; dot?: boolean }) {
  return (
    <div>
      <dt className="text-sm text-secondary-foreground mb-1">{label}</dt>
      <dd>
        {dot && <i className="fa-solid fa-circle text-amber-400 text-xs mr-2" />}
        {value}
      </dd>
    </div>
  )
}
function Badge({ children }: { children: string }) {
  return <span className="rounded bg-primary/10 text-primary px-2 py-0.5">{children}</span>
}
function targetPresentation(report: ModerationReportDetail) {
  if (report.targetType === 'user')
    return {
      icon: 'fa-user',
      label: window._('moderationPortal.report.targetUser'),
      to: `/users/${report.targetId}`,
    }
  if (report.targetType === 'book')
    return {
      icon: 'fa-book',
      label: window._('moderationPortal.report.targetBook'),
      to: `/books/${report.targetId}`,
    }
  return {
    icon: 'fa-comment',
    label: window._('moderationPortal.report.targetComment'),
    to: `/comments/${report.targetId}`,
  }
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
