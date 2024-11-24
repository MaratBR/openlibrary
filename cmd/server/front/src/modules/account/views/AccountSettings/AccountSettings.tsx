import { useAuthState } from '@/modules/auth/state'
import {
  Bell,
  Book,
  ExternalLink,
  Gavel,
  LockKeyhole,
  MessageSquareX,
  Palette,
  User,
  VenetianMask,
} from 'lucide-react'
import { Navigate } from 'react-router'
import { NavLink } from 'react-router-dom'
import './AccountSettings.css'
import React from 'react'
import clsx from 'clsx'
import { useQueryParam } from '@/lib/router-utils'
import AboutMeSettings from './AboutMeSettings'
import CustomizationSettings from './CustomizationSettings'
import PrivacySettings from './PrivacySettings'
import ModerationSettings from './ModerationSettings'

const settingsComponents: Record<string, React.ComponentType<{}>> = {
  about: AboutMeSettings,
  customization: CustomizationSettings,
  privacy: PrivacySettings,
  moderation: ModerationSettings,
}

export default function AccountSettings() {
  const isAuthorized = useAuthState((x) => !!x.user)

  let [activeTab] = useQueryParam('tab')
  activeTab ??= 'about'
  const SettingsComponent = settingsComponents[activeTab]

  if (!isAuthorized) return <Navigate to="/login" />

  return (
    <div className="container-default account-settings">
      <aside className="account-settings__sidebar">
        <div className="account-settings-group">
          <TabLink tab="about">
            <User /> About me
          </TabLink>
          <TabLink tab="notifications">
            <Bell /> Notifications
          </TabLink>
          <TabLink tab="editor">
            <Book /> Editor
          </TabLink>
          <TabLink tab="customization">
            <Palette /> Customization
          </TabLink>
        </div>
        <hr />
        <div className="account-settings-group">
          <div className="account-settings-group__name">Access</div>
          <div className="account-settings-group__description">
            Here you can manage your account security, privacy settings and related settings.
          </div>

          <TabLink tab="security">
            <LockKeyhole /> Security
          </TabLink>
          <TabLink tab="privacy">
            <VenetianMask /> Privacy
          </TabLink>
          <TabLink tab="ban">
            <Gavel /> Ban history
          </TabLink>
          <TabLink tab="moderation">
            <MessageSquareX /> Moderation
          </TabLink>
        </div>
        <hr />
        <div className="account-settings-group">
          <div className="account-settings-group__name">Privileges</div>
          <div className="account-settings-group__description">
            You can see this section of settings because you have a special role.
          </div>

          <NavLink to="/admin" className="account-settings-menu-item">
            <ExternalLink /> Admin dashboard
          </NavLink>
        </div>
      </aside>

      <main className="account-settings__content">
        {SettingsComponent && <SettingsComponent />}
      </main>
    </div>
  )
}

function TabLink({ tab, children }: React.PropsWithChildren<{ tab: string }>) {
  const [activeTab] = useQueryParam('tab')

  const disabled = settingsComponents[tab] === undefined

  return (
    <NavLink
      aria-disabled={disabled}
      onClick={(e) => disabled && e.preventDefault()}
      to={`?tab=${tab}`}
      className={({ isActive }) =>
        clsx('account-settings-menu-item', {
          active: isActive && activeTab === tab,
          disabled: disabled,
        })
      }
    >
      {children}
    </NavLink>
  )
}
