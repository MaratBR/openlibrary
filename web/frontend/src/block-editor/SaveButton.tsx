import { useState } from 'react'
import {
  useBEState,
  useDraftHasChanges,
  useDraftHasNewerRevision,
  useDraftHasPendingChanges,
} from './state'
import { PublishChapterPopup } from './PublishChapterPopup'

export function SaveButton() {
  const draftHasPendingChanges = useDraftHasChanges()
  const hasNewerRevision = useDraftHasNewerRevision()
  const saving = useBEState((s) => s.saving)
  const [openPublishPopup, setOpenPublishPopup] = useState(false)

  function handleSaveDraft() {
    useBEState.getState().saveDraft()
  }

  return (
    <div className="flex gap-4">
      <button
        onClick={() => setOpenPublishPopup(true)}
        disabled={(!draftHasPendingChanges && !hasNewerRevision) || saving}
        className="btn btn--outline btn--lg flex justify-center items-center"
      >
        <DraftPendingChangesIndicator />
        {saving ? <span className="loader loader--dark" /> : window._('editor.publishDraft')}
      </button>
      <button
        onClick={handleSaveDraft}
        disabled={!draftHasPendingChanges || saving}
        className="btn btn--lg btn--default w-30 flex justify-center items-center"
      >
        {saving ? <span className="loader loader--dark" /> : window._('common.save')}
      </button>

      <PublishChapterPopup open={openPublishPopup} onClose={() => setOpenPublishPopup(false)} />
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
