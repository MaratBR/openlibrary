import { stringArrayToQueryParameterValue } from '@/modules/common/api'

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
