export interface IMicrotask {
  next(): boolean
}

export function queueMicrotasksWithBursts(
  task: IMicrotask,
  burstMaxDuration: number,
  maxDuration: number,
) {
  _queueMicrotasksWithBursts(task, burstMaxDuration, maxDuration, 0)
}

function _queueMicrotasksWithBursts(
  task: IMicrotask,
  burstMaxDuration: number,
  maxDuration: number,
  accumulatedDuration: number,
) {
  queueMicrotask(() => {
    let finished = false
    let reachedDeadline = false
    const startedAt = performance.now()
    while (!finished && !reachedDeadline) {
      finished = task.next()
      reachedDeadline = performance.now() - startedAt > burstMaxDuration
    }

    if (!finished) {
      if (accumulatedDuration > maxDuration) {
        console.warn(
          `microtask took ${accumulatedDuration}ms, more than alloted ${maxDuration}ms and was interrupted`,
        )
        return
      }
      console.log('rescheduling microtask')
      // reschedule for next frame
      requestAnimationFrame(() => {
        _queueMicrotasksWithBursts(
          task,
          burstMaxDuration,
          maxDuration,
          accumulatedDuration + performance.now() - startedAt,
        )
      })
    }
  })
}

export function getPreferredMicroTaskDuration(): number {
  return 1000 / 60 / 3 // 1/3rd of a frame
}
