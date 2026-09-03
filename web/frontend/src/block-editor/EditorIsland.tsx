import { useMemo, useState } from 'react'
import { useAtomValue } from 'jotai'
import { ReactIslandProps } from '@/islands/common/react-island'
import { jotaiStore } from '@/react'
import './BookManagerEditor.scss'
import type { DraftDto } from './contracts'
import { draftAtom, initializeDraftAtom } from './state'
import { EditorIframe } from './EditorIframe'
import { SaveButton } from './SaveButton'
import { CenterHeader } from './CenterHeader'
import { WidgetsMenu, WidgetsService } from './widgets'
import { MoreFonts } from './MoreFonts'

export default function EditorIslandComponent({ data }: ReactIslandProps) {
  const { draft } = useMemo(() => data as { draft: DraftDto }, [data])
  useMemo(() => {
    jotaiStore.set(initializeDraftAtom, draft)
  }, [draft])

  return <EditorIsland draft={draft} />
}

function EditorIsland({ draft }: { draft: DraftDto }) {
  const [leftOpen, setLeftOpen] = useState(true)
  const [rightOpen, setRightOpen] = useState(true)

  return (
    <>
      <div className="BeLayout">
        <div className="BeLayout-header">
          <header className="BeHeader">
            <div className="BeHeader-left">
              <a
                href={`/books-manager#/books/${draft.book.id}?t=chapters`}
                className="BeHeader-logo"
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
            <div className="BeHeader-center">
              <CenterHeader />
            </div>
            <div className="BeHeader-right">
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
          className="BeLayout-body"
          style={{
            gridTemplateColumns: `${leftOpen ? '260px' : '0'} minmax(0, 1fr) ${rightOpen ? '260px' : '0'}`,
          }}
        >
          <div className="BeLayout-sidebarClip" aria-hidden={!leftOpen} inert={!leftOpen}>
            <div className="BeLayout-left">
              <WidgetsMenu service={WidgetsService.instance()} />
            </div>
          </div>
          <div className="BeLayout-center">
            <EditorIframe initialContent={draft.content} />
          </div>
          <div
            className="BeLayout-sidebarClip BeLayoutSidebarClip--right"
            aria-hidden={!rightOpen}
            inert={!rightOpen}
          >
            <div className="BeLayout-right">
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
    <button type="button" className="Btn Btn--icon Btn--ghost" onClick={onClick} aria-label={label}>
      <i className={`fa-solid ${icon}`} />
    </button>
  )
}

function ChapterDetails() {
  const draft = useAtomValue(draftAtom)
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
