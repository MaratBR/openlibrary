import { useEffect, useRef, useState } from 'preact/hooks'

export type BookCardPreviewContentProps = {
  bookId: string
}

export default function BookCardPreviewContent({ bookId }: BookCardPreviewContentProps) {
  const [content, setContent] = useState('')
  const cache = useRef<Record<string, string>>({})

  useEffect(() => {
    if (cache.current[bookId]) {
      setContent(cache.current[bookId])
      return
    }
    fetch(`/book/${bookId}/__fragment/preview-card`)
      .then((res) => res.text())
      .then((text) => {
        setContent(text)
        cache.current[bookId] = text
      })
  }, [bookId])

  if (content) {
    // eslint-disable-next-line react/no-danger
    return <div class="contents" dangerouslySetInnerHTML={{ __html: content }} />
  }

  return (
    <div class="w-full h-full flex items-center justify-center">
      <div class="loader" />
    </div>
  )
}
