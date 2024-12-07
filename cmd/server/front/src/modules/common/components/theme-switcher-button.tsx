import { setTheme, Theme, useTheme } from '@/lib/dark-mode/dark-mode'
import clsx from 'clsx'
import { MonitorCog, Moon, Sun } from 'lucide-react'
import React from 'react'

export type ThemeSwitcherButtonProps = {
  className?: string
}

const ICONS: Record<Theme, React.ComponentType> = {
  light: Sun,
  dark: Moon,
  system: MonitorCog,
}

export default function ThemeSwitcherButton({ className }: ThemeSwitcherButtonProps) {
  const theme = useTheme()

  const Icon = ICONS[theme]

  function handleClick(event: React.MouseEvent) {
    switch (theme) {
      case 'system':
        setTheme('dark')
        break
      case 'dark':
        setTheme('light')
        break
      case 'light':
        setTheme('system')
        break
    }
    // spark(120, 500, 500, 100, 10, event.clientX)
  }

  return (
    <div
      onClick={handleClick}
      role="button"
      className={clsx('menu-item flex gap-2 w-full', className)}
    >
      <Icon /> {getThemeName(theme)}
    </div>
  )
}

function getThemeName(theme: Theme) {
  switch (theme) {
    case 'system':
      return 'System theme'
    case 'dark':
      return 'Dark'
    case 'light':
      return 'Light'
  }
}
