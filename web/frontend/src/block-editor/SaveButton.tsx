import { useRef, useState } from 'react'
import { useAtomValue, useSetAtom } from 'jotai'
import {
  chapterNameIsValidAtom,
  draftAtom,
  draftHasNewerRevisionAtom,
  draftHasPendingChangesAtom,
  saveDraftAtom,
  savingAtom,
  useDraftHasChanges,
} from './state'
import { PublishChapterPopup } from './PublishChapterPopup'
import { ScheduleChapterPopup } from './ScheduleChapterPopup'
import Popper from '@/components/Popper'
import { appRuntime } from '@/effect/runtime'

export function SaveButton() {
  const draftHasPendingChanges = useDraftHasChanges()
  const hasNewerRevision = useAtomValue(draftHasNewerRevisionAtom)
  const saving = useAtomValue(savingAtom)
  const chapterNameIsValid = useAtomValue(chapterNameIsValidAtom)
  const saveDraft = useSetAtom(saveDraftAtom)
  const [openPublishPopup, setOpenPublishPopup] = useState(false)
  const [openSchedulePopup, setOpenSchedulePopup] = useState(false)
  const [openMenu, setOpenMenu] = useState(false)
  const menuButton = useRef<HTMLButtonElement | null>(null)
  const scheduledAt = useAtomValue(draftAtom)?.scheduledAt

  function handleSaveDraft() {
    void appRuntime.runPromise(saveDraft()).catch((error) => window.toast.error(error))
  }

  return (
    <div className="flex gap-3">
      <button
        type="button"
        onClick={handleSaveDraft}
        disabled={!chapterNameIsValid || !draftHasPendingChanges || saving}
        className="btn btn--lg btn--outline min-w-28 flex justify-center items-center"
      >
        {saving ? <span className="loader" /> : window._('common.save')}
      </button>

      <div className="btn-group">
        <button
          type="button"
          onClick={() => setOpenPublishPopup(true)}
          disabled={!chapterNameIsValid || (!draftHasPendingChanges && !hasNewerRevision) || saving}
          className="btn btn--primary btn--lg flex justify-center items-center"
        >
          <DraftPendingChangesIndicator />
          {window._('editor.publishDraft')}
        </button>
        <button
          ref={menuButton}
          type="button"
          className="btn btn--primary btn--lg btn--icon"
          disabled={!chapterNameIsValid || saving}
          aria-label={window._('editor.morePublishingOptions')}
          aria-expanded={openMenu}
          onClick={() => setOpenMenu((value) => !value)}
        >
          <i className="fa-solid fa-chevron-down" />
        </button>
      </div>

      <Popper
        anchorEl={menuButton}
        open={openMenu}
        onClose={() => setOpenMenu(false)}
        placement="bottom-end"
      >
        <div className="card card--elevated p-1 min-w-48">
          <button
            type="button"
            className="btn btn--ghost w-full justify-start"
            onClick={() => {
              setOpenMenu(false)
              setOpenSchedulePopup(true)
            }}
          >
            <i className="fa-regular fa-clock mr-2" />
            {window._(scheduledAt ? 'editor.reschedule' : 'editor.schedule')}
          </button>
        </div>
      </Popper>

      <PublishChapterPopup open={openPublishPopup} onClose={() => setOpenPublishPopup(false)} />
      {openSchedulePopup && <ScheduleChapterPopup onClose={() => setOpenSchedulePopup(false)} />}
    </div>
  )
}

function DraftPendingChangesIndicator() {
  const hasPendingChanges = useAtomValue(draftHasPendingChangesAtom)
  return (
    <div
      style={{
        transform: `scale(${hasPendingChanges ? 1 : 0})`,
      }}
      className="circular-indicator"
    />
  )
}
