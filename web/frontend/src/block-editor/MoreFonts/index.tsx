import Modal from '@/components/Modal'
import { appRuntime } from '@/effect/runtime'
import { create } from 'zustand'
import { FontsLoader } from '@/features/fonts-loader/loader'
import { Font } from '@/features/fonts-loader/api'
import { useVirtualizer } from '@tanstack/react-virtual'
import React, { useLayoutEffect, useMemo, useRef, useState } from 'react'

import './MoreFonts.scss'
import { ChapterContentEditor } from '../wysiwyg/editor'
import { useFavoriteFontState, useFonts } from '../fonts/state'
import { atom, useAtom, useAtomValue } from 'jotai'

const fontsLoader = appRuntime.runSync(FontsLoader)

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

  const { fonts } = useFonts()

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
          className: 'MoreFonts',
          role: 'dialog',
          'aria-modal': true,
          'aria-labelledby': 'more-fonts-title',
        },
      }}
    >
      <div className="MoreFonts-header">
        <div className="flex items-start justify-between gap-4">
          <div>
            <h2 id="more-fonts-title" className="text-xl font-semibold">
              {window._('editor.fonts.title')}
            </h2>
            <p className="mt-1 text-sm text-secondary-foreground">
              {window._('editor.fonts.description')}
            </p>
          </div>
          <button
            type="button"
            className="MoreFonts-iconButton"
            onClick={close}
            aria-label={window._('editor.fonts.close')}
          >
            <i className="fa-solid fa-xmark" aria-hidden="true" />
          </button>
        </div>
        <label className="MoreFonts-search">
          <i className="fa-solid fa-magnifying-glass" aria-hidden="true" />
          <input
            type="search"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="input"
            aria-label={window._('editor.fonts.search')}
            placeholder={window._('editor.fonts.search')}
          />
        </label>
        <Toggles />
      </div>
      <div className="MoreFonts-listHeading">
        <span>{window._('editor.fonts.available')}</span>
        <span>{filteredFonts.length}</span>
      </div>
      {filteredFonts.length === 0 && (
        <div className="MoreFonts-empty" role="status">
          <i className="fa-solid fa-magnifying-glass text-2xl" aria-hidden="true" />
          <p className="mt-3 font-medium">{window._('editor.fonts.noResults')}</p>
          <p className="mt-1 text-sm">{window._('editor.fonts.tryAnotherSearch')}</p>
        </div>
      )}
      <FontsList fonts={filteredFonts} />
    </Modal>
  )
}

function Toggles() {
  const [bold, setBold] = useAtom(boldEnabledAtom)
  const [italic, setItalic] = useAtom(italicEnabledAtom)
  const [testPhrase, setTestPhrase] = useAtom(testPhraseAtom)

  return (
    <div className="MoreFonts-preview">
      <label className="min-w-0 flex-1">
        <span className="MoreFonts-label">{window._('editor.fonts.preview')}</span>
        <input
          className="input"
          placeholder={window._('editor.fonts.testPhrase')}
          value={testPhrase}
          onChange={(e) => setTestPhrase(e.target.value)}
        />
      </label>
      <div className="MoreFonts-previewToggles">
        <button
          type="button"
          className="MoreFonts-iconButton"
          aria-pressed={bold}
          aria-label={window._('editor.bold')}
          onClick={() => setBold(!bold)}
        >
          <i className="fa-solid fa-bold" aria-hidden="true" />
        </button>
        <button
          type="button"
          className="MoreFonts-iconButton"
          aria-pressed={italic}
          aria-label={window._('editor.italic')}
          onClick={() => setItalic(!italic)}
        >
          <i className="fa-solid fa-italic" aria-hidden="true" />
        </button>
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
      className="MoreFonts-list"
      role="list"
      aria-label={window._('editor.fonts.available')}
      style={{ '--item-height': `${HEIGHT}px` } as React.CSSProperties}
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
      void appRuntime.runPromise(fontsLoader.ensureFontLoaded(font.name)).catch(() => undefined)
    }, 600)

    return () => clearTimeout(t)
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
          '--preview-weight': bold ? 'bold' : 'normal',
          '--preview-style': italic ? 'italic' : 'normal',
        } as React.CSSProperties
      }
    >
      <div className="BeFontPickerItem-aa apply-font">Aa</div>

      <div className="BeFontPickerItem-main">
        <span className="BeFontPickerItem-name apply-font">{testPhrase || font.name}</span>
        <span className="BeFontPickerItem-nameNormal">{font.name}</span>
      </div>

      <div className="BeFontPickerItem-star">
        <FavoriteFontButton font={font.name} />
      </div>
    </div>
  )
}

function FavoriteFontButton({ font }: { font: string }) {
  const { select, selected } = useFavoriteFontState(font)

  return (
    <button
      className="BeFontPickerItem-starButton"
      type="button"
      aria-pressed={selected}
      aria-label={`${window._(selected ? 'editor.fonts.removeFavorite' : 'editor.fonts.addFavorite')}: ${font}`}
      title={window._(selected ? 'editor.fonts.removeFavorite' : 'editor.fonts.addFavorite')}
      onClick={() => {
        select(!selected)
      }}
    >
      {selected ? <i className="fa-solid fa-star"></i> : <i className="fa-regular fa-star"></i>}
    </button>
  )
}
