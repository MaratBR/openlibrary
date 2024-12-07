import React from 'react'
import './site-header.css'
import { NavLink, NavLinkRenderProps } from 'react-router-dom'
import UserProfileButton from './user-profile-button'

export default function SiteHeader({ children }: React.PropsWithChildren) {
  return (
    <>
      <header id="site-header" className="sticky top-0 w-full z-50 border-border/40 site-header">
        <div className="container-default h-full grid grid-cols-[auto_1fr_auto] items-stretch">
          <div id="site-logo" className="mr-8 font-semibold text-4xl self-center">
            OpenLibrary
          </div>
          <div className="site-header__links">
            <SiteHeaderLink href="/home">Home</SiteHeaderLink>
            <SiteHeaderLink href="/manager/books">Your stories</SiteHeaderLink>
            <SiteHeaderLink href="/logout">Logout</SiteHeaderLink>
            <SiteHeaderLink href="/search">Search</SiteHeaderLink>
          </div>

          <div className="self-center">
            <UserProfileButton />
          </div>
        </div>
      </header>
    </>
  )
}

function SiteHeaderLink({ children, href }: React.PropsWithChildren<{ href: string }>) {
  return (
    <NavLink to={href} className={navLinkClass}>
      {children}
    </NavLink>
  )
}

const navLinkClass = ({ isActive }: NavLinkRenderProps) =>
  'site-header__nav-link' + (isActive ? ' site-header__nav-link--active' : '')
