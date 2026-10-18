import { ReactNode, useId } from 'react'
import { Form, Link, useLocation, useNavigation, useRouteError } from 'react-router'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { ErrorDisplay } from '@/components/error'
import { Pagination } from '@/components/Pagination'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components'

export const t = (key: string) => window._(`moderationPortal.workspace.${key}`)
export const date = (value: string) =>
  new Intl.DateTimeFormat(undefined, { dateStyle: 'medium', timeStyle: 'short' }).format(
    new Date(value),
  )
export function WorkspacePage({
  title,
  children,
  actions,
}: {
  title: string
  children: ReactNode
  actions?: ReactNode
}) {
  const navigation = useNavigation()
  return (
    <DashboardContent.Root>
      <DashboardContent.Header title={title}>{actions}</DashboardContent.Header>
      <div className="ModWorkspace" aria-busy={navigation.state !== 'idle'}>
        {navigation.state !== 'idle' && <p role="status">{t('loading')}</p>}
        {children}
      </div>
    </DashboardContent.Root>
  )
}
export function EmptyState({ children }: { children: ReactNode }) {
  return (
    <p className="ModWorkspace-empty" role="status">
      {children}
    </p>
  )
}
export function WorkspaceError() {
  const error = useRouteError()
  return (
    <WorkspacePage title={t('loadError')}>
      <div role="alert">
        <ErrorDisplay error={error} />
      </div>
      <div className="flex flex-wrap gap-3">
        <button className="Btn Btn--primary" onClick={() => window.location.reload()}>
          {t('retry')}
        </button>
        <Link className="Btn Btn--outline" to="/overview">
          {t('overview')}
        </Link>
      </div>
    </WorkspacePage>
  )
}
export function SearchForm({ children }: { children: ReactNode }) {
  const location = useLocation()
  return (
    <Form method="get" key={location.search} className="ModWorkspace-filters">
      {children}
      <button className="Btn Btn--primary" type="submit">
        {t('apply')}
      </button>
      <Link className="Btn Btn--ghost" to={location.pathname}>
        {t('reset')}
      </Link>
    </Form>
  )
}
export function FilterInput({
  name,
  value,
  label,
  type = 'text',
}: {
  name: string
  value?: string
  label: string
  type?: string
}) {
  const id = useId()
  return (
    <div className="ModWorkspace-field">
      <label htmlFor={id}>{label}</label>
      <input id={id} className="input" type={type} name={name} defaultValue={value ?? ''} />
    </div>
  )
}
export function FilterSelect({
  name,
  value,
  label,
  options,
}: {
  name: string
  value?: string
  label: string
  options: [string, string][]
}) {
  const id = useId()
  return (
    <div className="ModWorkspace-field">
      <label htmlFor={id}>{label}</label>
      <Select name={name} defaultValue={value || 'all'}>
        <SelectTrigger id={id}>
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {options.map(([v, text]) => (
            <SelectItem key={v} value={v}>
              {text}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  )
}
export function StateBadge({ value }: { value: string }) {
  return (
    <span className="ModWorkspace-status" data-state={value}>
      {t(value)}
    </span>
  )
}
export function PageNavigation({ page, totalPages }: { page: number; totalPages: number }) {
  const location = useLocation()
  return totalPages > 1 ? (
    <div className="flex justify-center">
      <Pagination.Facade
        page={page}
        totalPages={totalPages}
        size={5}
        getTo={(page) => {
          const params = new URLSearchParams(location.search)
          params.set('page', String(page))
          return { search: params.toString() }
        }}
      />
    </div>
  ) : null
}
export function filtersFrom(request: Request) {
  return Object.fromEntries(new URL(request.url).searchParams)
}
export const targetOptions = (): [string, string][] =>
  ['all', 'user', 'book', 'chapter', 'comment', 'report'].map((v) => [v, t(v)])
