import { render } from 'preact'
import type { OLNotification } from '@/http-client'
import SanitizeHTML from './SanitizeHTML'
import { Subject, useSubject } from './rx'

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
          <div key={i} class="ol-flash" data-type={notif.type}>
            <span>
              <SanitizeHTML value={notif.text} />
            </span>
            <div class="ol-flash__closeContainer">
              <button onClick={() => Notifications.instance.remove(notif)} class="ol-flash__close">
                <i class="fa-solid fa-xmark"></i>
              </button>
            </div>
          </div>
        )
      })}
    </>
  )
}

export function initFlashMessages() {
  const element = document.getElementById('client-flashes')

  if (!element) {
    throw new Error('cannot initialize flash messages: #client-flashes element not found')
  }

  render(<FlashesHost />, element)
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
