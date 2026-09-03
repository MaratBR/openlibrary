import { useMoreFontsState } from '@/block-editor/MoreFonts'
import { ChapterContentEditor } from '../../editor'
import { EditButtonMenuExpandableSection } from '../EditBubbleMenuExpandableSection'
import { useBubbleState } from '../state'

import './FontPicker.scss'

export function FontPicker({ editor }: { editor: ChapterContentEditor }) {
  const open = useBubbleState((x) => x.fontPickerOpen)

  return (
    <EditButtonMenuExpandableSection>
      {open && <FontPickerInternal editor={editor} />}
    </EditButtonMenuExpandableSection>
  )
}

function FontPickerInternal({ editor }: { editor: ChapterContentEditor }) {
  const favoriteFonts = ['Poppins', 'Merriweather', 'Literata']
  const toggle = useBubbleState((x) => x.toggleFontPicker)
  const openMoreFonts = useMoreFontsState((x) => x.open)

  function handleApplyFont(font: string) {
    editor.chain().setFontFamily(font).focus().run()
  }

  return (
    <div className="BeFontPickerMini">
      <ul className="BeFontPickerMini-list">
        {favoriteFonts.map((font) => {
          return (
            <li
              role="listitem"
              key={font}
              className="BeFontPickerMini-item"
              style={{ fontFamily: font }}
              onClick={() => handleApplyFont(font)}
            >
              {font}
            </li>
          )
        })}
      </ul>

      <hr />

      <button
        onClick={() => {
          openMoreFonts(editor)
          editor.chain().setTextSelection(editor.state.selection.from).blur().run()
          toggle()
        }}
        className="BeFontPickerMini-more"
      >
        {window._('common.more')}
      </button>
    </div>
  )
}
