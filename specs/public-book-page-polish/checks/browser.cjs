// Render fixtures first using the command in verification.md.
const { createRequire } = require('module')
const fs = require('fs'),
  path = require('path'),
  http = require('http'),
  assert = require('assert/strict')
const root = path.resolve(__dirname, '../../..')
const { chromium } = createRequire(
  fs.realpathSync(path.join(root, 'node_modules/@playwright/mcp/package.json')),
)('playwright')
const fixtureDir = process.env.BOOK_PAGE_FIXTURE_DIR || '/tmp/book-page-fixtures'
let releaseFragment
const server = http.createServer(async (req, res) => {
  const url = new URL(req.url, 'http://localhost')
  if (url.pathname.endsWith('/__fragment/toc')) {
    await new Promise((resolve) => {
      releaseFragment = resolve
    })
    res.setHeader('content-type', 'text/html')
    res.end(
      fs.readFileSync(
        path.join(
          fixtureDir,
          (req.headers.referer || '').includes('fixture=populated')
            ? 'toc-populated.html'
            : 'toc.html',
        ),
      ),
    )
    return
  }
  if (url.pathname === '/') {
    res.setHeader('content-type', 'text/html')
    // The standard shell currently omits this mount when there are no flashes.
    // Supply it in fixtures so that unrelated common.js initialization succeeds.
    res.end(
      fs
        .readFileSync(
          path.join(fixtureDir, (url.searchParams.get('fixture') || 'empty') + '.html'),
          'utf8',
        )
        .replace('</body>', '<div id="client-flashes"></div></body>'),
    )
    return
  }
  const file = url.pathname.startsWith('/_/assets/')
    ? path.join(root, 'dist', url.pathname.slice(10))
    : url.pathname.startsWith('/_/embed-assets/')
      ? path.join(root, 'web/frontend/embed-assets', url.pathname.slice(16))
      : null
  if (!file || !fs.existsSync(file) || fs.statSync(file).isDirectory()) {
    res.writeHead(404)
    res.end()
    return
  }
  res.setHeader(
    'content-type',
    file.endsWith('.css')
      ? 'text/css'
      : file.endsWith('.js')
        ? 'text/javascript'
        : file.endsWith('.svg')
          ? 'image/svg+xml'
          : file.endsWith('.woff2')
            ? 'font/woff2'
            : 'application/octet-stream',
  )
  res.end(fs.readFileSync(file))
})
let browser
;(async () => {
  await new Promise((resolve) => server.listen(0, '127.0.0.1', resolve))
  browser = await chromium.launch({
    headless: true,
    executablePath: process.env.PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH || '/usr/bin/chromium-browser',
    args: ['--no-sandbox'],
  })
  const page = await browser.newPage()
  const errors = []
  page.on('pageerror', (e) => errors.push(e.message))
  const results = []
  for (const fixture of ['empty', 'long', 'populated', 'authenticated', 'warning'])
    for (const dark of [false, true])
      for (const width of [360, 768, 1280]) {
        releaseFragment = null
        await page.setViewportSize({ width, height: 900 })
        await page.goto(`http://127.0.0.1:${server.address().port}/?fixture=${fixture}`)
        await page.evaluate((dark) => document.documentElement.classList.toggle('dark', dark), dark)
        if (fixture === 'warning') {
          const warning = await page.locator('[x-ref=card]').boundingBox()
          assert(warning.x >= 0 && warning.x + warning.width <= width)
          await page.getByRole('button', { name: 'Proceed', exact: true }).click()
          await page
            .getByRole('button', { name: 'Proceed', exact: true })
            .waitFor({ state: 'hidden' })
        }
        await page.waitForFunction(
          () => window.Alpine && document.querySelector('#book')._x_dataStack,
        )
        await page.evaluate((dark) => document.documentElement.classList.toggle('dark', dark), dark)
        await page.waitForFunction(
          () =>
            document.querySelector('#slot-book-toc .Loader') &&
            getComputedStyle(document.querySelector('#slot-book-toc')).display !== 'none',
        )
        const tabs = page.getByRole('tab')
        const before = await tabs.first().evaluate((el) => ({
          width: el.offsetWidth,
          height: el.offsetHeight,
          top: el.getBoundingClientRect().top + scrollY,
        }))
        await tabs.first().focus()
        await page.keyboard.press('ArrowRight')
        await assert.equal(await tabs.nth(1).getAttribute('aria-selected'), 'true')
        await page.locator('#book-panel-reviews').waitFor({ state: 'visible' })
        await page.keyboard.press('Home')
        assert.equal(await tabs.first().getAttribute('aria-selected'), 'true')
        await page.locator('#slot-book-toc').waitFor({ state: 'visible' })
        const after = await tabs.first().evaluate((el) => ({
          width: el.offsetWidth,
          height: el.offsetHeight,
          top: el.getBoundingClientRect().top + scrollY,
        }))
        assert.equal(before.height, after.height)
        assert.equal(before.width, after.width)
        assert.equal(before.top, after.top)
        await page.keyboard.press('End')
        await page.keyboard.press('ArrowLeft')
        assert.equal(await tabs.first().getAttribute('aria-selected'), 'true')
        for (let i = 0; !releaseFragment && i < 100; i++)
          await new Promise((resolve) => setTimeout(resolve, 20))
        assert(releaseFragment, 'chapter request missing')
        releaseFragment()
        const chapterContent = page.locator(fixture === 'populated' ? '.BookToc' : '#empty_toc')
        await chapterContent.waitFor()
        assert.equal(await page.locator('#slot-book-toc').getAttribute('role'), 'tabpanel')
        assert.equal(
          await page.locator('#slot-book-toc').getAttribute('aria-labelledby'),
          'book-tab-toc',
        )
        await tabs.nth(1).click()
        await page.locator('#slot-book-toc').waitFor({ state: 'hidden' })
        await tabs.first().click()
        await chapterContent.waitFor({ state: 'visible' })
        if (fixture === 'populated') {
          await tabs.nth(1).click()
          await page.locator('.BookPage-review').waitFor({ state: 'visible' })
        }
        if (fixture === 'authenticated' || fixture === 'warning') {
          const primary = page.locator('#reading-list button').first()
          const menuButton = page.locator('#reading-list button').nth(1)
          for (const status of ['want_to_read', 'reading', 'paused', 'read', 'dnf']) {
            await page.evaluate((status) => {
              document.querySelector('#book')._x_dataStack[0].rl = {
                status,
                chapterId: '1',
                chapterName: 'A chapter',
              }
            }, status)
            await page.waitForFunction(
              (status) =>
                document
                  .querySelector('#reading-list [role=menuitemradio][aria-checked=true]')
                  ?.getAttribute('@click')
                  ?.includes(status),
              status,
            )
            assert.equal(await primary.isDisabled(), ['read', 'dnf'].includes(status))
            const labelBounds = await primary
              .locator('[x-text="primaryLabel()"]')
              .evaluate((el) => ({
                width: el.clientWidth,
                textWidth: el.scrollWidth,
                whiteSpace: getComputedStyle(el).whiteSpace,
              }))
            assert.equal(labelBounds.whiteSpace, 'nowrap')
            assert(
              labelBounds.textWidth <= labelBounds.width + 1,
              JSON.stringify({ status, width, ...labelBounds }),
            )
            await menuButton.click()
            await page.locator('#reading-list [role=menu]').waitFor({ state: 'visible' })
            await page.keyboard.press('Escape')
            await page.locator('#reading-list [role=menu]').waitFor({ state: 'hidden' })
          }
          await page.evaluate(() => {
            document.querySelector('#book')._x_dataStack[0].rl = null
          })
          await tabs.nth(1).click()
          await page.locator('#current-review-rating').waitFor({ state: 'visible' })
        }
        if (fixture === 'long') await page.locator('.Tags [role=button]').click()
        const metrics = await page.evaluate(() => {
          const book = document.querySelector('#book'),
            cover = document.querySelector('.BookPage-cover').getBoundingClientRect()
          const elements = [...book.querySelectorAll('*')]
            .filter(
              (el) =>
                getComputedStyle(el).display !== 'none' &&
                el.getBoundingClientRect().width &&
                el.getBoundingClientRect().right > innerWidth + 1,
            )
            .map((el) => el.className)
          const canvas = document.createElement('canvas'),
            ctx = canvas.getContext('2d')
          canvas.width = canvas.height = 1
          const rgb = (color) => {
            ctx.clearRect(0, 0, 1, 1)
            ctx.fillStyle = color
            ctx.fillRect(0, 0, 1, 1)
            return [...ctx.getImageData(0, 0, 1, 1).data]
          }
          const lum = (color) => {
            const c = color.slice(0, 3).map((v) => {
              v /= 255
              return v <= 0.04045 ? v / 12.92 : ((v + 0.055) / 1.055) ** 2.4
            })
            return c[0] * 0.2126 + c[1] * 0.7152 + c[2] * 0.0722
          }
          const contrast = {}
          for (const selector of [
            '.PageHeader-title',
            '.PageHeader-meta a',
            '.BookPage-metadata',
            '.age-rating',
            '.BookPage-stat',
            '.BookPage-tab[aria-selected="false"]',
            '#reading-list .Btn',
          ]) {
            const el = document.querySelector(selector)
            let bg = el
            while (bg.parentElement && rgb(getComputedStyle(bg).backgroundColor)[3] === 0)
              bg = bg.parentElement
            const a = lum(rgb(getComputedStyle(el).color)),
              b = lum(rgb(getComputedStyle(bg).backgroundColor))
            contrast[selector] = (Math.max(a, b) + 0.05) / (Math.min(a, b) + 0.05)
          }
          return {
            contrast,
            scroll: document.documentElement.scrollWidth,
            width: innerWidth,
            cover: [cover.width, cover.height],
            overflowing: elements,
          }
        })
        assert(metrics.scroll <= width, JSON.stringify({ fixture, dark, width, ...metrics }))
        assert.deepEqual(metrics.cover, [200, 300])
        for (const [selector, ratio] of Object.entries(metrics.contrast))
          assert(ratio >= 4.5, JSON.stringify({ fixture, dark, width, selector, ratio }))
        results.push({ fixture, dark, width, ...metrics })
        if (fixture === 'empty' && width !== 768)
          await page.screenshot({
            path: path.join(fixtureDir, `book-${dark ? 'dark' : 'light'}-${width}.png`),
            fullPage: true,
          })
      }
  assert.deepEqual(errors, [])
  console.log(JSON.stringify(results, null, 2))
})()
  .catch((e) => {
    console.error(e)
    process.exitCode = 1
  })
  .finally(async () => {
    releaseFragment?.()
    await browser?.close()
    server.closeAllConnections()
    server.close()
  })
