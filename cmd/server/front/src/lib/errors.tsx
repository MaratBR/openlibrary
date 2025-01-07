import { HTTPError } from 'ky'
import React, { useEffect, useState } from 'react'

export function getErrorMessage(err: unknown): React.ReactNode {
  if (typeof err === 'string') {
    return err
  }

  if (err instanceof HTTPError) {
    return getErrorMessageFromResponse(err)
  }

  if (err instanceof Error) {
    return err.message
  }

  return 'unknown error'
}

function getErrorMessageFromResponse(httpError: HTTPError): React.ReactNode {
  return <ErrorMessageDisplay httpError={httpError} />
}

function ErrorMessageDisplay({ httpError }: { httpError: HTTPError }) {
  const [message, setMessage] = useState('')

  useEffect(() => {
    httpError.response.text().then(setMessage)
  }, [httpError])

  return (
    <div>
      <h3 className="text-md font-[500]">
        {httpError.response.status} {httpError.response.statusText}
      </h3>
      <p className="empty:hidden my-1">{message}</p>
      <p className="text-muted-foreground whitespace-pre leading-4 mt-1">
        Method: {httpError.request.method} <br />
        URL: {new URL(httpError.request.url).pathname}
      </p>
    </div>
  )
}
