import {
  getPreloadedData,
  httpClient,
  parseQueryStringArray,
  stringArrayToQueryParameterValue,
} from '@/modules/common/api'
import { ProtoAgeRating, ProtoSearchResult, ProtoTagsCategory } from '@/proto/search'
import { AgeRating, DefinedTagDto } from './api'
import { DateTime } from 'luxon'

export type BookSearchItem = {
  id: string
  name: string
  createdAt: string
  updatedAt: string
  ageRating: AgeRating
  words: number
  wordsPerChapter: number
  chapters: number
  favorites: number
  summary: string
  author: {
    id: string
    name: string
  }
  tags: string[]
  cover: string
}

export type SearchBooksResponse = {
  meta: {
    cacheHit: boolean
    cacheKey: string
    cacheTook: number
  }
  page: number
  pageSize: number
  totalPages: number
  booksTook: number
  books: BookSearchItem[]
  tags: DefinedTagDto[]
}

export type SearchBooksRequest = {
  'w.min'?: string
  'w.max'?: string
  'c.min'?: string
  'c.max'?: string
  'wc.min'?: string
  'wc.max'?: string
  'f.min'?: string
  'f.max'?: string
  it?: string[]
  et?: string[]
  iu?: string[]
  eu?: string[]
  p: number
}

export function isSearchBooksRequestEqual(req1: SearchBooksRequest, req2: SearchBooksRequest) {
  return (
    searchBooksRequestToURLSearchParams(req1).toString() ===
    searchBooksRequestToURLSearchParams(req2).toString()
  )
}

export function parseSearchBooksRequest(sp: URLSearchParams): SearchBooksRequest {
  let p = parseInt(sp.get('p') || '1')
  if (Number.isNaN(p)) p = 1
  else if (p < 1) p = 1

  return {
    'w.min': sp.get('w.min') || undefined,
    'w.max': sp.get('w.max') || undefined,
    'c.min': sp.get('c.min') || undefined,
    'c.max': sp.get('c.max') || undefined,
    'wc.min': sp.get('wc.min') || undefined,
    'wc.max': sp.get('wc.max') || undefined,
    'f.min': sp.get('f.min') || undefined,
    'f.max': sp.get('f.max') || undefined,
    it: parseQueryStringArray(sp.get('it')),
    et: parseQueryStringArray(sp.get('et')),
    iu: parseQueryStringArray(sp.get('iu')),
    p,
  }
}

export function searchBooksRequestToURLSearchParams(query: SearchBooksRequest): URLSearchParams {
  const urlSp = new URLSearchParams()

  if (query['w.max']) urlSp.set('w.max', query['w.max'])
  if (query['w.min']) urlSp.set('w.min', query['w.min'])
  if (query['c.max']) urlSp.set('c.max', query['c.max'])
  if (query['c.min']) urlSp.set('c.min', query['c.min'])
  if (query['wc.max']) urlSp.set('wc.max', query['wc.max'])
  if (query['wc.min']) urlSp.set('wc.min', query['wc.min'])
  if (query['f.max']) urlSp.set('f.max', query['f.max'])
  if (query['f.min']) urlSp.set('f.min', query['f.min'])
  if (query.it && query.it.length) urlSp.set('it', stringArrayToQueryParameterValue(query.it) || '')
  if (query.et && query.et.length) urlSp.set('et', stringArrayToQueryParameterValue(query.et) || '')
  if (query.iu && query.iu.length) urlSp.set('iu', stringArrayToQueryParameterValue(query.iu) || '')
  if (query.eu && query.eu.length) urlSp.set('eu', stringArrayToQueryParameterValue(query.eu) || '')
  if (query.p > 1) urlSp.set('p', query.p.toString())

  return urlSp
}

function protoAgeRating(r: ProtoAgeRating): AgeRating {
  switch (r) {
    case ProtoAgeRating.G:
      return 'G'
    case ProtoAgeRating.PG:
      return 'PG'
    case ProtoAgeRating.PG13:
      return 'PG-13'
    case ProtoAgeRating.R:
      return 'R'
    case ProtoAgeRating.NC17:
      return 'NC-17'
    case ProtoAgeRating.UNRECOGNIZED:
    case ProtoAgeRating.UNKNOWN:
    default:
      return '?'
  }
}

function protoTagsCategory(v: ProtoTagsCategory): DefinedTagDto['cat'] {
  switch (v) {
    case ProtoTagsCategory.OTHER:
      return 'other'
    case ProtoTagsCategory.REL:
      return 'rel'
    case ProtoTagsCategory.REL_TYPE:
      return 'reltype'
    case ProtoTagsCategory.FANDOM:
      return 'fandom'
    case ProtoTagsCategory.WARNING:
      return 'warning'
    case ProtoTagsCategory.UNRECOGNIZED:
    default:
      return 'other'
  }
}

export async function httpSearchBooks(query: SearchBooksRequest): Promise<SearchBooksResponse> {
  const sp = searchBooksRequestToURLSearchParams(query)

  function pbToDto(pbResult: ProtoSearchResult): SearchBooksResponse {
    const mappedResponse: SearchBooksResponse = {
      meta: {
        cacheHit: pbResult.cacheHit,
        cacheKey: pbResult.cacheKey,
        cacheTook: pbResult.cacheTook,
      },
      booksTook: pbResult.took,
      books: pbResult.items.map((item) => {
        const book: BookSearchItem = {
          id: item.id,
          words: item.words,
          chapters: item.chapters,
          favorites: item.favorites,
          cover: item.cover,
          name: item.name,
          createdAt: DateTime.fromMillis(item.createdAt * 1000).toISO(),
          updatedAt: DateTime.fromMillis(item.updatedAt * 1000).toISO(),
          ageRating: protoAgeRating(item.ageRating),
          wordsPerChapter: 0,
          summary: item.summary,
          author: {
            id: item.authorId,
            name: item.authorName,
          },
          tags: item.tagIds,
        }

        return book
      }),
      tags: pbResult.tags.map((tag) => {
        return {
          id: tag.id,
          name: tag.name,
          desc: tag.description,
          adult: tag.isAdult,
          spoiler: tag.isSpoiler,
          cat: protoTagsCategory(tag.category),
        }
      }),
      page: pbResult.page,
      pageSize: pbResult.pageSize,
      totalPages: pbResult.totalPages,
    }

    return mappedResponse
  }

  const get = async (): Promise<SearchBooksResponse> => {
    const response = await httpClient.get('/api/search', {
      searchParams: sp,
    })

    if (response.status > 299)
      throw new Error(`Unexpected status code: ${response.status} ${response.statusText}`)
    const contentType = response.headers.get('Content-Type')

    if (contentType === 'application/json') {
      return await response.json()
    }

    if (contentType !== 'application/vnd.google.protobuf') {
      throw new Error(`Unexpected Content-Type: ${contentType}`)
    }

    if (!response.body) throw new Error('No response body')

    const reader = response.body.getReader()
    const binaryContent = await reader.read()

    if (!binaryContent.value) throw new Error("Failed to ready response's binary content")

    const pbResult = ProtoSearchResult.decode(binaryContent.value)
    return pbToDto(pbResult)
  }

  const cacheKey = `/api/search?${sp.toString()}`

  if (__server__._preload && __server__._preload[cacheKey]) {
    const value = __server__._preload[cacheKey]
    if (typeof value !== 'string') {
      return await get()
    } else {
      delete __server__._preload[cacheKey]
      const arr = Uint8Array.from(window.atob(value), (v) => v.charCodeAt(0))
      const pbResult = ProtoSearchResult.decode(arr)
      return pbToDto(pbResult)
    }
  } else {
    return await get()
  }
}

export function getPreloadedBookSearchResult(query: SearchBooksRequest) {
  return getPreloadedData<SearchBooksResponse>(
    `/api/search?${searchBooksRequestToURLSearchParams(query).toString()}`,
  )
}
