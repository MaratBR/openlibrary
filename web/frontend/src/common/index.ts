import '../lib/island'
import '../lib/ScrollBlocker'
import '../toast/toast'
import './__server__'
import './cookies'
import './delay'
import './error'
import './flashes'
import './i18n'
import './nav'
import './scrollable'
import './theme'
import './util'
import '@/http-client'

import './style'
import './theme.css'

import { initAfterDOMReady } from './links'
import { initFirstActivityEvent } from '@/lib/user-activity-detector'
initAfterDOMReady()
initFirstActivityEvent()

requestIdleCallback(() => {
  document.dispatchEvent(new CustomEvent('ol:jsready'))
})
