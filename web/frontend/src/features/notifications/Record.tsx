import sanitizeHtml from 'sanitize-html'
import type { NotificationEvent } from './api'
import './record.scss'

/** Keep the client allowlist in sync with notifications.SanitizeContent. */
function safeContent(content: string) {
  return sanitizeHtml(content, {
    allowedTags: ['b', 'strong', 'i', 'em', 'u', 's', 'del', 'sub', 'sup', 'br', 'code', 'a'],
    allowedAttributes: { a: ['href', 'title', 'rel'] },
    allowedSchemes: ['http', 'https', 'mailto'],
    allowProtocolRelative: false,
  })
}

export function NotificationRecord({ event }: { event: NotificationEvent }) {
  const icon =
    event.icon ??
    (event.category === 'moderation'
      ? { type: 1, payload: 'fa-solid fa-shield' }
      : event.category === 'admin'
        ? { type: 1, payload: 'fa-solid fa-gauge-high' }
        : undefined)

  return (
    <article className={`NotificationRecord${event.urgent ? ' NotificationRecord--urgent' : ''}`}>
      {icon && (
        <span className="NotificationRecord-icon" aria-hidden="true">
          {icon.type === 2 ? (
            icon.payload
          ) : icon.type === 1 ? (
            <i className={icon.payload} />
          ) : icon.type === 0 && /^(https?:\/\/|\/(?!\/))/.test(icon.payload) ? (
            <img src={icon.payload} alt="" loading="lazy" referrerPolicy="no-referrer" />
          ) : null}
        </span>
      )}
      <div className="NotificationRecord-details">
        {event.urgent && (
          <span className="NotificationRecord-urgent">{window._('notifications.urgent')}</span>
        )}
        <h3 className="NotificationRecord-title">{event.title}</h3>
        <div
          className="NotificationRecord-content"
          dangerouslySetInnerHTML={{ __html: safeContent(event.content) }}
        />
        <time className="NotificationRecord-time" dateTime={event.createdAt}>
          {new Intl.DateTimeFormat(undefined, { dateStyle: 'medium', timeStyle: 'short' }).format(
            new Date(event.createdAt),
          )}
        </time>
      </div>
    </article>
  )
}
