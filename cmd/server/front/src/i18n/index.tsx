import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

import en from './translations/en.json'

const resources = {
  en: {
    translation: en,
  },
}

export function initI18n() {
  i18n
    .use(initReactI18next) // passes i18n down to react-i18next
    .init({
      resources,
      lng: 'en', // language to use, more information here: https://www.i18next.com/overview/configuration-options#languages-namespaces-resources
      interpolation: {
        escapeValue: false, // react already safes from xss
      },
    })
}
