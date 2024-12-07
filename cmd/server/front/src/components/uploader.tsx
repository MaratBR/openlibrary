import { cn } from '@/lib/utils'
import React from 'react'

export type UploaderProps = {
  subtitle?: string
  className?: string
  accept?: string
  onFile?: (file: File) => void
}

export default function Uploader({ subtitle, className, accept, onFile }: UploaderProps) {
  const inputRef = React.useRef<HTMLInputElement | null>(null)

  function handleDragOver(e: React.DragEvent<HTMLDivElement>) {
    if (e.dataTransfer.files.length === 0) {
      return
    }
    onFile?.(e.dataTransfer.files[0])
  }

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    if (e.target.files && e.target.files.length > 0) {
      onFile?.(e.target.files[0])
    }
  }

  return (
    <div
      role="button"
      className={cn(
        'flex flex-col items-center justify-center rounded-md p-8 text-muted-foreground border-2 border-dashed ' +
          'hover:text-primary',
        className,
      )}
      onClick={() => inputRef.current?.click()}
      onDragOver={handleDragOver}
    >
      <input
        ref={inputRef}
        aria-hidden="true"
        className="hidden"
        type="file"
        accept={accept}
        onChange={handleChange}
      />
      <span className="font-semibold">Drag 'n' drop file here, or click to select file</span>
      {subtitle && <span className="text-sm">{subtitle}</span>}
    </div>
  )
}
