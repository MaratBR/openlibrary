import { DateTime } from 'luxon'
import React from 'react'

export type TimestampProps = {
  value: Date | string | number
}

export default function Timestamp({ value }: TimestampProps) {
  const dt = React.useMemo(() => {
    let dt: DateTime
    if (value instanceof Date) {
      dt = DateTime.fromJSDate(value)
    } else if (typeof value === 'string') {
      dt = DateTime.fromISO(value)
    } else {
      dt = DateTime.fromMillis(value)
    }

    return dt
  }, [value])

  return (
    <time dateTime={dt.toISO()} className="timestamp">
      {dt.toRelative()}
    </time>
  )
}
