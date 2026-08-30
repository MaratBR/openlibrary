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

  return (
    <div className="BeFontPickerMini">
      {favoriteFonts.map((font) => {
        return (
          <button key={font} className="BeFontPickerMini-item" style={{ fontFamily: font }}>
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
        className="BeFontPickerMini-more"
      >
        {window._('common.more')}
      </button>
    </div>
  )
}
