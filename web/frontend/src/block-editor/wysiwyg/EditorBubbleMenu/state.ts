import { create } from 'zustand'

export type BubbleState = {
  colorPickerOpen: boolean

  toggleColorPicker(): void
}

export const useBubbleState = create<BubbleState>()((set) => ({
  colorPickerOpen: false,

  toggleColorPicker() {
    set((x) => ({ ...x, colorPickerOpen: !x.colorPickerOpen }))
  },
}))
