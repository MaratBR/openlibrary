import { useMemo, useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { useAtomValue, useSetAtom } from 'jotai'
import { draftAtom, scheduleDraftAtom } from './state'
import { appRuntime } from '@/effect/runtime'

export function ScheduleChapterPopup({ onClose }: { onClose: () => void }) {
  const scheduledAt = useAtomValue(draftAtom)?.scheduledAt
  const scheduleDraft = useSetAtom(scheduleDraftAtom)
  const [value, setValue] = useState(() => toLocalInputValue(scheduledAt ?? oneHourFromNow()))
  const selectedDate = useMemo(() => new Date(value), [value])
  const valid = value !== '' && !Number.isNaN(selectedDate.getTime()) && selectedDate > new Date()

  const mutation = useMutation({
    mutationFn: () => appRuntime.runPromise(scheduleDraft(selectedDate)),
    onSuccess() {
      window.toast({
        title: window._(scheduledAt ? 'editor.chapterRescheduled' : 'editor.chapterScheduled'),
      })
      onClose()
    },
    onError(error) {
      window.toast.error(error)
    },
  })

  return (
    <div className="be-publish-popup" role="dialog" aria-modal="true">
      <h2 className="text-xl font-semibold">
        {window._(scheduledAt ? 'editor.rescheduleChapter' : 'editor.scheduleChapter')}
      </h2>
      <label className="label mt-4" htmlFor="chapter-scheduled-at">
        {window._('editor.publishAt')}
      </label>
      <input
        id="chapter-scheduled-at"
        type="datetime-local"
        className="input mt-1"
        value={value}
        min={toLocalInputValue(new Date())}
        onChange={(event) => setValue(event.currentTarget.value)}
      />
      <div className="flex gap-2 mt-4">
        <button
          type="button"
          className="Btn Btn--primary"
          disabled={!valid || mutation.isPending}
          onClick={() => mutation.mutate()}
        >
          {window._(scheduledAt ? 'editor.reschedule' : 'editor.schedule')}
        </button>
        <button type="button" className="Btn Btn--ghost" onClick={onClose}>
          {window._('common.cancel')}
        </button>
      </div>
    </div>
  )
}

function oneHourFromNow() {
  return new Date(Date.now() + 60 * 60 * 1000)
}

function toLocalInputValue(value: string | Date) {
  const date = new Date(value)
  return new Date(date.getTime() - date.getTimezoneOffset() * 60_000).toISOString().slice(0, 16)
}
