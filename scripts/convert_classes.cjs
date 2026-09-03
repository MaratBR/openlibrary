const fs = require('fs')
const path = require('path')

const ROOT = process.argv[2] || '.'

const FILE_EXTENSIONS = new Set(['.ts', '.tsx', '.templ', '.html'])

const TYPE_CLASSES = ['primary', 'secondary', 'destructive']
const BTN_TYPE_CLASSES = TYPE_CLASSES.map((c) => `Btn--${c}`)
const STYLE_CLASSES = ['btn--ghost', 'btn--outline', 'btn--icon']

function processClassList(classString) {
  let classes = classString.trim().split(/\s+/)

  if (!classes.includes('btn')) {
    return classString
  }

  // Find first plain color class
  const plainType = TYPE_CLASSES.find((c) => classes.includes(c))

  if (plainType) {
    // Remove all plain color classes and add the converted one
    classes = classes.filter((c) => !TYPE_CLASSES.includes(c))
    if (!classes.includes(`Btn--${plainType}`)) {
      classes.push(`Btn--${plainType}`)
    }
  }

  // Ensure there is a btn--color
  const hasBtnType = BTN_TYPE_CLASSES.some((c) => classes.includes(c))
  if (!hasBtnType) {
    classes.push('btn--primary')
  }

  // Ensure there is a style
  const hasStyle = STYLE_CLASSES.some((c) => classes.includes(c))
  if (!hasStyle && !classes.includes('btn--solid')) {
    classes.push('btn--solid')
  }

  // Move all btn classes first
  const btnClasses = classes.filter((c) => c === 'btn' || c.startsWith('btn--'))
  const otherClasses = classes.filter((c) => !(c === 'btn' || c.startsWith('btn--')))

  return [...btnClasses, ...otherClasses].join(' ')
}

function processFile(file) {
  let text = fs.readFileSync(file, 'utf8')

  const updated = text.replace(
    /(class(?:Name)?\s*=\s*["'])([^"']*)(["'])/g,
    (_, start, classes, end) => {
      return start + processClassList(classes) + end
    },
  )

  if (updated !== text) {
    fs.writeFileSync(file, updated)
    console.log('Updated:', file)
  }
}

function walk(dir) {
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const full = path.join(dir, entry.name)

    if (entry.isDirectory()) {
      walk(full)
      continue
    }

    if (FILE_EXTENSIONS.has(path.extname(entry.name))) {
      processFile(full)
    }
  }
}

walk(ROOT)
