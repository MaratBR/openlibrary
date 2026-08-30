import { useState } from 'react'
import { useAtomValue, useSetAtom } from 'jotai'
import { draftAtom, publishDraftAtom } from './state'
import { useMutation } from '@tanstack/react-query'

import { ModalAnimation, useAnimation } from '@/lib/animate'
import Switch from '@/components/Switch'
import { createRoot } from 'react-dom/client'
import { appRuntime } from '@/effect/runtime'

export function PublishChapterPopup({ onClose, open }: { onClose: () => void; open: boolean }) {
  const isHidden = useAtomValue(draftAtom)?.isChapterPubliclyAvailable === false
  const publishDraft = useSetAtom(publishDraftAtom)
  const [makePublic, setMakePublic] = useState(true)

  const publishMutation = useMutation({
    mutationFn: async () => {
      const draft = await appRuntime.runPromise(publishDraft(makePublic))
      onClose()

      window.toast({
        title: window._('editor.chapterPublished'),
        duration: 15000,
        customContent(element) {
          const root = createRoot(element)
          root.render(
            <a className="Link" href={`/book/${draft.book.id}/chapters/${draft.chapter.id}`}>
              {window._('editor.viewChapter')}
              &nbsp;
              <i className="fa-solid fa-arrow-up-right-from-square" />
            </a>,
          )
          return () => root.unmount()
        },
      })
    },
    onError(error) {
      window.toast.error(error)
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
          className="Btn Btn--outline w-32"
          onClick={() => publishMutation.mutate()}
        >
          {publishMutation.isPending ? (
            <span className="Loader" />
          ) : (
            window._('editor.publishDraft')
          )}
        </button>
        <button
          disabled={publishMutation.isPending}
          className="Btn Btn--ghost"
          onClick={() => onClose()}
        >
          {window._('common.cancel')}
        </button>
      </div>
    </div>
  )
}
