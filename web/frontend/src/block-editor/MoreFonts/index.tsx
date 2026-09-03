import Modal from '@/components/Modal'
import { appRuntime } from '@/effect/runtime'
import { create } from 'zustand'
import { FontsLoader } from '@/features/fonts-loader/loader'
import { Font } from '@/features/fonts-loader/api'
import { useVirtualizer } from '@tanstack/react-virtual'
import React, { useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react'

import './MoreFonts.scss'
import { ChapterContentEditor } from '../wysiwyg/editor'
import { useFonts } from '../fonts/state'
import { atom, useAtom, useAtomValue } from 'jotai'
import { jotaiStore } from '@/react'
import EditorToggleButton from '../wysiwyg/EditorBubbleMenu/EditorToggleButton'

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

const boldEnabledAtom = atom(false)
const italicEnabledAtom = atom(false)
const testPhraseAtom = atom('')

export function MoreFonts() {
  const opened = useMoreFontsState((x) => x.opened)
  const close = useMoreFontsState((x) => x.close)

  const [search, setSearch] = useState('')

  const { fonts, init } = useFonts()

  useEffect(() => {
    if (opened) void appRuntime.runPromise(init())
  }, [opened, init])

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

        <Toggles />
      </div>
      <FontsList fonts={filteredFonts} />
    </Modal>
  )
}

function Toggles() {
  const [bold, setBold] = useAtom(boldEnabledAtom)
  const [italic, setItalic] = useAtom(italicEnabledAtom)
  const [testPhrase, setTestPhrase] = useAtom(testPhraseAtom)

  return (
    <div className="BeToggleGroup pl-8 pb-1 gap-1">
      <EditorToggleButton active={bold} onClick={() => setBold(!bold)}>
        <i className="fa-solid fa-bold" />
      </EditorToggleButton>
      <EditorToggleButton active={italic} onClick={() => setItalic(!italic)}>
        <i className="fa-solid fa-italic" />
      </EditorToggleButton>
      <div>
        <input
          placeholder={window._('editor.fonts.testPhrase')}
          value={testPhrase}
          onChange={(e) => setTestPhrase(e.target.value)}
        />
      </div>
    </div>
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

  const bold = useAtomValue(boldEnabledAtom)
  const italic = useAtomValue(italicEnabledAtom)
  const testPhrase = useAtomValue(testPhraseAtom).trim()

  return (
    <div
      className="BeFontPickerItem"
      data-font={font.name}
      role="listitem"
      style={
        {
          '--font-family': font.name,
          fontWeight: bold ? 'bold' : 'normal',
          fontStyle: italic ? 'italic' : 'normal',
        } as React.CSSProperties
      }
    >
      <div className="BeFontPickerItem-aa apply-font">Aa</div>

      <div className="BeFontPickerItem-main">
        <span className="BeFontPickerItem-name apply-font">{testPhrase || font.name}</span>
        <span className="BeFontPickerItem-nameNormal">{font.name}</span>
      </div>
    </div>
  )
}
