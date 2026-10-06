import assert from 'node:assert/strict'
import test from 'node:test'
import { build } from 'esbuild'

// Bundle the actual browser API schema so this check needs no TS test runner.
const result = await build({
  stdin: {
    contents: `export { DefinedTagDto } from './web/frontend/src/features/search/api';
      export { Schema } from 'effect';`,
    resolveDir: process.cwd(),
    loader: 'ts',
  },
  bundle: true,
  write: false,
  format: 'esm',
  platform: 'node',
})
const previousWindow = globalThis.window
// The HTTP client installs browser fetch and DOMContentLoaded handlers at import.
globalThis.window = { fetch: globalThis.fetch, addEventListener() {} }
let DefinedTagDto, Schema
try {
  ;({ DefinedTagDto, Schema } = await import(
    `data:text/javascript;base64,${Buffer.from(result.outputFiles[0].text).toString('base64')}`
  ))
} finally {
  if (previousWindow === undefined) delete globalThis.window
  else globalThis.window = previousWindow
}

test('tag search decodes the Martial Arts genre response', () => {
  const response = [
    {
      id: '6704893821110461035',
      name: 'Martial Arts',
      desc: 'Stories centered on combat disciplines, physical training, tournaments, and skilled fighters.',
      adult: false,
      spoiler: false,
      cat: 'genre',
    },
  ]
  assert.deepEqual(Schema.decodeUnknownSync(Schema.Array(DefinedTagDto))(response), response)
})

test('tag search accepts all server categories and the unknown fallback', () => {
  for (const cat of ['other', 'warning', 'fandom', 'rel', 'reltype', 'genre', 'unknown']) {
    assert.equal(
      Schema.decodeUnknownSync(DefinedTagDto)({
        id: '1',
        name: 'Tag',
        desc: '',
        adult: false,
        spoiler: false,
        cat,
      }).cat,
      cat,
    )
  }
})
