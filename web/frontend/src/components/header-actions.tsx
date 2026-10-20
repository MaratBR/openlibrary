import { useEffect, useState } from 'react'
import { useUserSelfData } from '@/api/auth/user'
import { Bell, NotificationsMenu } from '@/features/notifications/Menu'
import { UserMenu, type UserMenuLink } from '@/islands/header-user-menu/UserMenu'

/** Account and activity share one mount while keeping independent panels. */
export function HeaderActions({ hiddenLinks }: { hiddenLinks?: UserMenuLink[] }) {
  // Match navUser's initial links before attaching Radix popovers and their IDs.
  const [ready, setReady] = useState(false)
  const user = useUserSelfData()
  useEffect(() => setReady(true), [])

  return (
    <div className="flex items-center gap-2">
      {ready ? (
        <UserMenu hiddenLinks={hiddenLinks} />
      ) : (
        <a
          className="HeaderPopover-trigger UserMenu-trigger"
          href={`/users/${user.id}`}
          aria-label={window._('common.profile')}
        >
          <img className="avatar" src={user.avatar.md} alt="" />
        </a>
      )}{' '}
      {ready ? (
        <NotificationsMenu />
      ) : (
        <a
          className="HeaderPopover-trigger NotificationMenu-trigger"
          href="/notifications"
          aria-label={window._('notifications.title')}
        >
          <Bell />
        </a>
      )}
    </div>
  )
}
