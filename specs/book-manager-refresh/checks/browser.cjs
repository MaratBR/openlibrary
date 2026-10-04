const { createRequire } = require('module')
const fs = require('fs'),
  path = require('path')
const root = path.resolve(__dirname, '../../..')
const repoRequire = createRequire(path.join(root, 'package.json'))
let playwright
try {
  playwright = repoRequire('playwright')
} catch {
  const mcpRequire = createRequire(
    fs.realpathSync(path.join(root, 'node_modules/@playwright/mcp/package.json')),
  )
  playwright = mcpRequire('playwright')
}
const { chromium } = playwright
const assert = require('assert/strict')
const server = require('./server.cjs')
const screenshots = fs.mkdtempSync(path.join(require('os').tmpdir(), 'bm-refresh-check-'))
let browser
;(async () => {
  await new Promise((resolve) => server.listen(0, '127.0.0.1', resolve))
  const baseURL = `http://127.0.0.1:${server.address().port}`
  browser = await chromium.launch({
    headless: true,
    executablePath: process.env.PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH,
    args: ['--no-sandbox'],
  })
  const page = await browser.newPage()
  const errors = []
  page.on('pageerror', (e) => {
    errors.push(e.message)
    console.error('PAGEERROR', e.message)
  })
  let fail = false,
    delay = 0,
    empty = false,
    totalPages = 1,
    listFail = false,
    detailFail = false,
    mutations = 0,
    creations = 0,
    many = false
  const chapter = {
    id: '11',
    name: 'Opening chapter',
    words: 1,
    summary: '<p>Chapter summary</p>',
    isPubliclyVisible: false,
    scheduledAt: '2026-11-03T08:00:00Z',
  }
  const book = {
    id: '1',
    name: 'The Glass Orchard',
    ageRating: 'PG',
    adult: false,
    tags: [],
    words: 1,
    wordsPerChapter: 1,
    isPubliclyVisible: true,
    isBanned: false,
    isTrashed: false,
    summary: '',
    cover: { url: '' },
    chapters: [chapter],
  }
  await page.route('**/_api/**', async (route) => {
    const r = route.request(),
      url = new URL(r.url())
    if (r.method() === 'POST') {
      mutations++
      if (delay) await new Promise((r) => setTimeout(r, delay))
      if (fail) return route.fulfill({ status: 400, json: { message: 'Fixture failure' } })
      if (url.pathname.endsWith('/trash')) {
        book.isTrashed = url.searchParams.get('trash') === 'true'
        return route.fulfill({ json: 'ok' })
      }
      if (url.pathname.endsWith('/create-chapter')) {
        book.chapters.push({ ...chapter, id: '12', name: JSON.parse(r.postData()).name })
        return route.fulfill({ json: '12' })
      }
      const data = r.postDataJSON()
      if (url.pathname.includes('/chapter/')) {
        Object.assign(chapter, data)
        return route.fulfill({ json: 'ok' })
      }
      Object.assign(book, data)
      return route.fulfill({ json: book })
    }
    if (url.pathname.endsWith('/tags')) return route.fulfill({ json: [] })
    if (url.pathname.endsWith('/books')) {
      if (listFail) return route.fulfill({ status: 500, json: { message: 'Fixture load error' } })
      return route.fulfill({
        json: {
          books: empty
            ? []
            : Array.from({ length: many ? 20 : 1 }, (_, i) => ({
                ...book,
                id: String(i + 1),
                chapters: book.chapters.length,
              })),
          page: Number(url.searchParams.get('page')),
          totalPages,
        },
      })
    }
    if (detailFail) return route.fulfill({ status: 404, json: { message: 'Unavailable book' } })
    return route.fulfill({ json: book })
  })
  await page.route('**/books-manager/new', async (route) => {
    creations++
    await route.fulfill({ contentType: 'text/html', body: '<h1>Creation received</h1>' })
  })
  const go = async (hash) => {
    await page.goto(baseURL + '/#' + hash)
    await page
      .locator('.BM')
      .waitFor()
      .catch(async (e) => {
        console.log('FAILED PAGE', hash, await page.locator('body').innerText())
        throw e
      })
  }
  const fit = async (label) => {
    const d = await page.evaluate(() => ({
      w: innerWidth,
      sw: document.documentElement.scrollWidth,
    }))
    assert.ok(d.sw <= d.w, `${label}: ${JSON.stringify(d)}`)
  }
  for (const width of [360, 768, 1280]) {
    await page.setViewportSize({ width, height: 900 })
    for (const dark of [false, true]) {
      for (const route of [
        '/books',
        '/books/1',
        '/books/1?t=analytics',
        '/books/1/edit',
        '/books/1?t=chapters',
        '/books/new',
      ]) {
        await go(route)
        await page.evaluate((d) => document.documentElement.classList.toggle('dark', d), dark)
        await page.waitForTimeout(100)
        await fit(`${width} ${dark} ${route}`)
      }
    }
  }
  console.log('PASS responsive routes in both themes at 360/768/1280')
  await go('/books/1')
  await page.getByText('No summary yet.', { exact: true }).waitFor()
  await go('/books/1?t=analytics')
  await page.getByText('No summary yet.', { exact: true }).waitFor()
  await go('/books')
  assert.equal(await page.locator('.Pagination').count(), 0)
  await page.locator('summary[aria-label="Actions for The Glass Orchard"]').click()
  await page.getByRole('button', { name: 'Trash', exact: true }).click()
  fail = true
  await page.getByRole('dialog').getByRole('button', { name: 'Trash', exact: true }).click()
  await page.getByRole('alert').waitFor()
  assert.equal(book.isTrashed, false)
  fail = false
  delay = 350
  await page.getByRole('dialog').getByRole('button', { name: 'Trash', exact: true }).click()
  assert.equal(
    await page.getByRole('dialog').getByRole('button', { name: 'Saving…' }).isDisabled(),
    true,
  )
  await page.getByRole('dialog').waitFor({ state: 'detached' })
  await page.getByText('Trashed · PG').waitFor()
  delay = 0
  const beforeCancel = mutations
  await page.locator('summary[aria-label="Actions for The Glass Orchard"]').click()
  await page.getByRole('button', { name: 'Restore', exact: true }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Cancel', exact: true }).click()
  assert.equal(mutations, beforeCancel)
  assert.equal(book.isTrashed, true)
  await page.locator('summary[aria-label="Actions for The Glass Orchard"]').click()
  await page.getByRole('button', { name: 'Restore', exact: true }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Restore', exact: true }).click()
  await page.getByRole('dialog').waitFor({ state: 'detached' })
  assert.equal(book.isTrashed, false)
  console.log('PASS trash/restore, cancel, error/pending/success and default/old overview routes')
  await go('/books/1/edit')
  await page.getByRole('textbox', { name: 'Name', exact: true }).fill('Changed title')
  const adultField = page.getByRole('checkbox', { name: 'Adult', exact: true })
  await adultField.focus()
  await page.waitForTimeout(100)
  const focusClear = await adultField.evaluate(
    (e) =>
      e.getBoundingClientRect().bottom <=
      document.querySelector('.BM-actionBar').getBoundingClientRect().top,
  )
  assert.equal(focusClear, true)

  await page.getByRole('link', { name: 'Overview', exact: true }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Stay', exact: true }).click()
  assert.equal(
    await page.getByRole('textbox', { name: 'Name', exact: true }).inputValue(),
    'Changed title',
  )
  fail = true
  await page.getByRole('button', { name: 'Save changes', exact: true }).click()
  await page.getByRole('alert').waitFor()
  assert.equal(book.name, 'The Glass Orchard')
  fail = false
  await page.getByRole('button', { name: 'Save changes', exact: true }).click()
  await page.getByText('Changes saved.', { exact: true }).waitFor()
  await page.getByRole('button', { name: 'Cancel', exact: true }).click()
  await page.getByRole('heading', { name: 'Changed title', exact: true }).waitFor()
  assert.equal(await page.getByRole('dialog').count(), 0)
  console.log('PASS metadata dirty navigation, failure retained, save and clean cancel')
  await go('/books/1?t=chapters')
  await page.getByRole('button', { name: 'Edit details', exact: true }).click()
  await page.getByRole('dialog', { name: 'Chapter details', exact: true }).waitFor()
  assert.equal(
    await page.locator('#chapter-name-input').evaluate((e) => e === document.activeElement),
    true,
  )
  await page.keyboard.press('Escape')
  await page.getByRole('dialog').waitFor({ state: 'detached' })
  assert.equal(
    await page
      .getByRole('button', { name: 'Edit details', exact: true })
      .evaluate((e) => e === document.activeElement),
    true,
  )
  await page.getByRole('button', { name: 'Edit details', exact: true }).click()
  await page.locator('#chapter-name-input').fill('Changed chapter')
  await page.keyboard.press('Escape')
  await page.getByRole('button', { name: 'Stay', exact: true }).click()
  assert.equal(await page.locator('#chapter-name-input').inputValue(), 'Changed chapter')
  fail = true
  await page.getByRole('button', { name: 'Save changes', exact: true }).click()
  await page.getByRole('alert').waitFor()
  fail = false
  await page.getByRole('button', { name: 'Save changes', exact: true }).click()
  await page.getByRole('dialog').waitFor({ state: 'detached' })
  await page.getByRole('heading', { name: 'Changed chapter', exact: true }).waitFor()
  console.log('PASS chapter focus, Escape, dirty prompt and failed/successful save')
  await page.getByRole('button', { name: 'Add chapter', exact: false }).click()
  await page.locator('#add-chapter-name').fill('   ')
  const before = mutations
  await page.locator('#add-chapter-name').press('Enter')
  assert.equal(mutations, before)
  await page.locator('#add-chapter-name').fill('New chapter')
  fail = true
  await page.locator('#add-chapter-name').press('Enter')
  await page.getByRole('alert').waitFor()
  assert.equal(await page.locator('#add-chapter-name').inputValue(), 'New chapter')
  fail = false
  await page.locator('#add-chapter-name').press('Enter')
  await page.locator('#add-chapter-name').waitFor({ state: 'detached' })
  await page.getByRole('heading', { name: 'New chapter' }).waitFor()
  console.log('PASS Add chapter Enter validation and failed/successful save')
  await go('/books/new')
  await page.locator('#new-book-name').fill(' ')
  await page.locator('#new-book-name').press('Enter')
  await page.getByText('Step 1 of 4').waitFor()
  await page.locator('#new-book-name').fill('Wizard book')
  await page.locator('#new-book-name').press('Enter')
  await page.getByText('Step 2 of 4').waitFor()
  await page.locator('label[for="new-book-PG"]').click()
  assert.equal(await page.getByLabel('PG', { exact: true }).isChecked(), true)
  await page.getByRole('button', { name: 'Next', exact: true }).click()
  await page.getByText('Tags (optional)', { exact: true }).waitFor()
  await page.getByRole('button', { name: 'Next', exact: true }).click()
  await page
    .getByText('No tags selected', { exact: true })
    .waitFor()
    .catch(async (e) => {
      console.log('WIZARD FAILURE', page.url(), await page.locator('body').innerText())
      throw e
    })
  await page.getByRole('button', { name: '1. Name', exact: true }).click()
  assert.equal(await page.locator('#new-book-name').inputValue(), 'Wizard book')
  await page.getByRole('link', { name: 'Back to books', exact: false }).click()
  await page.getByRole('button', { name: 'Discard', exact: true }).click()
  await page.getByRole('heading', { name: 'Your books' }).waitFor()
  assert.equal(creations, 0)
  console.log('PASS wizard validation, step navigation, review and discard; no implicit creation')
  empty = true
  await go('/books')
  await page.getByText('Your library starts here.', { exact: false }).waitFor()
  empty = false
  totalPages = 3
  await go('/books')
  await page.locator('.Pagination').waitFor()
  await page.getByRole('link', { name: '2', exact: true }).click()
  await page.waitForURL('**page=2')
  assert.ok(page.url().includes('page=2'))
  listFail = true
  await go('/books')
  await page.getByRole('heading', { name: 'Unable to load this view' }).waitFor()
  listFail = false
  await page.getByRole('button', { name: 'Try again' }).click()
  await page.getByRole('heading', { name: 'Your books' }).waitFor()
  await go('/')
  await page.getByRole('heading', { name: 'Your books' }).waitFor()
  await go('/unknown-route')
  await page.getByText('404', { exact: true }).waitFor()
  detailFail = true
  await go('/books/1')
  await page.getByRole('heading', { name: 'Unable to load this view' }).waitFor()
  detailFail = false
  await page.getByRole('link', { name: 'Back to books', exact: true }).click()
  await page.getByRole('heading', { name: 'Your books' }).waitFor()
  console.log('PASS empty library, pagination, default/catch-all routes and list/detail recovery')

  for (const width of [360, 768, 1280]) {
    await page.setViewportSize({ width, height: 800 })
    await go('/books/1?t=chapters')
    await page.getByRole('button', { name: 'Edit details', exact: true }).first().click()
    await page.getByRole('dialog').waitFor()
    await fit('chapter panel ' + width)
    const rect = await page.getByRole('dialog').boundingBox()
    assert.ok(rect.x >= 0 && rect.x + rect.width <= width)
    await page.keyboard.press('Escape')
    await page.getByRole('button', { name: 'Add chapter', exact: false }).click()
    await page.locator('#add-chapter-name').waitFor()
    await fit('add chapter ' + width)
    await page.getByRole('button', { name: 'Cancel', exact: true }).click()
  }
  many = true
  book.name = 'An exceptionally long book title '.repeat(12)
  book.cover.url = '/missing-cover.png'
  for (const width of [360, 768, 1280]) {
    await page.setViewportSize({ width, height: 900 })
    await go('/books')
    await page.locator('.BM-bookCard').first().waitFor()
    assert.equal(await page.locator('.BM-bookCard').count(), 20)
    await fit('20 long-title cards ' + width)
  }
  console.log('PASS bounded panels/popovers, 20-card library, long titles and broken covers')
  book.name = 'The Glass Orchard'
  many = false
  book.cover.url = ''
  await go('/books')
  await page.screenshot({ path: path.join(screenshots, 'library.png'), fullPage: true })
  await go('/books/1/edit')
  await page.screenshot({ path: path.join(screenshots, 'details.png'), fullPage: true })
  await page.getByRole('checkbox', { name: 'Adult', exact: true }).focus()
  await page.keyboard.press('Space')
  await page.getByRole('button', { name: 'Save changes', exact: true }).click()
  await page.getByRole('alert').waitFor()
  assert.equal(await page.getByRole('checkbox', { name: 'Adult', exact: true }).isChecked(), true)
  await page.getByRole('checkbox', { name: 'Adult', exact: true }).focus()
  await page.keyboard.press('Space')
  await go('/books/new')
  await page.locator('#new-book-name').fill('Explicit creation')
  await page.getByRole('button', { name: 'Next', exact: true }).click()
  await page.locator('label[for="new-book-PG"]').click()
  await page.getByRole('button', { name: 'Next', exact: true }).click()
  await page.getByText('Tags (optional)', { exact: true }).waitFor()
  await page.getByRole('button', { name: 'Next', exact: true }).click()
  await page.getByText('No tags selected', { exact: true }).waitFor()
  assert.equal(creations, 0)
  await page.getByRole('button', { name: 'Create book', exact: true }).click()
  await page.getByRole('heading', { name: 'Creation received' }).waitFor()
  assert.equal(creations, 1)
  console.log('PASS unconfirmed Adult setting retained and explicit native creation')

  assert.deepEqual(errors, [])
  console.log('ALL CHECKS PASSED; screenshots:', screenshots)
})()
  .catch((e) => {
    console.error(e)
    process.exitCode = 1
  })
  .finally(async () => {
    await browser?.close()
    server.close()
  })
