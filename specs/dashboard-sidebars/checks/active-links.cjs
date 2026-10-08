const assert = require('node:assert/strict')
const fs = require('node:fs')
const path = require('node:path')
const { createRequire } = require('node:module')
const root = path.resolve(__dirname, '../../..')
const { chromium } = createRequire(
  fs.realpathSync(path.join(root, 'node_modules/@playwright/mcp/package.json')),
)('playwright')
const { build } = require('esbuild')
let browser
;(async () => {
  const bundle = await build({
    entryPoints: [path.join(root, 'web/frontend/src/common/links.ts')],
    bundle: true,
    write: false,
    format: 'iife',
    globalName: 'ActiveLinks',
    platform: 'browser',
  })
  browser = await chromium.launch({
    executablePath: process.env.PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH,
    args: ['--no-sandbox'],
  })
  const page = await browser.newPage()
  await page.route('http://fixture.test/**', (route) =>
    route.fulfill({
      contentType: 'text/html',
      body: '<a class="nav-link active" href="/">Home</a><a class="nav-link" href="/admin">Admin</a><a class="nav-link active" href="https://external.test/admin">External</a><a class="nav-link active" href="/admin#/users">Hash</a><a class="NavSubmenu-link" href="/admin">Submenu</a>',
    }),
  )
  const workspace =
    '<nav class="DashboardShell-nav"><a href="#/" class="active">Home</a><a href="#/users">Users</a><a href="#/tags">Tags</a></nav>'
  for (const islandFirst of [true, false]) {
    for (let reload = 0; reload < 5; reload++) {
      await page.goto('http://fixture.test/admin#/')
      await page.reload()
      await page.addScriptTag({ content: bundle.outputFiles[0].text })
      if (islandFirst)
        await page.evaluate(
          (html) => document.body.insertAdjacentHTML('beforeend', html),
          workspace,
        )
      await page.evaluate(() => ActiveLinks.initActiveLinks())
      if (!islandFirst)
        await page.evaluate(
          (html) => document.body.insertAdjacentHTML('beforeend', html),
          workspace,
        )
      await page.evaluate(() => ActiveLinks.initActiveLinks())
      assert.equal(await page.locator('.DashboardShell-nav .active').count(), 1)
      assert.equal(await page.locator('.nav-link.active').count(), 1)
      assert.equal(await page.locator('.nav-link.active').getAttribute('href'), '/admin')
      assert.equal(await page.locator('.NavSubmenu-link.active').count(), 1)
      // Simulate a router update, then rerun the common initializer.
      await page.evaluate(() => {
        const links = document.querySelectorAll('.DashboardShell-nav a')
        links[0].classList.remove('active')
        links[1].classList.add('active')
        history.replaceState(null, '', '#/users')
        ActiveLinks.initActiveLinks()
      })
      assert.equal(await page.locator('.DashboardShell-nav .active').count(), 1)
      assert.equal(
        await page.locator('.DashboardShell-nav .active').getAttribute('href'),
        '#/users',
      )
    }
  }
  console.log(
    'PASS: navbar active state, stale class removal, origin/hash checks, and workspace ownership in both initialization orders across repeated reloads.',
  )
})()
  .catch((error) => {
    console.error(error)
    process.exitCode = 1
  })
  .finally(async () => {
    if (browser) await browser.close()
  })
