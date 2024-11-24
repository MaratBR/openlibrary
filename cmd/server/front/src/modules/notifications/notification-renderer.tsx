import { NotificationContent } from './notification-content'
import { Notification, useNotificationsState } from './state'
import './notification-renderer.css'

export function NotificationRenderer() {
  const notifications = useNotificationsState((s) => s.notifications)

  if (notifications.length === 0) return null

  return (
    <div className="container-default space-y-0.5 px-1 md:px-8" data-testid="notifications">
      {notifications.map((notif) => (
        <RenderNotification key={notif.id} notification={notif} />
      ))}
    </div>
  )
}

function RenderNotification({ notification }: { notification: Notification }) {
  return (
    <div
      id={`notification-${notification.id}`}
      data-testid={`notification-${notification.id}`}
      className="notification p-2 bg-yellow-400/50 border border-yellow-600"
    >
      <div>
        <NotificationContent content={notification.text} />
      </div>
    </div>
  )
}
