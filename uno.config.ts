import { defineConfig, transformerDirectives, transformerVariantGroup } from 'unocss'
import presetWind4 from '@unocss/preset-wind4'

// const paths = [
//   `${path.resolve(__dirname, './internal/olhttp')}/*.templ`,
//   `${path.resolve(__dirname, './web/public/templates')}/*.templ`,
//   `${path.resolve(__dirname, './web/admin/templates')}/*.templ`,
//   // `${path.resolve(__dirname, './web/frontend', 'src')}/**/*.{js,ts,jsx,tsx}`,
// ]

// console.log('UnoCSS paths:', paths)

export default defineConfig({
  shortcuts: [
    // ...
  ],
  content: {
    filesystem: [
      'web/public/templates/*.templ',
      'web/admin/templates/*.templ',
      'internal/olhttp/*.templ',
      'web/frontend/src/**/*.{js,ts,jsx,tsx,css,scss}',
    ],

    pipeline: {
      include: [/\.([jt]sx|vine.ts|mdx?|html|templ)($|\?)/],
    },
  },

  extractors: [
    // {
    //   name: 'templ-extractor',
    //   order: -1,
    //   async extract(code) {
    //     console.log(code.id)
    //     return new Set()
    //   },
    // },
  ],

  // extractorDefault: false,

  theme: {
    animation: {
      keyframes: {
        'collapsible-down': 'collapsible-down 0.2s ease-out',
        'collapsible-up': 'collapsible-up 0.2s ease-out',
      },
    },
    radius: {
      lg: 'var(--radius)',
      md: 'calc(var(--radius) - 2px)',
      sm: 'calc(var(--radius) - 4px)',
    },
    font: {
      title: 'var(--font-title)',
      text: 'var(--font-text)',
      book: 'var(--font-book)',
    },
    colors: {
      background: 'hsl(var(--background))',
      foreground: 'hsl(var(--foreground))',
      card: {
        DEFAULT: 'hsl(var(--card))',
        foreground: 'hsl(var(--card-foreground))',
      },
      popover: {
        DEFAULT: 'hsl(var(--popover))',
        foreground: 'hsl(var(--popover-foreground))',
      },
      primary: {
        DEFAULT: 'hsl(var(--primary))',
        foreground: 'hsl(var(--primary-foreground))',
      },
      secondary: {
        DEFAULT: 'hsl(var(--secondary))',
        foreground: 'hsl(var(--secondary-foreground))',
      },
      muted: {
        DEFAULT: 'hsl(var(--muted))',
        foreground: 'hsl(var(--muted-foreground))',
      },
      'muted-2': {
        DEFAULT: 'hsl(var(--muted-2))',
        foreground: 'hsl(var(--muted-2-foreground))',
      },
      highlight: {
        DEFAULT: 'rgba(var(--highlight))',
        foreground: 'hsl(var(--foreground))',
      },
      accent: {
        DEFAULT: 'hsl(var(--accent))',
        foreground: 'hsl(var(--accent-foreground))',
      },
      destructive: {
        DEFAULT: 'hsl(var(--destructive))',
        foreground: 'hsl(var(--destructive-foreground))',
      },
      border: 'hsl(var(--border))',
      input: 'hsl(var(--input))',
      ring: 'hsl(var(--ring))',
      chart: {
        1: 'hsl(var(--chart-1))',
        2: 'hsl(var(--chart-2))',
        3: 'hsl(var(--chart-3))',
        4: 'hsl(var(--chart-4))',
        5: 'hsl(var(--chart-5))',
      },
    },
  },
  presets: [
    presetWind4(),
    // presetAttributify(),
    // presetIcons(),
    // presetTypography(),
    // presetWebFonts({
    //   fonts: {
    //     // ...
    //   },
    // }),
  ],
  transformers: [transformerDirectives(), transformerVariantGroup()],
})
