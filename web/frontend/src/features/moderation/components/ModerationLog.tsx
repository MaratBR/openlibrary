import { Pagination } from '@/components/Pagination'
import { ReactNode } from 'react'
import { NavLink } from 'react-router'

export type ModerationLogRecord = {
  id?: string
  time: string
  action: string
  targetType: string
  targetId: string
  reason: string
  payload: unknown
  actorUserId: string
  actorUserName: string
}
type PayloadRenderer = (payload: unknown) => ReactNode

// Action-specific payload renderers can be registered here as audit payload schemas evolve.
const payloadRenderers: Record<string, PayloadRenderer> = {
  rename: (payload) => <PayloadFields payload={payload} />,
  change_about: (payload) => <PayloadFields payload={payload} />,
  change_summary: (payload) => <PayloadFields payload={payload} />,
  change_age_rating: (payload) => <PayloadFields payload={payload} />,
}

export function ModerationLog({
  entries,
  emptyText,
  page,
  totalPages,
  showTarget = true,
}: {
  entries: ModerationLogRecord[]
  emptyText: string
  page?: number
  totalPages?: number
  showTarget?: boolean
}) {
  return (
    <section className="Card Card--nopad overflow-hidden">
      {entries.length ? (
        entries.map((entry, index) => (
          <ModerationLogEntry
            key={entry.id ?? `${entry.time}-${index}`}
            entry={entry}
            showTarget={showTarget}
          />
        ))
      ) : (
        <div className="p-8 text-center text-secondary-foreground">{emptyText}</div>
      )}
      {page && totalPages && totalPages > 1 ? (
        <div className="p-4 flex justify-center border-t border-border">
          <Pagination.Facade page={page} totalPages={totalPages} size={7} />
        </div>
      ) : null}
    </section>
  )
}

export function ModerationLogEntry({
  entry,
  showTarget = true,
}: {
  entry: ModerationLogRecord
  showTarget?: boolean
}) {
  return (
    <article className="p-4 sm:p-5 border-b border-border last:border-0">
      <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-2">
        <div>
          <div className="flex flex-wrap items-center gap-2">
            <strong>{formatAction(entry.action)}</strong>
            {showTarget && (
              <>
                <span className="Chip">{entry.targetType}</span>
                <TargetLink type={entry.targetType} id={entry.targetId} />
              </>
            )}
          </div>
          <div className="text-sm text-secondary-foreground mt-1">
            {entry.actorUserName || window._('moderationPortal.user.system')} ·{' '}
            {new Intl.DateTimeFormat(undefined, { dateStyle: 'medium', timeStyle: 'short' }).format(
              new Date(entry.time),
            )}
          </div>
        </div>
      </div>
      {entry.reason && <p className="mt-3">{entry.reason}</p>}
      <Payload action={entry.action} payload={entry.payload} />
    </article>
  )
}

function TargetLink({ type, id }: { type: string; id: string }) {
  const routes: Record<string, string> = {
    user: `/users/${id}`,
    book: `/books/${id}`,
    chapter: `/chapters/${id}`,
    comment: `/comments/${id}`,
  }
  return routes[type] ? (
    <NavLink className="Link text-sm" to={routes[type]}>
      {id}
    </NavLink>
  ) : (
    <span className="text-sm font-mono">{id}</span>
  )
}
function Payload({ action, payload }: { action: string; payload: unknown }) {
  if (payload === null || payload === undefined || payload === '') return null
  const renderer = payloadRenderers[action]
  return (
    <div className="mt-3">{renderer ? renderer(payload) : <PayloadPre payload={payload} />}</div>
  )
}
function PayloadFields({ payload }: { payload: unknown }) {
  if (!payload || typeof payload !== 'object' || Array.isArray(payload))
    return <PayloadPre payload={payload} />
  return (
    <dl className="grid sm:grid-cols-2 gap-2 rounded-lg bg-secondary/50 p-3">
      {Object.entries(payload).map(([key, value]) => (
        <div key={key}>
          <dt className="text-xs text-secondary-foreground">{formatAction(key)}</dt>
          <dd className="break-words">
            {typeof value === 'string' ? value : JSON.stringify(value)}
          </dd>
        </div>
      ))}
    </dl>
  )
}
function PayloadPre({ payload }: { payload: unknown }) {
  let text: string
  if (typeof payload === 'string') {
    try {
      text = JSON.stringify(JSON.parse(payload), null, 2)
    } catch {
      text = payload
    }
  } else {
    try {
      text = JSON.stringify(payload, null, 2)
    } catch {
      text = String(payload)
    }
  }
  return (
    <pre className="overflow-auto whitespace-pre-wrap break-words rounded-lg bg-secondary/50 p-3 text-sm">
      {text}
    </pre>
  )
}
function formatAction(value: string) {
  return value.replace(/_/g, ' ').replace(/\b\w/g, (letter) => letter.toUpperCase())
}
