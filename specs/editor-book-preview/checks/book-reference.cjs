const assert = require('node:assert/strict')
const { mkdtempSync, rmSync } = require('node:fs')
const { tmpdir } = require('node:os')
const { join } = require('node:path')
const { buildSync } = require('esbuild')

const directory = mkdtempSync(join(tmpdir(), 'openlibrary-book-reference-'))
try {
  const outfile = join(directory, 'reference.cjs')
  buildSync({
    entryPoints: ['web/frontend/src/features/book-preview/index.ts'],
    bundle: true,
    platform: 'node',
    format: 'cjs',
    outfile,
  })
  const { parseBookReference } = require(outfile)
  const origin = 'http://localhost:8080'
  const cases = [
    ['00123', '123'],
    ['9223372036854775807', '9223372036854775807'],
    ['/book/123-a-title', '123'],
    [origin + '/book/123?from=editor', '123'],
    [' 123 ', '123'],
    ['0', null],
    ['-1', null],
    ['9223372036854775808', null],
    ['99999999999999999999999999', null],
    ['https://other.test/book/123', null],
    ['//other.test/book/123', null],
    ['/book/123/chapters/1', null],
    ['javascript:alert(1)', null],
    ['1.2', null],
    ['', null],
  ]
  for (const [input, expected] of cases) {
    assert.equal(parseBookReference(input, origin), expected, input)
  }
  console.log(`${cases.length} book reference validation checks passed`)
} finally {
  rmSync(directory, { recursive: true, force: true })
}
