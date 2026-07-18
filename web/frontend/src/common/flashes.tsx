
import type { OLNotification } from '@/http-client'
import SanitizeHTML from './SanitizeHTML'
import { Subject, useSubject } from './rx'
import { createRoot } from 'react-dom/client';

class Notifications extends Subject<OLNotification[]> {
  constructor() {
    super([])
  }

  remove(notification: OLNotification) {
    const value = this.get()
    const idx = value.indexOf(notification)
    if (idx !== -1) {
      const newValue = [...value]
      newValue.splice(idx, 1)
      this.set(newValue)
    }
  }

  add(notification: OLNotification) {
    this.set([...this.get(), notification])
  }

  public static instance: Notifications = new Notifications()
}

function FlashesHost() {
  const notifications = useSubject(Notifications.instance)

  return (
    <>
      {notifications.map((notif, i) => {
        return (
          // TODO proper key value
          <div key={i} className="ol-flash" data-type={notif.type}>
            <span>
              <SanitizeHTML value={notif.text} />
            </span>
            <div className="ol-flash__closeContainer">
              <button onClick={() => Notifications.instance.remove(notif)} className="ol-flash__close">
                <i className="fa-solid fa-xmark"></i>
              </button>
            </div>
          </div>
        )
      })}
    </>
  )
}

export function initflashes() {
  const element = document.getElementById('client-flashes')

  if (!element) {
    throw new Error('cannot initialize flash messages: #client-flashes element not found')
  }

  const root = createRoot(element)
  root.render(<FlashesHost />)
}

declare global {
  function flash(notif: OLNotification): void
  function flash(text: string): void
  function flash(text: string, type: OLNotification['type']): void

  interface Window {
    flash: typeof flash
  }
}

// Implementation of flash
const flashFunc = (...args: [string] | [OLNotification] | [string, OLNotification['type']]) => {
  if (args.length >= 2 && typeof args[0] === 'string' && typeof args[1] === 'string') {
    // Handle (text: string, type: OLNotification['type'])
    const [text, type] = args
    Notifications.instance.add({
      text,
      type,
    })
    return
  }

  if (args.length >= 1) {
    if (typeof args[0] === 'string') {
      // Handle (text: string)
      Notifications.instance.add({
        text: args[0],
        type: 'info', // Default to 'info' type
      })
    } else if (typeof args[0] === 'object') {
      // Handle (notif: OLNotification)
      Notifications.instance.add(args[0])
    }
  }
} // Type assertion to ensure we match the overloads

window.flash = flashFunc
