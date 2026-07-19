#!/usr/bin/env node

import { readFile, readdir, writeFile } from 'node:fs/promises'
import { extname, join, relative, resolve } from 'node:path'
import process from 'node:process'

const repoRoot = resolve(import.meta.dirname, '..')
const defaultCatalog = join(repoRoot, 'translations', 'en.toml')
const sourceExtensions = new Set(['.go', '.js', '.jsx', '.mjs', '.cjs', '.ts', '.tsx', '.templ'])
const ignoredDirectories = new Set(['.git', 'node_modules', 'dist', 'build', 'coverage'])

const args = new Set(process.argv.slice(2))
const remove = args.has('--remove')

if (args.has('--help') || args.has('-h')) {
  console.log(`Usage: node scripts/check-unused-translations.mjs [--remove]

Find translation keys from translations/en.toml that are not referenced by:
  - _("key") / window._("key") calls in JS and TypeScript
  - l.T("key") and _tt(l, "key") calls in templ files
  - literal/template-literal keys with a static prefix
  - i18nExtractKeys(...) lists in templ files
  - i18nExtractKeysByPrefix(...) calls in templ files

The command is read-only by default. Pass --remove to delete unused entries.
Formatting, comments, and the order of the TOML catalog are preserved.`)
  process.exit(0)
}

const catalogText = await readFile(defaultCatalog, 'utf8')
const catalog = parseCatalog(catalogText)
const sourceFiles = await collectSourceFiles(repoRoot)
const exactKeys = new Set()
const usedPrefixes = new Set()

for (const file of sourceFiles) {
  const source = await readFile(file, 'utf8')
  collectTranslationCalls(source, exactKeys, usedPrefixes)
  collectExtractionHelpers(source, exactKeys, usedPrefixes)
}

const unused = catalog.entries.filter(({ key }) => !isUsed(key, exactKeys, usedPrefixes))

if (unused.length === 0) {
  console.log(`No unused translations found in ${relative(repoRoot, defaultCatalog)}.`)
  process.exit(0)
}

console.log(`Found ${unused.length} unused translation key${unused.length === 1 ? '' : 's'}:`)
for (const { key } of unused) console.log(`  ${key}`)

if (!remove) {
  console.log('\nDry run only. Re-run with --remove to delete these entries.')
  process.exitCode = 1
} else {
  const unusedLines = new Set(unused.map(({ lineIndex }) => lineIndex))
  const nextText = catalog.lines.filter((_, index) => !unusedLines.has(index)).join('\n')
  await writeFile(defaultCatalog, nextText, 'utf8')
  console.log(`\nRemoved ${unused.length} entries from ${relative(repoRoot, defaultCatalog)}.`)
}

async function collectSourceFiles(directory) {
  const files = []

  for (const entry of await readdir(directory, { withFileTypes: true })) {
    if (entry.isDirectory()) {
      if (!ignoredDirectories.has(entry.name)) files.push(...(await collectSourceFiles(join(directory, entry.name))))
      continue
    }

    const path = join(directory, entry.name)
    if (
      sourceExtensions.has(extname(entry.name)) &&
      !entry.name.endsWith('_templ.go') &&
      path !== defaultCatalog
    ) {
      files.push(path)
    }
  }

  return files
}

function collectTranslationCalls(source, exactKeys, prefixes) {
  const callPatterns = [
    { pattern: /(?:\bwindow\s*\.\s*)?\b_\s*\(\s*(['"`])([^'"`]*?)\1/g },
    { pattern: /\b[A-Za-z_$][\w$]*\s*\.\s*T\s*\(\s*(['"`])([^'"`]*?)\1/g, subtree: true },
    { pattern: /\b_tt\s*\(\s*[^,]+,\s*(['"`])([^'"`]*?)\1/g },
  ]

  for (const { pattern, subtree } of callPatterns) {
    for (const match of source.matchAll(pattern)) {
      const key = match[2]
      const interpolation = key.indexOf('${')
      const afterMatch = source.slice(match.index + match[0].length).trimStart()
      const isDynamic = (match[1] === '`' && interpolation !== -1) || afterMatch.startsWith('+')

      if (isDynamic) {
        const staticPart = interpolation === -1 ? key : key.slice(0, interpolation)
        const prefix = staticPart.replace(/\.$/, '')
        if (prefix) prefixes.add(prefix)
      } else {
        exactKeys.add(key)
        if (subtree) prefixes.add(key)
      }
    }
  }

  const formattedKeyPattern = /\bfmt\.Sprintf\s*\(\s*(['"])([^'"]*?)%[^'"]*\1/g
  for (const match of source.matchAll(formattedKeyPattern)) {
    const prefix = match[2].replace(/\.$/, '')
    if (prefix) prefixes.add(prefix)
  }
}

function collectExtractionHelpers(source, exactKeys, prefixes) {
  const prefixPattern = /\bi18nExtractKeysByPrefix\s*\(\s*[^,]+,\s*(['"])([^'"]+)\1\s*\)/g
  for (const match of source.matchAll(prefixPattern)) prefixes.add(match[2].replace(/\.$/, ''))

  for (const call of findCalls(source, 'i18nExtractKeys')) {
    const stringPattern = /(['"])([^'"]+)\1/g
    for (const match of call.matchAll(stringPattern)) exactKeys.add(match[2])
  }
}

function findCalls(source, functionName) {
  const calls = []
  const pattern = new RegExp(`\\b${functionName}\\s*\\(`, 'g')

  for (const match of source.matchAll(pattern)) {
    let depth = 1
    let quote = null
    let escaped = false
    let index = match.index + match[0].length

    for (; index < source.length && depth > 0; index++) {
      const char = source[index]
      if (quote) {
        if (escaped) escaped = false
        else if (char === '\\') escaped = true
        else if (char === quote) quote = null
      } else if (char === '"' || char === "'" || char === '`') quote = char
      else if (char === '(') depth++
      else if (char === ')') depth--
    }

    if (depth === 0) calls.push(source.slice(match.index, index))
  }

  return calls
}

function parseCatalog(text) {
  const lines = text.split('\n')
  const entries = []
  let section = ''

  lines.forEach((line, lineIndex) => {
    const sectionMatch = line.match(/^\s*\[([^\]]+)]\s*(?:#.*)?$/)
    if (sectionMatch) {
      section = sectionMatch[1].trim()
      return
    }

    const keyMatch = line.match(/^\s*([A-Za-z0-9_-]+)\s*=/)
    if (!keyMatch) return
    entries.push({ key: section ? `${section}.${keyMatch[1]}` : keyMatch[1], lineIndex })
  })

  return { entries, lines }
}

function isUsed(key, exactKeys, prefixes) {
  if (exactKeys.has(key)) return true
  for (const prefix of prefixes) {
    if (key === prefix || key.startsWith(`${prefix}.`)) return true
  }
  return false
}
