import { formatFileSize } from '@/lib/format'
import { ChangeEvent, useEffect, useRef } from 'preact/compat'

export type ImageUploadProps = {
  onFileSelected: (file: File) => void
  file: File | null
  displayedFile?: string | null
}

export function ImageUpload({ onFileSelected, file, displayedFile }: ImageUploadProps) {
  const ref = useRef<HTMLInputElement | null>(null)

  function handleChange(e: ChangeEvent) {
    const input = e.target as HTMLInputElement
    if (input.files && input.files.length >= 1) {
      onFileSelected(input.files[0])
    }
  }

  useEffect(() => {
    if (!file) {
      if (ref.current) ref.current.value = ''
    }
  }, [file])

  return (
    <div class="flex">
      <label class="bg-secondary rounded-r-2xl h-[200px] w-[133px] relative block cursor-pointer group overflow-hidden">
        <input
          ref={ref}
          onChange={handleChange}
          class="hidden"
          accept="image/png, image/jpeg, image/webp"
          type="file"
        />

        {displayedFile ? (
          <img src={displayedFile} />
        ) : (
          <div class="absolute inset-0 flex justify-center items-center group-hover:bg-highlight">
            <i class="fa-solid fa-image" />
          </div>
        )}
      </label>

      <div class="ml-4">
        {!file && displayedFile && (
          <span class="ml-2">{window._('bookManager.edit.clickOnImageToChange')}</span>
        )}

        {file && <FileInfo file={file} />}
      </div>
    </div>
  )
}

function FileInfo({ file }: { file: File }) {
  return (
    <div>
      {file.name}
      <br />
      {formatFileSize(file.size)}
    </div>
  )
}
