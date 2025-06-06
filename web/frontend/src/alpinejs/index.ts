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

// eslint-disable-next-line @typescript-eslint/no-explicit-any
;(window as any).Alpine = Alpine

Alpine.plugin(ajax)
Alpine.plugin(collapse)
Alpine.start()
