import { useLayoutEffect, useMemo, useState } from 'react'
import { ReactIslandProps } from '@/islands/common/react-island'
import './BookManagerEditor.scss'
import type { DraftDto } from './contracts'
import { useBEState } from './state'
import { EditorIframe } from './EditorIframe'
import { SaveButton } from './SaveButton'
import { CenterHeader } from './CenterHeader'
import { WidgetsMenu, WidgetsService } from './widgets'
import { MoreFonts } from './MoreFonts'

export default function EditorIslandComponent({ data }: ReactIslandProps) {
  const { draft } = useMemo(() => data as { draft: DraftDto }, [data])
  const [leftOpen, setLeftOpen] = useState(true)
  const [rightOpen, setRightOpen] = useState(true)

  useLayoutEffect(() => {
    useBEState.getState().init(draft)
  }, [draft])

  return (
    <>
      <div className="be-layout">
        <div className="be-layout__header">
          <header className="be-header">
            <div className="be-header__left">
              <a
                href={`/books-manager#/books/${draft.book.id}?t=chapters`}
                className="be-header__logo"
                aria-label={window._('editor.backToChapters')}
              >
                <img src="/_/embed-assets/logo.svg" alt="OpenLibrary" />
              </a>
              <SidebarToggle
                open={leftOpen}
                side="left"
                onClick={() => setLeftOpen((value) => !value)}
              />
            </div>
            <div className="be-header__center">
              <CenterHeader />
            </div>
            <div className="be-header__right">
              <SidebarToggle
                open={rightOpen}
                side="right"
                onClick={() => setRightOpen((value) => !value)}
              />
              <SaveButton />
            </div>
          </header>
        </div>
        <div
          className="be-layout__body"
          style={{
            gridTemplateColumns: `${leftOpen ? '260px' : '0'} minmax(0, 1fr) ${rightOpen ? '260px' : '0'}`,
          }}
        >
          <div className="be-layout__sidebar-clip" aria-hidden={!leftOpen} inert={!leftOpen}>
            <div className="be-layout__left">
              <WidgetsMenu service={WidgetsService.instance()} />
            </div>
          </div>
          <div className="be-layout__center">
            <EditorIframe initialContent={draft.content} />
          </div>
          <div
            className="be-layout__sidebar-clip be-layout__sidebar-clip--right"
            aria-hidden={!rightOpen}
            inert={!rightOpen}
          >
            <div className="be-layout__right">
              <ChapterDetails />
            </div>
          </div>
        </div>
      </div>

      <MoreFonts />
    </>
  )
}

function SidebarToggle({
  open,
  side,
  onClick,
}: {
  open: boolean
  side: 'left' | 'right'
  onClick: () => void
}) {
  const label = window._(open ? 'editor.collapseSidebar' : 'editor.expandSidebar')
  const icon = side === 'left' ? 'fa-table-columns' : 'fa-table-columns fa-flip-horizontal'
  return (
    <button type="button" className="btn btn--icon btn--ghost" onClick={onClick} aria-label={label}>
      <i className={`fa-solid ${icon}`} />
    </button>
  )
}

function ChapterDetails() {
  const draft = useBEState((state) => state.draft)
  if (!draft) return null
  return (
    <section className="p-4 space-y-4">
      <h2 className="font-semibold">{window._('editor.chapterDetails')}</h2>
      <div>
        <div className="text-sm text-secondary-foreground">{window._('editor.book')}</div>
        <div>{draft.book.name}</div>
      </div>
      <div>
        <div className="text-sm text-secondary-foreground">{window._('editor.visibility')}</div>
        <div>
          {window._(draft.isChapterPubliclyAvailable ? 'editor.publiclyVisible' : 'editor.hidden')}
        </div>
      </div>
    </section>
  )
}
