const assert = require('node:assert/strict')
const fs = require('node:fs')
const path = require('node:path')
const { createRequire } = require('node:module')
const root = path.resolve(__dirname, '../../..')
const { chromium } = createRequire(
  fs.realpathSync(path.join(root, 'node_modules/@playwright/mcp/package.json')),
)('playwright')
const servers = [
  require('../../admin-spa/checks/server.cjs'),
  require('../../book-manager-refresh/checks/server.cjs'),
]
let browser
;(async () => {
  browser = await chromium.launch({
    executablePath: process.env.PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH,
    args: ['--no-sandbox'],
  })
  const page = await browser.newPage()
  await page.route('**/_api/**', (route) =>
    route.fulfill({ json: { books: [], page: 1, totalPages: 1 } }),
  )
  for (const [index, server] of servers.entries()) {
    await new Promise((resolve) => server.listen(0, '127.0.0.1', resolve))
    for (const width of [360, 760, 768, 1024, 1280, 2560]) {
      await page.setViewportSize({ width, height: 900 })
      await page.goto(
        `http://127.0.0.1:${server.address().port}/${index === 0 ? 'admin#/' : '#/books'}`,
      )
      await page.locator('.DashboardShell-nav').waitFor()
      for (const dark of [false, true]) {
        await page.evaluate((dark) => document.documentElement.classList.toggle('dark', dark), dark)
        const geometry = await page.evaluate(() => {
          const sidebar = document.querySelector('.DashboardShell-sidebar').getBoundingClientRect()
          const content = document.querySelector('.DashboardShell-content').getBoundingClientRect()
          return {
            sidebarRight: sidebar.right,
            sidebarBottom: sidebar.bottom,
            contentLeft: content.left,
            contentTop: content.top,
            overflow: document.documentElement.scrollWidth > innerWidth,
            position: getComputedStyle(document.querySelector('.DashboardShell-sidebar')).position,
          }
        })
        if (index === 0 && width === 2560) {
          const widths = await page.evaluate(() => ({
            frame: document.querySelector('.Admin-frame').getBoundingClientRect().width,
            main: document.querySelector('.Admin-main').getBoundingClientRect().width,
            footer: document.querySelector('.Admin-footer').getBoundingClientRect().width,
          }))
          assert.ok(widths.frame > 2000)
          assert.equal(widths.main, 1344)
          assert.equal(widths.footer, 1344)
        }
        assert.equal(geometry.overflow, false, `workspace ${index}, width ${width}, dark ${dark}`)
        if (width > 760) {
          assert.ok(geometry.contentLeft >= geometry.sidebarRight)
          assert.equal(geometry.position, 'sticky')
        } else {
          assert.ok(geometry.contentTop >= geometry.sidebarBottom)
        }
        assert.equal(await page.locator('.DashboardShell-nav [aria-current="page"]').count(), 1)
      }
    }
    if (index === 1) {
      await page.locator('.DashboardShell-nav a[href="#/books/new"]').click()
      await page.waitForURL('**/#/books/new')
      await page.locator('.DashboardShell-nav a[href="#/books/new"][aria-current="page"]').waitFor()
      assert.equal(await page.locator('.DashboardShell-nav [aria-current="page"]').count(), 1)
      assert.equal(
        await page.locator('.DashboardShell-nav [aria-current="page"]').getAttribute('href'),
        '#/books/new',
      )
    }
  }
  console.log('PASS: both sidebar layouts, 24 responsive/theme combinations, active navigation.')
})()
  .catch((error) => {
    console.error(error)
    process.exitCode = 1
  })
  .finally(async () => {
    if (browser) await browser.close()
    servers.forEach((server) => server.close())
  })
