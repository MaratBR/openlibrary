import { Button } from '@/components/ui/button'
import Uploader from '@/components/uploader'
import { fileSize, getDataUrl } from '@/lib/files'
import { useMutation } from '@tanstack/react-query'
import { CircleX, Upload } from 'lucide-react'
import React, { useMemo } from 'react'
import { httpUploadBookCover } from '../api'
import { toast } from 'sonner'

export type BookCoverUploaderProps = {
  maxSize: number
  currentImage: string | null
  bookId: string
}

export default function BookCoverUploader({
  bookId,
  maxSize,
  currentImage,
}: BookCoverUploaderProps) {
  const [file, setFile] = React.useState<File | null>(null)
  const [url, setUrl] = React.useState<string | null>(null)

  React.useEffect(() => {
    if (!file) return
    setUrl(null)
    getDataUrl(file).then(setUrl)
  }, [file])

  const error = useMemo(() => {
    if (!file) return null

    if (file.size > maxSize) {
      return 'File is too big'
    }

    return null
  }, [file, maxSize])

  const upload = useMutation({
    mutationFn: () => {
      if (!file) throw new Error('No file')
      return httpUploadBookCover(bookId, file)
    },
    onSuccess() {
      toast(
        <div className="flex gap-2 items-center">
          <Upload />
          Cover uploaded
        </div>,
      )
    },
    onError(err) {
      console.error(err)
      toast(
        <div className="flex gap-2 items-center">
          <CircleX />
          Failed to upload cover
        </div>,
      )
    },
  })

  return (
    <section className="page-section">
      <div className="grid grid-cols-[300px_1fr] gap-4">
        <div className="flex items-center justify-center overflow-hidden rounded-xl bg-muted h-[300px] min-w-[300px] max-w-[300px]">
          {(url || currentImage) && (
            <img className="w-full h-full object-contain" src={url || currentImage || ''} alt="" />
          )}
        </div>

        <div className="flex flex-col gap-4">
          {error && (
            <div className="p-4 border border-red-600 bg-red-600/15 rounded-md">{error}</div>
          )}
          <Uploader
            onFile={setFile}
            accept="image/jpeg,image/png,image/webp"
            className="w-full h-full"
            subtitle={`You can upload jpeg or png image up to ${fileSize(maxSize)}`}
          />
        </div>
      </div>

      <div className="mt-8">
        <Button disabled={!file || !!error || upload.isPending} onClick={() => upload.mutate()}>
          <Upload /> Upload
        </Button>
      </div>
    </section>
  )
}
