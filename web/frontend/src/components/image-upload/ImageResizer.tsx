import Modal from '@/components/Modal'
import { useEffect, useRef, useState } from 'preact/hooks'
import { cropImage, loadImageAsElement, resizeImage } from './image-util'
import { initDraggable, ResizableImageBounds } from './draggable'

export function ImageResizer({
  file,
  expectedHeight,
  expectedWidth,
  width,
  height,
  onClose,
  handleUpload,
}: {
  file: File
  expectedHeight: number
  expectedWidth: number
  height: number
  width: number
  onClose: () => void
  handleUpload: (file: File, fileCropped: boolean) => Promise<void>
}) {
  const imgRef = useRef<HTMLImageElement | null>(null)
  const settings = useRef({
    width: expectedWidth,
    height: expectedHeight,
    scaleRatio: 0,
    originalWidth: 0,
    originalHeight: 0,
    x: 0,
    y: 0,
  })

  const [size, setSize] = useState({ width: 0, height: 0 })
  const [bounds, setBounds] = useState<ResizableImageBounds>({ x0: 0, y0: 0, x1: 0, y1: 0 })
  const [src, setSrc] = useState('')
  const [loading, setLoading] = useState(false)

  const xOffset = (width - expectedWidth) / 2
  const yOffset = (height - expectedHeight) / 2

  useEffect(() => {
    if (!file) return
    ;(async () => {
      const { dataUrl, width: actualWidth, height: actualHeight } = await analyzeImage(file)

      settings.current.originalHeight = actualHeight
      settings.current.originalWidth = actualWidth

      let w = actualWidth,
        h = actualHeight,
        mult = 1

      const isWider = w / h > expectedWidth / expectedHeight

      if (isWider) {
        mult = expectedHeight / h
      } else {
        mult = expectedWidth / w
      }

      settings.current.scaleRatio = mult

      w *= mult
      h *= mult

      if (isWider) {
        setBounds({
          y0: 0 + yOffset,
          y1: 0 + yOffset,
          x0: expectedWidth - w + xOffset,
          x1: 0 + xOffset,
        })
      } else {
        setBounds({
          x0: 0 + xOffset,
          x1: 0 + xOffset,
          y0: expectedHeight - h + yOffset,
          y1: 0 + yOffset,
        })
      }

      setSize({ width: w, height: h })
      setSrc(dataUrl)
    })()
  }, [file, expectedHeight, expectedWidth, height, width, yOffset, xOffset])

  useEffect(() => {
    if (imgRef.current) {
      const dispose = initDraggable({
        element: imgRef.current,
        bounds,
        onDrag(event) {
          settings.current.x = event.x - xOffset
          settings.current.y = event.y - yOffset
        },
      })
      return dispose
    }
  }, [bounds, expectedHeight, expectedWidth, height, width, xOffset, yOffset])

  async function handleResizeAndCrop() {
    if (loading) return
    setLoading(true)
    let finalFile: File
    let fileCropped = false

    try {
      finalFile = await transformImage(file, {
        scaleRatio: settings.current.scaleRatio,
        width: settings.current.width / settings.current.scaleRatio,
        height: settings.current.height / settings.current.scaleRatio,
        x: -settings.current.x / settings.current.scaleRatio,
        y: -settings.current.y / settings.current.scaleRatio,
        disableResize: true,
      })
      fileCropped = true
    } catch (e: unknown) {
      console.error('[ImageResizer] transformImage failed', e)
      finalFile = file
    }

    try {
      await handleUpload(finalFile, fileCropped)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Modal open>
      <div
        class="relative overflow-hidden flex items-center justify-center"
        style={{ height, width }}
      >
        <button
          class="z-10 size-10 flex items-center justify-center absolute top-1 right-1 hover:bg-white/20"
          onClick={() => onClose()}
        >
          <i class="fa-solid fa-xmark text-white" />
        </button>

        <Canvas cutoutHeight={expectedHeight} cutoutWidth={expectedWidth} />
        <div
          class="box-content overflow-hidden bg-background"
          style={{ width: expectedWidth, height: expectedHeight }}
        >
          <img class="max-w-none max-h-none select-none" ref={imgRef} src={src} style={size} />
        </div>

        <div class="absolute right-4 bottom-4 flex z-10">
          <button class="btn btn--secondary" onClick={handleResizeAndCrop}>
            {loading ? <span class="loader" /> : window._('bookManager.edit.cropAndUploadCover')}
          </button>
        </div>
      </div>
    </Modal>
  )
}

async function transformImage(
  file: File,
  options: {
    width: number
    height: number
    x: number
    y: number
    scaleRatio: number
    disableResize: boolean
  },
): Promise<File> {
  file = await cropImage(file, {
    x: options.x,
    y: options.y,
    width: options.width,
    height: options.height,
  })
  if (!options.disableResize) {
    file = await resizeImage(file, options.scaleRatio)
  }
  return file
}

function Canvas({ cutoutHeight, cutoutWidth }: { cutoutHeight: number; cutoutWidth: number }) {
  const canvasRef = useRef<HTMLCanvasElement | null>(null)

  useEffect(() => {
    if (canvasRef.current) {
      return initCutoutCanvas(canvasRef.current, cutoutHeight, cutoutWidth)
    }
  }, [cutoutHeight, cutoutWidth])

  return (
    <canvas
      style={{ zIndex: 1 }}
      class="w-full h-full absolute inset-0 pointer-events-none"
      ref={canvasRef}
    />
  )
}

function initCutoutCanvas(
  canvas: HTMLCanvasElement,
  cutoutHeight: number,
  cutoutWidth: number,
): () => void {
  const ctx = canvas.getContext('2d')
  if (!ctx) {
    throw new Error('Could not get 2D rendering context')
  }

  // Function to draw the cutout overlay
  const drawCutout = () => {
    // Clear the entire canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height)

    // Set the overlay color (semi-transparent black)
    ctx.fillStyle = 'rgba(0, 0, 0, 0.5)'

    // Fill the entire canvas
    ctx.fillRect(0, 0, canvas.width, canvas.height)

    // Create a "hole" in the center
    ctx.globalCompositeOperation = 'destination-out'

    // Calculate center position for the cutout
    const centerX = (canvas.width - cutoutWidth) / 2
    const centerY = (canvas.height - cutoutHeight) / 2

    // Draw a transparent rectangle in the center
    ctx.fillStyle = 'rgba(0, 0, 0, 1)'
    ctx.fillRect(centerX, centerY, cutoutWidth, cutoutHeight)

    // Reset composite operation
    ctx.globalCompositeOperation = 'source-over'
  }

  // Initial draw
  drawCutout()

  // Create a resize observer to handle canvas size changes
  const resizeObserver = new ResizeObserver(() => {
    // Sync canvas size with its display size
    canvas.width = canvas.clientWidth
    canvas.height = canvas.clientHeight

    // Redraw the cutout
    drawCutout()
  })

  // Start observing the canvas
  resizeObserver.observe(canvas)

  // Return a cleanup function
  return () => {
    // Stop observing the canvas
    resizeObserver.unobserve(canvas)
  }
}
type ImageAnalysisResult = {
  width: number
  height: number
  dataUrl: string
}

async function analyzeImage(file: File): Promise<ImageAnalysisResult> {
  const img = await loadImageAsElement(file)
  return {
    width: img.width,
    height: img.height,
    dataUrl: img.src,
  }
}
