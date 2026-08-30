#!/usr/bin/env node

import { readFile, readdir, stat, writeFile } from 'node:fs/promises'
import { extname, join, relative, resolve } from 'node:path'
import { pathToFileURL } from 'node:url'
import process from 'node:process'

const repoRoot = resolve(import.meta.dirname, '..')
const sourceExtensions = new Set([
  '.css',
  '.scss',
  '.html',
  '.templ',
  '.js',
  '.jsx',
  '.mjs',
  '.cjs',
  '.ts',
  '.tsx',
])
const stylesheetExtensions = new Set(['.css', '.scss'])
const ignoredDirectories = new Set(['.git', 'node_modules', 'dist', 'build', 'coverage'])

const explicitBemPattern =
  /(?<![A-Za-z0-9_-])([A-Za-z][A-Za-z0-9-]*(?:__[A-Za-z0-9_-]+|--[A-Za-z0-9_-]+))(?=$|[^A-Za-z0-9_-])/g
const classSelectorPattern = /\.([A-Za-z_][A-Za-z0-9_-]*)/g

export function componentName(name) {
  if (/^[A-Z][A-Za-z0-9]*$/.test(name)) return name
  return name
    .split('-')
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join('')
}

export function descendantName(name, descendantCase = 'camel') {
  if (descendantCase === 'kebab') return name
  const [first = '', ...rest] = name.split('-')
  return first + rest.map((part) => part.charAt(0).toUpperCase() + part.slice(1)).join('')
}

export function parseBemToken(token) {
  const elementAt = token.indexOf('__')
  const modifierAt = token.indexOf('--')

  if (elementAt === -1 && modifierAt === -1) return null

  if (elementAt !== -1 && (modifierAt === -1 || elementAt < modifierAt)) {
    const block = token.slice(0, elementAt)
    const elementAndModifier = token.slice(elementAt + 2)
    const elementModifierAt = elementAndModifier.indexOf('--')
    return {
      block,
      element:
        elementModifierAt === -1
          ? elementAndModifier
          : elementAndModifier.slice(0, elementModifierAt),
      modifier: elementModifierAt === -1 ? null : elementAndModifier.slice(elementModifierAt + 2),
    }
  }

  return {
    block: token.slice(0, modifierAt),
    element: null,
    modifier: token.slice(modifierAt + 2),
  }
}

export function convertBemToken(token, options = {}) {
  const parsed = parseBemToken(token)
  if (!parsed) return token

  const descendantCase = options.descendantCase ?? 'camel'
  const components = options.components ?? new Map()
  const block = components.get(parsed.block) ?? componentName(parsed.block)
  const element = parsed.element ? `-${descendantName(parsed.element, descendantCase)}` : ''
  const modifier = parsed.modifier ? `--${descendantName(parsed.modifier, descendantCase)}` : ''

  return `${block}${element}${modifier}`
}

export function discoverComponents(sources) {
  const components = new Set()

  for (const source of sources) {
    for (const match of source.matchAll(explicitBemPattern)) {
      const parsed = parseBemToken(match[1])
      if (parsed?.block) components.add(parsed.block)
    }
  }

  return components
}

function replaceExplicitBem(source, options) {
  return source.replace(explicitBemPattern, (match) => convertBemToken(match, options))
}

function replaceClassTokens(value, components, options) {
  return value.replace(/(^|\s)([^\s]+)(?=\s|$)/g, (match, whitespace, token) => {
    const converted = convertClassToken(token, components, options)
    return `${whitespace}${converted}`
  })
}

function convertClassToken(token, components, options) {
  const importantPrefix = token.startsWith('!') ? '!' : ''
  const bareToken = importantPrefix ? token.slice(1) : token
  if (parseBemToken(bareToken)) return importantPrefix + convertBemToken(bareToken, options)
  if (components.has(bareToken)) {
    return importantPrefix + (options.components.get(bareToken) ?? componentName(bareToken))
  }
  return token
}

function replaceStaticClassAttributes(source, components, options) {
  return source.replace(
    /(\bclass(?:Name)?\s*=\s*)(["'])([^"']*)(\2)/g,
    (match, prefix, quote, value) =>
      `${prefix}${quote}${replaceClassTokens(value, components, options)}${quote}`,
  )
}

function replaceClassBearingStringLiterals(source, components, options) {
  const stringPattern = /(["'`])([^"'`\n]*)(\1)/g

  return source.replace(stringPattern, (match, quote, value, endQuote, offset) => {
    const context = source.slice(Math.max(0, offset - 160), offset)
    if (
      !/(?:className|classList|clsx|\bcn|querySelector(?:All)?|closest|matches)[\s\S]*$/m.test(
        context,
      )
    ) {
      return match
    }

    // Selector APIs use .Component; class composition APIs use whitespace-separated tokens.
    let converted = value.replace(classSelectorPattern, (selector, name) => {
      if (parseBemToken(name)) return `.${convertBemToken(name, options)}`
      if (!components.has(name)) return selector
      return `.${options.components.get(name) ?? componentName(name)}`
    })
    converted = replaceClassTokens(converted, components, options)
    return `${quote}${converted}${endQuote}`
  })
}

export function transformSource(source, extension, components, options = {}) {
  const normalizedOptions = {
    components: options.components ?? new Map(),
    descendantCase: options.descendantCase ?? 'camel',
  }
  let result = replaceExplicitBem(source, normalizedOptions)

  if (stylesheetExtensions.has(extension)) {
    result = result.replace(classSelectorPattern, (selector, name) => {
      if (!components.has(name)) return selector
      return `.${normalizedOptions.components.get(name) ?? componentName(name)}`
    })
    result = result.replace(
      /&__([A-Za-z0-9_-]+?)(?:--([A-Za-z0-9_-]+))?(?=$|[^A-Za-z0-9_-])/g,
      (_, name, modifier) =>
        `&-${descendantName(name, normalizedOptions.descendantCase)}` +
        (modifier ? `--${descendantName(modifier, normalizedOptions.descendantCase)}` : ''),
    )
    result = result.replace(
      /&--([A-Za-z0-9_-]+)/g,
      (_, name) => `&--${descendantName(name, normalizedOptions.descendantCase)}`,
    )
    return result
  }

  result = replaceStaticClassAttributes(result, components, normalizedOptions)
  return replaceClassBearingStringLiterals(result, components, normalizedOptions)
}

async function main() {
  const config = parseArguments(process.argv.slice(2))
  if (config.help) {
    printHelp()
    return
  }

  const componentOverrides = await readComponentOverrides(config.mapPath)
  const targetPaths = config.paths.length === 0 ? [join(repoRoot, 'web')] : config.paths
  const files = []
  for (const target of targetPaths)
    files.push(...(await collectSourceFiles(resolve(repoRoot, target))))

  const records = await Promise.all(
    [...new Set(files)]
      .sort()
      .map(async (file) => ({ file, source: await readFile(file, 'utf8') })),
  )
  const components = discoverComponents(records.map(({ source }) => source))
  for (const component of componentOverrides.keys()) components.add(component)

  const options = {
    components: componentOverrides,
    descendantCase: config.descendantCase,
  }
  validateMappings(records, components, options)
  const changed = []

  for (const record of records) {
    const updated = transformSource(record.source, extname(record.file), components, options)
    if (updated === record.source) continue
    changed.push({ ...record, updated })
    if (config.write) await writeFile(record.file, updated, 'utf8')
  }

  if (config.printMap) printComponentMap(components, componentOverrides)

  if (changed.length === 0) {
    console.log('No BEM-to-SUIT changes found.')
    return
  }

  const verb = config.write ? 'Updated' : 'Would update'
  for (const { file } of changed) console.log(`${verb}: ${relative(repoRoot, file)}`)
  console.log(
    `\n${config.write ? 'Updated' : 'Found'} ${changed.length} file${changed.length === 1 ? '' : 's'} ` +
      `across ${components.size} component famil${components.size === 1 ? 'y' : 'ies'}.`,
  )

  const dynamicWarnings = findDynamicWarnings(
    changed.map(({ file, updated }) => ({ file, source: updated })),
  )
  if (dynamicWarnings.length > 0) {
    console.log('\nReview possible dynamically constructed BEM names:')
    for (const warning of dynamicWarnings) console.log(`  ${warning}`)
  }

  if (!config.write) console.log('\nDry run only. Re-run with --write to apply the changes.')
  if (config.check && changed.length > 0) process.exitCode = 1
}

function parseArguments(argv) {
  const config = {
    check: false,
    descendantCase: 'camel',
    help: false,
    mapPath: null,
    paths: [],
    printMap: false,
    write: false,
  }

  for (let index = 0; index < argv.length; index++) {
    const arg = argv[index]
    if (arg === '--write') config.write = true
    else if (arg === '--check') config.check = true
    else if (arg === '--print-map') config.printMap = true
    else if (arg === '--help' || arg === '-h') config.help = true
    else if (arg === '--map') config.mapPath = argv[++index]
    else if (arg.startsWith('--map=')) config.mapPath = arg.slice('--map='.length)
    else if (arg.startsWith('--descendant-case=')) {
      config.descendantCase = arg.slice('--descendant-case='.length)
    } else if (arg.startsWith('-')) throw new Error(`Unknown option: ${arg}`)
    else config.paths.push(arg)
  }

  if (!['camel', 'kebab'].includes(config.descendantCase)) {
    throw new Error('--descendant-case must be "camel" or "kebab"')
  }
  if (config.mapPath === undefined) throw new Error('--map requires a JSON file path')
  return config
}

async function readComponentOverrides(mapPath) {
  if (!mapPath) return new Map()
  const parsed = JSON.parse(await readFile(resolve(repoRoot, mapPath), 'utf8'))
  if (!parsed || Array.isArray(parsed) || typeof parsed !== 'object') {
    throw new Error(
      'The component map must be a JSON object of BEM block names to SUIT component names',
    )
  }

  const overrides = new Map(Object.entries(parsed))
  for (const [from, to] of overrides) {
    if (!/^[A-Za-z][A-Za-z0-9-]*$/.test(from) || !/^[A-Z][A-Za-z0-9]*$/.test(to)) {
      throw new Error(`Invalid component mapping: ${JSON.stringify(from)} -> ${JSON.stringify(to)}`)
    }
  }
  return overrides
}

async function collectSourceFiles(path) {
  const pathStat = await stat(path)
  if (!pathStat.isDirectory()) return isSourceFile(path) ? [path] : []

  const files = []
  for (const entry of await readdir(path, { withFileTypes: true })) {
    if (entry.isDirectory()) {
      if (!ignoredDirectories.has(entry.name))
        files.push(...(await collectSourceFiles(join(path, entry.name))))
    } else {
      const file = join(path, entry.name)
      if (isSourceFile(file)) files.push(file)
    }
  }
  return files
}

function validateMappings(records, components, options) {
  const componentTargets = new Map()
  for (const block of components) {
    addMapping(componentTargets, options.components.get(block) ?? componentName(block), block)
  }

  const classTargets = new Map()
  for (const { source } of records) {
    for (const match of source.matchAll(explicitBemPattern)) {
      addMapping(classTargets, convertBemToken(match[1], options), match[1])
    }
  }
}

function addMapping(targets, target, source) {
  const previous = targets.get(target)
  if (previous && previous !== source) {
    throw new Error(
      `Naming collision: ${JSON.stringify(previous)} and ${JSON.stringify(source)} both map to ${JSON.stringify(target)}`,
    )
  }
  targets.set(target, source)
}

function isSourceFile(file) {
  return sourceExtensions.has(extname(file)) && !file.endsWith('_templ.go')
}

function findDynamicWarnings(records) {
  const warnings = []
  for (const { file, source } of records) {
    source.split('\n').forEach((line, index) => {
      if (/(?:[A-Za-z][A-Za-z0-9-]*|&)__(?:\$\{|[A-Za-z0-9-]*\s*\+)/.test(line)) {
        warnings.push(`${relative(repoRoot, file)}:${index + 1}`)
      }
    })
  }
  return warnings
}

function printComponentMap(components, overrides) {
  const map = Object.fromEntries(
    [...components].sort().map((name) => [name, overrides.get(name) ?? componentName(name)]),
  )
  console.log(`${JSON.stringify(map, null, 2)}\n`)
}

function printHelp() {
  console.log(`Usage: node scripts/bem-to-suit.mjs [options] [path ...]

Convert BEM class names to SUIT CSS names. The default target is web/ and the
command is read-only unless --write is supplied.

  block-name                   -> BlockName
  block-name__child-name       -> BlockName-childName
  block-name--feature-name     -> BlockName--featureName
  &__child-name                -> &-childName
  &--feature-name              -> &--featureName

Options:
  --write                      Apply changes (otherwise perform a dry run)
  --check                      Exit 1 when the dry run finds changes
  --print-map                  Print the discovered component mapping
  --map <file>                 Override component names with a JSON object
  --descendant-case=camel      Use canonical SUIT childName casing (default)
  --descendant-case=kebab      Preserve child-name casing after the SUIT dash
  -h, --help                   Show this help

Example override file:
  { "ol-comments": "OLComments", "btn": "Button" }

After writing, review the diff and run templ generate, pnpm run build, affected
Go tests, and git diff --check.`)
}

if (process.argv[1] && pathToFileURL(resolve(process.argv[1])).href === import.meta.url) {
  main().catch((error) => {
    console.error(error.message)
    process.exitCode = 1
  })
}
