import { ErrorDisplay } from '@/components/error'
import { ImageResizer } from '@/components/image-upload'
import Alpine from 'alpinejs'
import { render } from 'preact'

// class ImageUploaderFileEvent extends CustomEvent<{ file: File; fileCropped: boolean }> {
//   fileCropped: boolean
//   file: File
//   promise?: Promise<unknown>

//   constructor(file: File, fileCropped: boolean) {
//     super('image-uploader:file', { detail: { file, fileCropped } })

//     this.file = file
//     this.fileCropped = fileCropped
//   }

//   setPromise(promise: Promise<unknown>) {
//     this.promise = promise
//   }
// }

Alpine.data(
  'ImageUploader',
  (params?: {
    crop?: { width?: number; height?: number }
    callback?: (file: File) => void | Promise<unknown>
  }) => ({
    _id: Math.random().toString().substring(2),
    _disposeErr: () => {},
    uploading: false,
    hasError: false,

    init() {
      const { crop } = params ?? {}

      const input = this.$refs.input
      if (!(input instanceof HTMLInputElement)) {
        return
      }

      input.addEventListener('change', () => {
        const file = input.files && input.files.length > 0 ? input.files[0] : null
        if (file) {
          if (crop && crop.width && crop.height) {
            resizeImage(file, crop.width, crop.height)
              .then((resizedFile) => {
                this.handleFile(resizedFile)
              })
              .catch((err) => {
                this.handleError(err)
              })
          } else {
            this.handleFile(file)
          }
        }
      })
    },

    handleError(err: unknown) {
      this.hasError = true
      const errorContainer = this.$refs.err

      if (errorContainer instanceof HTMLElement) {
        this._disposeErr()
        render(<ErrorDisplay error={err} />, errorContainer)
        this._disposeErr = () => {
          render(null, errorContainer)
          this._disposeErr = () => {}
        }
      }
    },

    handleFile(file: File) {
      const callback = params?.callback

      if (callback) {
        const result = callback(file)
        if (result instanceof Promise) {
          this.uploading = true
          result
            .catch((err) => this.handleError(err))
            .finally(() => {
              this.uploading = false
            })
        }
      }
    },

    destroy() {
      this._disposeErr()
    },
  }),
)

const RESIZE_IMAGE_ERR_CLOSED = new Error('closed')

function resizeImage(file: File, width: number, height: number) {
  return new Promise<File>((resolve, reject) => {
    const container = document.createElement('div')
    container.style.display = 'contents'

    const unmount = () =>
      requestAnimationFrame(() => {
        render(null, container)
        container.remove()
      })

    const onClose = () => {
      unmount()
      reject(RESIZE_IMAGE_ERR_CLOSED)
    }

    const handleUpload = async (file: File, fileCropped: boolean) => {
      unmount()
      // TODO handle error proper
      if (!fileCropped) throw new Error('failed to crop the file')
      resolve(file)
    }

    requestAnimationFrame(() => {
      document.body.appendChild(container)
      render(
        <ImageResizer
          handleUpload={handleUpload}
          onClose={onClose}
          file={file}
          expectedHeight={height}
          expectedWidth={width}
          height={600}
          width={600}
        />,
        container,
      )
    })
  })
}
