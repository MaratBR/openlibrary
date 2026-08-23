import { create } from 'zustand'

export type BubbleState = {
  colorPickerOpen: boolean
  fontPickerOpen: boolean

  toggleColorPicker(): void
  toggleFontPicker(): void
}

export const useBubbleState = create<BubbleState>()((set) => ({
  colorPickerOpen: false,
  fontPickerOpen: false,

  toggleColorPicker() {
    set((x) => ({ ...x, colorPickerOpen: !x.colorPickerOpen, fontPickerOpen: false }))
  },

  toggleFontPicker() {
    set((x) => ({ ...x, colorPickerOpen: false, fontPickerOpen: !x.fontPickerOpen }))
  },
}))
