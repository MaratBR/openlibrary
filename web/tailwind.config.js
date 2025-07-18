const path = require('path');

/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'selector',
  content: [
    path.resolve(__dirname, 'olresponse') + '/*.templ',
    path.resolve(__dirname, 'public/templates') + '/*.templ',
    path.resolve(__dirname, 'admin/templates') + '/*.templ',
    path.resolve(__dirname, 'frontend', 'src') + '/**/*.{js,ts,jsx,tsx,svelte}',
  ],
  theme: {
      extend: {
        keyframes: {
          'collapsible-down': {
            from: { height: '0' },
            to: { height: 'var(--radix-collapsible-content-height)' },
          },
          'collapsible-up': {
            from: { height: 'var(--radix-collapsible-content-height)' },
            to: { height: '0' },
          },
        },
        animation: {
          'collapsible-down': 'collapsible-down 0.2s ease-out',
          'collapsible-up': 'collapsible-up 0.2s ease-out',
        },
        borderRadius: {
          lg: 'var(--radius)',
          md: 'calc(var(--radius) - 2px)',
          sm: 'calc(var(--radius) - 4px)',
        },
        fontFamily: {
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
  },
  plugins: [],
}