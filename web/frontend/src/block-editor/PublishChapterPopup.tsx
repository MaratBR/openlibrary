import { useState } from 'preact/hooks'
import { useBEState } from './state'
import { useMutation } from '@tanstack/react-query'
import { render } from 'preact'
import { AnimationWrapper, ModalAnimation } from '@/lib/animate'
import Switch from '@/components/Switch'

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
            render(
              <a class="link" href={`/book/${draft.book.id}/chapters/${draft.chapterId}`}>
                {window._('editor.viewChapter')}
                &nbsp;
                <i class="fa-solid fa-arrow-up-right-from-square" />
              </a>,
              element,
            )
            return () => render(null, element)
          }
        },
      })
    },
  })

  return (
    <AnimationWrapper show={open} animation={ModalAnimation.factory(150)}>
      <div class="be-publish-popup">
        <header class="text-xl font-semibold">{window._('editor.publishAreYouSure')}</header>

        <p>{window._('editor.publishWarning')}</p>

        {isHidden && (
          <div class="mt-4 flex gap-2">
            <Switch
              name="makePublic"
              id="editor-makePublic"
              value={makePublic}
              onChange={setMakePublic}
            />
            <label class="label" for="editor-makePublic">
              {window._('editor.makeChapterVisible')}
            </label>
          </div>
        )}

        <div class="mt-4 flex gap-1">
          <button
            disabled={publishMutation.isPending}
            class="btn btn--outline w-32"
            onClick={() => publishMutation.mutate()}
          >
            {publishMutation.isPending ? <span class="loader" /> : window._('editor.publishDraft')}
          </button>
          <button
            disabled={publishMutation.isPending}
            class="btn btn--ghost"
            onClick={() => onClose()}
          >
            {window._('common.cancel')}
          </button>
        </div>
      </div>
    </AnimationWrapper>
  )
}
