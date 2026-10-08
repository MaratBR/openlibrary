import { ReactNode, useEffect } from 'react'
import { useAtom } from 'jotai'
import { recentItemsAtom, rememberDashboardItem } from './dashboard-recent-items'
import { NavLink, useLocation, useMatches } from 'react-router'

// Route handles supply labels from successfully loaded detail data.
export function dashboardItem<T>(label: (data: T) => string) {
  return { dashboardItem: (data: unknown) => label(data as T) }
}

export function DashboardNavItem({
  to,
  dashboard,
  children,
  end = false,
  exclude,
}: {
  to: string
  dashboard: 'admin' | 'moderation' | 'bookmanager'
  children: ReactNode
  end?: boolean
  exclude?: string
}) {
  const location = useLocation()
  const matches = useMatches()
  const [recent, setRecent] = useAtom(recentItemsAtom(`${dashboard}:${to}`))
  const segments = location.pathname.split('/')
  const itemPath = segments[2] ? `${to}/${segments[2]}` : undefined
  const match = matches.at(-1)
  const handle = match?.handle as ReturnType<typeof dashboardItem> | undefined
  const label =
    location.pathname.startsWith(`${to}/`) &&
    match?.loaderData != null &&
    typeof handle?.dashboardItem === 'function'
      ? handle.dashboardItem(match.loaderData)
      : undefined

  useEffect(() => {
    if (!itemPath || !label) return
    setRecent((items) => rememberDashboardItem(items, { to: itemPath, label }))
  }, [itemPath, label, location.pathname, setRecent])

  return (
    <div className="DashboardShell-navGroup">
      <NavLink
        to={to}
        end={end}
        className={
          exclude
            ? ({ isActive }) => (isActive && location.pathname !== exclude ? 'active' : '')
            : undefined
        }
        aria-current={location.pathname === exclude ? false : undefined}
      >
        {children}
      </NavLink>
      {recent.length > 0 && (
        <ul className="DashboardShell-submenu">
          {recent.map((item) => (
            <li key={item.to}>
              <NavLink to={item.to} title={item.label}>
                <span>{item.label}</span>
              </NavLink>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
