import { ChapterContentEditor } from '../../editor'
import { useMoreFontsState } from '../../MoreFonts'
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

  return (
    <div className="be-font-picker-mini">
      {favoriteFonts.map((font) => {
        return (
          <button key={font} className="be-font-picker-mini__item" style={{ fontFamily: font }}>
            {font}
          </button>
        )
      })}

      <button
        onClick={() => {
          openMoreFonts(editor)
          editor.chain().setTextSelection(editor.state.selection.from).blur().run()
          toggle()
        }}
        className="be-font-picker-mini__more"
      >
        {window._('common.more')}
      </button>
    </div>
  )
}
