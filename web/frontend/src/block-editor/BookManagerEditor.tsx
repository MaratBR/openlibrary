import { useLayoutEffect, useMemo, useState } from 'preact/hooks'
import { PreactIslandProps } from '@/islands/common/preact-island'
import { z } from 'zod'
import './BookManagerEditor.scss'
import { DraftDtoSchema } from './contracts'
import { EditorIframe } from './wysiwyg'
import { useWYSIWYG } from './wysiwyg/state'
import { useBEState, useDraftHasChanges, useDraftHasNewerRevision } from './state'
import { createPortal } from 'preact/compat'
import Switch from '@/components/Switch'
import { AnimationWrapper, ModalAnimation } from '@/lib/animate'
import { useMutation } from '@tanstack/react-query'
import { render } from 'preact'

const dataSchema = z.object({
  bookId: z.string(),
  draft: DraftDtoSchema,
})

export default function BookManagerEditor({ data }: PreactIslandProps) {
  const { draft } = useMemo(() => dataSchema.parse(data), [data])

  useLayoutEffect(() => {
    useBEState.getState().init(draft)
  }, [draft])

  return (
    <div class="be-layout">
      <div class="be-layout__header">
        <header class="be-header">
          <div />
          <div class="be-header__left">Left</div>
          <div class="be-header__center">Center</div>
          <div class="be-header__right">
            <SaveButton />
          </div>
          <div />
        </header>
      </div>
      <div class="be-layout__body">
        <div class="be-layout__left">Left</div>
        <div class="be-layout__center">
          <ChapterNameInput />
          <EditorIframe />
        </div>
        <div class="be-layout__right">Right</div>
      </div>
    </div>
  )
}

function SaveButton() {
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
        class="btn btn--ghost btn--lg btn--sq flex justify-center items-center"
      >
        {saving ? <span class="loader loader--dark" /> : window._('editor.publishDraft')}
      </button>
      <button
        onClick={handleSaveDraft}
        disabled={!draftHasPendingChanges || saving}
        class="btn btn--lg btn--sq w-30 flex justify-center items-center"
      >
        {saving ? <span class="loader loader--dark" /> : window._('common.save')}
      </button>

      <PublishChapterPopup open={openPublishPopup} onClose={() => setOpenPublishPopup(false)} />
    </div>
  )
}

function PublishChapterPopup({ onClose, open }: { onClose: () => void; open: boolean }) {
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

function ChapterNameInput() {
  const chapterName = useBEState((s) => s.chapterName)
  const elements = useWYSIWYG((s) => s.initData?.elements)

  if (!elements) return null

  return createPortal(
    <div class="my-4">
      <span class="text-muted-foreground">Chapter name</span>
      <input
        name="chapterName"
        class="be-chapter-name-input"
        value={chapterName}
        onChange={(e) => {
          useBEState.getState().setChapterName((e.target as HTMLInputElement).value)
        }}
      />
    </div>,
    elements.contentWrapperHeader,
  )
}
