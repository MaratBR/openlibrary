import React, { useRef, useState } from 'react'
import { create } from 'zustand'

export type NotificationDescriptor = {
  text: string
  id: string
}

export type Notification = {
  text: string
  id: string
  addedAt: number
}

export type NotificationsState = {
  notifications: Notification[]

  add(descriptor: NotificationDescriptor): void
  remove(...ids: string[]): void
}

function combine(oldList: Notification[], newList: Notification[]): Notification[] {
  return [
    ...oldList.filter((old) => !newList.some((newNotif) => newNotif.id === old.id)),
    ...newList,
  ]
}

export const useNotificationsState = create<NotificationsState>()((set) => ({
  notifications: [],

  add(descriptor) {
    set((state) => ({
      notifications: combine(state.notifications, [{ ...descriptor, addedAt: Date.now() }]),
    }))
  },

  remove(...ids) {
    set((state) => ({
      notifications: state.notifications.filter((n) => !ids.includes(n.id)),
    }))
  },
}))

export function useNotificationsSlot() {
  const [notifications, setNotifications] = useState([] as NotificationDescriptor[])
  const activeNotifications = useRef<string[]>([])

  React.useEffect(() => {
    const state = useNotificationsState.getState()

    for (const notif of notifications) {
      state.add(notif)
    }
    activeNotifications.current = notifications.map((x) => x.id)

    return () => {
      if (activeNotifications.current.length) {
        state.remove(...activeNotifications.current)
      }
    }
  }, [notifications])

  return setNotifications
}
