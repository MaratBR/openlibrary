import assert from 'node:assert/strict'
import test from 'node:test'
import { build } from 'esbuild'

const result = await build({
  entryPoints: ['web/frontend/src/features/search/model.ts'],
  bundle: true,
  write: false,
  format: 'esm',
  platform: 'node',
})
const { getDefaultDetailedBookSearchQuery, getQueryParams } = await import(
  `data:text/javascript;base64,${Buffer.from(result.outputFiles[0].text).toString('base64')}`
)

test('search defaults to 20 and serializes HTTP pagination and index filters', () => {
  const query = getDefaultDetailedBookSearchQuery()
  assert.equal(query.pageSize, 20)
  assert.equal(getQueryParams(query).toString(), '')
  Object.assign(query, {
    page: 3,
    pageSize: 50,
    rating: 'PG',
    createdAt: { min: '2026-01-01', max: '' },
    popularity: { min: 1.25, max: null },
    views: { min: 0, max: 100 },
  })
  assert.deepEqual(Object.fromEntries(getQueryParams(query)), {
    rating: 'PG',
    'views.min': '0',
    'views.max': '100',
    'popularity.min': '1.25',
    'createdAt.min': '2026-01-01',
    p: '3',
    ps: '50',
  })
})
