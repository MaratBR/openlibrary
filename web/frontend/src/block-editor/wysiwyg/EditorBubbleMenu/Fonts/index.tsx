import { useBubbleState } from '../state'

import { FontPicker } from './FontPicker'

export { FontPicker }

export function Fonts() {
  const toggle = useBubbleState((x) => x.toggleFontPicker)

  return (
    <button className="be-listitem be-listitem--btn" onClick={toggle}>
      <i className="fa-solid fa-font"></i>{' '}
    </button>
  )
}
