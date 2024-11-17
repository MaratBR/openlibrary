import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import './lib/dark-mode/dark-mode.ts'

import './index.css'
import { Settings } from 'luxon'

const root = document.getElementById('root')!

Settings.throwOnInvalid = true

declare module 'luxon' {
  interface TSSettings {
    throwOnInvalid: true
  }
}

createRoot(root).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
