import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './lib/dark-mode/dark-mode.ts'
import './index.css'
import { Settings } from 'luxon'
import './server-data.ts'
import App, { staticInitApp } from './App.tsx'

const root = document.getElementById('root')!

Settings.throwOnInvalid = true

declare module 'luxon' {
  interface TSSettings {
    throwOnInvalid: true
  }
}

if (__server__.iframeAllowed !== true && self !== top) {
  console.error('React app will not be initialized: iframe is not allowed')
} else {
  staticInitApp()
  createRoot(root).render(
    <StrictMode>
      <App />
    </StrictMode>,
  )
}
