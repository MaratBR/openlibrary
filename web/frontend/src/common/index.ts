import './i18n'

import { OverlayScrollbars } from 'overlayscrollbars'
import 'overlayscrollbars/overlayscrollbars.css'

// eslint-disable-next-line @typescript-eslint/no-explicit-any
;(window as any).OverlayScrollbars = OverlayScrollbars
document.dispatchEvent(new CustomEvent('overlayscrollbars-ready'))

import './ol-global'
import './delay'
import '../lib/island'
import '../toast/toast'

import '@/http-client'
import './common.css'

import { initAfterDOMReady } from './links'
initAfterDOMReady()
