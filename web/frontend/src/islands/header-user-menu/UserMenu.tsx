import { useEffect, useRef, useState } from 'react'
import { useUserSelfData } from '@/api/auth/user'
import './UserMenu.scss'
import { ModalAnimation, useAnimation } from '@/lib/animate'
import { useForkRef } from '@/lib/ref'

export function UserMenu() {
  const ref = useRef<HTMLDivElement | null>(null)
  const [open, setOpen] = useState(false)

  useEffect(() => {
    const $btn = document.querySelector('#nav-user > button')
    if (!$btn) {
      console.error('Cannot find user header button')
      return
    }

    const onClick = () => {
      setOpen(true)
    }
    $btn.addEventListener('click', onClick)
    return () => {
      $btn.removeEventListener('click', onClick)
    }
  }, [])

  useEffect(() => {
    if (!open) return

    const onClickOutside = (event: MouseEvent) => {
      console.log(ref.current)
      if (event.target instanceof Node && ref.current?.contains(event.target)) {
        return
      }

      setOpen(false)
    }

    window.addEventListener('click', onClickOutside)
    return () => window.removeEventListener('click', onClickOutside)
  }, [open])

  const { ref: animationRef } = useAnimation({
    animation: ModalAnimation.default,
    show: open,
  })

  const finalRef = useForkRef(animationRef, ref)

  return (
    <div className="UserMenu" ref={finalRef}>
      <UserMenuBody />
    </div>
  )
}

function UserMenuBody() {
  const user = useUserSelfData()

  return (
    <div>
      <div className="flex justify-center py-4">
        <img className="avatar size-16 shadow-2xl" src={user.avatar.md} />
      </div>

      <ul className="UserMenu__items">
        <a className="UserMenu__item" role="listitem" href={`/user/${user.id}`}>
          {window._('common.profile')}
        </a>

        <a className="UserMenu__item" role="listitem" href="/profile/settings">
          {window._('common.settings')}
        </a>

        <a className="UserMenu__item" role="listitem" href="/books-manager">
          {window._('common.bookManager')}
        </a>
      </ul>
    </div>
  )
}
