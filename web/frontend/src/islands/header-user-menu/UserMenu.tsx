import { AnchorHTMLAttributes, useEffect, useLayoutEffect, useRef, useState } from 'react'
import { useUserSelfData } from '@/api/auth/user'
import './UserMenu.scss'
import { ModalAnimation, useAnimation } from '@/lib/animate'
import { useForkRef } from '@/lib/ref'

export function UserMenu() {
  const ref = useRef<HTMLDivElement | null>(null)
  const triggerRef = useRef<HTMLButtonElement | null>(null)
  const [open, setOpen] = useState(false)

  useLayoutEffect(() => {
    const trigger = document.querySelector<HTMLButtonElement>('#nav-user > button')
    if (!trigger) {
      console.error('Cannot find user header button')
      return
    }

    triggerRef.current = trigger
    trigger.setAttribute('aria-haspopup', 'menu')
    const onClick = () => setOpen((value) => !value)
    trigger.addEventListener('click', onClick)
    return () => {
      trigger.removeEventListener('click', onClick)
      triggerRef.current = null
    }
  }, [])

  useLayoutEffect(() => {
    if (!open) return

    const onClickOutside = (event: MouseEvent) => {
      if (
        event.target instanceof Node &&
        (ref.current?.contains(event.target) || triggerRef.current?.contains(event.target))
      ) {
        return
      }

      setOpen(false)
    }

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setOpen(false)
        triggerRef.current?.focus()
      }
    }

    window.addEventListener('click', onClickOutside)
    window.addEventListener('keydown', onKeyDown)
    return () => {
      window.removeEventListener('click', onClickOutside)
      window.removeEventListener('keydown', onKeyDown)
    }
  }, [open])

  useLayoutEffect(() => {
    triggerRef.current?.setAttribute('aria-expanded', String(open))
  }, [open])

  const { ref: animationRef } = useAnimation({
    animation: ModalAnimation.default,
    show: open,
  })

  const finalRef = useForkRef(animationRef, ref)

  return (
    <div className="UserMenu" ref={finalRef} role="menu" aria-hidden={!open}>
      <UserMenuBody onNavigate={() => setOpen(false)} />
    </div>
  )
}

function UserMenuBody({ onNavigate }: { onNavigate: () => void }) {
  const user = useUserSelfData()
  const isAdmin = user.role === 'admin' || user.role === 'system'
  const isModerator = isAdmin || user.role === 'moderator'

  return (
    <div className="UserMenu-body">
      <div className="UserMenu-identity">
        <img className="avatar size-14" src={user.avatar.md} alt="" />
        <div className="min-w-0">
          <div className="font-semibold truncate">{user.name}</div>
          <div className="text-sm text-secondary-foreground truncate">{user.email}</div>
        </div>
      </div>

      <ul className="UserMenu-items">
        <MenuItem icon="fa-regular fa-user" href={`/users/${user.id}`} onClick={onNavigate}>
          {window._('common.profile')}
        </MenuItem>
        <MenuItem icon="fa-solid fa-gear" href="/account/settings" onClick={onNavigate}>
          {window._('common.settings')}
        </MenuItem>
        <MenuItem icon="fa-solid fa-book-open" href="/books-manager" onClick={onNavigate}>
          {window._('common.bookManager')}
        </MenuItem>
        {isModerator && (
          <MenuItem icon="fa-solid fa-shield" href="/moderation/" onClick={onNavigate}>
            {window._('common.moderationPortal')}
          </MenuItem>
        )}
        {isAdmin && (
          <MenuItem
            icon="fa-solid fa-shield-halved"
            href="/admin"
            target="_blank"
            rel="noreferrer"
            onClick={onNavigate}
          >
            {window._('common.adminDashboard')}
            <i className="fa-solid fa-arrow-up-right-from-square ml-auto text-xs" />
          </MenuItem>
        )}
        <li className="UserMenu-separator" aria-hidden="true" />
        <MenuItem icon="fa-solid fa-right-from-bracket" href="/logout" onClick={onNavigate}>
          {window._('common.logout')}
        </MenuItem>
      </ul>
    </div>
  )
}

function MenuItem({
  icon,
  children,
  ...props
}: { icon: string } & AnchorHTMLAttributes<HTMLAnchorElement>) {
  return (
    <li>
      <a className="UserMenu-item" role="menuitem" {...props}>
        <i className={`${icon} UserMenu-itemIcon`} aria-hidden="true" />
        {children}
      </a>
    </li>
  )
}
