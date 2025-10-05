import { ImageResizer, ImageUpload } from '@/components/image-upload'
import { httpUploadCover, ManagerBookDetailsDto } from './api'
import { useState } from 'preact/hooks'

export default function BookCover({ book }: { book: ManagerBookDetailsDto }) {
  const [file, setFile] = useState<File | null>(null)
  const [displayedFile, setDisplayedFile] = useState(book.cover)

  function handleFile(file: File) {
    setFile(file)
  }

  async function handleFileUpload(file: File, fileCropped: boolean) {
    const newUrl = await httpUploadCover({
      clientCropped: fileCropped,
      file,
      bookId: book.id,
    })
    const url = new URL(newUrl)
    url.searchParams.set('__cacheBuster', Date.now().toString())
    setDisplayedFile(url.toString())
    setFile(null)
  }

  return (
    <div class="p-4">
      <ImageUpload displayedFile={displayedFile} file={file} onFileSelected={handleFile} />
      {file && (
        <ImageResizer
          width={600}
          height={600}
          file={file}
          expectedHeight={300}
          expectedWidth={200}
          onClose={() => setFile(null)}
          handleUpload={handleFileUpload}
        />
      )}
    </div>
  )
}
