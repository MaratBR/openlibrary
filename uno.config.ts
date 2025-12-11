import { defineConfig, transformerDirectives, transformerVariantGroup } from 'unocss'
import presetWind4 from '@unocss/preset-wind4'

export default defineConfig({
  content: {
    filesystem: [
      'web/public/templates/*.templ',
      'web/admin/templates/*.templ',
      'internal/olhttp/*.templ',
      'internal/olhttp/webcomponents/*.templ',
      'web/frontend/src/**/*.{js,ts,jsx,tsx,css,scss}',
    ],

    pipeline: {
      include: [/\.([jt]sx|mdx?|html|templ)($|\?)/],
    },
  },

  safelist: ['bg-primary', 'bg-destructive', 'bg-background', 'bg-foreground'],

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
      background: 'var(--background)',
      foreground: 'var(--foreground)',
      card: {
        DEFAULT: 'var(--card)',
        foreground: 'var(--card-foreground)',
      },
      popover: {
        DEFAULT: 'var(--popover)',
        foreground: 'var(--popover-foreground)',
      },
      primary: {
        DEFAULT: 'var(--primary)',
        foreground: 'var(--primary-foreground)',
      },
      secondary: {
        DEFAULT: 'var(--secondary)',
        foreground: 'var(--secondary-foreground)',
      },
      muted: {
        DEFAULT: 'var(--muted)',
        foreground: 'var(--muted-foreground)',
      },
      'muted-2': {
        DEFAULT: 'var(--muted-2)',
        foreground: 'var(--muted-2-foreground)',
      },
      highlight: {
        DEFAULT: 'var(--highlight)',
        foreground: 'var(--foreground)',
      },
      accent: {
        DEFAULT: 'var(--accent)',
        foreground: 'var(--accent-foreground)',
      },
      destructive: {
        DEFAULT: 'var(--destructive)',
        foreground: 'var(--destructive-foreground)',
      },
      border: 'var(--border)',
      input: 'var(--input)',
      ring: 'var(--ring)',
      chart: {
        1: 'var(--chart-1)',
        2: 'var(--chart-2)',
        3: 'var(--chart-3)',
        4: 'var(--chart-4)',
        5: 'var(--chart-5)',
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
