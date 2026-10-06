const { createRequire } = require('module')
const fs = require('fs'),
  path = require('path'),
  assert = require('assert/strict')
const root = path.resolve(__dirname, '../../..')
const mcpRequire = createRequire(
  fs.realpathSync(path.join(root, 'node_modules/@playwright/mcp/package.json')),
)
const { chromium } = mcpRequire('playwright')
const server = require('./server.cjs')
let browser
;(async () => {
  await new Promise((resolve) => server.listen(0, '127.0.0.1', resolve))
  const base = `http://127.0.0.1:${server.address().port}`
  browser = await chromium.launch({
    headless: true,
    executablePath: process.env.PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH,
    args: ['--no-sandbox'],
  })
  const page = await browser.newPage()
  const errors = []
  page.on('pageerror', (e) => errors.push(e.message))
  page.on('console', (m) => {
    if (m.type() === 'error' && !/Failed to load resource:.*status of (500|401)/.test(m.text()))
      errors.push(m.text())
  })
  const user = {
    id: '416f17a0-1740-4e61-99e0-e9e5309c8030',
    name: 'Reader',
    role: 'user',
    isBanned: false,
    avatar: '',
    joinedAt: '2026-08-07',
    bio: 'Original biography',
    gender: 'nonbinary',
  }
  const tag = {
    id: '6704893821110461031',
    name: 'Action',
    description: 'A tag description',
    category: 'genre',
    adult: false,
    spoiler: false,
    createdAt: '2026-08-07',
    isDefault: true,
    synonymOf: null,
  }
  const parent = { ...tag, id: '6704893821110461032', name: 'Adventure' }
  let listFail = false,
    saveFail = false,
    delay = 0,
    empty = false,
    posts = [],
    debugCalls = 0,
    unauthorized = false
  await page.route('**/admin/api/**', async (route) => {
    const req = route.request(),
      url = new URL(req.url()),
      pathname = url.pathname
    if (delay) await new Promise((resolve) => setTimeout(resolve, delay))
    let status = 200,
      body
    if (unauthorized) {
      status = 401
      body = { error: 'Unauthorized' }
    } else if (req.method() === 'POST') {
      const fields = Object.fromEntries(new URLSearchParams(req.postData()))
      posts.push({ path: pathname, fields, csrf: req.headers()['x-csrf-token'] })
      if (saveFail) {
        status = 500
        body = { error: 'Failed' }
      } else {
        if (pathname.endsWith('/debug')) debugCalls++
        else if (pathname.includes('/users/'))
          Object.assign(user, { bio: fields.about, gender: fields.gender, role: fields.role })
        else
          Object.assign(tag, {
            name: fields.name,
            description: fields.description,
            category: fields.type,
            adult: fields.adult === 'on',
            spoiler: fields.spoiler === 'on',
            synonymOf: fields.synonymOf ? { id: parent.id, name: parent.name } : null,
          })
        body = { ok: true }
      }
    } else if (listFail) {
      status = 500
      body = { error: 'Failed' }
    } else if (pathname === '/admin/api/users') {
      const pageNo = Number(url.searchParams.get('p') || 1)
      body = {
        users: empty ? [] : [user],
        page: pageNo,
        totalPages: empty ? 0 : 3,
        total: empty ? 0 : 41,
      }
    } else if (pathname.startsWith('/admin/api/users/')) body = user
    else if (pathname === '/admin/api/tags')
      body = {
        tags: empty ? [] : [tag, parent],
        page: Number(url.searchParams.get('p') || 1),
        totalPages: empty ? 0 : 2,
        total: empty ? 0 : 51,
      }
    else if (pathname.startsWith('/admin/api/tags/'))
      body = pathname.endsWith(parent.id) ? parent : tag
    else if (pathname === '/admin/api/debug') body = ['books:elastic:reindex']
    else {
      status = 404
      body = { error: 'Missing' }
    }
    await route.fulfill({ status, contentType: 'application/json', body: JSON.stringify(body) })
  })
  async function go(hash, title) {
    await page.goto(`${base}/admin#${hash}`)
    await page.getByRole('heading', { name: title, level: 1, exact: true }).waitFor()
  }
  const nav = (name) =>
    page.getByRole('navigation', { name: 'Workspace' }).getByRole('link', { name, exact: true })
  const snapshots = fs.mkdtempSync(path.join(require('os').tmpdir(), 'admin-spa-check-'))
  // The real admin Alpine bundle must mount Island successfully, including deep links.
  await go('/', 'Library administration')
  await page.evaluate(() => {
    window.__adminDocumentMarker = 'persistent'
  })
  await nav('Users').click()
  await page.getByRole('heading', { name: 'Users', level: 1 }).waitFor()
  assert.equal(await page.evaluate(() => window.__adminDocumentMarker), 'persistent')
  await page.getByLabel('Search users').fill('reader')
  await page.getByLabel('Role', { exact: true }).selectOption('moderator')
  await page.getByRole('button', { name: 'Search', exact: true }).click()
  await page.waitForURL('**#/users?q=reader&usersFilter.role=moderator')
  await page.getByRole('link', { name: 'Next', exact: true }).click()
  await page.waitForURL('**p=2')
  assert.ok(page.url().includes('q=reader&usersFilter.role=moderator'))
  await page.getByRole('link', { name: 'Edit', exact: true }).click()
  await page.getByRole('heading', { name: 'Reader', level: 1 }).waitFor()
  assert.equal(await page.getByLabel('Custom gender').inputValue(), 'nonbinary')
  await page.getByLabel('About you').fill('Edited biography')
  await nav('Tags').click()
  await page.getByRole('dialog', { name: 'Discard unsaved changes?' }).waitFor()
  await page.getByRole('button', { name: 'Stay', exact: true }).click()
  assert.equal(await page.getByLabel('About you').inputValue(), 'Edited biography')
  await page.getByRole('button', { name: 'Reset password', exact: true }).click()
  assert.ok((await page.getByLabel('New password', { exact: true }).inputValue()).length >= 18)
  await page.getByRole('button', { name: 'Cancel password reset', exact: true }).click()
  saveFail = true
  await page.getByRole('button', { name: 'Save', exact: true }).click()
  await page.getByRole('alert').waitFor()
  assert.equal(await page.getByLabel('About you').inputValue(), 'Edited biography')
  saveFail = false
  delay = 150
  await page.getByRole('button', { name: 'Save', exact: true }).click()
  await page.getByRole('button', { name: 'Saving…', exact: true }).waitFor()
  assert.equal(await page.getByLabel('About you').isDisabled(), true)
  await page.getByText('Changes saved.', { exact: true }).waitFor()
  delay = 0
  const userPost = posts.at(-1)
  assert.equal(userPost.fields.gender, 'nonbinary')
  assert.equal(userPost.fields.password, '')
  assert.equal(userPost.csrf, 'fixture-csrf')
  await page.getByRole('link', { name: /Back to users/ }).click()
  await page.waitForURL('**p=2')
  assert.ok(page.url().includes('q=reader'))
  await nav('Tags').click()
  await page.getByRole('heading', { name: 'Tags', level: 1 }).waitFor()
  await page.getByLabel('Search tags').fill('action')
  await page.getByLabel('Adult tags only').check()
  await page.getByRole('button', { name: 'Apply', exact: true }).click()
  await page.waitForURL('**#/tags?q=action&onlyAdultTags=true')
  await page.getByRole('link', { name: 'Action', exact: true }).click()
  await page.getByRole('heading', { name: 'Action', level: 1 }).waitFor()
  assert.ok((await page.locator('main').innerText()).includes('Genre'))
  await page.getByRole('link', { name: 'Edit', exact: true }).click()
  await page.getByRole('heading', { name: 'Edit tag', level: 1 }).waitFor()
  await page.getByRole('button', { name: 'Change', exact: true }).click()
  await page.getByRole('dialog', { name: 'Choose a parent tag' }).waitFor()
  await page.getByRole('button', { name: 'Adventure', exact: false }).click()
  await page.getByLabel('Description', { exact: true }).fill('')
  await page.getByLabel('Is adult', { exact: false }).check()
  await page.getByLabel('Is adult', { exact: false }).uncheck()
  saveFail = true
  await page.getByRole('button', { name: 'Save', exact: true }).click()
  await page.getByRole('alert').waitFor()
  assert.equal(await page.getByLabel('Description', { exact: true }).inputValue(), '')
  saveFail = false
  await page.getByRole('button', { name: 'Save', exact: true }).click()
  await page.getByText('Changes saved.', { exact: true }).waitFor()
  const tagPost = posts.at(-1)
  assert.equal(tagPost.fields.description, '')
  assert.equal(tagPost.fields.adult, '')
  assert.equal(tagPost.fields.spoiler, '')
  assert.equal(tagPost.fields.synonymOf, '6704893821110461032')
  await page.getByRole('link', { name: 'Cancel', exact: true }).click()
  await page.getByRole('heading', { name: 'Action', level: 1 }).waitFor()
  await page.getByRole('link', { name: /Back to tags/ }).click()
  await page.waitForURL('**#/tags?q=action&onlyAdultTags=true')
  await nav('Debug actions').click()
  await page.getByRole('heading', { name: 'Debug actions', level: 1 }).waitFor()
  await page.getByRole('button', { name: 'Run action', exact: true }).click()
  await page.getByRole('button', { name: 'Cancel', exact: true }).click()
  assert.equal(debugCalls, 0)
  await page.getByRole('button', { name: 'Run action', exact: true }).click()
  await page.getByRole('dialog').getByRole('button', { name: 'Run action', exact: true }).click()
  await page
    .getByText('Reindexing has been scheduled. It may take some time to finish.', { exact: true })
    .waitFor()
  assert.equal(debugCalls, 1)
  // Empty, failure/retry, unauthorized and invalid route states.
  empty = true
  await go('/users', 'Users')
  await page.getByText('No users match your filters.', { exact: true }).waitFor()
  assert.equal(await page.getByRole('link', { name: 'Next', exact: true }).count(), 0)
  empty = false
  listFail = true
  await go('/tags', 'Page unavailable')
  await page.getByRole('alert').waitFor()
  listFail = false
  await page.getByRole('button', { name: 'Try again', exact: true }).click()
  await page.getByRole('heading', { name: 'Tags', level: 1 }).waitFor()
  unauthorized = true
  await go('/users', 'Page unavailable')
  await page.getByRole('link', { name: 'Sign in', exact: true }).waitFor()
  unauthorized = false
  await go('/users/invalid', 'Page unavailable')
  await go('/missing', 'Page not found')
  await go(`/users/${user.id}`, 'Reader')
  await page.getByLabel('About you').fill('Discard this edit')
  await nav('Books').click()
  await page.getByRole('dialog', { name: 'Discard unsaved changes?' }).waitFor()
  await page.getByRole('button', { name: 'Discard changes', exact: true }).click()
  await page.getByRole('heading', { name: 'Books', level: 1 }).waitFor()
  assert.equal(user.bio, 'Edited biography')
  await go('/books', 'Books')
  await page.getByLabel('Book ID', { exact: true }).fill('0')
  await page.getByRole('button', { name: 'Open', exact: true }).click()
  await page.getByRole('alert').waitFor()
  // Every existing route, both themes, three viewport sizes; catch raw keys/overflow.
  const routes = [
    ['/', 'Library administration'],
    ['/users', 'Users'],
    [`/users/${user.id}`, 'Reader'],
    ['/books', 'Books'],
    ['/tags', 'Tags'],
    [`/tags/${tag.id}`, 'Action'],
    [`/tags/${tag.id}/edit`, 'Edit tag'],
    ['/debug', 'Debug actions'],
  ]
  for (const width of [360, 768, 1280]) {
    await page.setViewportSize({ width, height: 900 })
    for (const [hash, title] of routes) {
      await go(hash, title)
      for (const dark of [false, true]) {
        await page.evaluate((dark) => document.documentElement.classList.toggle('dark', dark), dark)
        assert.equal(
          await page.evaluate(() => document.documentElement.scrollWidth > innerWidth),
          false,
          `${hash} overflows at ${width}`,
        )
        assert.equal(
          /\b(?:admin|tag|role|common)\.[a-zA-Z]/.test(await page.locator('main').innerText()),
          false,
          `${hash} untranslated`,
        )
      }
    }
    await go('/', 'Library administration')
    await page.screenshot({ path: path.join(snapshots, `home-${width}.png`), fullPage: true })
  }
  // Pending navigation has a status and does not expose the old route's actions.
  delay = 200
  await nav('Users').click()
  await page.getByText('Loading…', { exact: true }).waitFor()
  await page.getByRole('heading', { name: 'Users', level: 1 }).waitFor()
  delay = 0
  await page.reload()
  await page.getByRole('heading', { name: 'Users', level: 1 }).waitFor()
  await nav('Tags').click()
  await page.getByRole('heading', { name: 'Tags', level: 1 }).waitFor()
  await page.goBack()
  await page.getByRole('heading', { name: 'Users', level: 1 }).waitFor()
  await page.goForward()
  await page.getByRole('heading', { name: 'Tags', level: 1 }).waitFor()
  assert.deepEqual(errors, [])
  console.log(
    `PASS: SPA registration, navigation, filters, mutations, failures, permissions UI, dialogs, 48 themed route layouts. Screenshots: ${snapshots}`,
  )
})()
  .catch((error) => {
    console.error(error)
    process.exitCode = 1
  })
  .finally(async () => {
    await browser?.close()
    server.close()
  })
