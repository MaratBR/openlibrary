export function loadImageAsDataUrl(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (event) => {
      resolve(event.target?.result as string)
    }
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

// Create a promise-based image loader
export function loadImageAsElement(file: File): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    loadImageAsDataUrl(file)
      .then((url) => {
        const img = new Image()
        img.onload = () => resolve(img)
        img.onerror = reject
        img.src = url
      })
      .catch(reject)
  })
}

/**
 * Resizes an image by a specific scale ratio
 * @param imageFile Input image file
 * @param ratio Scale ratio (e.g., 0.5 for half size, 2 for double size)
 * @param quality JPEG compression quality (default 0.85)
 * @returns Promise resolving to resized image file
 */
export async function resizeImage(
  imageFile: File,
  ratio: number,
  quality: number = 0.85,
): Promise<File> {
  // Load the image
  const img = await loadImageAsElement(imageFile)

  // Calculate new dimensions
  const newWidth = Math.round(img.width * ratio)
  const newHeight = Math.round(img.height * ratio)

  // Create canvas with new dimensions
  const canvas = document.createElement('canvas')
  canvas.width = newWidth
  canvas.height = newHeight

  // Get 2D rendering context
  const ctx = canvas.getContext('2d')
  if (!ctx) {
    throw new Error('Could not get 2D rendering context')
  }

  // Set background to white for transparency handling
  ctx.fillStyle = 'white'
  ctx.fillRect(0, 0, newWidth, newHeight)

  // Draw the resized image
  ctx.drawImage(
    img,
    0,
    0, // Destination canvas start coordinates
    newWidth,
    newHeight, // Destination canvas dimensions
  )

  // Convert canvas to blob
  return new Promise((resolve, reject) => {
    canvas.toBlob(
      (blob) => {
        if (!blob) {
          reject(new Error('Failed to create blob'))
          return
        }

        const resizedFile = new File([blob], imageFile.name, {
          type: imageFile.type || 'image/jpeg',
          lastModified: Date.now(),
        })

        resolve(resizedFile)
      },
      imageFile.type || 'image/jpeg',
      quality,
    )
  })
}

/**
 * Crops an image to specified coordinates
 * @param imageFile Input image file
 * @param cropOptions Cropping coordinates and dimensions
 * @returns Promise resolving to cropped image file
 */
export async function cropImage(
  imageFile: File,
  cropOptions: {
    x: number
    y: number
    width: number
    height: number
    quality?: number
  },
): Promise<File> {
  // Destructure crop options with defaults
  const { x, y, width, height, quality = 1 } = cropOptions

  // Load the image
  const img = await loadImageAsElement(imageFile)

  // Create canvas with cropped dimensions
  const canvas = document.createElement('canvas')
  canvas.width = width
  canvas.height = height

  // Get 2D rendering context
  const ctx = canvas.getContext('2d')
  if (!ctx) {
    throw new Error('Could not get 2D rendering context')
  }

  // Draw the cropped portion of the image
  ctx.drawImage(
    img,
    x,
    y, // Source image start coordinates
    width,
    height, // Source image crop dimensions
    0,
    0, // Destination canvas start coordinates
    width,
    height, // Destination canvas dimensions
  )

  // Convert canvas to blob
  return new Promise((resolve, reject) => {
    canvas.toBlob(
      (blob) => {
        if (!blob) {
          reject(new Error('Failed to create blob'))
          return
        }

        const croppedFile = new File([blob], imageFile.name, {
          type: imageFile.type || 'image/jpeg',
          lastModified: Date.now(),
        })

        resolve(croppedFile)
      },
      imageFile.type || 'image/jpeg',
      quality,
    )
  })
}
