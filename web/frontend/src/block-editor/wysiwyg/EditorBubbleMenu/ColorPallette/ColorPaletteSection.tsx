import { ChapterContentEditor } from '../../editor'
import { EditButtonMenuExpandableSection } from '../EditBubbleMenuExpandableSection'
import { useBubbleState } from '../state'
import { ColorPallettePicker } from './ColorPalettePicker'

export function ColorPalletteSection({ editor }: { editor: ChapterContentEditor }) {
  const open = useBubbleState((x) => x.colorPickerOpen)

  return (
    <EditButtonMenuExpandableSection>
      {open && (
        <div className="p2">
          <ColorPallettePicker onColorChosen={handleColor} />
        </div>
      )}
    </EditButtonMenuExpandableSection>
  )

  function handleColor(color: { r: number; g: number; b: number } | null) {
    if (color) {
      editor
        .chain()
        .focus()
        .setColor(rgbToHex(color.r, color.g, color.b))
        .run()
    } else {
      editor.chain().focus().unsetColor().run()
    }
  }
}

function rgbToHex(r: number, g: number, b: number): string {
  return '#' + [r, g, b].map((value) => value.toString(16).padStart(2, '0')).join('')
}
