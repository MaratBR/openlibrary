/* eslint-disable no-undef */
import { resolve } from 'node:path'
import { defineConfig, Plugin } from 'vite'
import glob from 'fast-glob'
import UnoCSS from 'unocss/vite'
import react from '@vitejs/plugin-react'

const SOURCEMAP = true


// function esbuildMinifyPlugin(): Plugin {
//   return {
//     name: 'esbuild-minify-post',
//     apply: 'build',
//     async closeBundle() {
//       // adjust this to match your outDir
//       const outDir = resolve(process.cwd(), 'dist')
//       // get all .js files from dist
//       const files = await glob('**/*.js', { cwd: outDir, absolute: true })

//       await Promise.all(
//         files.map(async (file) => {
//           const code = await readFile(file, 'utf8')
//           const result = await esbuild({
//             stdin: {
//               contents: code,
//               resolveDir: dirname(file),
//               sourcefile: file,
//               loader: 'js',
//             },
//             outfile: file,
//             write: true,
//             bundle: false,
//             minify: true,
//             sourcemap: SOURCEMAP,
//             allowOverwrite: true,
//           })

//           if (result.errors.length) {
//             console.error(`esbuild failed on ${file}`, result.errors)
//           }
//         }),
//       )
//     },
//   }
// }

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

const ENTRIES = [
  'common',
  'alpinejs',
  'http-client',
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

  'islands/mod',
]

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
  ],

  resolve: {
    alias: {
      '@': resolve(__dirname, './web/frontend/src'),
    },
  },


  build: {
    minify: false, // TODO toggle depending on env
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
        ENTRIES.map((entry) => [entry, resolve(__dirname, 'web/frontend/src', entry, 'index.ts')]),
      ),
    },
    sourcemap: SOURCEMAP,
  },
}))
