import { useRef, useState } from 'react'
import { NavLink } from 'react-router'
import { httpClient } from '@/features/http-client'
import { roles } from '@/features/admin/api'
import { Card, Notice, PageHeader, errorMessage, t } from '@/features/admin/components/ui'

export function DebugNavigation() {
  return (
    <nav className="Tabs mb-6" aria-label={t('admin.sidebar.debugActions')}>
      <NavLink className="Tab" to="/debug" end>
        {t('admin.notifications.actions')}
      </NavLink>
      <NavLink className="Tab" to="/debug/notifications">
        {t('admin.notifications.title')}
      </NavLink>
    </nav>
  )
}

export default function SendNotifications() {
  const [target, setTarget] = useState('user')
  const [iconType, setIconType] = useState('emoji')
  const [pending, setPending] = useState(false)
  const [count, setCount] = useState<number>()
  const [error, setError] = useState<unknown>()
  const locked = useRef(false)
  async function send(form: HTMLFormElement) {
    if (locked.current) return
    locked.current = true
    setPending(true)
    setCount(undefined)
    setError(undefined)
    try {
      const data = new FormData(form)
      const body = new URLSearchParams()
      for (const [key, value] of data) body.set(key, String(value))
      const result = await httpClient
        .post('/admin/api/debug/notifications', { body, retry: 0 })
        .json<{ count: number }>()
      setCount(result.count)
    } catch (error) {
      setError(error)
    } finally {
      locked.current = false
      setPending(false)
    }
  }
  return (
    <>
      <PageHeader title={t('admin.sidebar.debugActions')} />
      <DebugNavigation />
      <Card
        title={t('admin.notifications.title')}
        description={t('admin.notifications.description')}
      >
        <form
          className="grid gap-5 max-w-160"
          onSubmit={(event) => {
            event.preventDefault()
            void send(event.currentTarget)
          }}
        >
          <label className="grid gap-2">
            {t('admin.notifications.target')}
            <select
              className="input"
              value={target}
              disabled={pending}
              onChange={(event) => setTarget(event.target.value)}
            >
              <option value="user">{t('admin.notifications.user')}</option>
              <option value="role">{t('admin.notifications.role')}</option>
            </select>
          </label>
          {target === 'user' ? (
            <label className="grid gap-2">
              {t('admin.notifications.userId')}
              <input
                className="input"
                name="userId"
                required
                pattern="[0-9a-fA-F-]{36}"
                disabled={pending}
              />
            </label>
          ) : (
            <label className="grid gap-2">
              {t('admin.notifications.role')}
              <select className="input" name="role" disabled={pending}>
                {roles.map((role) => (
                  <option key={role} value={role}>
                    {t(`role.${role}`)}
                  </option>
                ))}
              </select>
            </label>
          )}
          <label className="grid gap-2">
            {t('admin.notifications.subject')}
            <input className="input" name="title" required maxLength={100} disabled={pending} />
          </label>
          <label className="grid gap-2">
            {t('admin.notifications.content')}
            <textarea
              className="input"
              name="content"
              rows={5}
              maxLength={1000}
              disabled={pending}
            />
          </label>
          <p className="text-sm text-secondary-foreground">{t('admin.notifications.formatHint')}</p>
          <div className="grid gap-5 sm:grid-cols-2">
            <label className="grid gap-2">
              {t('admin.notifications.category')}
              <select className="input" name="category" defaultValue="system" disabled={pending}>
                {['system', 'moderation', 'admin'].map((category) => (
                  <option key={category} value={category}>
                    {t(`notifications.category.${category}`)}
                  </option>
                ))}
              </select>
            </label>
            <label className="grid gap-2">
              {t('admin.notifications.type')}
              <input className="input" name="type" maxLength={100} disabled={pending} />
            </label>
          </div>
          <label className="flex items-center gap-3">
            <input type="checkbox" name="urgent" value="true" disabled={pending} />
            {t('admin.notifications.urgent')}
          </label>
          <div className="grid gap-5 sm:grid-cols-2">
            <label className="grid gap-2">
              {t('admin.notifications.iconType')}
              <select
                className="input"
                name="iconType"
                value={iconType}
                onChange={(event) => setIconType(event.target.value)}
                disabled={pending}
              >
                {['emoji', 'fontAwesome', 'image'].map((type) => (
                  <option key={type} value={type}>
                    {t(`admin.notifications.icon.${type}`)}
                  </option>
                ))}
              </select>
            </label>
            <label className="grid gap-2">
              {t('admin.notifications.iconOverride')}
              <input
                className="input"
                name="icon"
                maxLength={iconType === 'emoji' ? 32 : 256}
                disabled={pending}
              />
            </label>
          </div>
          <p className="text-sm text-secondary-foreground">
            {t(`admin.notifications.iconHint.${iconType}`)}
          </p>
          {target === 'role' && (
            <p className="text-secondary-foreground">{t('admin.notifications.broadcastHint')}</p>
          )}
          {!!error && <Notice error>{errorMessage(error)}</Notice>}
          {count !== undefined && (
            <div role="status">
              <Notice>
                {t(count ? 'admin.notifications.sent' : 'admin.notifications.noRecipients')}{' '}
                {count > 0 && count}
              </Notice>
            </div>
          )}
          <div>
            <button className="Btn Btn--primary" disabled={pending}>
              {t(pending ? 'admin.notifications.sending' : 'admin.notifications.send')}
            </button>
          </div>
        </form>
      </Card>
    </>
  )
}
