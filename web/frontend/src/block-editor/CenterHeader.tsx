import { useAtomValue } from 'jotai'
import { draftAtom } from './state'

export function CenterHeader() {
  const scheduledAt = useAtomValue(draftAtom)?.scheduledAt

  return (
    <div className="flex justify-center items-center h-full">
      {scheduledAt && (
        <span className="text-sm text-secondary-foreground">
          {window._('editor.publishingScheduledFor', {
            date: new Date(scheduledAt).toLocaleString(),
          })}
        </span>
      )}
    </div>
  )
}
