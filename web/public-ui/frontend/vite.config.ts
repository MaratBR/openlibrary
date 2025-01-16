import { resolve } from 'node:path'
import { defineConfig, Plugin } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

type AutoInjectCSSAsLinkOptions = {
  baseUrl: string;
}

function autoInjectCSSAsLinkTagPlugin({
  baseUrl
}: AutoInjectCSSAsLinkOptions): Plugin {
  return {
    name: 'auto-inject-css-as-link-tag',
    enforce: 'post',
    apply: 'build',
  
    async generateBundle(_, bundle) {
      const chunksAssets: Record<string, string[]> = {};
  
      for (const [fileName, chunk] of Object.entries(bundle)) {
        if (chunk.type !== 'chunk' || !/\.js$/.test(fileName)) {
          continue;
        }
  
        const assets = new Set<string>();
  
        if (chunk.viteMetadata) {
          console.log(chunk.viteMetadata)
          for (const importedAsset of chunk.viteMetadata.importedCss) {
            assets.add(importedAsset);
          }
        }
  
        // // Recursively collect assets for a chunk
        // const collectAssets = (chunkName) => {
        //   const entry = bundle[chunkName];
        //   if (!entry) return;
  
        //   if (entry.type === 'asset' && entry) {
        //     assets.add(entry.fileName);
        //   } else if (entry.type === 'chunk') {
        //     entry.imports.forEach(collectAssets);
        //     entry.dynamicImports.forEach(collectAssets);
  
        //     // Add directly imported CSS/other assets
        //     for (const dep of entry.viteMetadata?.importedAssets || []) {
        //       assets.add(dep);
        //     }
        //   }
        // };
  
        // collectAssets(fileName);
  
        chunksAssets[fileName] = [...assets];
      }
  
      // Append asset data to each chunk
      for (const [fileName, chunk] of Object.entries(bundle)) {
        if (chunk.type !== 'chunk') continue;
  
        const assets = chunksAssets[fileName];
        if (!assets || assets.length === 0) {
          continue;
        }

        const funcName = `__cssInject${Math.random().toString(36).substring(2)}`
        const func = `function ${funcName}(url){${baseUrl ? `url=${JSON.stringify(baseUrl)}+url;` : ''}if(!Array.from(document.head.querySelectorAll('link[rel="stylesheet"]')).some(link=>link.getAttribute('href')===url)){const link=document.createElement("link");link.rel="stylesheet";link.href=url;document.head.appendChild(link);}}`
        const injectedCode = `(()=>{${func};${JSON.stringify([...assets])}.forEach(${funcName})})();\n` 
        chunk.code = injectedCode + chunk.code
      }
            
    }
  }
}

const ENTRIES = [
  'common',
  'alpinejs',
  'http-client'
]


export default defineConfig({

  plugins: [
    svelte({}),
    autoInjectCSSAsLinkTagPlugin({
      baseUrl: '/_/assets/'
    }),
  ],

  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
    }
  },

  build: {
    rollupOptions: {
      output: {
        chunkFileNames: 'chunks/[hash].js',
        // Put chunk styles at <output>/assets
        assetFileNames: (assetInfo) => {

          if (assetInfo.names.length === 1 && assetInfo.names[0].endsWith('.css') && ENTRIES.includes(assetInfo.names[0].substring(0, assetInfo.names[0].length - 4))) {
            return '[name][extname]'
          }

          return 'assets/[name]-[hash][extname]'
        },
        entryFileNames: '[name].js',
      }
    },
    cssCodeSplit: true,
    lib: {
      name: 'ol-public-ui',
      formats: ['es'],
      entry: Object.fromEntries(ENTRIES.map(entry => [
        entry,
        resolve(__dirname, 'src', entry, 'index.ts')
      ]))
    },
  }
})