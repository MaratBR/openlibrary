const http = require('http'),
  fs = require('fs'),
  path = require('path')
const root = path.resolve(__dirname, '../../..')
const { execFileSync } = require('child_process')
const translations = execFileSync(
  'python3',
  [
    '-c',
    `
import json, sys, tomllib
result = {}
def flatten(data, prefix=''):
    for key, value in data.items():
        if isinstance(value, dict): flatten(value, prefix + key + '.')
        elif isinstance(value, str): result[prefix + key] = value
with open(sys.argv[1], 'rb') as source: flatten(tomllib.load(source))
print(json.dumps(result))
`,
    path.join(root, 'translations/en.toml'),
  ],
  { encoding: 'utf8' },
)
module.exports = http.createServer((req, res) => {
  if (req.url === '/' || req.url.startsWith('/admin')) {
    res.setHeader('content-type', 'text/html')
    res.end(
      `<!doctype html><html class="theme-default"><head><meta name="viewport" content="width=device-width, initial-scale=1"><link rel="stylesheet" href="/dist/common.css"><link rel="stylesheet" href="/dist/admin-common.css"><link rel="stylesheet" href="/_/embed-assets/fonts.css"><link rel="stylesheet" href="/_/embed-assets/css/fa/all.min.css"></head><body><div x-data="Island({name:'admin/App',data:null})"><span>Loading…</span></div><script>window.i18n=${translations};window.__server__={ageRatings:['G','PG','PG-13','R','NC-17','?']};window.getCookie=()=> 'fixture-csrf';window.OLTheme={toggle(){document.documentElement.classList.toggle('dark')}};window.toast=Object.assign(()=>{}, {error:()=>{}});window._=(key,args={})=>Object.entries(args).reduce((s,[k,v])=>s.replace('{{'+k+'}}',v).replace('{{.'+k+'}}',v),window.i18n[key]||key);</script><script type="module" src="/dist/admin-alpinejs.js"></script></body></html>`,
    )
    return
  }
  const urlPath = decodeURIComponent(req.url.split('?')[0])
  const file = urlPath.startsWith('/_/embed-assets/')
    ? path.join(root, 'web/frontend/embed-assets', urlPath.slice('/_/embed-assets/'.length))
    : path.join(root, urlPath)
  if (!file.startsWith(root + '/dist/') && !file.startsWith(root + '/web/frontend/embed-assets/')) {
    res.writeHead(404)
    res.end()
    return
  }
  try {
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
  } catch {
    res.writeHead(404)
    res.end()
  }
})
