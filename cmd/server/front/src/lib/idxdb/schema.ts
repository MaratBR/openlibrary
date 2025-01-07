import { z } from 'zod'

export function getIDBSchemaFromZodObject(zodObject: z.AnyZodObject) {
  const allFields = Object.keys(zodObject.shape)
  if (!allFields.includes('_id')) {
    throw new Error('_id field is required')
  }

  allFields.sort((a, b) => a.localeCompare(b))
  return allFields.map((x) => (x === '_id' ? '++_id' : x)).join(', ')
}
