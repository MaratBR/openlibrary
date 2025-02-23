import { httpClient } from '@/http-client'
import { z } from 'zod'

const tagCategorySchema = z.enum(['other', 'warning', 'fandom', 'rel', 'reltype', 'unknown'])

export type TagsCategory = z.infer<typeof tagCategorySchema>

export const definedTagDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  desc: z.string(),
  adult: z.boolean(),
  spoiler: z.boolean(),
  cat: tagCategorySchema,
})

export type DefinedTagDto = z.infer<typeof definedTagDtoSchema>

const userDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  avatar: z.string(),
})

const numberRangeSchema = z.object({
  min: z.number().int().nullable(),
  max: z.number().int().nullable(),
})

export const detailedBookSearchQuerySchema = z.object({
  words: numberRangeSchema,
  chapters: numberRangeSchema,
  wordsPerChapter: numberRangeSchema,
  includeTags: z.array(definedTagDtoSchema),
  excludeTags: z.array(definedTagDtoSchema),
  includeUsers: z.array(userDtoSchema),
  excludeUsers: z.array(userDtoSchema),
  includeBanned: z.boolean(),
  includeHidden: z.boolean(),
  includeEmpty: z.boolean(),
  page: z.number().int(),
  pageSize: z.number().int(),
})

export type DetailedBookSearchQuery = z.infer<typeof detailedBookSearchQuerySchema>

export function getDefaultDetailedBookSearchQuery(): DetailedBookSearchQuery {
  return {
    words: { min: null, max: null },
    chapters: { min: null, max: null },
    wordsPerChapter: { min: null, max: null },
    includeTags: [],
    excludeTags: [],
    includeUsers: [],
    excludeUsers: [],
    includeBanned: false,
    includeHidden: false,
    includeEmpty: false,
    page: 1,
    pageSize: 10,
  }
}

export function getQueryParams(query: DetailedBookSearchQuery): URLSearchParams {
  const params = new URLSearchParams()

  if (query.words.min !== null) params.set('w.min', query.words.min.toString())
  if (query.words.max !== null) params.set('w.max', query.words.max.toString())
  if (query.chapters.min !== null) params.set('c.min', query.chapters.min.toString())
  if (query.chapters.max !== null) params.set('c.max', query.chapters.max.toString())
  if (query.wordsPerChapter.min !== null) params.set('wc.min', query.wordsPerChapter.min.toString())
  if (query.wordsPerChapter.max !== null) params.set('wc.max', query.wordsPerChapter.max.toString())
  if (query.includeTags.length > 0) params.set('it', query.includeTags.map((x) => x.id).join(','))
  if (query.excludeTags.length > 0) params.set('et', query.excludeTags.map((x) => x.id).join(','))
  if (query.includeUsers.length > 0) params.set('iu', query.includeUsers.map((x) => x.id).join(','))
  if (query.excludeUsers.length > 0) params.set('eu', query.excludeUsers.map((x) => x.id).join(','))

  if (query.page > 1) params.set('page', query.page.toString())
  if (query.pageSize !== 20) params.set('pageSize', query.pageSize.toString())

  return params
}

export async function searchTags(query: string): Promise<DefinedTagDto[]> {
  const response = await httpClient.get('/_api/tags', { searchParams: { q: query } })
  const json = await response.json()

  return z.array(definedTagDtoSchema).parse(json)
}
