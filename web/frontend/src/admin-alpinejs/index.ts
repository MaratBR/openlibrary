import Alpine from 'alpinejs'
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
//@ts-expect-error
import ajax from '@imacrayon/alpine-ajax'
// eslint-disable-next-line @typescript-eslint/no-explicit-any
;(window as any).Alpine = Alpine
Alpine.plugin(ajax)

if (!new URLSearchParams(window.location.search).has('debug.disableAlpineJS')) {
  Alpine.start()
}
