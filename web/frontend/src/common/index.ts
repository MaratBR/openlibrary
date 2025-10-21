import './i18n'
import './util'
import './cookies'
import './theme'
import './scrollable'
import './nav'

import './__server__'
import './delay'
import './flashes'
import '../lib/island'
import '../lib/ScrollBlocker'
import '../toast/toast'

import '@/http-client'
import 'virtual:uno.css'
import './common.css'

import { initAfterDOMReady } from './links'
import { initFirstActivityEvent } from '@/lib/user-activity-detector'
initAfterDOMReady()
initFirstActivityEvent()

requestIdleCallback(() => {
  document.dispatchEvent(new CustomEvent('ol:jsready'))
})
