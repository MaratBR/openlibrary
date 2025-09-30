import './i18n'
import './util'
import './cookies'
import './theme'

import { OverlayScrollbars } from 'overlayscrollbars'
import 'overlayscrollbars/overlayscrollbars.css'

// eslint-disable-next-line @typescript-eslint/no-explicit-any
;(window as any).OverlayScrollbars = OverlayScrollbars
document.dispatchEvent(new CustomEvent('overlayscrollbars-ready'))

import './__server__'
import './delay'
import './flashes'
import '../lib/island'
import '../toast/toast'

import '@/http-client'
import './common.css'

import { initAfterDOMReady } from './links'
import { initFirstActivityEvent } from '@/lib/user-activity-detector'
initAfterDOMReady()
initFirstActivityEvent()
