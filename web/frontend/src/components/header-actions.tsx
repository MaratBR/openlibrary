import { NotificationsMenu } from '@/features/notifications/Menu'
import { UserMenu, type UserMenuLink } from '@/islands/header-user-menu/UserMenu'

/** Account and activity share one mount while keeping independent panels. */
export function HeaderActions({ hiddenLinks }: { hiddenLinks?: UserMenuLink[] }) {
  return (
    <div className="flex items-center gap-2">
      <UserMenu hiddenLinks={hiddenLinks} />
      <NotificationsMenu />
    </div>
  )
}
