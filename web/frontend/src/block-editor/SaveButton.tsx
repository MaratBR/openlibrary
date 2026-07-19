import { useRef, useState } from 'react'
import {
  useBEState,
  useDraftHasChanges,
  useDraftHasNewerRevision,
  useDraftHasPendingChanges,
  useChapterNameIsValid,
} from './state'
import { PublishChapterPopup } from './PublishChapterPopup'
import { ScheduleChapterPopup } from './ScheduleChapterPopup'
import Popper from '@/components/Popper'

export function SaveButton() {
  const draftHasPendingChanges = useDraftHasChanges()
  const hasNewerRevision = useDraftHasNewerRevision()
  const saving = useBEState((s) => s.saving)
  const chapterNameIsValid = useChapterNameIsValid()
  const [openPublishPopup, setOpenPublishPopup] = useState(false)
  const [openSchedulePopup, setOpenSchedulePopup] = useState(false)
  const [openMenu, setOpenMenu] = useState(false)
  const menuButton = useRef<HTMLButtonElement | null>(null)
  const scheduledAt = useBEState((state) => state.draft?.scheduledAt)

  function handleSaveDraft() {
    void useBEState
      .getState()
      .saveDraft()
      .catch((error) => window.toast.error(error))
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
          className="btn btn--default btn--lg flex justify-center items-center"
        >
          <DraftPendingChangesIndicator />
          {window._('editor.publishDraft')}
        </button>
        <button
          ref={menuButton}
          type="button"
          className="btn btn--default btn--lg btn--icon"
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
  const hasPendingChanges = useDraftHasPendingChanges()
  return (
    <div
      style={{
        transform: `scale(${hasPendingChanges ? 1 : 0})`,
      }}
      className="circular-indicator"
    />
  )
}
