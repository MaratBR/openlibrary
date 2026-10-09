import { useEffect, useMemo } from 'react'
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/Select'
import { useAtomValue } from 'jotai'
import type { Level } from '@tiptap/extension-heading'
import { wysiwygEditorAtom } from './state'
import { ChapterContentEditor, useEditorToolbarState } from './editor'
import { useFavoriteFonts } from '../fonts/state'
import { useMoreFontsState } from '../MoreFonts'
import { appRuntime } from '@/effect/runtime'
import { FontsLoader } from '@/features/fonts-loader/loader'
import './EditorToolbar.scss'

const fontsLoader = appRuntime.runSync(FontsLoader)

export function EditorToolbar({ allowedFontSizes }: { allowedFontSizes: string[] }) {
  const editor = useAtomValue(wysiwygEditorAtom)
  return (
    <div className="BeToolbar" role="toolbar" aria-label={window._('editor.formattingToolbar')}>
      {editor && <ToolbarControls editor={editor} allowedFontSizes={allowedFontSizes} />}
    </div>
  )
}

function ToolbarControls({
  editor,
  allowedFontSizes,
}: {
  editor: ChapterContentEditor
  allowedFontSizes: string[]
}) {
  const state = useEditorToolbarState(editor)
  const favorites = useFavoriteFonts()
  const moreFonts = useMoreFontsState((x) => x.open)
  const fonts = useMemo(
    () => [
      ...new Set([
        'Poppins',
        'Merriweather',
        'Literata',
        ...favorites,
        ...(state.font ? [state.font] : []),
      ]),
    ],
    [favorites, state.font],
  )

  useEffect(() => {
    void appRuntime
      .runPromise(fontsLoader.addFonts(fonts))
      .then(() =>
        Promise.all(fonts.map((font) => document.fonts.load(`20px ${JSON.stringify(font)}`))),
      )
      .catch(() => undefined)
  }, [fonts])

  const fontSizes = allowedFontSizes

  function button(label: string, icon: string, active: boolean, action: () => void) {
    return (
      <button
        type="button"
        className="BeToolbar-button"
        title={window._(`editor.${label}`)}
        aria-label={window._(`editor.${label}`)}
        aria-pressed={active}
        onMouseDown={(event) => event.preventDefault()}
        onClick={action}
      >
        <i className={`fa-solid ${icon}`} aria-hidden="true" />
      </button>
    )
  }

  return (
    <>
      <div className="BeToolbar-group" role="group" aria-label={window._('editor.fontTools')}>
        <div className="BeToolbar-controls">
          <Select
            value={state.font ?? '__default__'}
            onValueChange={(font) => {
              if (font === '__more__') {
                moreFonts(editor)
                return
              }
              if (font !== '__default__') {
                void appRuntime
                  .runPromise(fontsLoader.ensureFontLoaded(font))
                  .catch(() => undefined)
                editor.chain().focus().setFontFamily(font).run()
              } else editor.chain().focus().unsetFontFamily().run()
            }}
          >
            <SelectTrigger
              className="BeToolbar-font"
              aria-label={window._('editor.fontFamily')}
              style={state.font ? { fontFamily: state.font } : undefined}
            >
              <SelectValue />
            </SelectTrigger>
            <SelectContent onCloseAutoFocus={(event) => event.preventDefault()}>
              <SelectItem value="__default__">{window._('editor.defaultFont')}</SelectItem>
              {fonts.map((font) => (
                <SelectItem
                  key={font}
                  value={font}
                  className="text-xl py-2"
                  style={{ fontFamily: font }}
                >
                  <span style={{ fontFamily: font }}>{font}</span>
                </SelectItem>
              ))}
              <SelectItem value="__more__">{window._('editor.moreFonts')}</SelectItem>
            </SelectContent>
          </Select>
          <Select
            value={state.fontSize ?? '__default__'}
            onValueChange={(size) => {
              if (size !== '__default__') editor.chain().focus().setFontSize(size).run()
              else editor.chain().focus().unsetFontSize().run()
            }}
          >
            <SelectTrigger className="BeToolbar-size" aria-label={window._('editor.fontSize')}>
              <SelectValue />
            </SelectTrigger>
            <SelectContent onCloseAutoFocus={(event) => event.preventDefault()}>
              <SelectItem value="__default__">{window._('editor.defaultFontSize')}</SelectItem>
              {fontSizes.map((size) => (
                <SelectItem key={size} value={size}>
                  {size}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {button('bold', 'fa-bold', state.bold, () => editor.chain().focus().toggleBold().run())}
          {button('italic', 'fa-italic', state.italic, () =>
            editor.chain().focus().toggleItalic().run(),
          )}
          {button('underline', 'fa-underline', state.underline, () =>
            editor.chain().focus().toggleUnderline().run(),
          )}
          {button('strike', 'fa-strikethrough', state.strikethrough, () =>
            editor.chain().focus().toggleStrike().run(),
          )}
        </div>
        <span className="BeToolbar-caption">{window._('editor.fontTools')}</span>
      </div>
      <div className="BeToolbar-group" role="group" aria-label={window._('editor.paragraphTools')}>
        <div className="BeToolbar-controls">
          <Select
            value={
              state.textType === 'ul' || state.textType === 'ol'
                ? 'p'
                : state.textType === 'text'
                  ? 'p'
                  : state.textType
            }
            onValueChange={(style) => {
              const chain = editor.chain().focus()
              if (editor.isActive('bulletList')) chain.toggleBulletList()
              if (editor.isActive('orderedList')) chain.toggleOrderedList()
              if (style === 'p') chain.setParagraph().run()
              else chain.setHeading({ level: Number(style.slice(1)) as Level }).run()
            }}
          >
            <SelectTrigger
              className="BeToolbar-style"
              aria-label={window._('editor.paragraphStyle')}
            >
              <SelectValue />
            </SelectTrigger>
            <SelectContent onCloseAutoFocus={(event) => event.preventDefault()}>
              {['p', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6'].map((style) => (
                <SelectItem key={style} value={style}>
                  {window._(`editor.${style}`)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {button('ul', 'fa-list-ul', state.textType === 'ul', () =>
            editor.chain().focus().toggleBulletList().run(),
          )}
          {button('ol', 'fa-list-ol', state.textType === 'ol', () =>
            editor.chain().focus().toggleOrderedList().run(),
          )}
          <span className="BeToolbar-divider" />
          {button('textAlignLeft', 'fa-align-left', state.textAlign === 'left', () =>
            editor.chain().focus().setTextAlign('left').run(),
          )}
          {button('textAlignCenter', 'fa-align-center', state.textAlign === 'center', () =>
            editor.chain().focus().setTextAlign('center').run(),
          )}
          {button('textAlignRight', 'fa-align-right', state.textAlign === 'right', () =>
            editor.chain().focus().setTextAlign('right').run(),
          )}
          {button('textAlignJustify', 'fa-align-justify', state.textAlign === 'justify', () =>
            editor.chain().focus().setTextAlign('justify').run(),
          )}
        </div>
        <span className="BeToolbar-caption">{window._('editor.paragraphTools')}</span>
      </div>
    </>
  )
}
