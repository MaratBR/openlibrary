import { useState } from 'react'
import { useBEState } from './state'
import { useMutation } from '@tanstack/react-query'

import { ModalAnimation, useAnimation } from '@/lib/animate'
import Switch from '@/components/Switch'
import { createRoot } from 'react-dom/client'

export function PublishChapterPopup({ onClose, open }: { onClose: () => void; open: boolean }) {
  const isHidden = useBEState((s) => s.draft?.isChapterPubliclyAvailable === false)
  const [makePublic, setMakePublic] = useState(true)

  const publishMutation = useMutation({
    mutationFn: async () => {
      await useBEState.getState().saveAndPublishDraft(makePublic)
      onClose()

      window.toast({
        title: window._('editor.chapterPublished'),
        duration: 15000,
        customContent(element) {
          const { draft } = useBEState.getState()

          if (!draft) {
            element.innerText = 'ERROR: no draft in state, cannot display toast message'
          } else {
            const root = createRoot(element)
            root.render(
              <a className="link" href={`/book/${draft.book.id}/chapters/${draft.chapter.id}`}>
                {window._('editor.viewChapter')}
                &nbsp;
                <i className="fa-solid fa-arrow-up-right-from-square" />
              </a>,
            )
            return () => root.unmount()
          }
        },
      })
    },
  })

  const { ref } = useAnimation({
    show: open,
    animation: ModalAnimation.default,
  })

  return (
    <div ref={ref} className="be-publish-popup">
      <header className="text-xl font-semibold">{window._('editor.publishAreYouSure')}</header>

      <p>{window._('editor.publishWarning')}</p>

      {isHidden && (
        <div className="mt-4 flex gap-2">
          <Switch
            name="makePublic"
            id="editor-makePublic"
            value={makePublic}
            onChange={setMakePublic}
          />
          <label className="label" htmlFor="editor-makePublic">
            {window._('editor.makeChapterVisible')}
          </label>
        </div>
      )}

      <div className="mt-4 flex gap-1">
        <button
          disabled={publishMutation.isPending}
          className="btn btn--outline w-32"
          onClick={() => publishMutation.mutate()}
        >
          {publishMutation.isPending ? (
            <span className="loader" />
          ) : (
            window._('editor.publishDraft')
          )}
        </button>
        <button
          disabled={publishMutation.isPending}
          className="btn btn--ghost"
          onClick={() => onClose()}
        >
          {window._('common.cancel')}
        </button>
      </div>
    </div>
  )
}
