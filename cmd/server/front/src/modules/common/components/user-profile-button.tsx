import { useAuthState } from '@/modules/auth/state'
import './user-profile-button.css'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { NavLink } from 'react-router-dom'
import { LogIn, LogOut, Settings, User } from 'lucide-react'
import React from 'react'
import ThemeSwitcherButton from './theme-switcher-button'

export default function UserProfileButton() {
  const user = useAuthState((s) => s.user)

  const [open, setOpen] = React.useState(false)

  const closePopover = () => setOpen(false)

  if (!user) {
    return (
      <NavLink to="/login" className="menu-item">
        <LogIn /> Log in
      </NavLink>
    )
  }

  return (
    <Popover modal onOpenChange={setOpen} open={open}>
      <PopoverTrigger asChild>
        <button className="user-profile-button">
          <img className="user-profile-button__img" src={user.avatar.md} />
        </button>
      </PopoverTrigger>
      <PopoverContent className="p-0" align="end">
        <div className="p-2 space-y-2">
          <NavLink to={`/users/${user.id}`} className="menu-item" onClick={closePopover}>
            <User /> Profile
          </NavLink>
          <NavLink to={`/account/settings`} className="menu-item" onClick={closePopover}>
            <Settings /> Account settings
          </NavLink>
          <ThemeSwitcherButton />
        </div>
        <hr />
        <div className="p-2 space-y-2">
          <NavLink to={`/logout`} className="menu-item" onClick={closePopover}>
            <LogOut /> Logout
          </NavLink>
        </div>
      </PopoverContent>
    </Popover>
  )
}
