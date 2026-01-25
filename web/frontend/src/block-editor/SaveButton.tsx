import { useState } from 'preact/hooks'
import { useBEState, useDraftHasChanges, useDraftHasNewerRevision } from './state'
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
    <div class="flex gap-4">
      <button
        onClick={() => setOpenPublishPopup(true)}
        disabled={(!draftHasPendingChanges && !hasNewerRevision) || saving}
        class="btn btn--ghost btn--lg flex justify-center items-center"
      >
        {saving ? <span class="loader loader--dark" /> : window._('editor.publishDraft')}
      </button>
      <button
        onClick={handleSaveDraft}
        disabled={!draftHasPendingChanges || saving}
        class="btn primary btn--lg w-30 flex justify-center items-center"
      >
        {saving ? <span class="loader loader--dark" /> : window._('common.save')}
      </button>

      <PublishChapterPopup open={openPublishPopup} onClose={() => setOpenPublishPopup(false)} />
    </div>
  )
}
