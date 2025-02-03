import { resolve } from 'node:path'
import { defineConfig, Plugin } from 'vite'
import preact from '@preact/preset-vite'
import type { OutputAsset } from 'rollup'

type AutoInjectCSSAsLinkOptions = {
  baseUrl: string
}

function autoInjectCSSAsLinkTagPlugin({ baseUrl }: AutoInjectCSSAsLinkOptions): Plugin {
  return {
    name: 'auto-inject-css-as-link-tag',
    enforce: 'post',
    apply: 'build',

    async generateBundle(_, bundle) {
      const chunksAssets: Record<string, string[]> = {}

      for (const [fileName, chunk] of Object.entries(bundle)) {
        if (chunk.type !== 'chunk' || !/\.js$/.test(fileName)) {
          continue
        }

        const assets = new Set<string>()

        if (chunk.viteMetadata) {
          for (const importedAsset of chunk.viteMetadata.importedCss) {
            assets.add(importedAsset)
          }
        }

        chunksAssets[fileName] = [...assets]
      }

      // Append asset data to each chunk
      for (const [fileName, chunk] of Object.entries(bundle)) {
        if (chunk.type !== 'chunk') continue

        const assets = chunksAssets[fileName]
        if (!assets || assets.length === 0) {
          continue
        }

        const funcName = `__cssInject${Math.random().toString(36).substring(2)}`
        const func = `function ${funcName}(url){${baseUrl ? `url=${JSON.stringify(baseUrl)}+url;` : ''}if(!Array.from(document.head.querySelectorAll('link[rel="stylesheet"]')).some(link=>link.getAttribute('href')===url)){const link=document.createElement("link");link.rel="stylesheet";link.href=url;document.head.appendChild(link);}}`
        const injectedCode = `(()=>{${func};${JSON.stringify([...assets])}.forEach(${funcName})})();\n`
        chunk.code = injectedCode + chunk.code
      }

      const cssManifest: OutputAsset = {
        fileName: '__injectedCSS.js',
        needsCodeReference: false,
        name: '__injectedCSS.js',
        names: ['__injectedCSS.js'],
        originalFileName: null,
        originalFileNames: [],
        source: `window.__injectedCSS=${JSON.stringify(chunksAssets)}`,
        type: 'asset',
      }

      bundle['__injectedCSS.json'] = cssManifest
    },
  }
}

const ENTRIES = [
  'common',
  'alpinejs',
  'http-client',

  'admin-common',
  'admin-alpinejs',

  'public.api',

  'book-reader',

  'islands/book-card-preview',
  'islands/review-editor',
  'islands/search-filters',
  'islands/admin-password-reset',
]

export default defineConfig({
  plugins: [
    preact({
      devToolsEnabled: true,
      prefreshEnabled: true,
      babel: {
        plugins: [
          [
            '@babel/plugin-proposal-decorators',
            {
              decoratorsBeforeExport: true,
              version: '2023-05',
            },
          ],
          '@babel/plugin-transform-class-static-block',
          '@babel/plugin-transform-class-properties',
        ],
      },
    }),
    autoInjectCSSAsLinkTagPlugin({
      baseUrl: '/_/assets/',
    }),
  ],

  resolve: {
    alias: {
      // eslint-disable-next-line no-undef
      '@': resolve(__dirname, './src'),
    },
  },

  build: {
    rollupOptions: {
      output: {
        chunkFileNames: 'chunks/[hash].js',
        // Put chunk styles at <output>/assets
        assetFileNames: (assetInfo) => {
          console.log(assetInfo.names)
          if (
            assetInfo.names.length === 1 &&
            assetInfo.names[0].endsWith('.css') &&
            ENTRIES.includes(assetInfo.names[0].substring(0, assetInfo.names[0].length - 4))
          ) {
            return '[name][extname]'
          }

          return '[name]-[hash][extname]'
        },
        entryFileNames: '[name].js',
      },
    },
    cssCodeSplit: true,
    lib: {
      name: 'ol-public-ui',
      formats: ['es'],
      entry: Object.fromEntries(
        // eslint-disable-next-line no-undef
        ENTRIES.map((entry) => [entry, resolve(__dirname, 'src', entry, 'index.ts')]),
      ),
    },
  },
})
