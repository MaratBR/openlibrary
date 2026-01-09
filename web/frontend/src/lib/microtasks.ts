export interface IMicrotask {
  next(): boolean
}

export type MicrotaskCallbacks = {
  onSuccess?: () => void
  onTimeout?: () => void
  onError?: (err: unknown) => void
}

export function queueMicrotasksWithBursts(
  task: IMicrotask,
  burstMaxDuration: number,
  maxDuration: number,
  callbacks?: MicrotaskCallbacks,
) {
  _queueMicrotasksWithBursts(task, burstMaxDuration, maxDuration, 0, callbacks)
}

function _queueMicrotasksWithBursts(
  task: IMicrotask,
  burstMaxDuration: number,
  maxDuration: number,
  accumulatedDuration: number,
  { onError, onSuccess, onTimeout }: MicrotaskCallbacks = {},
) {
  queueMicrotask(() => {
    let finished = false
    let reachedDeadline = false
    const startedAt = performance.now()
    while (!finished && !reachedDeadline) {
      try {
        finished = task.next()
      } catch (e: unknown) {
        if (onError) {
          onError(e)
          return
        }
        throw e
      }
      reachedDeadline = performance.now() - startedAt > burstMaxDuration
    }

    if (!finished) {
      if (accumulatedDuration > maxDuration) {
        console.warn(
          `microtask took ${accumulatedDuration}ms, more than alloted ${maxDuration}ms and was interrupted`,
        )
        if (onTimeout) onTimeout()
        return
      }
      // reschedule for next frame
      requestAnimationFrame(() => {
        _queueMicrotasksWithBursts(
          task,
          burstMaxDuration,
          maxDuration,
          accumulatedDuration + performance.now() - startedAt,
        )
      })
    } else if (onSuccess) onSuccess()
  })
}

export function getPreferredMicroTaskDuration(): number {
  return 1000 / 60 / 3 // 1/3rd of a frame
}
