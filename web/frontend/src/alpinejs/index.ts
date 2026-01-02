import Alpine from 'alpinejs'
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
//@ts-expect-error
import ajax from '@imacrayon/alpine-ajax'

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
//@ts-expect-error
import collapse from '@alpinejs/collapse'

// alpinejs components
import './BookReader'
import './rating-input'
import './Collapse'
import './page-progress'
import './passwordInput'
import './navbar'
import './SimpleEditor'
import './TagsAutocomplete'
import './Tabs'
import './ImageUploader'
import './Popover'

import './island'
import './islands'

Alpine.plugin(ajax)
Alpine.plugin(collapse)
if (!new URLSearchParams(window.location.search).has('debug.disableAlpineJS')) {
  Alpine.start()
}
;(window as unknown as { Alpine: typeof Alpine }).Alpine = Alpine
