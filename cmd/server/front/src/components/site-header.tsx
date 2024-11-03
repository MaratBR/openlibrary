import React from 'react'
import './site-header.css'
import { NavLink, NavLinkRenderProps } from 'react-router-dom'

export default function SiteHeader({ children }: React.PropsWithChildren<{}>) {
  return (
    <header id="site-header" className="sticky top-0 z-50 w-full border-border/40 site-header">
      <div className="container-default h-14 flex items-center">
        <div id="site-logo" className="mr-8 font-semibold text-4xl">
          OpenLibrary
        </div>
        <div className="space-x-2 flex">
          <SiteHeaderLink href="/home">Home</SiteHeaderLink>
          <SiteHeaderLink href="/my-books">Your stories</SiteHeaderLink>
          <SiteHeaderLink href="/logout">Logout</SiteHeaderLink>
          <SiteHeaderLink href="/about">About</SiteHeaderLink>
        </div>
      </div>
    </header>
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
