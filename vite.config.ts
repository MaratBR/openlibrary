/* eslint-disable no-undef */
import { resolve } from 'node:path'
import { transform } from 'esbuild'
import { defineConfig, Plugin } from 'vite'
import glob from 'fast-glob'
import UnoCSS from 'unocss/vite'
import react from '@vitejs/plugin-react'

const SOURCEMAP = true

function watchExternalPlugin(options: { paths?: string[] } = {}): Plugin {
  const { paths = [] } = options

  return {
    name: 'watch-external',

    buildStart() {
      for (const file of glob.sync(paths, { absolute: true })) {
        this.addWatchFile(file)
      }
    },
  }
}

function minifyLibraryPlugin(enabled: boolean): Plugin {
  return {
    name: 'minify-library',
    renderChunk: {
      order: 'post',
      async handler(code, chunk, outputOptions) {
        if (!enabled || outputOptions.format !== 'es') return null

        return transform(code, {
          format: 'esm',
          minify: true,
          sourcemap: SOURCEMAP,
          sourcefile: chunk.fileName,
        })
      },
    },
  }
}

const ENTRIES = [
  'common',
  'alpinejs',
  'block-editor',

  // admin stuff
  'admin-common',
  'admin-alpinejs',

  // moderation
  'mod',

  // global public API (remove?)
  'public.api',

  // specifically stuff for book-reader
  'book-reader',

  // bookmanager stuff

  // islands
  'islands/public', // all islands available in public pages
  'islands/signup',
  'islands/review-editor',
  'islands/admin-password-reset',

  'islands/bookmanager',

  'islands/admin',
]

const ENTRY_PATHS: Partial<Record<(typeof ENTRIES)[number], string>> = {
  mod: 'features/moderation/styles',
}

export default defineConfig((env) => ({
  cacheDir: '.vite',
  define: {
    'process.env.NODE_ENV': JSON.stringify(env.mode),
  },
  plugins: [
    watchExternalPlugin({
      paths: [
        'web/public/templates/*.templ',
        'web/admin/templates/*.templ',
        'internal/olhttp/*.templ',
        'internal/olhttp/webcomponents/*.templ',
      ],
    }),

    UnoCSS(),

    react({}),

    // Vite deliberately preserves whitespace in ES library builds, even with
    // build.minify enabled. Run a final minification pass for production mode.
    minifyLibraryPlugin(env.mode === 'production'),
  ],

  resolve: {
    alias: {
      '@': resolve(__dirname, './web/frontend/src'),
    },
  },

  build: {
    minify: env.mode === 'production' ? 'esbuild' : false,
    rollupOptions: {
      output: {
        chunkFileNames: 'chunks/[hash].js',
        // Put chunk styles at <output>/assets
        assetFileNames: (assetInfo) => {
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
        ENTRIES.map((entry) => [
          entry,
          resolve(__dirname, 'web/frontend/src', ENTRY_PATHS[entry] ?? entry, 'index.ts'),
        ]),
      ),
    },
    sourcemap: SOURCEMAP,
  },
}))
