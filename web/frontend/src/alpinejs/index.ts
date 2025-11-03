import Alpine from 'alpinejs'
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
//@ts-expect-error
import ajax from '@imacrayon/alpine-ajax'

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
//@ts-expect-error
import collapse from '@alpinejs/collapse'

// alpinejs components
import './book-reader'
import './rating-input'
import './collapse'
import './page-progress'
import './navbar'
;(window as unknown as { Alpine: typeof Alpine }).Alpine = Alpine

Alpine.plugin(ajax)
Alpine.plugin(collapse)
if (!new URLSearchParams(window.location.search).has('debug.disableAlpineJS')) {
  Alpine.start()
}
