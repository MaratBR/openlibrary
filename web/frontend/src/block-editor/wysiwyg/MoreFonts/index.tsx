import Modal from '@/components/Modal'
import { create } from 'zustand'
import { ChapterContentEditor } from '../editor'
import { useQuery } from '@tanstack/react-query'
import { FontsLoader } from '@/features/fonts-loader/loader'
import { Font } from '@/features/fonts-loader/api'
import { useVirtualizer } from '@tanstack/react-virtual'
import React, { useLayoutEffect, useMemo, useRef, useState } from 'react'

import './MoreFonts.scss'

export type MoreFontsState = {
  opened: boolean
  editorState: null | {
    editor: ChapterContentEditor
    selection: { from: number; to: number }
  }

  open(editor: ChapterContentEditor): void
  close(): void
}

export const useMoreFontsState = create<MoreFontsState>()((set) => ({
  opened: false,
  editorState: null,

  open(editor) {
    set({
      opened: true,
      editorState: {
        editor,
        selection: { from: editor.state.selection.from, to: editor.state.selection.to },
      },
    })
  },

  close() {
    set({ opened: false })
  },
}))

export function MoreFonts() {
  const opened = useMoreFontsState((x) => x.opened)
  const close = useMoreFontsState((x) => x.close)

  const [search, setSearch] = useState('')

  const { data: fonts } = useQuery({
    queryFn: () => FontsLoader.fetchFonts(),
    queryKey: ['FontLoader-loadFonts'],
    staleTime: 0,
    gcTime: Infinity,
    initialData: [],
  })

  const filteredFonts = useMemo(() => {
    const trimmed = search.trim().toLowerCase()

    if (trimmed === '') return fonts

    return fonts.filter((font) => font.name.toLowerCase().includes(trimmed))
  }, [fonts, search])

  return (
    <Modal
      onClose={() => close()}
      open={opened}
      slotProps={{
        content: {
          className: 'p-0 rounded-none',
        },
      }}
    >
      <div className="bg-surface border-b">
        <input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="text-3xl pl-8 pt-8 pb-4 outline-none"
          placeholder="Search fonts"
        />
      </div>
      <FontsList fonts={filteredFonts} />
    </Modal>
  )
}

function FontsList({ fonts }: { fonts: ReadonlyArray<Readonly<Font>> }) {
  const parentRef = useRef<HTMLDivElement | null>(null)

  const HEIGHT = 80

  const virtualizer = useVirtualizer({
    count: fonts.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => HEIGHT,
  })

  return (
    <div
      ref={parentRef}
      style={
        { height: '500px', overflow: 'auto', '--item-height': `${HEIGHT}px` } as React.CSSProperties
      }
    >
      <div
        className="relative w-full"
        style={{
          height: `${virtualizer.getTotalSize()}px`,
        }}
      >
        {virtualizer.getVirtualItems().map((virtualItem) => {
          const font = fonts[virtualItem.index]

          return (
            <div
              key={virtualItem.key}
              style={{
                position: 'absolute',
                top: 0,
                left: 0,
                width: '100%',
                height: `${virtualItem.size}px`,
                transform: `translateY(${virtualItem.start}px)`,
              }}
            >
              <FontRow font={font} />
            </div>
          )
        })}
      </div>
    </div>
  )
}

function FontRow({ font }: { font: Readonly<Font> }) {
  useLayoutEffect(() => {
    const t = setTimeout(() => {
      FontsLoader.ensureFontLoaded(font.name)
    }, 600)

    return () => clearInterval(t)
  }, [font.name])

  return (
    <div
      className="be-font-picker-item"
      data-font={font.name}
      style={{ '--font-family': font.name } as React.CSSProperties}
    >
      <div className="be-font-picker-item__aa apply-font">Aa</div>

      <div className="be-font-picker-item__main">
        <span className="be-font-picker-item__name apply-font">{font.name}</span>
        <span className="be-font-picker-item__name-normal">{font.name}</span>
      </div>
    </div>
  )
}
