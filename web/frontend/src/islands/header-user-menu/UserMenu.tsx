import { AnchorHTMLAttributes } from 'react'
import { Popover } from 'radix-ui'
import { useUserSelfData } from '@/api/auth/user'
import '@/components/header-popover.scss'
import './UserMenu.scss'

export type UserMenuLink = 'bookManager' | 'moderation' | 'admin'

type UserMenuProps = {
  hiddenLinks?: UserMenuLink[]
}

export function UserMenu({ hiddenLinks }: UserMenuProps) {
  const user = useUserSelfData()

  return (
    <Popover.Root>
      <Popover.Trigger
        className="HeaderPopover-trigger UserMenu-trigger"
        aria-label={window._('common.profile')}
      >
        <img className="avatar" src={user.avatar.md} alt="" />
      </Popover.Trigger>
      <Popover.Portal>
        <Popover.Content
          className="HeaderPopover UserMenu"
          align="end"
          sideOffset={10}
          collisionPadding={12}
          aria-label={window._('common.profile')}
        >
          <UserMenuBody hiddenLinks={hiddenLinks} />
        </Popover.Content>
      </Popover.Portal>
    </Popover.Root>
  )
}

function UserMenuBody({ hiddenLinks = [] }: { hiddenLinks?: UserMenuLink[] }) {
  const user = useUserSelfData()
  const isAdmin = user.role === 'admin' || user.role === 'system'
  const isModerator = isAdmin || user.role === 'moderator'

  const pathname = window.location.pathname
  const visible = (key: UserMenuLink, path: string) =>
    !hiddenLinks.includes(key) && !(pathname === path || pathname.startsWith(`${path}/`))

  const showWorkspaces =
    visible('bookManager', '/books-manager') ||
    (isModerator && visible('moderation', '/moderation')) ||
    (isAdmin && visible('admin', '/admin'))

  return (
    <div className="UserMenu-body">
      <div className="UserMenu-identity">
        <div className="UserMenu-account">
          <div className="UserMenu-name">{user.name}</div>
          {user.email && <div className="UserMenu-email">{user.email}</div>}
        </div>
        <img className="avatar UserMenu-avatar" src={user.avatar.md} alt="" />
      </div>

      <ul className="UserMenu-shortcuts">
        <MenuItem icon="fa-solid fa-user" href={`/users/${user.id}`}>
          {window._('common.profile')}
        </MenuItem>
        <MenuItem icon="fa-solid fa-gear" href="/account/settings">
          {window._('common.settings')}
        </MenuItem>
      </ul>
      {showWorkspaces && (
        <ul className="UserMenu-items">
          {visible('bookManager', '/books-manager') && (
            <MenuItem icon="fa-solid fa-book-open" href="/books-manager">
              {window._('common.bookManager')}
            </MenuItem>
          )}
          {isModerator && visible('moderation', '/moderation') && (
            <MenuItem icon="fa-solid fa-shield" href="/moderation/">
              {window._('common.moderationPortal')}
            </MenuItem>
          )}
          {isAdmin && visible('admin', '/admin') && (
            <MenuItem icon="fa-solid fa-shield-halved" href="/admin">
              {window._('common.adminDashboard')}
            </MenuItem>
          )}
        </ul>
      )}
      <ul className="UserMenu-footer">
        <MenuItem icon="fa-solid fa-right-from-bracket" href="/logout" logout>
          {window._('common.logout')}
        </MenuItem>
      </ul>
    </div>
  )
}

function MenuItem({
  icon,
  children,
  logout = false,
  ...props
}: { icon: string; logout?: boolean } & AnchorHTMLAttributes<HTMLAnchorElement>) {
  return (
    <li>
      <Popover.Close asChild>
        <a className={`UserMenu-item${logout ? ' UserMenu-item--logout' : ''}`} {...props}>
          <i className={`${icon} UserMenu-itemIcon`} aria-hidden="true" />
          <span className="UserMenu-itemLabel">{children}</span>
          {!logout && (
            <i className="fa-solid fa-chevron-right UserMenu-itemArrow" aria-hidden="true" />
          )}
        </a>
      </Popover.Close>
    </li>
  )
}
